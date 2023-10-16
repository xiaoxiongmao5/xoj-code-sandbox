/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-09-27 14:46:54
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-16 13:17:51
 */
package utils

import (
	"encoding/json"
	"errors"
	"reflect"
	"strconv"

	"github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/mylog"
)

/** 字符串切片转成int切片
 */
func StrSlice2IntSlice(strSlice []string) ([]int, error) {
	if strSlice == nil {
		return nil, errors.New("param is nil")
	}
	intSlice := make([]int, len(strSlice))
	for i, str := range strSlice {
		intVal, err := strconv.Atoi(str)
		if err != nil {
			return intSlice, err
		}
		intSlice[i] = intVal
	}
	return intSlice, nil
}

/** int切片转字符串切片
 */
func IntSlice2StrSlice(intSlice []int) []string {
	if intSlice == nil {
		return nil
	}
	strSlice := make([]string, len(intSlice))
	for i, intVal := range intSlice {
		strSlice[i] = strconv.Itoa(intVal)
	}
	return strSlice
}

/** 前端分页参数转数据库查询的limit和offser
 */
func CalculateLimitOffset[T int | int32 | int64](current, pageSize T) (limit, offset T) {
	if current < 1 {
		current = 1
	}
	offset = (current - 1) * pageSize
	limit = pageSize
	return limit, offset
}

/** string 转 int64
 */
func String2Int64(str string) (num int64, err error) {
	return strconv.ParseInt(str, 10, 64)
}

// 将源结构的字段值复制到目标结构中
func CopyStructFields(src, dst interface{}) bool {
	//将 src 转换为反射值reflect.Value,以便后续可以通过反射操作源结构体的字段。
	srcValue := reflect.ValueOf(src)
	// 将 dst 转换为反射值,dst 通常是目标结构体的指针，因此需要通过 Elem() 方法获取指针指向的结构体。
	dstValue := reflect.ValueOf(dst)

	// 检查 src 是否是结构体类型
	if srcValue.Kind() != reflect.Struct {
		mylog.Log.Error("[CopyStructFields] src 不是结构体类型")
		return false
	}
	// 检查 dst 是否是指针类型
	if dstValue.Kind() != reflect.Pointer {
		mylog.Log.Error("[CopyStructFields] dst 不是指针类型")
		return false
	}
	// dstValue.Elem() 返回指针指向的值，然后检查它是否是结构体类型。
	// 检查目标结构体是否有效
	if dstValue.Elem().Kind() != reflect.Struct {
		mylog.Log.Error("[CopyStructFields] dst 结构体无效")
		return false
	}

	dstElem := dstValue.Elem()
	srcType := srcValue.Type()
	dstType := dstElem.Type()

	for i := 0; i < srcType.NumField(); i++ {
		srcField := srcType.Field(i)
		dstField, found := dstType.FieldByName(srcField.Name)
		if !found {
			continue
		}

		if srcField.Type == dstField.Type {
			dstFieldValue := dstElem.FieldByName(srcField.Name)
			srcFieldValue := srcValue.Field(i)
			dstFieldValue.Set(srcFieldValue)
		}
	}

	return true
}

func JsonMarshal(data interface{}, defaultStr string) (string, error) {
	jsonStr, err := json.Marshal(data)
	if err != nil {
		mylog.Log.Error("json.Marshal failed, err=", err)
		return defaultStr, err
	}
	return string(jsonStr), nil
}
