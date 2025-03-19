package async

// AsyncRunCoroutine 启动一个协程,执行无参数函数
func AsyncRunCoroutine(cb func()) {
	go func() {
		defer TryException()
		cb()
	}()
}

type FuncWithAny[T any] func(in T)

// AsyncRunWithAny 启动一个协程,执行有一个泛型参数的函数
func AsyncRunWithAny[T any](cb FuncWithAny[T], in T) {
	go func(in T) {
		defer TryException()
		cb(in)
	}(in)
}

type FuncWithAnyMultiParam[T1 any, T2 any] func(int1 T1, in2 T2)

// AsyncRunWithAnyMulti 启动一个协程,执行有两个泛型参数的函数
func AsyncRunWithAnyMulti[T1 any, T2 any](cb FuncWithAnyMultiParam[T1, T2], in1 T1, in2 T2) {
	go func(in1 T1, in2 T2) {
		defer TryException()
		cb(in1, in2)
	}(in1, in2)
}
