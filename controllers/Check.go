/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-03 19:42:49
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-04 20:30:12
 * @FilePath: /xoj-code-sandbox/controllers/Check.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package controllers

import (
	"errors"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/google/uuid"
	"github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/model"
	"github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/mylog"
	"github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/utils"
)

// Operations about object
type CheckController struct {
	beego.Controller
}

func (this CheckController) Health() {
	this.Data["json"] = map[string]int32{"ok": 0}
	this.ServeJSON()
}

func (this CheckController) CheckExec() {
	// compileCmd := "/usr/local/go/bin/go env"
	// compileCmd := "pwd"

	compileProcess := exec.Command("go", "env")
	compileProcess.Stderr = os.Stderr
	compileProcess.Stdout = os.Stdout

	// 启动编译进程
	if err := compileProcess.Start(); err != nil {
		mylog.Log.Error("启动编译进程[compileProcess.Start] 失败, err=", err.Error())
		this.Data["json"] = map[string]int32{"ok": 1}
		this.ServeJSON()
		return
	}

	this.Data["json"] = map[string]int32{"ok": 0}
	this.ServeJSON()
}

func (this CheckController) Do() {
	code := `
	package main

	import (
		"fmt"
		"os"
		"strconv"
	)

	func main() {
		// 使用 os.Args 来获取命令行参数，os.Args[0] 是程序名称，os.Args[1:] 包含所有的命令行参数
		args := os.Args[1:]
		a, _ := strconv.Atoi(args[0])
		b, _ := strconv.Atoi(args[1])
		fmt.Printf("结果: %d", a+b)
	}
	`
	inputList := []string{"1 2", "1 3"}

	// 1. 保存用户代码到文件中
	userCodePath, err := SaveCodeToFile(code)
	if err != nil {
		mylog.Log.Error("保存代码失败, err=", err.Error())
		return
	}

	// 2. 编译代码文件 go build -o xxx
	if err := CompileFile(userCodePath); err != nil {
		mylog.Log.Error("编译错误, err=", err.Error())
		return
	}

	// 3. 运行执行文件，得到输出结果
	executeMessageList, err := RunFile(userCodePath, inputList)
	if err != nil {
		mylog.Log.Error("执行错误, err=", err.Error())
		return
	}

	// 4. 整理结果
	executeCodeResponse := GetOutputResponse(executeMessageList)

	// 5. 清理文件
	if err := DeleteFile(userCodePath); err != nil {
		mylog.Log.Error("清理文件错误, err=", err.Error())
		return
	}

	this.Data["json"] = executeCodeResponse
	this.ServeJSON()

	// 6. 错误处理
}

const (
	GLOBAL_CODE_DIR_NAME   = "tmpcode"
	GLOBAL_GO_FILE_NAME    = "main.go"
	GLOBAL_GO_BINARY_NAME  = "main"
	TIME_OUT               = 5 * time.Second
	COMPILE_TIMEOUT_ERROR  = "编译超时"
	RUN_TIMEOUT_ERROR      = "运行超时"
	RUN_ERROR              = "运行错误"
	COMPILE_ERROR          = "编译错误"
	CODE_SAND_BOX_ERROR    = "代码沙箱错误"
	EXECUTION_SUCCESS      = 1 //正常运行完成
	EXECUTION_COMPILE_FAIL = 2
	EXECUTION_RUNTIME_FAIL = 3 //用户提交的代码执行中存在错误
)

// 1. 把用户的代码保存为文件
func SaveCodeToFile(code string) (string, error) {
	userDir, err := os.Getwd()
	if err != nil {
		mylog.Log.Error("获取[os.Getwd] 失败, err=", err.Error())
		return "", err
	}
	mylog.Log.Info("userDir: ", userDir)

	// 判断全局代码目录是否存在，没有则创建
	globalCodePathName := fmt.Sprintf("%s/%s", userDir, GLOBAL_CODE_DIR_NAME)
	// Stat返回一个描述命名文件的FileInfo。如果出现错误，则其类型为*PathError。
	// IsNotExist返回一个布尔值，指示是否已知报告文件或目录不存在的错误。
	if _, err := os.Stat(globalCodePathName); os.IsNotExist(err) {
		// os.ModePerm 511
		if err := os.Mkdir(globalCodePathName, os.ModePerm); err != nil {
			mylog.Log.Errorf("创建目录[os.Mkdir(%s)] 失败, err=%v", globalCodePathName, err.Error())
			return "", err
		}
	}

	// 把用户的代码隔离存放
	userCodeParentPath := fmt.Sprintf("%s/%s", globalCodePathName, uuid.New().String())
	if err := os.Mkdir(userCodeParentPath, os.ModePerm); err != nil {
		mylog.Log.Errorf("创建目录[os.Mkdir(%s)] 失败, err=%v", userCodeParentPath, err.Error())
		return "", err
	}

	userCodePath := fmt.Sprintf("%s/%s", userCodeParentPath, GLOBAL_GO_FILE_NAME)
	if err := os.WriteFile(userCodePath, []byte(code), 0644); err != nil {
		mylog.Log.Error("保存代码到文件[os.WriteFile] 失败, err=", err.Error())
		return "", err
	}

	return userCodePath, nil
}

// 2. 编译代码
func CompileFile(userCodePath string) error {
	userCodeParentPath := filepath.Dir(userCodePath)
	// fmt.Println("userCodeParentPath=", userCodeParentPath)

	// 构建编译命令
	var compileCmd string
	if runtime.GOOS == "windows" {
		compileCmd = fmt.Sprintf("go build -o %s\\%s.exe %s", userCodeParentPath, GLOBAL_GO_BINARY_NAME, userCodePath)
	} else {
		compileCmd = fmt.Sprintf("go build -o %s/%s %s", userCodeParentPath, GLOBAL_GO_BINARY_NAME, userCodePath)
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
	case <-time.After(TIME_OUT):
		// 超时
		compileProcess.Process.Kill()
		return errors.New(COMPILE_TIMEOUT_ERROR)
	case err := <-done:
		if err != nil {
			mylog.Log.Errorf("%v : err= %v", COMPILE_ERROR, err.Error())
			return err
		}
	}
	return nil
}

// 3. 运行编译后的可执行文件, 获得执行结果列表
func RunFile(userCodePath string, inputList []string) ([]model.ExecuteMessage, error) {
	userCodeParentPath := filepath.Dir(userCodePath)

	// 运行编译后的可执行文件
	var runCmd string
	if runtime.GOOS == "windows" {
		runCmd = fmt.Sprintf("%s\\%s.exe", userCodeParentPath, GLOBAL_GO_BINARY_NAME)
	} else {
		runCmd = fmt.Sprintf("%s/%s", userCodeParentPath, GLOBAL_GO_BINARY_NAME)
	}

	runCmd = strings.TrimSpace(runCmd)

	var executeMessageList []model.ExecuteMessage

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
				executeMessageList = append(executeMessageList, model.ExecuteMessage{
					ErrorMessage: RUN_TIMEOUT_ERROR,
					Time:         latencyTm,
				})
				return executeMessageList, err
			} else {
				executeMessageList = append(executeMessageList, model.ExecuteMessage{
					ErrorMessage: err.Error(),
					Time:         latencyTm,
				})
				return executeMessageList, err
			}
		} else {
			executeMessageList = append(executeMessageList, model.ExecuteMessage{
				Message: string(output),
				Time:    latencyTm,
			})
		}
	}
	return executeMessageList, nil
}

// 4. 获取输出结果
func GetOutputResponse(executeMessageList []model.ExecuteMessage) model.ExecuteCodeResponse {
	var executeCodeResponse model.ExecuteCodeResponse
	var outputList []string
	// 取用时最大值，便于判断是否超时
	var maxTime int64

	for _, executeMessage := range executeMessageList {
		errorMessage := executeMessage.ErrorMessage
		if utils.IsNotBlank(errorMessage) {
			executeCodeResponse.Message = errorMessage
			// 用户提交的代码执行中存在错误
			executeCodeResponse.Status = EXECUTION_RUNTIME_FAIL
			break
		}
		outputList = append(outputList, executeMessage.Message)
		maxTime = int64(math.Max(float64(maxTime), float64(executeMessage.Time)))
	}

	// 正常运行完成
	if utils.CheckSame[int]("判断outputList和executeMessageList长度一致", len(outputList), len(executeMessageList)) {
		executeCodeResponse.Status = EXECUTION_SUCCESS
	}

	executeCodeResponse.OutputList = outputList

	// 要借助第三方库来获取内存占用，非常麻烦，此处不做实现
	executeCodeResponse.JudgeInfo = model.JudgeInfo{
		Time: maxTime,
	}

	return executeCodeResponse
}

// 5. 删除文件
func DeleteFile(userCodePath string) error {
	return os.RemoveAll(filepath.Dir(userCodePath))
}

// 6. 获取错误响应
func GetErrorResponse() {}
