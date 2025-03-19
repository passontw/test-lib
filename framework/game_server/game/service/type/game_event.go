package types

import "sl.framework.com/trace"

type EventDTO struct {
	GameRoomId      int64 // 房间ID
	GameRoundId     int64 // 当前局信息
	GameId          int64
	NextGameRoundId int64 // 下一局局信息
	GameRoundNo     string
	Command         string
	Time            int64 //时间，采用时间戳毫秒级
	ReceiveTime     int64 //接收到时间的时间 采用时间戳毫秒级
	Payload         any
}
type EventBase struct {
	Dto            *EventDTO
	RoundDTO       *GameRoundDTO
	TraceId        string // traceId 跟踪流程
	RequestId      string //请求id
	RetHandleEvent *int   // 消息处理返回值
	MsgHeader      string //具体类对应的消息 用于打印日志
}

/**
 * HandleEvent
 * 处理游戏事件函数
 *
 * @param traceId string - 跟踪id
 * @return RETURN
 */

func (e *EventBase) HandleRondEvent() {
	trace.Info("HandleEvent traceId:%s", e.TraceId)
}
