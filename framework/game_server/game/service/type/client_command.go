package types

// MessageCommandType 消息命令类型
type MessageCommandType int
type GameEventType int

const (
	MessageCommandTypeChat MessageCommandType = 28 //聊天消息 指令

	/*
		游戏流程 相关指令
	*/

	MessageCommandTypeInvalid   MessageCommandType = -1 //投注开始 指令 (开局)
	MessageCommandTypeBetStart  MessageCommandType = 2  //投注开始 指令 (开局)
	MessageCommandTypeBetStop   MessageCommandType = 3  //投注结束 指令
	MessageCommandTypeGameDraw  MessageCommandType = 4  //奖 指令 (局结束)
	MessageCommandTypeGameData  MessageCommandType = 8  //游戏数据，如：发牌，骰宝开牌，轮盘结果等等
	MessageCommandTypeGamePause MessageCommandType = 9  //游戏暂停,如换靴，维护等
	MessageCommandTypeGameEnd   MessageCommandType = 25 //游戏结束，本局结束

	/*
		用户游戏行为 相关指令
	*/

	MessageCommandTypeBet         MessageCommandType = 12 //投注 指令
	MessageCommandTypeBetCancel   MessageCommandType = 14 //取消投注 指令
	MessageCommandTypeBetReceipt  MessageCommandType = 26 //投注小票 指令
	MessageCommandTypeDynamicOdds MessageCommandType = 29 //动态赔率 指令

	GAMEEVENT_PAUSE     GameEventType = 1
	GAMEEVENT_GAMESTART GameEventType = 2
	GAMEEVENT_STOP      GameEventType = 3
	GAMEEVENT_GAMEDATA  GameEventType = 4
	GAMEEVENT_GAMEDRAW  GameEventType = 5
	GAMEEVENT_GAMEEND   GameEventType = 6
)

type GameEvent struct {
	GameRoomId     int64              `json:"gameRoomId"`     //游戏房间id
	RoundNo        string             `json:"roundNo"`        //游戏局号，不可以超过32长度字符串
	MessageCommand MessageCommandType `json:"messageCommand"` //指令
	Payload        interface{}        `json:"payload"`        //承载数据
	Time           int64              `json:"time"`           //时间，采用时间戳毫秒级
}

type GameRound struct {
	RoundId string `json:"id"`      //游戏局id
	RoundNo string `json:"roundNo"` //游戏局号，不可以超过32长度字符串
}
