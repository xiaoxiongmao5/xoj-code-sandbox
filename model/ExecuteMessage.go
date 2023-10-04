/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-04 11:15:40
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-04 11:21:29
 */
package model

// 进程执行信息
type ExecuteMessage struct {
	ExitValue    int    `json:"exitValue"`
	Message      string `json:"message"`
	ErrorMessage string `json:"errorMessage"`
	Time         int64  `json:"time"`
	Memory       int64  `json:"memory"`
}
