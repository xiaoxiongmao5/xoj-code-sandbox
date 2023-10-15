/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-08 11:37:18
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-15 20:30:41
 */
package service

import (
	"context"

	"github.com/docker/docker/client"
	"github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/model"
	commonservice "github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/service/commonService"
)

type GoCodeSandboxByDockerNative struct {
	Ctx context.Context
	Cli *client.Client
}

// 1. 把用户的代码保存为文件
func (this GoCodeSandboxByDockerNative) SaveCodeToFile(code string) (string, error) {
	return commonservice.SaveCodeToFile(code)
}

// 2. 编译代码
func (this GoCodeSandboxByDockerNative) CompileFile(userCodePath string) error {
	return GoCodeSandboxByNative{}.CompileFile(userCodePath)
}

// 3. 运行编译后的可执行文件, 获得执行结果列表
func (this GoCodeSandboxByDockerNative) RunFile(userCodePath string, inputList []string) (execResultList []model.ExecResult, err error) {
	return GoCodeSandboxByDocker{
		Ctx: this.Ctx,
		Cli: this.Cli,
	}.RunFile(userCodePath, inputList)
}

// 4. 获取输出结果
func (this GoCodeSandboxByDockerNative) GetOutputResponse(execResultList []model.ExecResult) model.ExecuteCodeResponse {
	return commonservice.GetOutputResponse(execResultList)
}

// 5. 删除文件
func (this GoCodeSandboxByDockerNative) DeleteFile(userCodePath string) error {
	return commonservice.DeleteFile(userCodePath)
}
