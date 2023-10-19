package comhelper

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

/**
 * 数字相关，类型转换
 */

// 十六进制字符串转换成十进制整数
func HexToDec(param string) int64 {
	dec, _ := strconv.ParseInt(param, 16, 64)
	return dec
}

// 十六进制字符串转换成二进制字符串
func HexToBinStr(param, bit_size string) string {
	dec, _ := strconv.ParseInt(param, 16, 64)
	bin := fmt.Sprintf("%0"+bit_size+"b", dec)
	return bin
}

// 十六进制字符串转换成二进制数组
func HexToBinArr(param string, bit_size int) []int64 {
	dec, _ := strconv.ParseInt(param, 16, 64)
	bits := []int64{}
	for i := 0; i < bit_size; i++ {
		bits = append([]int64{dec & 0x1}, bits...)
		dec = dec >> 1
	}
	return bits
}

// 二进制加法
func AddBinary(a string, b string) string {
	var carry, sum int
	i, j := len(a), len(b)
	if i < j {
		i, j = j, i
		a, b = b, a
	}

	res := make([]byte, i+1)

	for j > 0 {
		j--
		i--
		sum = int(a[i]-'0') + int(b[j]-'0') + carry
		carry = sum / 2
		sum = sum % 2
		res[i+1] = byte(sum + '0')
	}

	for i > 0 {
		i--
		sum = int(a[i]-'0') + carry
		carry = sum / 2
		sum = sum % 2
		res[i+1] = byte(sum + '0')
	}

	res[0] = byte(carry + '0')

	for i < len(res)-1 {
		if res[i] == '0' {
			i++
		} else {
			break
		}
	}
	return string(res[i:])
}

// 二进制取反
func TransBinary(s string) string {
	bits := ""
	for i := 0; i < len(s); i++ {
		bit := StringToInt(s[i:i+1]) ^ 1
		bits += IntToString(bit)
	}
	return bits
}

// 十进制转换成2,8,16进制
func DecConvertToX(n, num int) (string, error) {
	if n < 0 {
		return strconv.Itoa(n), errors.New("只支持正整数")
	}
	if num != 2 && num != 8 && num != 16 {
		return strconv.Itoa(n), errors.New("只支持二、八、十六进制的转换")
	}
	result := ""
	h := map[int]string{
		0:  "0",
		1:  "1",
		2:  "2",
		3:  "3",
		4:  "4",
		5:  "5",
		6:  "6",
		7:  "7",
		8:  "8",
		9:  "9",
		10: "A",
		11: "B",
		12: "C",
		13: "D",
		14: "E",
		15: "F",
	}
	for ; n > 0; n /= num {
		lsb := h[n%num]
		result = lsb + result
	}
	return result, nil
}

// 16进制转换成字符串
func HexToString(hex string) string {
	var dst []byte
	fmt.Sscanf(hex, "%X", &dst)
	return string(dst)
}

// 字符串转换成16进制
func StringToHex(str string) string {
	d := fmt.Sprintf("%X", str)
	return d
}

// 16进制位数补0
func HexFillZero(data string, need_len int) string {
	if l := len(data); l < need_len {
		b := ""
		for x := need_len - l; x > 0; x-- {
			b += "0"
		}
		data = b + data
	}
	return data
}

// 获取随机数，可做验证码
func GetRandNumString(length int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())
	var sb strings.Builder
	for i := 0; i < length; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}
