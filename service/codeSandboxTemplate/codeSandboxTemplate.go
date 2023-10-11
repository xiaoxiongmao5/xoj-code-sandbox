/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-08 11:22:12
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-11 21:50:45
 * @FilePath: /xoj-code-sandbox/service/CodeSandboxTemplate.go
 * @Description: 代码沙箱-模版方法
 */
package codesandboxtemplate

import (
	"errors"

	"github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/model"
	codeexecstatusenum "github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/model/enums/CodeExecStatusEnum"
	"github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/mylog"
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
		executeCodeResponse.Message = codeexecstatusenum.SYSTEM_ERROR.GetText() + ", err: " + err.Error()
		// 系统错误(保存代码失败)
		executeCodeResponse.Status = codeexecstatusenum.SYSTEM_ERROR.GetValue()
		return executeCodeResponse, err
	}

	defer c.DeleteFile(userCodePath)

	err = c.CompileFile(userCodePath)
	if err != nil {
		var e model.ErrTimeOut
		if errors.As(err, &e) {
			mylog.Log.Error("编译超时, err=", err.Error())
			executeCodeResponse.Message = codeexecstatusenum.COMPILE_TIMEOUT_ERROR.GetText()
			executeCodeResponse.Status = codeexecstatusenum.COMPILE_TIMEOUT_ERROR.GetValue()
			return executeCodeResponse, err
		}
		mylog.Log.Error("编译失败,err=", err.Error())
		executeCodeResponse.Message = codeexecstatusenum.COMPILE_FAIL.GetText() + ", err: " + err.Error()
		executeCodeResponse.Status = codeexecstatusenum.COMPILE_FAIL.GetValue()
		return executeCodeResponse, nil
	}

	execResultList, err := c.RunFile(userCodePath, param.InputList)
	if err != nil {
		var e model.ErrTimeOut
		if errors.As(err, &e) {
			mylog.Log.Error("运行用户代码超时, err=", err.Error())
			executeCodeResponse.Message = codeexecstatusenum.RUN_TIMEOUT_ERROR.GetText()
			executeCodeResponse.Status = codeexecstatusenum.RUN_TIMEOUT_ERROR.GetValue()
			return executeCodeResponse, err
		}
		mylog.Log.Error("运行用户代码失败,err=", err.Error())
		executeCodeResponse.Message = codeexecstatusenum.RUN_FAIL.GetText() + ", err: " + err.Error()
		executeCodeResponse.Status = codeexecstatusenum.RUN_FAIL.GetValue()
		return executeCodeResponse, nil
	}

	executeCodeResponse = c.GetOutputResponse(execResultList)

	return executeCodeResponse, nil
}
