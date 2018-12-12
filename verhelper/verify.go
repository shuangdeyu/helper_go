package verhelper

import (
	"fmt"
	"regexp"
)

/**
 * 获取变量类型
 */
func Typeof(v interface{}) string {
	return fmt.Sprintf("%T", v)
}

/**
 * 校验是否为手机号
 */
func IsMobile(mobile string) bool {
	if m, _ := regexp.MatchString(`^1[3,4,5,6,7,8,9]\d{9}$`, mobile); !m {
		return false
	}
	return true
}

/**
 * 邮箱格式是否正确
 */
func IsEmail(email string) bool {
	if m, _ := regexp.MatchString(`^([\w\.\_]{2,10})@(\w{1,}).([a-z]{2,4})$`, email); !m {
		return false
	}
	return true
}

/**
 * 判断字符串是否为整数
 */
func IsOnlyInteger(i string) bool {
	if m, _ := regexp.MatchString(`^[0-9]\d*$`, i); !m {
		return false
	}
	return true
}

/**
 * 判断字符串是否为数学类型
 */
func IsOnlyNum(num string) bool {
	if m, _ := regexp.MatchString(`^(\-|\+)?\d+(\.\d+)?$`, num); !m {
		return false
	}
	return true
}

/**
 * 是否包含字母以及数字
 */
func IsCharOrNum(str string) bool {
	h_char := false
	h_num := false
	if m, _ := regexp.MatchString(`.*[a-zA-Z]+.*`, str); m {
		h_char = true
	}
	if m, _ := regexp.MatchString(`.*[0-9]+.*`, str); m {
		h_num = true
	}
	return h_num && h_char
}

/**
 * 密码验证
 * 必须包含字母以及数字，8~20位，
 * 允许包含 ~!@#$%^&*()_+-=*-+.,'; 特殊字符
 */
func PasswordReg(pwd string) bool {
	m1, _ := regexp.MatchString(`^[\w~!@#$%^&*()_+-=*-+.,';]{8,20}$`, pwd)
	m2, _ := regexp.MatchString(`^[0-9]{1,}$`, pwd)
	m3, _ := regexp.MatchString(`^[a-zA-Z]{1,}$`, pwd)
	m4, _ := regexp.MatchString(`^[~!@#$%^&*()_+-=*-+.,';]{1,}$`, pwd)
	m5 := IsCharOrNum(pwd)
	if !m1 || m2 || m3 || m4 || !m5 {
		return false
	}
	return true
}

/**
 * 隐藏身份证号码
 */
func HideCard(card string) string {
	ret := ""
	if card != "" {
		len := len(card)
		if len < 6 {
			ret = "******"
		} else {
			ret = card[0:6] + "******"
			ret += card[(len - 4):len]
		}
	}
	return ret
}

/**
 * 隐藏电话号码
 */
func HideMobile(mobile string) string {
	len := len(mobile)
	if len < 3 {
		return "***"
	}
	ret := mobile[0:3] + "****"
	if len > 7 {
		ret += mobile[len-4 : len]
	}
	return ret
}

/**
 * 隐藏真实姓名
 */
func HideTrueName(true_name string) string {
	name := ""
	if true_name != "" {
		rs_name := []rune(true_name)
		len := len(rs_name)
		if len > 2 {
			name = string(rs_name[0:1])
			for i := 0; i < len-2; i++ {
				name += "*"
			}
			name += string(rs_name[len-1 : len])
		} else {
			name = string(rs_name[0:1]) + "*"
		}
	}
	return name
}
