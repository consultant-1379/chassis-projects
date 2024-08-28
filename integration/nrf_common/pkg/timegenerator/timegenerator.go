package timegenerator

import (
	"time"
)

var localtimeInSecond int64

// GetLocalTime is to get the local time
func GetLocalTime() int64 {
	return localtimeInSecond
}

// GenerateLocalTime is to generate the local time every second
func GenerateLocalTime() {
	ticker := time.NewTicker(time.Second * time.Duration(1))
	go func() {
		for {
			select {
			case t := <-ticker.C:
				localtimeInSecond = t.Unix()
			}
		}
	}()
}
