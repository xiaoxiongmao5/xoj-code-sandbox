/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-02 12:27:09
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-04 11:17:33
 */
package model

type ExecuteCodeResponse struct {
	OutputList []string
	// 接口信息
	Message string `json:"message"`
	// 执行状态
	Status int32 `json:"status"`
	// 判题信息
	JudgeInfo JudgeInfo `json:"judgeInfo"`
}
