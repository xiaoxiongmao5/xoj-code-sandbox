/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-04 11:15:40
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-08 10:49:04
 */
package model

// 进程执行信息
type ExecResult struct {
	ExitCode int    `json:"exitCode"`
	StdOut   string `json:"stdOut"`
	StdErr   string `json:"stdErr"`
	Time     int64  `json:"time"`
	Memory   int64  `json:"memory"`
}
