package apihelper

import (
	"encoding/json"
	"helper_go/nethelper"
)

/**
 * 高德地图api相关
 */

// 坐标转换
type CoordinateConvertArgs struct {
	Key       string // 应用key
	Locations string // 经纬度
	Coordsys  string // 原坐标系类型
}

func CoordinateConvert(param *CoordinateConvertArgs) map[string]interface{} {
	url := "https://restapi.amap.com/v3/assistant/coordinate/convert?"
	url += "key=" + param.Key + "&locations=" + param.Locations + "&coordsys=" + param.Coordsys

	ret := nethelper.HttpGet(url)
	var data map[string]interface{}
	err := json.Unmarshal([]byte(ret), &data)
	if err != nil {
		return map[string]interface{}{}
	}
	return data
}

// 逆地理编码
type GeocodeRegeoArgs struct {
	Key        string // 应用key
	Locations  string // 经纬度
	Poitype    string // 返回附近POI类型
	Radius     string // 搜索半径
	Extensions string // 返回结果控制
	Batch      string // 批量查询控制
	Roadlevel  string // 道路等级
}

func GeocodeRegeo(param *GeocodeRegeoArgs) map[string]interface{} {
	url := "https://restapi.amap.com/v3/geocode/regeo?"
	url += "key=" + param.Key + "&location=" + param.Locations
	url += "&poitype=" + param.Poitype + "&radius=" + param.Radius
	url += "&extensions=" + param.Extensions + "&batch=" + param.Batch
	url += "&roadlevel=" + param.Roadlevel

	ret := nethelper.HttpGet(url)
	var data map[string]interface{}
	err := json.Unmarshal([]byte(ret), &data)
	if err != nil {
		return map[string]interface{}{}
	}
	return data
}
