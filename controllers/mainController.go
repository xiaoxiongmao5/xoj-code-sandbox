/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-08 13:10:34
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-18 00:49:57
 * @FilePath: /xoj-code-sandbox/controllers/controller.go
 */
package controllers

import (
	"context"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/config"
	"github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/model"
	"github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/mydocker"
	"github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/myresq"
	"github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/service"
	codesandboxtemplate "github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/service/codeSandboxTemplate"
)

type MainController struct {
	beego.Controller
}

//	@Summary		执行代码
//	@Description	执行代码
//	@Tags			代码沙箱
//	@Accept			application/json
//	@Produce		application/json
//	@Param			request	body		model.ExecuteCodeRequest		true	"请求参数"
//	@Success		200		{object}	swagtype.ExecuteCodeResponse	"响应数据"
//	@Router			/executeCode [post]
func (this MainController) ExecuteCode() {
	var params model.ExecuteCodeRequest
	if err := this.BindJSON(&params); err != nil {
		myresq.Abort(this.Ctx, myresq.PARAMS_ERROR, "")
		return
	}

	var executeCodeResponse model.ExecuteCodeResponse
	var err error
	var goCodeSandbox codesandboxtemplate.CodeSandboxInterface

	switch config.AppConfigDynamic.CodeSandboxType {
	case "docker":
		goCodeSandbox = service.GoCodeSandboxByDocker{
			Ctx: context.Background(),
			Cli: mydocker.Cli,
		}
	case "native":
		goCodeSandbox = service.GoCodeSandboxByNative{}
	case "dockerAndNative":
		goCodeSandbox = service.GoCodeSandboxByDockerNative{
			Ctx: context.Background(),
			Cli: mydocker.Cli,
		}
	default:
		goCodeSandbox = service.GoCodeSandboxByDocker{
			Ctx: context.Background(),
			Cli: mydocker.Cli,
		}
	}

	executeCodeResponse, err = codesandboxtemplate.CodeSandboxTemplate(goCodeSandbox, params)

	if err != nil {
		myresq.AbortWithData(this.Ctx, myresq.EXECUTE_CODE_ERROR, err.Error(), executeCodeResponse)
		return
	}
	myresq.Success(this.Ctx, executeCodeResponse)
}
