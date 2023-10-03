package utils

import (
	"fmt"
	"time"
)

// 获得当前时间的毫秒级时间戳
func GetCurrentTimeMillis() int64 {
	unixTimestampNano := time.Now().UnixNano()
	return unixTimestampNano / int64(time.Millisecond)
}

// 获取当前时间的秒级时间戳
func GetCurrentTimeUnixSecond() int64 {
	return time.Now().Unix()
}

// 获取当前时间的纳秒级时间戳
func GetCurrentTimeUnixNano() int64 {
	return time.Now().UnixNano()
}

func ShowTime() {
	// 微秒（us）
	us := 1000 // 1毫秒 = 1000微秒
	durationUs := time.Duration(us) * time.Microsecond
	fmt.Printf("%d us = %v\n", us, durationUs)

	// 毫秒（ms）
	ms := 1000 // 1秒 = 1000毫秒
	durationMs := time.Duration(ms) * time.Millisecond
	fmt.Printf("%d ms = %v\n", ms, durationMs)

	// 秒
	seconds := 5
	durationSeconds := time.Duration(seconds) * time.Second
	fmt.Printf("%d seconds = %v\n", seconds, durationSeconds)

	// 纳秒（ns）
	ns := 1000 // 1微秒 = 1000纳秒
	durationNs := time.Duration(ns) * time.Nanosecond
	fmt.Printf("%d ns = %v\n", ns, durationNs)
}
