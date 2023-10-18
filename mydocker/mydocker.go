/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-04 20:03:09
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-18 00:11:06
 */
package mydocker

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/sirupsen/logrus"
	"github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/model"
	"github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/myerror"
	"github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/mylog"
	commonservice "github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/service/commonService"
	"github.com/xiaoxiongmao5/xoj/xoj-code-sandbox/utils"
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

// 下载镜像
func PullGolangImage(ctx context.Context, cli *client.Client, image string) error {
	var tmPull int64
	tmPullStart := time.Now()
	reader, err := cli.ImagePull(ctx, image, types.ImagePullOptions{})
	tmPull = time.Since(tmPullStart).Milliseconds()
	if err != nil {
		mylog.Log.Infof("下载[%s]镜像失败, err=%v", image, err.Error())
		return err
	}
	io.Copy(os.Stdout, reader)

	defer mylog.Log.WithFields(logrus.Fields{
		"下载镜像耗时": tmPull,
		"镜像":     image,
	}).Info("Docker-下载镜像-统计")
	return nil
}

// 创建并启动Docker容器：使用Docker客户端创建并启动一个Docker容器。需要指定容器的镜像、卷、工作目录等信息。
func CreateAndStartContainer(ctx context.Context, cli *client.Client, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *v1.Platform, containerName string) (containerID string, err error) {
	var tmCreate, tmStart int64

	// 创建容器
	tmCreateStart := time.Now()
	resp, err := cli.ContainerCreate(
		ctx,
		config,
		hostConfig,
		networkingConfig,
		platform,
		containerName,
	)

	tmCreate = time.Since(tmCreateStart).Milliseconds()
	if err != nil {
		mylog.Log.Error("创建容器失败, err=", err.Error())
		return "", err
	}
	containerID = resp.ID

	// 启动容器
	tmStartStart := time.Now()
	err = cli.ContainerStart(ctx, containerID, types.ContainerStartOptions{})
	tmStart = time.Since(tmStartStart).Milliseconds()
	if err != nil {
		mylog.Log.Error("启动容器失败, err=", err.Error())
		return "", err
	}

	defer mylog.Log.WithFields(logrus.Fields{
		"创建容器耗时": tmCreate,
		"启动容器耗时": tmStart,
		"容器ID":   containerID,
	}).Info("Docker-创建启动-统计")

	return containerID, nil
}

// 在Docker容器内执行命令：使用 cli.ContainerExecCreate 和 cli.ContainerExecStart 函数在Docker容器内执行编译和运行Go代码的命令。
// command := "go build -o main /app/main.go"
func ExecuteInContainer(ctx context.Context, cli *client.Client, containerID string, command string) (execResult model.ExecResult, err error) {
	var tmExecCreate, tmExecAttach, tmIORead, tmGetExecInspect int64

	// 创建一个通道用于共享数据
	memMaxChannel := make(chan uint64)

	defer func() {
		memory := <-memMaxChannel
		mylog.Log.Infof("====================在defer中，获取到[%s]的memMaxChannel=%v", containerID, memory)
		execResult.Memory = int64(math.Ceil(float64(memory))) //向上取整
	}()

	// 创建一个 done 通道
	doneOfWatchMemory := make(chan bool)
	defer close(doneOfWatchMemory)

	command = strings.TrimSpace(command)
	commandParts := strings.Split(command, " ")

	// 创建一个在容器内执行的命令
	tmExecCreateStart := time.Now()
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
	tmExecCreate = time.Since(tmExecCreateStart).Milliseconds()
	if err != nil {
		mylog.Log.Error("创建在容器内执行的命令失败, err=", err.Error())
		return execResult, err
	}

	// 等待命令执行完成并获取输出
	tmExecAttachStart := time.Now()
	resp, err := cli.ContainerExecAttach(ctx, createResp.ID, types.ExecStartCheck{})
	tmExecAttach = time.Since(tmExecAttachStart).Milliseconds()
	if err != nil {
		mylog.Log.Error("等待命令执行完成并获取输出失败, err=", err.Error())
		return execResult, err
	}
	defer resp.Close()

	// 开始监视容器内存
	go func(ctx context.Context, cli *client.Client, containerID string, doneOfWatchMemory chan bool) {
		var maxMemory uint64
		for {
			// 检查是否需要停止协程
			select {
			case <-doneOfWatchMemory:
				memMaxChannel <- maxMemory
				return
			default:
			}

			// 获取内存
			memory, err := GetContainerMemoryUsage(ctx, cli, containerID)
			if err != nil {
				mylog.Log.Error("获取内存失败, err=", err)
			} else {
				if memory > maxMemory {
					maxMemory = memory
				}
			}
			mylog.Log.Info("监视容器内存的协程, memory=", memory)
			// time.Sleep(10 * time.Millisecond)
		}
	}(ctx, cli, containerID, doneOfWatchMemory)

	done := make(chan error, 1)
	go func() {
		tmIOReadStart := time.Now()
		// read the output
		var outBuf, errBuf bytes.Buffer
		_, err = stdcopy.StdCopy(&outBuf, &errBuf, resp.Reader)
		stdout, err := io.ReadAll(&outBuf)
		execResult.StdOut = string(stdout)
		if err != nil {
			mylog.Log.Error("读取[Exec]执行的[outBuf]失败,err=", err.Error())
			done <- err
			return
		}
		stderr, err := io.ReadAll(&errBuf)
		execResult.StdErr = string(stderr)
		if err != nil {
			mylog.Log.Error("读取[Exec]执行的[errBuf]失败,err=", err.Error())
			done <- err
			return
		}
		tmIORead = time.Since(tmIOReadStart).Milliseconds()
		execResult.Time = tmIORead
		done <- nil
	}()

	select {
	case <-time.After(commonservice.TIME_OUT):
		msg := fmt.Sprintf("执行超时，限制为 %v", commonservice.TIME_OUT)
		mylog.Log.Errorf("ExecuteInContainer : %s", msg)
		return execResult, myerror.ErrTimeOut{Message: msg} //exec 超时
	case err = <-done:
		if err != nil {
			return execResult, err
		}
	}

	// 获取exec执行信息
	tmGetExecInspectStart := time.Now()
	containerExecInspect, err := cli.ContainerExecInspect(ctx, createResp.ID)
	tmGetExecInspect = time.Since(tmGetExecInspectStart).Milliseconds()
	if err != nil {
		mylog.Log.Error("获取[Exec]执行信息失败,err=", err.Error())
		return execResult, err
	}
	exitCode := containerExecInspect.ExitCode
	execResult.ExitCode = exitCode
	if !utils.CheckSame[int]("检查[Exec]执行的退出码是否为0", exitCode, 0) {
		mylog.Log.Errorf("运行进程的退出码=[%d], StdErr=[%s], StdOut=[%s]", exitCode, execResult.StdErr, execResult.StdOut)

		if utils.CheckSame[int]("检查[Exec]执行的退出码是否为137", exitCode, 137) {
			execResult.StdErr = "操作系统的内存不足"
			return execResult, myerror.ErrMemoryFullOut{Message: execResult.StdErr}
		}

		return execResult, errors.New(execResult.StdErr)
	}

	// scanner := bufio.NewScanner(resp.Reader)
	// var outputLines []string
	// for scanner.Scan() {
	// 	// 逐行读取输出并将每一行的内容存储在 outputLines 切片中
	// 	line := scanner.Text()
	// 	outputLines = append(outputLines, line)
	// }
	// if err = scanner.Err(); err != nil {
	// 	mylog.Log.Error("获取输出中的错误内容err=", err.Error())
	// 	return
	// }
	// output = strings.Join(outputLines, "\n")
	// outputArr, err := io.ReadAll(resp.Reader)

	defer mylog.Log.WithFields(logrus.Fields{
		"创建exec":          tmExecCreate,
		"运行exec":          tmExecAttach,
		"读取stdout、stderr": tmIORead,
		"获取exec执行信息":      tmGetExecInspect,
	}).Info("Docker-运行命令-耗时统计")

	return execResult, nil
}

// 清理容器：在使用完Docker容器后，清理容器，以释放资源。
func StopAndRemoveContainer(ctx context.Context, cli *client.Client, containerID string) error {
	if cli == nil {
		mylog.Log.Error("StopAndRemoveContainer cli is nil")
		return errors.New("docker cli is nil")
	}

	if utils.IsEmpty(containerID) {
		mylog.Log.Error("StopAndRemoveContainer containerID is empty")
		return errors.New("docker containerID is nil")
	}

	var tmStop, tmRemove int64

	errMsg := ""
	// 停止容器
	tmStopStart := time.Now()
	err := cli.ContainerStop(ctx, containerID, container.StopOptions{})
	tmStop = time.Since(tmStopStart).Milliseconds()
	if err != nil {
		errMsg += err.Error()
	}

	// 删除容器
	tmRemoveStart := time.Now()
	err = cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{
		RemoveVolumes: true, Force: true,
	})
	tmRemove = time.Since(tmRemoveStart).Milliseconds()
	if err != nil {
		errMsg += err.Error()
	}

	defer mylog.Log.WithFields(logrus.Fields{
		"停止容器耗时": tmStop,
		"删除容器耗时": tmRemove,
		"容器ID":   containerID,
	}).Info("Docker-停止删除-统计")

	if utils.IsEmpty(errMsg) {
		return nil
	}

	mylog.Log.Infof("容器[%s]清理失败", containerID)

	return errors.New(errMsg)
}

// 获取内存消耗
func GetContainerMemoryUsage(ctx context.Context, cli *client.Client, containerID string) (uint64, error) {
	if cli == nil {
		mylog.Log.Error("GetContainerMemoryUsage cli is nil")
		return 0, errors.New("docker cli is nil")
	}

	if utils.IsEmpty(containerID) {
		mylog.Log.Error("GetContainerMemoryUsage containerID is empty")
		return 0, errors.New("docker containerID is nil")
	}
	var tmGetStat int64
	var memoryUsage uint64

	// ContainerStats返回给定容器的近乎实时的统计信息。这取决于caller关闭io.ReadCloser返回。
	// 查询容器的统计数据（stats）
	tmGetStatStart := time.Now()
	stats, err := cli.ContainerStats(ctx, containerID, false)
	tmGetStat = time.Since(tmGetStatStart).Milliseconds()
	if err != nil {
		mylog.Log.Error("查询容器的统计数据失败, err=", err.Error())
		return 0.0, err
	}
	defer stats.Body.Close()

	// 解析统计数据中的内存使用信息
	var memUsage types.StatsJSON
	if err := json.NewDecoder(stats.Body).Decode(&memUsage); err != nil {
		mylog.Log.Error("解析统计数据中的内存使用信息失败, err=", err.Error())
		return 0.0, err
	}

	// 获取内存使用信息（以字节为单位）
	memoryUsage = memUsage.MemoryStats.Usage //内存的当前res_counter使用情况
	// memoryUsage = memUsage.MemoryStats.MaxUsage //有记录以来的最大使用量

	// 如果需要，你可以将内存使用信息转换为其他单位，例如MB或GB
	// memoryUsageInKB := memoryUsage / 1024

	defer mylog.Log.WithFields(logrus.Fields{
		"获取容器内存耗时(ms)": tmGetStat,
		"容器内存消耗(byte)": memoryUsage,
		"容器ID":         containerID,
	}).Info("Docker-获取内存-统计")

	return memoryUsage, nil
}
