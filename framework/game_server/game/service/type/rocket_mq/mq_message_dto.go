package rocket_mq

import "encoding/json"

// rocketMq 消息封装
type RocketMqMessageDTO struct {
	MessageType    string          `json:"message_type"`    //消息类型 Game_Draw 开奖消息  Bet_confirm 提交注单消息
	MessagePayload json.RawMessage `json:"message_payload"` //消息负载 json.RawMessage 允许先将 MessagePayload 存为 []byte，然后再根据 MessageType 进行反序列化
}
