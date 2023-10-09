/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-08 11:37:18
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-09 17:12:16
 */
package service

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/model"
	"github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/mydocker"
	"github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/mylog"
	commonservice "github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/service/commonService"
)

const (
	COMPILE_GO_IMAGE   = "golang:1.20.8-alpine"
	COMPILE_COMMAND    = "go build -o main /app/main.go"
	RUN_GO_IMAGE       = "alpine:latest"
	RUN_COMMAND_PREFIX = "./main "
	CONTAINER_WORK_DIR = "/app" //容器的工作目录
)

type GoCodeSandboxByDocker struct {
	Ctx context.Context
	Cli *client.Client
}

// 1. 把用户的代码保存为文件
func (this GoCodeSandboxByDocker) SaveCodeToFile(code string) (string, error) {
	return commonservice.SaveCodeToFile(code)
}

// 2. 编译代码
func (this GoCodeSandboxByDocker) CompileFile(userCodePath string) error {
	userCodeParentPath := filepath.Dir(userCodePath)

	// 创建并启动容器
	containerID, err := this.CreateContainerCfgOfDefault(this.Cli, COMPILE_GO_IMAGE, userCodeParentPath)
	if err != nil {
		return err
	}
	// defer 关闭并删除容器
	defer mydocker.StopAndRemoveContainer(this.Ctx, this.Cli, containerID)

	command := COMPILE_COMMAND

	// 编译代码
	execResult, err := mydocker.ExecuteInContainer(this.Ctx, this.Cli, containerID, command)
	if err != nil {
		mylog.Log.Error("编译失败, err=", err.Error())
		return err
	}

	mylog.Log.WithFields(logrus.Fields{
		"内存":     execResult.Memory,
		"耗时":     execResult.Time,
		"StdOut": execResult.StdOut,
		"StdErr": execResult.StdErr,
	}).Info("编译-资源消耗统计")

	return nil
}

// 3. 运行编译后的可执行文件, 获得执行结果列表
func (this GoCodeSandboxByDocker) RunFile(userCodePath string, inputList []string) (execResultList []model.ExecResult, err error) {
	userCodeParentPath := filepath.Dir(userCodePath)

	// 创建并启动容器
	containerID, err := this.CreateContainerCfgOfRunExec(this.Cli, RUN_GO_IMAGE, userCodeParentPath)
	if err != nil {
		return execResultList, err
	}
	// defer 关闭并删除容器
	defer mydocker.StopAndRemoveContainer(this.Ctx, this.Cli, containerID)

	for i, input := range inputList {

		command := RUN_COMMAND_PREFIX + strings.TrimSpace(input)

		// 根据每条输入用例，运行代码
		execResult, err := mydocker.ExecuteInContainer(this.Ctx, this.Cli, containerID, command)
		if err != nil {
			mylog.Log.Errorf("运行用户代码,输入示例[%d]失败,err=%s", i, err.Error())
			return execResultList, err
		}
		mylog.Log.WithFields(logrus.Fields{
			"耗时":     execResult.Time,
			"内存":     execResult.Memory,
			"StdOut": execResult.StdOut,
			"StdErr": execResult.StdErr,
		}).Infof("运行用例[%v]-资源和输出-统计", i)

		// // 运行用户代码，存在错误输出
		// if utils.IsNotBlank(execResult.StdErr) {
		// 	errMsg := fmt.Sprintf("运行用户代码,输入示例[%d]失败,错误输出[StdErr=%s]", i, execResult.StdErr)
		// 	mylog.Log.Error(errMsg)
		// 	return execResultList, errors.New(errMsg)
		// }

		execResultList = append(execResultList, execResult)
	}

	return execResultList, nil
}

// 4. 获取输出结果
func (this GoCodeSandboxByDocker) GetOutputResponse(execResultList []model.ExecResult) model.ExecuteCodeResponse {
	return commonservice.GetOutputResponse(execResultList)
}

// 5. 删除文件
func (this GoCodeSandboxByDocker) DeleteFile(userCodePath string) error {
	return commonservice.DeleteFile(userCodePath)
}

// 创建执行用户代码的容器（有资源限制：内存、运行时间、CPU）
func (this GoCodeSandboxByDocker) CreateContainerCfgOfRunExec(cli *client.Client, image string, userCodeParentPath string) (containerID string, err error) {
	containerConfig := &container.Config{
		Image:        image,
		Tty:          true,
		AttachStdout: true,
		AttachStderr: true,
		WorkingDir:   CONTAINER_WORK_DIR, // 设置工作目录，可以根据需要修改
	}

	// 设置容器的卷，用于共享代码和编译结果
	codeVolume := fmt.Sprintf("%s:%s", userCodeParentPath, CONTAINER_WORK_DIR)

	containerHostConfig := &container.HostConfig{
		Binds: []string{codeVolume}, //此容器的卷绑定列表
		Resources: container.Resources{
			Memory: commonservice.MEMORY_LIMIT, //内存限制（字节）
			// CPUShares: 1,                 //CPU份额（相对于其他容器的相对重量）
		},
	}

	// 创建一个随机的容器名
	containerName := fmt.Sprintf("my-container-runexec-%s", uuid.New().String())

	return mydocker.CreateAndStartContainer(this.Ctx, cli, containerConfig, containerHostConfig, nil, nil, containerName)
}

// 创建默认配置的容器
func (this GoCodeSandboxByDocker) CreateContainerCfgOfDefault(cli *client.Client, image string, userCodeParentPath string) (containerID string, err error) {

	// 设置容器的配置
	containerConfig := &container.Config{
		Image:        image,
		Tty:          true,
		AttachStdout: true,
		AttachStderr: true,
		WorkingDir:   CONTAINER_WORK_DIR, // 设置工作目录，可以根据需要修改
	}

	// 设置容器的卷，用于共享代码和编译结果
	codeVolume := fmt.Sprintf("%s:%s", userCodeParentPath, CONTAINER_WORK_DIR)

	containerHostConfig := &container.HostConfig{
		Binds:     []string{codeVolume}, //此容器的卷绑定列表
		Resources: container.Resources{
			// CPUShares: 1, //CPU份额（相对于其他容器的相对重量）
		},
	}

	// 创建一个随机的容器名
	containerName := fmt.Sprintf("my-container-default-%s", uuid.New().String())

	return mydocker.CreateAndStartContainer(this.Ctx, cli, containerConfig, containerHostConfig, nil, nil, containerName)
}
