package pwdhelper

import (
	"crypto/sha512"
	"encoding/base64"
	"golang.org/x/crypto/pbkdf2"
	"helper_go/comhelper"
)

/**
 * 一种密码加密方式
 * mode只决定password和salt的组合方式
 */
func PasswordHash(password string, salt string, mode int) string {
	toHAshString := ""
	if mode == 0 {
		toHAshString = salt + comhelper.Md5(password+salt) + salt
	} else if mode == 1 {
		toHAshString = salt + comhelper.Md5(salt+SubStr(password, 1, 10)) + password
		salt = salt + SubStr(password, 1, 5)
	} else {
		toHAshString = salt + comhelper.Md5(password+salt) + salt
		salt = SubStr(salt, 1, 10) + SubStr(password, 4, 4) + salt
	}
	pwd := []byte(toHAshString)
	saltByte := []byte(salt)
	dk := pbkdf2.Key(pwd, saltByte, 15000, 128, sha512.New)
	return base64.StdEncoding.EncodeToString(dk)
}
func SubStr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0
	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length
	if start > end {
		start, end = end, start
	}
	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}
	return string(rs[start:end])
}
