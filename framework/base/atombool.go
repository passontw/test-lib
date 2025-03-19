package base

import "sync/atomic"

// 封装bool类型的原子变量
type AtomBool struct {
	val int32
}

func (a *AtomBool) Set(flag bool) {
	var i int32 = 0
	if flag {
		i = 1
	}
	atomic.StoreInt32(&a.val, i)
}

func (a *AtomBool) Get() bool {
	i := atomic.LoadInt32(&a.val)
	if 0 == i {
		return false
	}
	return true
}
