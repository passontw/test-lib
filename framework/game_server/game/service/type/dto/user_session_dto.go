package dto

type UserSessionDTO struct {
	NameSpace  string `json:"name_space"`   //聊天室所在命名空间，如:站点id
	RoomId     string `json:"room_id"`      //房间号
	UserId     int64  `json:"user_id"`      //用户Id
	GameRoomId int64  `json:"game_room_id"` //游戏房间id
	GameId     int64  `json:"game_id"`      //游戏id
	SessionId  string `json:"session_id"`   //会话id
	ClientIp   string `json:"client_ip"`    //客户ip
	Currency   string `json:"currency"`     //币种
	Online     bool   `json:"online"`       //是否在线，是 true,否 false
}
