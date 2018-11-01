package timehelper

import (
	"helper_go/comhelper"
	"strconv"
	"time"
)

/**
 * 获取当前时间
 */
func CurrentTime() string {
	now := time.Now().Format("2006-01-02 15:04:05")
	return now
}

/**
 * 获取当前时间戳
 */
func CurrentTimeUnix() string {
	now := time.Now().Unix()
	str := strconv.FormatInt(now, 10)
	return str
}

/**
 * 获取当天零点时间
 */
func CurrentTimeZero() string {
	todayDateStr := time.Now().Format("2006-01-02")
	t, _ := time.ParseInLocation("2006-01-02", todayDateStr, time.Local)
	return t.Format("2006-01-02 15:04:05")
}

/**
 * 获取当天末点时间
 */
func CurrentTimeEnd() string {
	todayDateStr := time.Now().Format("2006-01-02")
	return todayDateStr + " 23:59:59"
}

/**
 * 时间转换成时间戳
 */
func TimeToUnix(t string) string {
	loc, _ := time.LoadLocation("Local")                              // 获取时区
	theTime, _ := time.ParseInLocation("2006-01-02 15:04:05", t, loc) // 使用模板在对应时区转化为time.time类型
	sr := theTime.Unix()                                              // 转化为时间戳
	str := strconv.FormatInt(sr, 10)
	return str
}

/**
 * 时间戳转换成时间
 */
func UnixToTime(u string) string {
	t := time.Unix(comhelper.StringToInt64(u), 0).Format("2006-01-02 15:04:05")
	return t
}

/**
 * 时间转换成日期
 */
func TimeToDate(t string) string {
	tu := TimeToUnix(t)
	d := time.Unix(comhelper.StringToInt64(tu), 0).Format("2006-01-02")
	return d
}

/**
 * 获取时间所对应的天
 */
func TimeToDay(t string) int {
	tu := TimeToUnix(t)
	day := time.Unix(comhelper.StringToInt64(tu), 0).Day()
	return day
}

/**
 * 获取两时间相差天数
 */
func GetDayDiffer(end_time, start_time string) int64 {
	var day int64
	t1, _ := time.ParseInLocation("2006-01-02 15:04:05", start_time, time.Local)
	t2, _ := time.ParseInLocation("2006-01-02 15:04:05", end_time, time.Local)
	diff := t2.Unix() - t1.Unix() // 时间戳差
	day = diff / 86400
	return day
}

/**
 * 比较两个时间的大小
 */
func CompareTwoTime(t1, t2 string) bool {
	tn1 := comhelper.StringToInt64(TimeToUnix(t1))
	tn2 := comhelper.StringToInt64(TimeToUnix(t2))
	if tn1 > tn2 {
		return true
	} else {
		return false
	}
}

/**
 * 根据时间获取其零点值
 */
func TimeToTimeZero(t string) string {
	// 先转换成时间戳
	unix := TimeToUnix(t)
	zt := time.Unix(comhelper.StringToInt64(unix), 0).Format("2006-01-02")
	return zt + " 00:00:00"
}
