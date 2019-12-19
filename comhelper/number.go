package comhelper

import (
	"fmt"
	"strconv"
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
