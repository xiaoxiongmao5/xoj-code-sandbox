/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-09 20:45:35
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-09 20:59:39
 */
package codeexecstatusenum

type CodeExecStatusEnum int32

func (this CodeExecStatusEnum) GetValue() int32 {
	return int32(this)
}

func (this CodeExecStatusEnum) GetText() string {
	return CodeExecStatusEnumName[this]
}

// 代码沙箱操作代码的状态
const (
	SUCCEED               CodeExecStatusEnum = 1
	COMPILE_FAIL          CodeExecStatusEnum = 2
	COMPILE_TIMEOUT_ERROR CodeExecStatusEnum = 3
	RUN_FAIL              CodeExecStatusEnum = 4
	RUN_TIMEOUT_ERROR     CodeExecStatusEnum = 5
	SYSTEM_ERROR          CodeExecStatusEnum = 6
)

var CodeExecStatusEnumName = map[CodeExecStatusEnum]string{
	SUCCEED:               "正常运行完成",
	COMPILE_FAIL:          "编译失败",
	COMPILE_TIMEOUT_ERROR: "编译超时",
	RUN_FAIL:              "用户提交的代码执行中存在错误",
	RUN_TIMEOUT_ERROR:     "运行超时",
	SYSTEM_ERROR:          "系统错误",
}
