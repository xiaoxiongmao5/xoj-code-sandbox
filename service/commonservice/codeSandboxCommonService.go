/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-08 11:33:07
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-18 00:35:59
 * @FilePath: /xoj-code-sandbox/service/codeSandboxCommon.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package commonservice

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/model"
	codeexecstatusenum "github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/model/enums/CodeExecStatusEnum"
	"github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/mylog"
	"github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/utils"
)

const (
	GLOBAL_CODE_DIR_NAME  = "tmpcode"
	GLOBAL_GO_FILE_NAME   = "main.go"
	GLOBAL_GO_BINARY_NAME = "main"
	TIME_OUT              = 10 * time.Second
	MEMORY_LIMIT          = 10 * 1024 * 1024 //内存限制（字节）(docker容器要求最低位6M)
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

// 4. 获取输出结果
func GetOutputResponse(execResultList []model.ExecResult) model.ExecuteCodeResponse {
	var executeCodeResponse model.ExecuteCodeResponse
	if utils.IsEmpty(execResultList) {
		return executeCodeResponse
	}
	var outputList []string
	// 取用时最大值，便于判断是否超时
	var maxTime int64
	var maxMemory int64

	for _, execResult := range execResultList {
		outputList = append(outputList, execResult.StdOut)
		maxTime = int64(math.Max(float64(maxTime), float64(execResult.Time)))
		maxMemory = int64(math.Max(float64(maxMemory), float64(execResult.Memory)))

		stdErr := execResult.StdErr
		if utils.IsNotBlank(stdErr) {
			executeCodeResponse.Message = codeexecstatusenum.RUN_FAIL.GetText() + " : " + stdErr
			// 用户提交的代码执行中存在错误
			executeCodeResponse.Status = codeexecstatusenum.RUN_FAIL.GetValue()
			break
		}
	}

	// 正常运行完成
	if utils.CheckSame[int]("判断outputList和executeMessageList长度一致", len(outputList), len(execResultList)) {
		executeCodeResponse.Message = codeexecstatusenum.SUCCEED.GetText()
		executeCodeResponse.Status = codeexecstatusenum.SUCCEED.GetValue()
	}

	executeCodeResponse.OutputList = outputList

	executeCodeResponse.JudgeInfo = model.JudgeInfo{
		Time:   maxTime,
		Memory: maxMemory,
	}

	return executeCodeResponse
}

// 5. 删除文件
func DeleteFile(userCodePath string) error {
	return os.RemoveAll(filepath.Dir(userCodePath))
}
