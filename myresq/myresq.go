/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-15 21:29:16
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-16 13:20:16
 * @FilePath: /xoj-code-sandbox/myresq/myresq.go
 */
package myresq

import (
	"github.com/beego/beego/v2/server/web/context"
	"github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/mylog"
)

// 通用返回类
type BaseResponse struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"` //当 data 为 nil 时，将不会包含 data 字段
	Message string      `json:"message"`
}

// 创建新的通用返回对象
func NewBaseResponse(code RespCode, message string, data interface{}) *BaseResponse {
	return &BaseResponse{
		Code:    int(code),
		Data:    data,
		Message: message,
	}
}

func Abort(ctx *context.Context, code RespCode, msg string) {
	if ctx == nil {
		mylog.Log.Error("Abort param ctx is nil")
		return
	}
	message := code.GetMessage()
	if msg != "" {
		message = msg
	}
	jsondata := NewBaseResponse(code, message, nil)
	ctx.Input.SetData("json", jsondata)
	ctx.Abort(200, "")
}

func Success(ctx *context.Context, data interface{}) {
	if ctx == nil {
		mylog.Log.Error("Success param ctx is nil")
		return
	}
	jsondata := NewBaseResponse(SUCCESS, SUCCESS.GetMessage(), data)
	ctx.Input.SetData("json", jsondata)
	ctx.Abort(200, "")
}

func AbortWithData(ctx *context.Context, code RespCode, msg string, data interface{}) {
	if ctx == nil {
		mylog.Log.Error("AbortWithData param ctx is nil")
		return
	}
	message := code.GetMessage()
	if msg != "" {
		message = msg
	}
	jsondata := NewBaseResponse(code, message, data)
	ctx.Input.SetData("json", jsondata)
	ctx.Abort(200, "")
}
