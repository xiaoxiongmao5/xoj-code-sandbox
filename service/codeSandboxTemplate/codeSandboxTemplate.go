/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-08 11:22:12
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-18 00:51:37
 * @FilePath: /xoj-code-sandbox/service/CodeSandboxTemplate.go
 * @Description: 代码沙箱-模版方法
 */
package codesandboxtemplate

import (
	"errors"

	"github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/model"
	codeexecstatusenum "github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/model/enums/CodeExecStatusEnum"
	"github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/myerror"
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

func CodeSandboxTemplate(c CodeSandboxInterface, param model.ExecuteCodeRequest) (model.ExecuteCodeResponse, error) {
	var executeCodeResponse model.ExecuteCodeResponse
	userCodePath, err := c.SaveCodeToFile(param.Code)
	if err != nil {
		executeCodeResponse.Message = codeexecstatusenum.SYSTEM_ERROR.GetText() + " : " + err.Error()
		// 系统错误(保存代码失败)
		executeCodeResponse.Status = codeexecstatusenum.SYSTEM_ERROR.GetValue()
		return executeCodeResponse, err
	}

	defer func(c CodeSandboxInterface, userCodePaths string) {
		go c.DeleteFile(userCodePath)
	}(c, userCodePath)

	err = c.CompileFile(userCodePath)
	if err != nil {
		var e myerror.ErrTimeOut
		if errors.As(err, &e) {
			mylog.Log.Error("编译超时, err=", err.Error())
			executeCodeResponse.Message = codeexecstatusenum.COMPILE_TIMEOUT_ERROR.GetText()
			executeCodeResponse.Status = codeexecstatusenum.COMPILE_TIMEOUT_ERROR.GetValue()
			return executeCodeResponse, err
		}
		var e2 myerror.ErrMemoryFullOut
		if errors.As(err, &e2) {
			mylog.Log.Error("编译导致内存不足, err=", err.Error())
			executeCodeResponse.Message = codeexecstatusenum.OUT_OF_MEMORY_ERROR.GetText()
			executeCodeResponse.Status = codeexecstatusenum.OUT_OF_MEMORY_ERROR.GetValue()
			return executeCodeResponse, err
		}
		mylog.Log.Error("编译失败, err=", err.Error())
		executeCodeResponse.Message = codeexecstatusenum.COMPILE_FAIL.GetText() + " : " + err.Error()
		executeCodeResponse.Status = codeexecstatusenum.COMPILE_FAIL.GetValue()
		return executeCodeResponse, nil
	}

	execResultList, err := c.RunFile(userCodePath, param.InputList)
	executeCodeResponse = c.GetOutputResponse(execResultList)
	if err != nil {
		var e myerror.ErrTimeOut
		if errors.As(err, &e) {
			mylog.Log.Error("运行代码超时, err=", err.Error())
			executeCodeResponse.Message = codeexecstatusenum.RUN_TIMEOUT_ERROR.GetText()
			executeCodeResponse.Status = codeexecstatusenum.RUN_TIMEOUT_ERROR.GetValue()
			return executeCodeResponse, err
		}
		var e2 myerror.ErrMemoryFullOut
		if errors.As(err, &e2) {
			mylog.Log.Error("运行代码导致内存不足, err=", err.Error())
			executeCodeResponse.Message = codeexecstatusenum.OUT_OF_MEMORY_ERROR.GetText()
			executeCodeResponse.Status = codeexecstatusenum.OUT_OF_MEMORY_ERROR.GetValue()
			return executeCodeResponse, err
		}
		mylog.Log.Error("运行代码失败, err=", err.Error())
		executeCodeResponse.Message = codeexecstatusenum.RUN_FAIL.GetText() + " : " + err.Error()
		executeCodeResponse.Status = codeexecstatusenum.RUN_FAIL.GetValue()
		return executeCodeResponse, nil
	}

	return executeCodeResponse, nil
}
