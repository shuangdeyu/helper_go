package leetcode

/*
69. x的平方根

给你一个非负整数 x ，计算并返回 x 的 算术平方根 。
由于返回类型是整数，结果只保留 整数部分 ，小数部分将被 舍去 。
注意：不允许使用任何内置指数函数和算符，例如 pow(x, 0.5) 或者 x ** 0.5 。

示例 1：
输入：x = 4
输出：2

示例 2：
输入：x = 8
输出：2
解释：8 的算术平方根是 2.82842..., 由于返回类型是整数，小数部分将被舍去。

提示：
0 <= x <= 231 - 1
*/

func mySqrt(x int) int {
	/*if x == 0 {
		return 0
	}
	if x == 1 {
		return 1
	}
	if x == 2 {
		return 1
	}
	for i := 1; i < x; i++ {
		if i*i > x {
			return i - 1
		}
	}
	return x*/

	// 二分法
	/*if x == 0 || x == 1 {
		return x
	}
	l := 0
	r := x
	for l <= r {
		mid := (l + r) / 2
		if mid*mid < x {
			l = mid + 1
		} else {
			r = mid - 1
		}
	}
	if l*l > x {
		return r
	} else {
		return l
	}*/

	// 牛顿迭代法 求根号 x，即求满足 x^2 - n = 0 方程的所有解
	r := x
	for r*r > x {
		r = (r + x/r) / 2
	}
	return r
}
