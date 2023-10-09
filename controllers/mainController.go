/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-08 13:10:34
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-09 19:15:38
 * @FilePath: /xoj-code-sandbox/controllers/controller.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package controllers

import (
	"context"

	beego "github.com/beego/beego/v2/server/web"
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

	goCodeSandboxByDocker := service.GoCodeSandboxByDocker{
		Ctx: context.Background(),
		Cli: mydocker.Cli,
	}
	// goCodeSandboxByNative := service.GoCodeSandboxByNative{}

	executeCodeResponse, err := codesandboxtemplate.CodeSandboxTemplate(goCodeSandboxByDocker, params)
	if err != nil {
		myresq.AbortWithData(this.Ctx, myresq.EXECUTE_CODE_ERROR, err.Error(), executeCodeResponse)
		return
	}
	myresq.Success(this.Ctx, executeCodeResponse)
}
