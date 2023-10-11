/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-08 11:34:56
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-11 21:15:19
 */
package service

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/model"
	codeexecstatusenum "github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/model/enums/CodeExecStatusEnum"
	"github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/mylog"
	commonservice "github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/service/commonService"
)

const (
	WINDOWS = "windows"
)

type GoCodeSandboxByNative struct {
}

// 1. 把用户的代码保存为文件
func (this GoCodeSandboxByNative) SaveCodeToFile(code string) (string, error) {
	return commonservice.SaveCodeToFile(code)
}

// 2. 编译代码
func (this GoCodeSandboxByNative) CompileFile(userCodePath string) error {
	userCodeParentPath := filepath.Dir(userCodePath)

	// 构建编译命令
	var compileCmd string
	if runtime.GOOS == WINDOWS {
		compileCmd = fmt.Sprintf("go build -o %s\\%s.exe %s", userCodeParentPath, commonservice.GLOBAL_GO_BINARY_NAME, userCodePath)
	} else {
		compileCmd = fmt.Sprintf("go build -o %s/%s %s", userCodeParentPath, commonservice.GLOBAL_GO_BINARY_NAME, userCodePath)
	}

	// TrimSpace返回字符串s的一个片段，去掉所有前导和尾随空格
	compileCmd = strings.TrimSpace(compileCmd)
	compileCmdParts := strings.Split(compileCmd, " ")

	// 编译代码
	compileProcess := exec.Command(compileCmdParts[0], compileCmdParts[1:]...)
	compileProcess.Stderr = os.Stderr
	compileProcess.Stdout = os.Stdout

	// 启动编译进程
	if err := compileProcess.Start(); err != nil {
		mylog.Log.Error("启动编译进程[compileProcess.Start] 失败, err=", err.Error())
		return err
	}

	// 等待编译完成或超时
	done := make(chan error, 1)
	go func() {
		done <- compileProcess.Wait()
	}()

	select {
	case <-time.After(commonservice.TIME_OUT):
		compileProcess.Process.Kill()
		return &model.ErrTimeOut{} //编译超时
	case err := <-done:
		if err != nil {
			mylog.Log.Errorf("%v : err= %v", codeexecstatusenum.COMPILE_FAIL.GetText(), err.Error()) //编译失败
			return err
		}
	}
	return nil
}

// 3. 运行编译后的可执行文件, 获得执行结果列表
func (this GoCodeSandboxByNative) RunFile(userCodePath string, inputList []string) ([]model.ExecResult, error) {
	userCodeParentPath := filepath.Dir(userCodePath)

	// 运行编译后的可执行文件
	var runCmd string
	if runtime.GOOS == WINDOWS {
		runCmd = fmt.Sprintf("%s\\%s.exe", userCodeParentPath, commonservice.GLOBAL_GO_BINARY_NAME)
	} else {
		runCmd = fmt.Sprintf("%s/%s", userCodeParentPath, commonservice.GLOBAL_GO_BINARY_NAME)
	}

	runCmd = strings.TrimSpace(runCmd)

	var execResultList []model.ExecResult

	for i, input := range inputList {
		inputParst := strings.Split(strings.TrimSpace(input), " ")
		runProcess := exec.Command(runCmd, inputParst...)
		// runProcess.Stdin = strings.NewReader(input)

		startTime := time.Now()
		// CombinedOutput运行该命令并返回其组合的标准输出和标准错误。
		output, err := runProcess.CombinedOutput()
		latencyTm := time.Since(startTime).Milliseconds()

		execResult := model.ExecResult{Time: latencyTm}
		if err != nil {
			if strings.Contains(err.Error(), "signal: killed") {
				return execResultList, &model.ErrTimeOut{} //运行超时
			}
			mylog.Log.Errorf("运行用户代码,输入示例[%d]失败,err=%s", i, err.Error())
			return execResultList, err
		}
		execResult.StdOut = string(output)
		execResultList = append(execResultList, execResult)
	}
	return execResultList, nil
}

// 4. 获取输出结果
func (this GoCodeSandboxByNative) GetOutputResponse(execResultList []model.ExecResult) model.ExecuteCodeResponse {
	return commonservice.GetOutputResponse(execResultList)
}

// 5. 删除文件
func (this GoCodeSandboxByNative) DeleteFile(userCodePath string) error {
	return commonservice.DeleteFile(userCodePath)
}
