package main

import (
	"net/http"

	"github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/mydocker"
	"github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/mylog"
	"github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/myresq"
	_ "github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/routers"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

func init() {
	mylog.Log.Info("init begin: main")

	var err error

	// 创建Docker客户端
	mydocker.Cli, err = mydocker.CreateDockerClient()
	if err != nil {
		panic(err)
	}

	mylog.Log.Info("init end  : main")
}

func main() {
	defer mylog.Log.Writer().Close()
	defer mydocker.Cli.Close()

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"

		// // 开启监控：Admin 管理后台
		// beego.BConfig.Listen.EnableAdmin = true
		// // 修改监听的地址和端口：
		// beego.BConfig.Listen.AdminAddr = "localhost"
		// beego.BConfig.Listen.AdminPort = 8089
	}

	// 全局异常捕获
	beego.BConfig.RecoverPanic = true
	beego.BConfig.RecoverFunc = func(ctx *context.Context, config *beego.Config) {
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

	beego.Run()
}
