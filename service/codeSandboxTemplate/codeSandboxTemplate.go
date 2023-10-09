/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-08 11:22:12
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-09 21:05:34
 * @FilePath: /xoj-code-sandbox/service/CodeSandboxTemplate.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package codesandboxtemplate

import (
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

	err = c.CompileFile(userCodePath)
	if err != nil {
		mylog.Log.Errorf("编译失败,但不影响成功的返回沙箱执行结果[err=%s]", err.Error())
		executeCodeResponse.Message = codeexecstatusenum.COMPILE_FAIL.GetText() + ", err: " + err.Error()
		// 编译失败
		executeCodeResponse.Status = codeexecstatusenum.COMPILE_FAIL.GetValue()
		return executeCodeResponse, nil
	}

	execResultList, err := c.RunFile(userCodePath, param.InputList)
	if err != nil {
		mylog.Log.Errorf("运行失败,但不影响成功的返回沙箱执行结果[err=%s]", err.Error())
		executeCodeResponse.Message = codeexecstatusenum.RUN_FAIL.GetText() + ", err: " + err.Error()
		// 用户提交的代码执行中存在错误
		executeCodeResponse.Status = codeexecstatusenum.RUN_FAIL.GetValue()
		return executeCodeResponse, nil
	}

	executeCodeResponse = c.GetOutputResponse(execResultList)

	err = c.DeleteFile(userCodePath)
	if err != nil {
		mylog.Log.Errorf("清理文件失败,但不影响成功的返回沙箱执行结果,[userCodePath=%s] [err=%s]", userCodePath, err.Error())
		return executeCodeResponse, nil
	}

	return executeCodeResponse, nil
}
