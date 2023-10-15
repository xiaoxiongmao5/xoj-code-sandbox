/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-09 20:24:21
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-15 22:05:50
 * @FilePath: /xoj-code-sandbox/main.go
 */
package main

import (
	"github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/config"
	"github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/middleware"
	"github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/mydocker"
	"github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/mylog"
	_ "github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/routers"

	beego "github.com/beego/beego/v2/server/web"
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

	// 启动动态配置文件加载协程
	go config.LoadAppDynamicConfigCycle()

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
	beego.BConfig.RecoverFunc = middleware.ExceptionHandingMiddleware

	beego.Run()
}
