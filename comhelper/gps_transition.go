package comhelper

import (
	"math"
)

const (
	X_PI   = math.Pi * 3000.0 / 180.0
	OFFSET = 0.00669342162296594323
	AXIS   = 6378245.0
)

// WGS84toGCJ02 WGS84坐标系->火星坐标系
func WGS84toGCJ02(lon, lat float64) (float64, float64) {
	if isOutOFChina(lon, lat) {
		return lon, lat
	}

	mgLon, mgLat := delta(lon, lat)

	return mgLon, mgLat
}

// 高德地图用的是GCJ02标准,我们后台用的是WGS84标准
// GCJ02toWGS84 火星坐标系->WGS84坐标系
func GCJ02toWGS84(lon, lat float64) (float64, float64) {
	if isOutOFChina(lon, lat) {
		return lon, lat
	}

	mgLon, mgLat := delta(lon, lat)

	return lon*2 - mgLon, lat*2 - mgLat
}
func delta(lon, lat float64) (float64, float64) {
	dlat := transformLat(lon-105.0, lat-35.0)
	dlon := transformLng(lon-105.0, lat-35.0)

	radlat := lat / 180.0 * math.Pi
	magic := math.Sin(radlat)
	magic = 1 - OFFSET*magic*magic
	sqrtmagic := math.Sqrt(magic)

	dlat = (dlat * 180.0) / ((AXIS * (1 - OFFSET)) / (magic * sqrtmagic) * math.Pi)
	dlon = (dlon * 180.0) / (AXIS / sqrtmagic * math.Cos(radlat) * math.Pi)

	mgLat := lat + dlat
	mgLon := lon + dlon

	return mgLon, mgLat
}

func transformLat(lon, lat float64) float64 {
	var ret = -100.0 + 2.0*lon + 3.0*lat + 0.2*lat*lat + 0.1*lon*lat + 0.2*math.Sqrt(math.Abs(lon))
	ret += (20.0*math.Sin(6.0*lon*math.Pi) + 20.0*math.Sin(2.0*lon*math.Pi)) * 2.0 / 3.0
	ret += (20.0*math.Sin(lat*math.Pi) + 40.0*math.Sin(lat/3.0*math.Pi)) * 2.0 / 3.0
	ret += (160.0*math.Sin(lat/12.0*math.Pi) + 320*math.Sin(lat*math.Pi/30.0)) * 2.0 / 3.0
	return ret
}

func transformLng(lon, lat float64) float64 {
	var ret = 300.0 + lon + 2.0*lat + 0.1*lon*lon + 0.1*lon*lat + 0.1*math.Sqrt(math.Abs(lon))
	ret += (20.0*math.Sin(6.0*lon*math.Pi) + 20.0*math.Sin(2.0*lon*math.Pi)) * 2.0 / 3.0
	ret += (20.0*math.Sin(lon*math.Pi) + 40.0*math.Sin(lon/3.0*math.Pi)) * 2.0 / 3.0
	ret += (150.0*math.Sin(lon/12.0*math.Pi) + 300.0*math.Sin(lon/30.0*math.Pi)) * 2.0 / 3.0
	return ret
}

func isOutOFChina(lon, lat float64) bool {
	return !(lon > 73.66 && lon < 135.05 && lat > 3.86 && lat < 53.55)
}

/**
 * 地理围栏计算
 * @param lon 		经度
 * @param lat 		纬度
 * @param points 	围栏经纬度
 * @param radius 	半径(误差距离)
 */
func IsInGeoFence(lon, lat float64, points []map[string]float64, radius float64) bool {
	iSum := 0
	var dLon1, dLon2, dLat1, dLat2, dLon float64

	// 判断围栏点数是否小于3
	iCount := len(points)
	if iCount < 3 {
		return false
	}

	// 开始计算
	for i := 0; i < iCount; i++ {
		if i == (iCount - 1) {
			dLon1 = points[i]["longitude"]
			dLat1 = points[i]["latitude"]
			dLon2 = points[0]["longitude"]
			dLat2 = points[0]["latitude"]
		} else {
			dLon1 = points[i]["longitude"]
			dLat1 = points[i]["latitude"]
			dLon2 = points[i+1]["longitude"]
			dLat2 = points[i+1]["latitude"]
		}
		// 以下语句判断A点是否在边的两端点的水平平行线之间，在则可能有交点，开始判断交点是否在左射线上
		if ((lat >= dLat1) && (lat < dLat2)) || ((lat >= dLat2) && (lat < dLat1)) {
			if math.Abs(dLat1-dLat2) > 0 {
				// 得到 A点向左射线与边的交点的x坐标：
				dLon = dLon1 - ((dLon1-dLon2)*(dLat1-lat))/(dLat1-dLat2)
				if dLon < lon {
					iSum++
				}
			}
		}
	}
	if iSum%2 != 0 {
		return true
	}

	// 在区域外的，判断误差，通过点和误差值构造圆形，判断圆形和区域是否相交
	if radius > 0 {
		for i := 0; i < iCount; i++ {
			if i == iCount-1 {
				dLon1 = points[i]["longitude"]
				dLat1 = points[i]["latitude"]
				dLon2 = points[0]["longitude"]
				dLat2 = points[0]["latitude"]
			} else {
				dLon1 = points[i]["longitude"]
				dLat1 = points[i]["latitude"]
				dLon2 = points[i+1]["longitude"]
				dLat2 = points[i+1]["latitude"]
			}
			// 获取点到线的距离，如果距离小于误差距离，则说明在区域内，退出
			dis := PointToLinenDistance(dLat1, dLon1, dLat2, dLon2, lat, lon)
			if dis <= radius {
				return true
			}
		}
	}
	return false
}

/**
 * 计算两经纬度点之间距离
 * @param lon1 经度
 * @param lat1 纬度
 * @param lon2 经度
 * @param lat2 纬度
 * 返回单位：米
 */
func TwoPointDistance(lon1, lat1, lon2, lat2 float64) float64 {
	radius := 6371000.0 // 6378137.0
	rad := math.Pi / 180.0
	lat1 = lat1 * rad
	lon1 = lon1 * rad
	lat2 = lat2 * rad
	lon2 = lon2 * rad
	theta := lon2 - lon1
	dist := math.Acos(math.Sin(lat1)*math.Sin(lat2) + math.Cos(lat1)*math.Cos(lat2)*math.Cos(theta))
	return dist * radius
}

/**
 * 获取不规则区域内的重心点
 * @param mPoints 区域经纬度数组
 * @return Gy 经度
 * @return Gx 纬度
 */
func GetCenterOfGravityPoint(mPoints []map[string]float64) (float64, float64) {
	area := 0.0        // 多边形面积
	Gx, Gy := 0.0, 0.0 // 重心的x、y
	length := len(mPoints)
	for i := 1; i <= length; i++ {
		iLat := mPoints[i%length]["latitude"]
		iLng := mPoints[i%length]["longitude"]
		nextLat := mPoints[i-1]["latitude"]
		nextLng := mPoints[i-1]["longitude"]
		temp := (iLat*nextLng - iLng*nextLat) / 2.0
		area += temp
		Gx += temp * (iLat + nextLat) / 3.0
		Gy += temp * (iLng + nextLng) / 3.0
	}
	Gx = Gx / area
	Gy = Gy / area
	return Gy, Gx
}

/**
 * 计算点到线段的距离，海伦公式
 * @param lat1 	线段点1的纬度
 * @param lng1 	线段点1的经度
 * @param lat2 	线段点2的纬度
 * @param lng2 	线段点2的经度
 * @param lat 	计算点的纬度
 * @param lng 	计算点的经度
 */
func PointToLinenDistance(lat1, lng1, lat2, lng2, lat, lng float64) float64 {
	a := TwoPointDistance(lng1, lat1, lng2, lat2)
	b := TwoPointDistance(lng2, lat2, lng, lat)
	c := TwoPointDistance(lng1, lat1, lng, lat)
	if b*b >= (c*c + a*a) {
		return c
	}
	if c*c >= (b*b + a*a) {
		return b
	}
	l := (a + b + c) / 2
	s := math.Sqrt(l * (l - a) * (l - b) * (l - c))
	return 2 * s / a
}
