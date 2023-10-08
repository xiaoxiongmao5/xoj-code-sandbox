/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-08 11:22:12
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-08 13:41:30
 * @FilePath: /xoj-code-sandbox/service/CodeSandboxTemplate.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package service

import (
	"github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/model"
)

type CodeSandboxInterface interface {
	// 1. 保存用户代码到文件中
	SaveCodeToFile(string) (string, error)

	// 2. 编译代码
	CompileFile(string) error

	// 3. 运行代码
	RunFile(string, []string) ([]model.ExecResult, error)

	// 4. 整理结果
	GetOutputResponse([]model.ExecResult) model.ExecuteCodeResponse

	// 5. 清理文件
	DeleteFile(string) error
}

func CodeSandboxTemplate(c CodeSandboxInterface, param model.ExecuteCodeRequest) (executeCodeResponse model.ExecuteCodeResponse, err error) {
	userCodePath, err := c.SaveCodeToFile(param.Code)
	if err != nil {
		return executeCodeResponse, err
	}

	err = c.CompileFile(userCodePath)
	if err != nil {
		return executeCodeResponse, err
	}

	execResultList, err := c.RunFile(userCodePath, param.InputList)
	if err != nil {
		return executeCodeResponse, err
	}

	executeCodeResponse = c.GetOutputResponse(execResultList)

	err = c.DeleteFile(userCodePath)
	if err != nil {
		return executeCodeResponse, err
	}

	return executeCodeResponse, nil
}
