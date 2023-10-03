/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-09-26 10:35:03
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-03 19:34:26
 */
package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/mylog"
)

// 检查是否为空字符串
func IsAnyBlank(values ...interface{}) bool {
	for _, value := range values {
		if IsEmpty(value) {
			return true
		}
	}
	return false
}

// 检查不为空
func IsNotBlank(value interface{}) bool {
	return !IsEmpty(value)
}

// 检查为空
func IsEmpty(value interface{}) bool {
	if value == nil {
		return true
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String:
		return v.Len() == 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Slice, reflect.Array:
		return v.Len() == 0
	case reflect.Map, reflect.Ptr, reflect.Interface:
		return v.IsNil()
	}

	return false
}

// 检查是否一样（使用 == 检查）
func CheckSame[T string | int | int8 | int16 | int32 | int64](desc string, str1 T, str2 T) bool {
	res := false
	if str1 == str2 {
		res = true
	} else {
		res = false
	}
	mylog.Log.WithFields(logrus.Fields{
		"got":    str1,
		"export": str2,
		"pass":   res,
	}).Info(desc)
	return res
}

// 检查字符串忽略大小写后是否一样（使用 EqualFold 检查）
func CheckSameStrFold(desc string, str1 string, str2 string) bool {
	res := false
	if strings.EqualFold(str1, str2) {
		res = true
	} else {
		res = false
	}
	mylog.Log.WithFields(logrus.Fields{
		"got":    str1,
		"export": str2,
		"pass":   res,
		"notes":  "已忽略大小写",
	}).Info(desc)
	return res
}

// 检查数组是否一样（使用 DeepEqual 检查）
func CheckSameArr[T string | int | []int](desc string, str1 T, str2 T) bool {
	res := false
	if reflect.DeepEqual(str1, str2) {
		res = true
	} else {
		res = false
	}
	mylog.Log.WithFields(logrus.Fields{
		"got":    fmt.Sprintf("%v", str1),
		"export": fmt.Sprintf("%v", str2),
		"pass":   res,
	}).Info(desc)
	return res
}
