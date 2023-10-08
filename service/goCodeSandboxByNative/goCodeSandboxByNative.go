/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-08 11:34:56
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-08 13:53:54
 * @FilePath: /xoj-code-sandbox/service/codeSandboxByNative/codeSandboxByNative.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package gocodesandboxbynative

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/model"
	"github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/mylog"
	"github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/service"
)

const (
	WINDOWS = "windows"
)

type GoCodeSandboxByNative struct {
}

// 1. 把用户的代码保存为文件
func (this GoCodeSandboxByNative) SaveCodeToFile(code string) (string, error) {
	return service.SaveCodeToFile(code)
}

// 2. 编译代码
func (this GoCodeSandboxByNative) CompileFile(userCodePath string) error {
	userCodeParentPath := filepath.Dir(userCodePath)

	// 构建编译命令
	var compileCmd string
	if runtime.GOOS == WINDOWS {
		compileCmd = fmt.Sprintf("go build -o %s\\%s.exe %s", userCodeParentPath, service.GLOBAL_GO_BINARY_NAME, userCodePath)
	} else {
		compileCmd = fmt.Sprintf("go build -o %s/%s %s", userCodeParentPath, service.GLOBAL_GO_BINARY_NAME, userCodePath)
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
	case <-time.After(service.TIME_OUT):
		// 超时
		compileProcess.Process.Kill()
		return errors.New(service.COMPILE_TIMEOUT_ERROR)
	case err := <-done:
		if err != nil {
			mylog.Log.Errorf("%v : err= %v", service.COMPILE_ERROR, err.Error())
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
		runCmd = fmt.Sprintf("%s\\%s.exe", userCodeParentPath, service.GLOBAL_GO_BINARY_NAME)
	} else {
		runCmd = fmt.Sprintf("%s/%s", userCodeParentPath, service.GLOBAL_GO_BINARY_NAME)
	}

	runCmd = strings.TrimSpace(runCmd)

	var execResultList []model.ExecResult

	for _, input := range inputList {
		inputParst := strings.Split(strings.TrimSpace(input), " ")
		runProcess := exec.Command(runCmd, inputParst...)
		// runProcess.Stdin = strings.NewReader(input)

		startTime := time.Now()

		// CombinedOutput运行该命令并返回其组合的标准输出和标准错误。
		output, err := runProcess.CombinedOutput()

		latencyTm := time.Since(startTime).Milliseconds()

		if err != nil {
			if strings.Contains(err.Error(), "signal: killed") {
				execResultList = append(execResultList, model.ExecResult{
					StdErr: service.RUN_TIMEOUT_ERROR,
					Time:   latencyTm,
				})
				return execResultList, err
			} else {
				execResultList = append(execResultList, model.ExecResult{
					StdErr: err.Error(),
					Time:   latencyTm,
				})
				return execResultList, err
			}
		} else {
			execResultList = append(execResultList, model.ExecResult{
				StdOut: string(output),
				Time:   latencyTm,
			})
		}
	}
	return execResultList, nil
}

// 4. 获取输出结果
func (this GoCodeSandboxByNative) GetOutputResponse(execResultList []model.ExecResult) model.ExecuteCodeResponse {
	return service.GetOutputResponse(execResultList)
}

// 5. 删除文件
func (this GoCodeSandboxByNative) DeleteFile(userCodePath string) error {
	return service.DeleteFile(userCodePath)
}
