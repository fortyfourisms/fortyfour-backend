package utils

import "time"

func TimeNowUnix() int64 {
	return time.Now().Unix()
}
