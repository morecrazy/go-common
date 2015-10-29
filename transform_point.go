package common

import (
	//"reflect"
	"math"
)

var a = 6378245.0
var ee = 0.00669342162296594323

func OutOfChina(lat, lon float64) bool {
	if (lat < 0.8293) || (lat > 55.8271) {
		return true
	}

	if (lon < 72.004) || (lon > 137.8347) {
		return true
	}

	return false
}

//World Geodetic System ==> Mars Geodetic System
func TransForm(wglat, wglon float64) (float64, float64) {
	//note: 原始点转高德
	//transform(latitude,longitude) , WGS84
	//return (latitude,longitude) , GCJ02
	if OutOfChina(wglat, wglon) {
		return wglat, wglon
	}

	dlat := TransFormLat(wglon-105.0, wglat-35.0)
	dlon := TransFormLon(wglon-105.0, wglat-35.0)
	rad_lat := wglat / 18.0 * math.Pi
	magic := math.Sin(rad_lat)
	magic = 1 - ee*magic*magic
	sqrt_magic := math.Sqrt(magic)
	dlat = (dlat * 180.0) / ((a * (1 - ee)) / (magic * sqrt_magic) * math.Pi)
	dlon = (dlon * 180.0) / (a / sqrt_magic * math.Cos(rad_lat) * math.Pi)
	mglat := wglat + dlat
	mglon := wglon + dlon
	return mglat, mglon
}

func TransFormLat(x, y float64) float64 {
	ret := -100.0 + 2.0*x + 3.0*y + 0.2*y*y + 0.1*x*y + 0.2*math.Sqrt(math.Abs(x))
	ret += (20.0*math.Sin(6.0*x*math.Pi) + 20.0*math.Sin(2.0*x*math.Pi)) * 2.0 / 3.0
	ret += (20.0*math.Sin(y*math.Pi) + 40.0*math.Sin(y/3.0*math.Pi)) * 2.0 / 3.0
	ret += (160.0*math.Sin(y/12.0*math.Pi) + 320*math.Sin(y*math.Pi/30.0)) * 2.0 / 3.0
	return ret
}

func TransFormLon(x, y float64) float64 {
	ret := 300.0 + x + 2.0*y + 0.1*x*x + 0.1*x*y + 0.1*math.Sqrt(math.Abs(x))
	ret += (20.0*math.Sin(6.0*x*math.Pi) + 20.0*math.Sin(2.0*x*math.Pi)) * 2.0 / 3.0
	ret += (20.0*math.Sin(x*math.Pi) + 40.0*math.Sin(x/3.0*math.Pi)) * 2.0 / 3.0
	ret += (150.0*math.Sin(x/12.0*math.Pi) + 300.0*math.Sin(x/30.0*math.Pi)) * 2.0 / 3.0
	return ret
}

//角度求弧度
//A2 = A1 * pi / 180.0
//A2:弧度, A1:角度
func Radians(x float64) float64 {
	return x * math.Pi / 180.0
}

func Distance(lat, lng, val_lat, val_lng float64) float64 {
	dlat := Radians(val_lat - lat)
	dlon := Radians(val_lng - lng)
	tmp := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(Radians(lat))*math.Cos(Radians(val_lat))*math.Sin(dlon/2)*math.Sin(dlon/2)

	return 6371 * 2 * math.Atan2(math.Sqrt(tmp), math.Sqrt(1-tmp))
}
