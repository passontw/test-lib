package dto

type GameCommandDTO struct {
	GameRoomId  int64  `json:"gameRoomId"`
	GameRoundId int64  `json:"gameRoundId"`
	Command     string `json:"command"`
	Payload     any    `json:"payload"`
	CreateTime  int64  `json:"createTime"`
}
