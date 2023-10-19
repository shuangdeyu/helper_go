package timehelper

import (
	"fmt"
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
 * 获取当前毫秒级时间戳
 */
func CurrentTimeUnixMilli() string {
	now := time.Now().UnixNano() / 1e6
	str := strconv.FormatInt(now, 10)
	return str
}

/**
 * 获取当前纳秒级时间戳
 */
func CurrentTimeUnixNano() string {
	now := time.Now().UnixNano()
	str := strconv.FormatInt(now, 10)
	return str
}

/**
 * 获取当天零点时间
 */
func CurrentTimeZero() string {
	todayDateStr := time.Now().In(time.Local).Format("2006-01-02")
	return todayDateStr + " 00:00:00"
	/*t, _ := time.Parse("2006-01-02", todayDateStr)
	return t.Format("2006-01-02 15:04:05")*/
}

/**
 * 获取当天末点时间
 */
func CurrentTimeEnd() string {
	todayDateStr := time.Now().In(time.Local).Format("2006-01-02")
	return todayDateStr + " 23:59:59"
}

/**
 * 获取当前日期
 */
func CurrentDay() string {
	now := time.Now().Format("2006-01-02")
	return now
}

/**
 * 获取本月的第一天日期
 */
func MonthFirstDay() string {
	now := time.Now()
	year, month, _ := now.Date()
	location := now.Location()
	day := time.Date(year, month, 1, 0, 0, 0, 0, location).Format("2006-01-02")
	return day
}

/**
 * 获取本月的最后一天日期
 */
func MonthEndDay() string {
	now := time.Now()
	year, month, _ := now.Date()
	location := now.Location()
	day := time.Date(year, month, 1, 0, 0, 0, 0, location)
	last_day := day.AddDate(0, 1, -1).Format("2006-01-02")
	return last_day
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
 * 获取两时间相差秒
 */
func GetDayDifferBySecond(end_time, start_time string) int64 {
	//var day int64
	t1, _ := time.ParseInLocation("2006-01-02 15:04:05", start_time, time.Local)
	t2, _ := time.ParseInLocation("2006-01-02 15:04:05", end_time, time.Local)
	diff := t2.Unix() - t1.Unix() // 时间戳差
	//day = diff / 60
	return diff
}

/**
 * 获取两时间相差分钟
 */
func GetDayDifferByMin(end_time, start_time string) int64 {
	var day int64
	t1, _ := time.ParseInLocation("2006-01-02 15:04:05", start_time, time.Local)
	t2, _ := time.ParseInLocation("2006-01-02 15:04:05", end_time, time.Local)
	diff := t2.Unix() - t1.Unix() // 时间戳差
	day = diff / 60
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
func TimeToTimeEnd(t string) string {
	// 先转换成时间戳
	unix := TimeToUnix(t)
	zt := time.Unix(comhelper.StringToInt64(unix), 0).Format("2006-01-02")
	return zt + " 23:59:59"
}

/**
 * 获取年月日时字符串
 */
func YearMonthDayHourStr() string {
	year, month, day := time.Now().Date()
	return strconv.Itoa(year) + time_zero(strconv.Itoa(int(month))) + time_zero(strconv.Itoa(day)) + time_zero(strconv.Itoa(time.Now().Hour()))
}

/**
 * 获取当前时间字符串
 */
func CurrentTimeStr() string {
	year, month, day := time.Now().Date()
	now := time.Now()
	data := strconv.Itoa(year) + time_zero(strconv.Itoa(int(month))) + time_zero(strconv.Itoa(day))
	data += time_zero(strconv.Itoa(now.Hour())) + time_zero(strconv.Itoa(now.Minute())) + time_zero(strconv.Itoa(now.Second()))
	return data
}

func time_zero(t string) string {
	if len(t) < 2 {
		return "0" + t
	} else {
		return t
	}
}

/**
 * 获取UTC时间
 */
func CurrentUtcTime() string {
	now := time.Now()
	year, mon, day := now.UTC().Date()
	hour, min, sec := now.UTC().Clock()
	t := fmt.Sprintf("%04d%02d%02d%02d%02d%02d", year, mon, day, hour, min, sec)
	return t
}

/**
 * 给定日期加指定天数
 */
func TimeAddDay(times string, day int) string {
	// 字符串时间转换成时间戳
	unix := comhelper.StringToInt64(TimeToUnix(times))
	// 加上指定天数
	ts := time.Unix(unix, 0)
	tst := ts.AddDate(0, 0, day)
	// 转换成字符串
	t := tst.Format("2006-01-02 15:04:05")
	return t
}

/**
 * 给定日期加指定时间
 */
func TimeAdd(times string, add_time time.Duration) string {
	// 字符串时间转换成时间戳
	unix := comhelper.StringToInt64(TimeToUnix(times))
	// 加上指定时间
	ts := time.Unix(unix, 0)
	tst := ts.Add(add_time)
	// 转换成字符串
	t := tst.Format("2006-01-02 15:04:05")
	return t
}

/**
 * 是否在工作时间
 */
func InWorkTime() bool {
	t := time.Now()
	wd := t.Weekday()
	if wd == time.Sunday || wd == time.Saturday {
		return false
	}
	hour, _, _ := t.Clock()
	if hour < 9 || hour > 18 {
		return false
	}
	return true
}

/**
 * 获取N天前的时间
 */
func DayAgoTime(day int64) string {
	tim := time.Now().Unix() - int64(86400)*day
	return UnixToTime(comhelper.Int64ToString(tim))
}
