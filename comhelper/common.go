package comhelper

import "encoding/json"

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
