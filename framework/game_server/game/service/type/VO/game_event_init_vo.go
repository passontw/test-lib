package VO

// 创建游戏事件使用的赋值对象
type GameEventInitVO struct {
	TraceId     string `json:"traceId"`
	RoomId      int64  `json:"room_id"`
	RoundId     int64  `json:"round_id"`
	NextRoundId int64  `json:"next_round_id"`
	RequestId   string `json:"request_id"`
	Code        *int   `json:"code"`
	Time        int64  `json:"time"`         //时间时间
	ReceiveTime int64  `json:"receive_time"` //时间接收时间
}
