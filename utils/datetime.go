package utils

import "time"

var (
	zone *time.Location
	err  error
)

func init() {
	zone, err = time.LoadLocation("Asia/Shanghai")
	if err != nil {
		zone = time.FixedZone("UTC", 8*3600)
	}
}

func LocalMillTime() int64 {
	return time.Now().In(zone).UnixMilli()
}

func ParseMillTimeToStr(datetime int64) string {
	return time.UnixMilli(datetime).In(zone).Format("2006-01-02 15:04:05.000")
}

func LocalTimeStr() string {
	return localTimeStr("2006-01-02 15:04:05")
}

func LocalMillTimeStr() string {
	return localTimeStr("2006-01-02 15:04:05.000")
}

func localTimeStr(format string) string {
	return time.Now().In(zone).Format(format)
}
