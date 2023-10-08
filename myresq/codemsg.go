/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-01 17:34:58
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-08 13:57:11
 */
package myresq

import "strconv"

// 自定义返回码的类型
type RespCode int

// 获取错误消息
func (ec RespCode) GetMessage() string {
	if respCodeMessages[ec] != "" {
		return respCodeMessages[ec]
	}
	return strconv.Itoa(int(ec))
}

const (
	// 定义枚举值
	SUCCESS RespCode = iota
	PARAMS_ERROR
	NOT_LOGIN_ERROR
	NO_AUTH_ERROR
	NOT_FOUND_ERROR
	FORBIDDEN_ERROR
	SYSTEM_ERROR
	OPERATION_ERROR
	GET_CONTEXT_ERROR
	API_REQUEST_ERROR
	GENERATE_RANDOMKEY_FAILED
	GENERATE_TOKEN_FAILED
)
const (
	USER_NOT_EXIST RespCode = iota + 3000
	USER_EXIST
	CREATE_USER_FAILED
	USER_PASSWORD_ERROR
)

const (
	EXECUTE_CODE_ERROR RespCode = iota + 4000
)

// 错误消息映射

var respCodeMessages = map[RespCode]string{
	SUCCESS:             "success",
	PARAMS_ERROR:        "请求参数错误",
	NOT_LOGIN_ERROR:     "未登录",
	NO_AUTH_ERROR:       "无权限",
	NOT_FOUND_ERROR:     "请求数据不存在",
	FORBIDDEN_ERROR:     "禁止访问",
	SYSTEM_ERROR:        "系统内部异常",
	OPERATION_ERROR:     "操作失败",
	GET_CONTEXT_ERROR:   "获取上下文信息失败",
	API_REQUEST_ERROR:   "接口调用失败",
	USER_NOT_EXIST:      "用户不存在",
	USER_EXIST:          "用户已存在",
	CREATE_USER_FAILED:  "账号创建失败",
	USER_PASSWORD_ERROR: "账号密码错误",
	EXECUTE_CODE_ERROR:  "代码执行错误",
}
