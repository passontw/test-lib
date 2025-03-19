package rpcreq

import (
	"fmt"
	"sl.framework.com/game_server/conf"
	"time"
)

/*
 GameCreate创建游戏事件
*/

type GameCreate struct {
	GameRoomId  int64      `json:"gameRoomId"`  //数据库字段:game_room_id 游戏房间id
	GameRoundId int64      `json:"gameRoundId"` //数据库字段:game_round_id 游戏局id
	Command     MessageCmd `json:"command"`     //数据库字段:command 指令
	Payload     string     `json:"payload"`     //数据库字段:payload 负载数据
	CreateTime  time.Time  `json:"createTime"`  //数据库字段:create_time 创建时间
}

/**
 * UserMessageRequest
 * 发送玩家专属消息
 *
 * @param traceId string - traceId 用于日志跟踪
 * @param gameMessage GameMessage - 所发送的玩家专属消息
 * @return int - 请求返回码
 */

func CreateGameRequest(traceId string, userMessage UserMessage) int {
	url := fmt.Sprintf("%v/feign/message/user/send", conf.GetPlatformInfoUrl())
	msg := fmt.Sprintf("UserMessageRequest traceId=%v, userMessage=%v, url=%v", traceId, userMessage, url)

	ret := runHttpPost(traceId, msg, url, userMessage, nil)
	return ret
}
