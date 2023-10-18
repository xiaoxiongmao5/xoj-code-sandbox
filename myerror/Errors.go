/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-15 21:29:16
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-17 13:01:19
 */
package myerror

var (
// ErrInTimeOut  = fmt.Errorf("timeout: %w", errors.New("超时错误"))
)

// 超时错误
type ErrTimeOut struct {
	Message string
}

func (this ErrTimeOut) Error() string {
	return this.Message
}

// 内存溢出错误
type ErrMemoryFullOut struct {
	Message string
}

func (this ErrMemoryFullOut) Error() string {
	return this.Message
}
