package model

var (
// ErrInTimeOut  = fmt.Errorf("timeout: %w", errors.New("超时错误"))
// MemoryFullOut = errors.New("内存溢出错误")
)

type ErrTimeOut struct {
	Message string
}

func (this ErrTimeOut) Error() string {
	return this.Message
}

// func ErrTimeOutFn() error {
// 	return ErrTimeOut{"超时"}
// }
