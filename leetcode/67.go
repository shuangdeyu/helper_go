package leetcode

import "strconv"

/*
67.二进制求和
给你两个二进制字符串 a 和 b ，以二进制字符串的形式返回它们的和。

示例 1：
输入:a = "11", b = "1"
输出："100"

示例 2：
输入：a = "1010", b = "1011"
输出："10101"

提示：
1 <= a.length, b.length <= 104
a 和 b 仅由字符 '0' 或 '1' 组成
字符串如果不是 "0" ，就不含前导零
*/
func addBinary(a string, b string) string {
	/*ret := ""
	num := 0
	lena := len(a)
	lenb := len(b)
	if lena == 0 && lenb == 0 {
		return ""
	} else if lena == 0 {
		return b
	} else if lenb == 0 {
		return a
	}

	for {
		i := ""
		if lena > 0 {
			i = a[lena-1 : lena]
			lena--
		}
		j := ""
		if lenb > 0 {
			j = b[lenb-1 : lenb]
			lenb--
		}
		if i == "" && j == "" && num == 0 {
			break
		}

		if i == "1" && j == "1" {
			if num == 1 {
				ret = "1" + ret
			} else {
				ret = "0" + ret
			}
			num = 1
		} else if i == "1" || j == "1" {
			if num == 1 {
				ret = "0" + ret
				num = 1
			} else {
				ret = "1" + ret
				num = 0
			}
		} else {
			if num == 1 {
				ret = "1" + ret
			} else {
				ret = "0" + ret
			}
			num = 0
		}
	}
	return ret*/

	ans := ""
	carry := 0
	lenA, lenB := len(a), len(b)
	n := max(lenA, lenB)
	for i := 0; i < n; i++ {
		if i < lenA {
			carry += int(a[lenA-i-1] - '0')
		}
		if i < lenB {
			carry += int(b[lenB-i-1] - '0')
		}
		ans = strconv.Itoa(carry%2) + ans
		carry /= 2
	}
	if carry > 0 {
		ans = "1" + ans
	}
	return ans
}
