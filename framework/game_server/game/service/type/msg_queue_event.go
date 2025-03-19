package types

// GameEventMessageTmp 从数据源接收数据 转为GameEventMessageV2供内部使用
type GameEventMessageTmp struct {
	GameRoomId     int64       `json:"gameRoomId"`  //游戏房间id
	RoundNo        string      `json:"roundNo"`     //游戏局号，不可以超过32长度字符串
	NextRoundNo    string      `json:"nextRoundNo"` //下一局局号,需要当前局给出下一局局号做预热数据
	MessageCommand string      `json:"command"`     //指令
	Time           int64       `json:"time"`        //时间，采用时间戳毫秒级
	Payload        interface{} `json:"payload"`     //承载数据
}

// GameEventVO 从数据源收到的数据对应的结构
type GameEventVO struct {
	GameRoomId      int64            `json:"gameRoomId"`  //游戏房间id
	GameRoundNo     string           `json:"roundNo"`     //游戏局号，不可以超过32长度字符串
	NextGameRoundNo string           `json:"nextRoundNo"` //下一局局号,需要当前局给出下一局局号做预热数据
	Command         GameEventCommand `json:"command"`     //字符串形式的指令
	Time            int64            `json:"time"`        //时间，采用时间戳毫秒级
	ReceiveTime     int64            `json:"receiveTime"` //接收到时间的时间 采用时间戳毫秒级
	Payload         interface{}      `json:"payload"`     //承载数据
}

// GameEventResultVO 返回给数据源接口的结构体
type GameEventResultVO struct {
	GameRoomId      int64            `json:"gameRoomId"`  //游戏房间id
	GameRoundNo     string           `json:"roundNo"`     //游戏局号，不可以超过32长度字符串
	NextGameRoundNo string           `json:"nextRoundNo"` //下一局局号,需要当前局给出下一局局号做预热数据
	CountDown       int64            `json:"countDown"`   //倒计时时长，秒
	Duration        int64            `json:"duration"`    //游戏时长，秒单位
	Command         GameEventCommand //字符串形式的指令
	Time            int64            `json:"time"`    //时间，采用时间戳毫秒级
	Payload         interface{}      `json:"payload"` //承载数据

}

/*游戏事件相关model*/

// GameEventMessageHeader 游戏事件头
type GameEventMessageHeader struct {
	Command         string `json:"command"`         //指令
	GameId          int64  `json:"gameId"`          //游戏id
	GameRoomId      int64  `json:"gameRoomId"`      //游戏房间id
	GameRoundId     int64  `json:"gameRoundId"`     //游戏局id
	NextGameRoundId int64  `json:"nextGameRoundId"` //游戏下一局Id
	GameRoundNo     string `json:"gameRoundNo"`     //游戏局号
}

type JoinLeaveType string

const (
	RoomActionJoin  JoinLeaveType = "Join"  //加入房间
	RoomActionLeave JoinLeaveType = "Leave" //离开房间
)

// JoinLeaveGameRoom 玩家进入房间 离开房间
type JoinLeaveGameRoom struct {
	Namespace  string        `json:"namespace"`  //命名空间
	RoomId     string        `json:"roomId"`     //房间号
	UserId     int64         `json:"userId"`     //用户Id
	GameRoomId int64         `json:"gameRoomId"` //游戏房间id
	GameId     int64         `json:"gameId"`     //游戏id
	SessionId  string        `json:"sessionId"`  //会话id
	ClientIp   string        `json:"clientIp"`   //客户ip
	Currency   string        `json:"currency"`   //币种
	Type       JoinLeaveType `json:"type"`       //加入 Join，离开 Leave
}
