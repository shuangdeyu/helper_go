package pwdhelper

import (
	"errors"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

/**
 * iban 编码解码，用于 ETH
 */

var (
	Base36Chars        = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	errICAPLength      = errors.New("invalid ICAP length")
	errICAPEncoding    = errors.New("invalid ICAP encoding")
	errICAPChecksum    = errors.New("invalid ICAP checksum")
	errICAPCountryCode = errors.New("invalid ICAP country code")
	errICAPAssetIdent  = errors.New("invalid ICAP asset identifier")
	errICAPInstCode    = errors.New("invalid ICAP institution code")
	errICAPClientIdent = errors.New("invalid ICAP client identifier")
)

var (
	Big1  = big.NewInt(1)
	Big0  = big.NewInt(0)
	Big36 = big.NewInt(36)
	Big97 = big.NewInt(97)
	Big98 = big.NewInt(98)
)

//export ConvertICAPToAddress
func ConvertICAPToAddress(s string) (common.Address, error) {
	switch len(s) {
	case 35: // "XE" + 2 digit checksum + 31 base-36 chars of address
		return parseICAP(s)
	case 34: // "XE" + 2 digit checksum + 30 base-36 chars of address
		return parseICAP(s)
	case 20: // "XE" + 2 digit checksum + 3-char asset identifier +
		// 4-char institution identifier + 9-char institution client identifier
		return parseIndirectICAP(s)
	default:
		return common.Address{}, errICAPLength
	}
}

func parseICAP(s string) (common.Address, error) {
	if !strings.HasPrefix(s, "XE") {
		return common.Address{}, errICAPCountryCode
	}
	if err := IbanvalidCheckSum(s); err != nil {
		return common.Address{}, err
	}
	// checksum is ISO13616, Ethereum address is base-36
	bigAddr, _ := new(big.Int).SetString(s[4:], 36)
	return common.BigToAddress(bigAddr), nil
}

func parseIndirectICAP(s string) (common.Address, error) {
	if !strings.HasPrefix(s, "XE") {
		return common.Address{}, errICAPCountryCode
	}
	if s[4:7] != "ETH" {
		return common.Address{}, errICAPAssetIdent
	}
	if err := IbanvalidCheckSum(s); err != nil {
		return common.Address{}, err
	}
	return common.Address{}, errors.New("not implemented")
}

//export ConvertAddressToICAP
func ConvertAddressToICAP(a common.Address) (string, error) {
	enc := base36Encode(a.Big())
	// zero padd encoded address to Direct ICAP length if needed
	if len(enc) < 30 {
		enc = join(strings.Repeat("0", 30-len(enc)), enc)
	}
	icap := join("XE", checkDigits(enc), enc)
	return icap, nil
}

// https://en.wikipedia.org/wiki/International_Bank_Account_Number#Validating_the_IBAN
func IbanvalidCheckSum(s string) error {
	s = join(s[4:], s[:4])
	expanded, err := iso13616Expand(s)
	if err != nil {
		return err
	}
	checkSumNum, _ := new(big.Int).SetString(expanded, 10)
	if checkSumNum.Mod(checkSumNum, Big97).Cmp(Big1) != 0 {
		return errICAPChecksum
	}
	return nil
}

func checkDigits(s string) string {
	expanded, _ := iso13616Expand(strings.Join([]string{s, "XE00"}, ""))
	num, _ := new(big.Int).SetString(expanded, 10)
	num.Sub(Big98, num.Mod(num, Big97))

	checkDigits := num.String()
	// zero padd checksum
	if len(checkDigits) == 1 {
		checkDigits = join("0", checkDigits)
	}
	return checkDigits
}

// not base-36, but expansion to decimal literal: A = 10, B = 11, ... Z = 35
func iso13616Expand(s string) (string, error) {
	var parts []string
	if !validBase36(s) {
		return "", errICAPEncoding
	}
	for _, c := range s {
		i := uint64(c)
		if i >= 65 {
			parts = append(parts, strconv.FormatUint(uint64(c)-55, 10))
		} else {
			parts = append(parts, string(c))
		}
	}
	return join(parts...), nil
}

func base36Encode(i *big.Int) string {
	var chars []rune
	x := new(big.Int)
	for {
		x.Mod(i, Big36)
		chars = append(chars, rune(Base36Chars[x.Uint64()]))
		i.Div(i, Big36)
		if i.Cmp(Big0) == 0 {
			break
		}
	}
	// reverse slice
	for i, j := 0, len(chars)-1; i < j; i, j = i+1, j-1 {
		chars[i], chars[j] = chars[j], chars[i]
	}
	return string(chars)
}

func validBase36(s string) bool {
	for _, c := range s {
		i := uint64(c)
		// 0-9 or A-Z
		if i < 48 || (i > 57 && i < 65) || i > 90 {
			return false
		}
	}
	return true
}

func join(s ...string) string {
	return strings.Join(s, "")
}
