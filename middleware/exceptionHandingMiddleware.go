/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-10 17:11:45
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-11 17:01:32
 * @Description: 捕获中断业务异常的中间件
 */
package middleware

import (
	"net/http"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/mylog"
	"github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/myresq"
)

func ExceptionHandingMiddleware(ctx *context.Context, config *beego.Config) {
	if err := recover(); err != nil {
		mylog.Log.Errorf("beego.BConfig.RecoverFunc err= %v \n", err)

		// 从 Context 中获取错误码和消息
		response, ok := ctx.Input.GetData("json").(*myresq.BaseResponse)
		if !ok {
			response = myresq.NewBaseResponse(500, "未知错误", nil)
		}

		// 将 JSON 响应写入 Context，并设置响应头
		ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
		ctx.Output.SetStatus(http.StatusOK)
		ctx.Output.JSON(response, false, false)
	}
}
