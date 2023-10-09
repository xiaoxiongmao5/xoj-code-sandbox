/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-02 12:27:09
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-09 17:10:37
 */
package model

type ExecuteCodeResponse struct {
	OutputList []string `json:"outputList"`
	// 接口信息(对应Status的信息描述，1:正常运行完成, 2.错误输出, 3. 错误输出)
	Message string `json:"message"`
	// 执行状态(1:正常运行完成 2:编译失败 3:用户提交的代码运行有错误 4:系统错误)
	Status int32 `json:"status"`
	// 判题信息
	JudgeInfo JudgeInfo `json:"judgeInfo"`
}
