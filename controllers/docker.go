/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-04 16:33:26
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-07 23:33:31
 */
package controllers

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/model"
	"github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/mydocker"
	"github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/mylog"
	"github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/utils"
)

type DockerController struct {
	beego.Controller
}

func (this DockerController) Do() {
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
		fmt.Printf("%d", a+b)
	}
	`
	inputList := []string{"2 2", "2 3"}

	// 1. 保存用户代码到文件中
	userCodePath, err := SaveCodeToFile(code)
	if err != nil {
		mylog.Log.Error("保存代码失败, err=", err.Error())
		return
	}

	ctx := context.Background()
	cli := mydocker.Cli

	userCodeParentPath := filepath.Dir(userCodePath)

	// 2. 编译代码
	if err := CompileFileInDocker(ctx, cli, userCodeParentPath); err != nil {
		mylog.Log.Error("编译代码失败, err=", err.Error())
		return
	}

	// 3. 运行代码
	executeMessageList, err := RunFileInDocker(ctx, cli, userCodeParentPath, inputList)
	if err != nil {
		mylog.Log.Error("运行可执行文件失败, err=", err.Error())
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

}

// 编译代码
func CompileFileInDocker(ctx context.Context, cli *client.Client, userCodeParentPath string) error {
	containerID, err := CreateContainerCfgOfDefault(ctx, cli, "golang:1.20.8-alpine", userCodeParentPath)
	if err != nil {
		return err
	}
	defer mydocker.StopAndRemoveContainer(ctx, cli, containerID)

	command := "go build -o main /app/main.go"

	execResult, err := mydocker.ExecuteInContainer(ctx, cli, containerID, command)
	if err != nil {
		mylog.Log.Error("编译失败, err=", err.Error())
		return err
	}

	// 获取内存
	memory, err := mydocker.GetContainerMemoryUsage(ctx, cli, containerID)
	if err != nil {
		mylog.Log.Error("获取编译内存失败, err=", err)
	}
	mylog.Log.WithFields(logrus.Fields{
		"内存":     memory,
		"耗时":     execResult.Tm,
		"StdOut": execResult.StdOut,
		"StdErr": execResult.StdErr,
	}).Info("编译-资源消耗统计")
	return nil
}

// 运行代码
func RunFileInDocker(ctx context.Context, cli *client.Client, userCodeParentPath string, inputList []string) (executeMessageList []model.ExecuteMessage, err error) {
	containerID, err := CreateContainerCfgOfRunExec(ctx, cli, "alpine:latest", userCodeParentPath)
	if err != nil {
		return
	}
	defer mydocker.StopAndRemoveContainer(ctx, cli, containerID)

	for i, input := range inputList {

		command := "./main " + strings.TrimSpace(input)

		// 根据每条输入用例，运行代码，拿到输出结果
		execResult, err := mydocker.ExecuteInContainer(ctx, cli, containerID, command)
		if err != nil {
			mylog.Log.Errorf("运行用户代码,输入示例[%d]失败,err=%s", i, err.Error())
			return executeMessageList, err
		}
		mylog.Log.WithFields(logrus.Fields{
			"耗时":     execResult.Tm,
			"StdOut": execResult.StdOut,
			"StdErr": execResult.StdErr,
		}).Infof("运行用例[%v]-资源和输出-统计", i)

		// 运行用户代码，存在错误输出
		if utils.IsNotBlank(execResult.StdErr) {
			errMsg := fmt.Sprintf("运行用户代码,输入示例[%d]失败,存在错误输出,StdErr=%s", i, execResult.StdErr)
			mylog.Log.Error(errMsg)
			return executeMessageList, errors.New(errMsg)
		}

		// 获取内存
		memory, err := mydocker.GetContainerMemoryUsage(ctx, cli, containerID)
		if err != nil {
			mylog.Log.Errorf("运行用户代码,输入示例[%d],获取内存失败,err=%s", i, err.Error())
			return executeMessageList, err
		}
		mylog.Log.Infof("运行用例[%v]-资源消耗-内存[%v]", i, memory)

		executeMessageList = append(executeMessageList, model.ExecuteMessage{
			Message: execResult.StdOut,
			Time:    execResult.Tm,
			Memory:  int64(memory),
		})
	}
	return executeMessageList, nil
}

// 创建执行用户代码的容器（有资源限制：内存、运行时间、CPU）
func CreateContainerCfgOfRunExec(ctx context.Context, cli *client.Client, image string, userCodeParentPath string) (containerID string, err error) {
	containerConfig := &container.Config{
		Image:        image,
		Tty:          true,
		AttachStdout: true,
		AttachStderr: true,
		WorkingDir:   "/app", // 设置工作目录，可以根据需要修改
	}

	// 设置容器的卷，用于共享代码和编译结果
	codeVolume := fmt.Sprintf("%s:/app", userCodeParentPath)

	containerHostConfig := &container.HostConfig{
		Binds: []string{codeVolume}, //此容器的卷绑定列表
		Resources: container.Resources{
			Memory: 7 * 1024 * 1024, //内存限制（字节）
			// CPUShares: 1,                 //CPU份额（相对于其他容器的相对重量）
		},
	}

	// 创建一个随机的容器名
	containerName := fmt.Sprintf("my-container-runexec-%s", uuid.New().String())

	return mydocker.CreateAndStartContainer(ctx, cli, containerConfig, containerHostConfig, nil, nil, containerName)
}

// 创建默认配置的容器
func CreateContainerCfgOfDefault(ctx context.Context, cli *client.Client, image string, userCodeParentPath string) (containerID string, err error) {

	// 设置容器的配置
	containerConfig := &container.Config{
		Image:        image,
		Tty:          true,
		AttachStdout: true,
		AttachStderr: true,
		WorkingDir:   "/app", // 设置工作目录，可以根据需要修改
	}

	// 设置容器的卷，用于共享代码和编译结果
	codeVolume := fmt.Sprintf("%s:/app", userCodeParentPath)

	containerHostConfig := &container.HostConfig{
		Binds:     []string{codeVolume}, //此容器的卷绑定列表
		Resources: container.Resources{
			// CPUShares: 1, //CPU份额（相对于其他容器的相对重量）
		},
	}

	// 创建一个随机的容器名
	containerName := fmt.Sprintf("my-container-default-%s", uuid.New().String())

	return mydocker.CreateAndStartContainer(ctx, cli, containerConfig, containerHostConfig, nil, nil, containerName)
}
