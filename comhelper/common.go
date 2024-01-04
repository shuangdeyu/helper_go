package comhelper

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/axgle/mahonia"
	"github.com/vmihailenco/msgpack"
	"github.com/yvasiyarov/php_session_decoder/php_serialize"
	"golang.org/x/crypto/bcrypt"
)

/**
 * 打印数据
 */
func Dump(data interface{}) {
	fmt.Println(JsonEncode(data))
}

/**
 * 取值，含默认值
 */
func DefaultParam(param string, def_param string) string {
	if param == "" {
		return strings.TrimSpace(def_param)
	}
	return strings.TrimSpace(param)
}

/**
 * ********************************************* 编码相关 **********************************************
 *
 */

/**
 * 转换成json
 */
func JsonEncode(v interface{}) string {
	ret, err := json.Marshal(v)
	if err != nil {
		return ""
	} else {
		return string(ret)
	}
}

/**
 * Msgpack序列化
 */
func MsgpackEncode(v interface{}) ([]byte, error) {
	ret, err := msgpack.Marshal(v)
	return ret, err
}

/**
 * md5 加密
 */
func Md5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has) // 将[]byte转成16进制
	return md5str
}

/**
 * bcrypt加密，可替代md5 ，更可靠
 */
func BcryptEncode(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	encodePW := string(hash) // 保存在数据库的密码，虽然每次生成都不同，只需保存一份即可
	return encodePW
}

/**
 * bcrypt 验证
 */
func BcryptDecode(password, enpassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(enpassword), []byte(password))
	if err != nil {
		return false
	} else {
		return true
	}
}

/**
 * php-UnSerialize 反序列化
 */
func UnSerialize(str string) (interface{}, error) {
	data, err := php_serialize.NewUnSerializer(str).Decode()
	return data, err
}

/**
 * ********************************************* 数组相关 **********************************************
 */

/**
 * 将数组转换到url形式字符串
 * @param remove string 需要去除的字符串
 */
func Array2UrlString(arr map[string]string, arr_key []string, remove string) string {
	content := ""
	if len(arr) > 0 && len(arr_key) > 0 {
		for _, v := range arr_key {
			if v != remove {
				content += v + "=" + arr[v] + "&"
			}
		}
		content = content[0 : len(content)-1]
	}
	content = strings.TrimSpace(content)
	return content
}

/**
 * 将数组转换到url形式字符串
 * @param remove string 需要去除的字符串
 */
func Array2UrlStringByEncode(arr map[string]string, arr_key []string, remove string) string {
	content := ""
	if len(arr) > 0 && len(arr_key) > 0 {
		for _, v := range arr_key {
			if v != remove {
				content += v + "=" + arr[v] + "&"
			}
		}
		content = content[0 : len(content)-1]
	}
	content = strings.TrimSpace(content)
	content = url.QueryEscape(content)
	return content
}

/**
 * 合并两个map数组
 */
func MergeMap(m1 map[string]interface{}, m2 map[string]interface{}) map[string]interface{} {
	for k, v := range m2 {
		m1[k] = v
	}
	return m1
}

/**
 * 合并两个数组 []string
 */
func MergeString(s ...[]string) []string {
	switch len(s) {
	case 0:
		return []string{}
	case 1:
		return s[0]
	default:
		var str []string
		str = append(s[0], s[1]...)
		return str
	}
}

/**
 * 判断数组是否包含某个元素
 */
func InArray(arr interface{}, val interface{}) bool {
	targetValue := reflect.ValueOf(arr)
	switch reflect.TypeOf(arr).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == val {
				return true
			}
		}
	case reflect.Map:
		// 匹配的是键值，不是值
		if targetValue.MapIndex(reflect.ValueOf(val)).IsValid() {
			return true
		}
	}
	return false
}

/**
 * 判断数组是否包含某个元素的部分
 */
func InArrayContains(arr []string, str string) bool {
	for _, v := range arr {
		if strings.Contains(str, v) {
			return true
		}
	}
	return false
}

/**
 * 字符串数组去重
 */
func DistinctArrString(arr []string) []string {
	tmpMap := make(map[string]interface{})
	ret := []string{}
	for _, val := range arr {
		if _, ok := tmpMap[val]; !ok && len(strings.TrimSpace(val)) > 0 {
			ret = append(ret, val)
			tmpMap[val] = struct{}{}
		}
	}
	return ret
}

/**
 * map[string]string to map[string]interface{}
 */
func MapStringToInterface(param map[string]string) (data map[string]interface{}) {
	if len(param) > 0 {
		ret := map[string]interface{}{}
		for k, v := range param {
			ret[k] = v
		}
		data = ret
	}
	return
}

/**
 * ********************************************* 字符串相关 **********************************************
 */

/**
 * 去除字符串中的html标签
 */
func TrimHtml(src string) string {
	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)
	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")
	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")
	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "\n")
	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "\n")
	return strings.TrimSpace(src)
}

/**
 * string转float
 */
func StringToFloat(str string, bit int) float64 {
	if bit != 32 && bit != 64 {
		bit = 64
	}
	f, _ := strconv.ParseFloat(strings.TrimSpace(str), bit)
	return f
}

/**
 * float64 转 string
 */
func Float64ToString(f float64) string {
	str := strconv.FormatFloat(f, 'g', -1, 64)
	return str
}

/**
 * string转int
 */
func StringToInt(str string) int {
	i, _ := strconv.Atoi(strings.TrimSpace(str))
	return i
}

/**
 * string转int64
 */
func StringToInt64(str string) int64 {
	i, _ := strconv.ParseInt(str, 10, 64)
	return i
}

/**
 * int转string
 */
func IntToString(i int) string {
	str := strconv.Itoa(i)
	return str
}

/**
 * int64转string
 */
func Int64ToString(i int64) string {
	str := strconv.FormatInt(i, 10)
	return str
}

/**
 * 任何类型转换成int
 */
func AnyToInt(i interface{}) int {
	switch i.(type) {
	case int:
		return i.(int)
	case int64:
		return int(i.(int64))
	case string:
		return StringToInt(i.(string))
	case float64:
		return int(i.(float64))
	}
	return 0
}

/**
 * 任何类型转换成string
 */
func AnyToString(i interface{}) string {
	switch i.(type) {
	case int:
		return IntToString(i.(int))
	case int64:
		return Int64ToString(i.(int64))
	case string:
		return i.(string)
	case float64:
		return Float64ToString(i.(float64))
	}
	return ""
}

/**
 * 任何类型转换成float64
 */
func AnyToFloat64(i interface{}) float64 {
	switch i.(type) {
	case int:
		return float64(i.(int))
	case int64:
		return float64(i.(int64))
	case string:
		return StringToFloat(i.(string), 64)
	case float64:
		return i.(float64)
	}
	return 0
}

// 判断变量是否为空
func Empty(params interface{}) bool {
	//初始化变量
	var (
		flag         bool = true
		defaultValue reflect.Value
	)

	r := reflect.ValueOf(params)

	//获取对应类型默认值
	defaultValue = reflect.Zero(r.Type())
	//由于params 接口类型 所以default_value也要获取对应接口类型的值 如果获取不为接口类型 一直为返回false
	if !reflect.DeepEqual(r.Interface(), defaultValue.Interface()) {
		flag = false
	}
	return flag
}

// 字符串编码转换
func ConvertEncoding(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}
