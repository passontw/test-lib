package log

// NoOpWriter 是一个不执行任何操作的 io.Writer 实现。
type NoOpWriter struct{}

// Write 方法满足 io.Writer 接口，但不执行任何实际写入操作。
func (w *NoOpWriter) Write(p []byte) (n int, err error) {
	// 返回输入长度和 nil 错误，表示成功处理了所有输入，但实际上没有做任何事。
	return len(p), nil
}
