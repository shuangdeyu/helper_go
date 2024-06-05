package encrypt

import (
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"encoding/base64"
	"strings"
)

/**
 * 用golang重构的一款java加密的样式
 * eg: encryptPasswd, err := EncryptText(*ftp.Password, GetEncryptKey(ftp.Name), DefaultSalt, 19)
 * eg: encryptPasswd, _ = DecryptText(*ftp.Password, GetEncryptKey(ftp.Name), DefaultSalt, 19)
 */

var DefaultSalt = []byte{0xb8, 0x95, 0x34, 0xc6, 0xed, 0xf1, 0x57, 0xa2}

func isAscii(str string) bool {
	for _, char := range str {
		if char > 0x7f {
			return false
		}
	}
	return true
}

func GetEncryptKey(name string) string {
	if isAscii(name) {
		return name
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(name))
	return encoded
}

func GetEncodingPassword(strKey string) []rune {
	r := NewRandomJava(1997)
	fixKey := []rune{'k', '&', 'm', '1', 'n', '-', 'k', '7'}
	l := len(strKey) + len(fixKey)
	key := make([]rune, l)
	for i := 0; i < l; i++ {
		n, _ := r.nextInt(l)
		if n < len(fixKey) {
			key[i] = fixKey[n]
		} else {
			key[i] = rune(strKey[n-len(fixKey)])
		}
	}
	return key
}

func getDerivedKey(password string, salt []byte, count int) ([]byte, []byte) {
	key := md5.Sum([]byte(password + string(salt)))
	for i := 0; i < count-1; i++ {
		key = md5.Sum(key[:])
	}
	return key[:8], key[8:]
}

func EncryptText(text, keyword string, salt []byte, iterations int) (string, error) {
	passwd := GetEncodingPassword(keyword)

	padNum := byte(8 - len(text)%8)
	for i := byte(0); i < padNum; i++ {
		text += string(padNum)
	}

	dk, iv := getDerivedKey(string(passwd), salt, iterations)

	block, err := des.NewCipher(dk)
	if err != nil {
		return "", err
	}

	encrypter := cipher.NewCBCEncrypter(block, iv)
	encrypted := make([]byte, len(text))
	encrypter.CryptBlocks(encrypted, []byte(text))

	return base64.StdEncoding.EncodeToString(encrypted), nil
}

func DecryptText(text, keyword string, salt []byte, iterations int) (string, error) {
	passwd := GetEncodingPassword(keyword)

	msgBytes, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return "", err
	}

	dk, iv := getDerivedKey(string(passwd), salt, iterations)
	block, err := des.NewCipher(dk)
	if err != nil {
		return "", err
	}

	decrypter := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(msgBytes))
	decrypter.CryptBlocks(decrypted, msgBytes)

	decryptedString := strings.TrimRight(string(decrypted), "\x01\x02\x03\x04\x05\x06\x07\x08")
	return decryptedString, nil
}
