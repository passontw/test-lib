package gamelogic

import "sl.framework.com/trace"

// IEvent 事件接口
type IEvent interface {
	//HandleEvent 创建事件后调用该函数处理事件
	HandleEvent()
}

// CEvent 事件父类,非必须组合该类
type CEvent struct {
}

// RedisEvent 设置或者更新redis事件
type RedisEvent struct {
	traceId string
}

func (r *RedisEvent) HandleEvent() {
	trace.Info("RedisEvent HandleEvent function")
}
