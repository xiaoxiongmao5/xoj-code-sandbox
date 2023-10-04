/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-04 20:03:09
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-04 22:02:26
 */
package mydocker

import (
	"context"
	"io"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
	"github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/mylog"
)

var Cli *client.Client

func init() {
	// fmt.Println("init begin: mydocker")

	// fmt.Println("init begin: mydocker")
}

// 创建Docker客户端
func CreateDockerClient() (*client.Client, error) {
	// 使用 client.NewClientWithOpts 函数来创建一个Docker客户端。
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return cli, nil
}

// 下载golang镜像
func PullGolangImage(cli *client.Client) error {
	ctx := context.Background()
	reader, err := cli.ImagePull(ctx, "golang:1.20.8-alpine", types.ImagePullOptions{})
	if err != nil {
		mylog.Log.Info("下载golang镜像失败, err=", err.Error())
		return err
	}
	io.Copy(os.Stdout, reader)
	return nil
}

// 在Docker容器内执行命令：使用 cli.ContainerExecCreate 和 cli.ContainerExecStart 函数在Docker容器内执行编译和运行Go代码的命令。
// command := "go build -o main /app/main.go"
func ExecuteInContainer(cli *client.Client, containerID string, command string) ([]byte, error) {
	ctx := context.Background()

	command = strings.TrimSpace(command)
	commandParts := strings.Split(command, " ")

	// 创建一个在容器内执行的命令
	timeExecCreateStart := time.Now()
	createResp, err := cli.ContainerExecCreate(
		ctx,
		containerID,
		types.ExecConfig{
			AttachStdout: true,
			AttachStderr: true,
			AttachStdin:  true,
			WorkingDir:   "/app", // 设置工作目录，可以根据需要修改
			Cmd:          commandParts,
		},
	)
	timeExecCreate := time.Since(timeExecCreateStart).Milliseconds()
	if err != nil {
		mylog.Log.Error("创建在容器内执行的命令失败, err=", err.Error())
		return nil, err
	}

	// 等待命令执行完成并获取输出
	timeExecAttachStart := time.Now()
	execResp, err := cli.ContainerExecAttach(ctx, createResp.ID, types.ExecStartCheck{})
	timeExecAttach := time.Since(timeExecAttachStart).Milliseconds()
	if err != nil {
		mylog.Log.Error("等待命令执行完成并获取输出失败, err=", err.Error())
		return nil, err
	}
	defer execResp.Close()

	timeIOReadStart := time.Now()
	output, err := io.ReadAll(execResp.Reader)
	timeIORead := time.Since(timeIOReadStart).Milliseconds()
	if err != nil {
		mylog.Log.Error("获取输出中的错误内容err=", err.Error())
		return nil, err
	}

	mylog.Log.WithFields(logrus.Fields{
		"timeExecCreate": timeExecCreate,
		"timeExecAttach": timeExecAttach,
		"timeIORead":     timeIORead,
	}).Info("Docker容器内执行命令-耗时统计")

	return output, nil

	// timeGetContainerLogsStart := time.Now()
	// var stdout bytes.Buffer
	// var stderr bytes.Buffer
	// out, err := cli.ContainerLogs(ctx, containerID, types.ContainerLogsOptions{
	// 	ShowStdout: true, ShowStderr: true, Timestamps: true,
	// })
	// timeGetContainerLogs := time.Since(timeGetContainerLogsStart).Milliseconds()
	// if err != nil {
	// 	mylog.Log.Error("cli.ContainerLogs 错误, err=", err.Error())
	// 	return nil, err
	// }
	// timeStdCopyStart := time.Now()
	// stdcopy.StdCopy(&stdout, &stderr, out)
	// timeStdCopy := time.Since(timeStdCopyStart).Milliseconds()

	// mylog.Log.WithFields(logrus.Fields{
	// 	"timeExecCreate": timeExecCreate,
	// 	"timeExecAttach": timeExecAttach,
	// 	"timeGetContainerLogs": timeGetContainerLogs,
	// 	"timeStdCopy":          timeStdCopy,
	// }).Info("Docker容器内执行命令-耗时统计")

	// return []byte(stdout.String()), errors.New(stderr.String())
}
