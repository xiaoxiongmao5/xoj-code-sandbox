/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-04 16:33:26
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-04 20:38:03
 */
package controllers

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/google/uuid"
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
		fmt.Printf("结果: %d", a+b)
	}
	`
	inputList := []string{"2 2", "2 3"}

	// 1. 保存用户代码到文件中
	userCodePath, err := SaveCodeToFile(code)
	if err != nil {
		mylog.Log.Error("保存代码失败, err=", err.Error())
		return
	}

	cli := mydocker.Cli

	userCodeParentPath := filepath.Dir(userCodePath)

	// 创建并启动Docker容器
	containerID, err := CreateAndStartContainer(cli, userCodeParentPath)
	if err != nil {
		mylog.Log.Error("创建并启动Docker容器失败, err=", err.Error())
		return
	}

	// 编译代码
	if err := CompileFileInDocker(cli, containerID); err != nil {
		mylog.Log.Error("编译代码失败, err=", err.Error())
		return
	}

	// 运行代码
	executeMessageList, err := RunFileInDocker(cli, containerID, inputList)
	if err != nil {
		mylog.Log.Error("运行可执行文件失败, err=", err.Error())
		return
	}

	// 4. 整理结果
	executeCodeResponse := GetOutputResponse(executeMessageList)

	// 5. 清理容器和文件
	if err := CleanContainerAndDeleteFile(cli, containerID, userCodePath); err != nil {
		mylog.Log.Error("清理容器和文件错误, err=", err.Error())
		return
	}

	this.Data["json"] = executeCodeResponse
	this.ServeJSON()

}

// 创建并启动Docker容器：使用Docker客户端创建并启动一个Docker容器，该容器用于编译和运行Go代码。需要指定容器的镜像、卷、工作目录等信息。
func CreateAndStartContainer(cli *client.Client, userCodeParentPath string) (string, error) {
	ctx := context.Background()
	// 创建一个随机的容器名
	containerName := fmt.Sprintf("my-go-container-%s", uuid.New().String())

	// 设置容器的配置
	containerConfig := &container.Config{
		Image:        "golang:1.20.8-alpine", // Go编译环境的镜像
		Tty:          true,
		AttachStdout: true,
		AttachStderr: true,
		WorkingDir:   "/app", // 设置工作目录，可以根据需要修改
	}

	// 设置容器的卷，用于共享代码和编译结果
	codeVolume := fmt.Sprintf("%s:/app", userCodeParentPath)

	containerHostConfig := &container.HostConfig{
		Binds:   []string{codeVolume},
		ShmSize: 1000,
		Runtime: "",
	}

	// 创建容器
	resp, err := cli.ContainerCreate(
		ctx,
		containerConfig,
		containerHostConfig,
		nil,
		nil,
		containerName,
	)
	if err != nil {
		mylog.Log.Error("创建容器失败, err=", err.Error())
		// panic(err)
		return "", err
	}

	// 启动容器
	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		mylog.Log.Error("启动容器失败, err=", err.Error())
		return "", err
	}

	return resp.ID, nil
}

// 编译代码
func CompileFileInDocker(cli *client.Client, containerID string) error {
	command := "go build -o main /app/main.go"
	output, err := mydocker.ExecuteInContainer(cli, containerID, command)
	if err != nil {
		mylog.Log.Error("编译失败, err=", err.Error())
		return err
	}
	mylog.Log.Info("编译成功, output=", string(output))
	return nil
}

// 运行代码
func RunFileInDocker(cli *client.Client, containerID string, inputList []string) ([]model.ExecuteMessage, error) {
	var executeMessageList []model.ExecuteMessage
	for _, input := range inputList {
		command := "./main " + strings.TrimSpace(input)

		startTime := time.Now()

		output, err := mydocker.ExecuteInContainer(cli, containerID, command)

		latencyTm := time.Since(startTime).Milliseconds()

		if err != nil {
			executeMessageList = append(executeMessageList, model.ExecuteMessage{
				ErrorMessage: err.Error(),
				Time:         latencyTm,
			})
			return executeMessageList, err
		}
		executeMessageList = append(executeMessageList, model.ExecuteMessage{
			Message: string(output),
			Time:    latencyTm,
		})
	}
	return executeMessageList, nil
}

// 清理容器：在使用完Docker容器后，清理容器，以释放资源。
func CleanContainerAndDeleteFile(cli *client.Client, containerID, userCodePath string) error {
	ctx := context.Background()
	errMsg := ""
	// 停止容器
	if err := cli.ContainerStop(ctx, containerID, container.StopOptions{}); err != nil {
		errMsg += err.Error()
	}

	// 删除容器
	if err := cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{
		RemoveVolumes: true, Force: true,
	}); err != nil {
		errMsg += err.Error()
	}

	// 删除宿主机文件
	if err := os.RemoveAll(filepath.Dir(userCodePath)); err != nil {
		errMsg += err.Error()
	}
	if utils.IsEmpty(errMsg) {
		return nil
	}
	return errors.New(errMsg)
}
