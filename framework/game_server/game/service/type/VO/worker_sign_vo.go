package VO

type WorkerSignVO struct {
	GameRoomId int64  `json:"gameRoomId"` //房间id
	UserName   string `json:"username"`   //员工用户名
}
