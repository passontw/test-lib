package types

type GameDrawMessage struct {
	GameRoomId         int64              `json:"gameRoomId"`
	GameRoundId        int64              `json:"gameRoundId"`
	GameId             int64              `json:"gameId"`
	GameRoundNo        string             `json:"gameRoundNo"`
	GameRoundResultDTO GameRoundResultDTO `json:"gameRoundResultDTOs"`
	OrderList          []int64            `json:"orderList"`
}
type MQMESSAGE_TOPICS string

const (
	GameDraw_0_OUT MQMESSAGE_TOPICS = "gameDraw-0-out-0"
	GameDraw_1_OUT MQMESSAGE_TOPICS = "gameDraw-1-out-0"
	GameDraw_2_OUT MQMESSAGE_TOPICS = "gameDraw-2-out-0"
	GameDraw_3_OUT MQMESSAGE_TOPICS = "gameDraw-3-out-0"
	GameDraw_4_OUT MQMESSAGE_TOPICS = "gameDraw-4-out-0"
	GameDraw_5_OUT MQMESSAGE_TOPICS = "gameDraw-5-out-0"
	GameDraw_6_OUT MQMESSAGE_TOPICS = "gameDraw-6-out-0"
	GameDraw_7_OUT MQMESSAGE_TOPICS = "gameDraw-7-out-0"
	GameDraw_8_OUT MQMESSAGE_TOPICS = "gameDraw-8-out-0"
	GameDraw_9_OUT MQMESSAGE_TOPICS = "gameDraw-9-out-0"
)
