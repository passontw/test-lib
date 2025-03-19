package rpcreq

import (
	"encoding/json"
	"fmt"
	"sl.framework.com/async"
	"sl.framework.com/game_server/conf"
	"sl.framework.com/trace"
	"time"
)

/*
	GameMessage游戏消息
*/

type GameMessage[T any] struct {
	GameRoomId     string  `json:"gameRoomId"`  //游戏房间id
	GameRoundId    string  `json:"gameRoundId"` //游戏局id
	MessageCommand string  `json:"command"`     //消息指令
	Date           int64   `json:"time"`        //UTC0时间毫秒
	UserIdList     []int64 `json:"userIdList"`  //客户端id,如用户id,此集合不为空则指定推送
	Body           T       `json:"body"`        //todo:使用interface{}
}

/**
 * AsyncGameMessageRequest
 * 发送游戏消息
 *
 * @param traceId string - traceId 用于日志跟踪
 * @param gameMessage GameMessage - 所发送的游戏消息
 * @return
 */

func AsyncGameMessageRequest[T any](traceId string, gameMessage GameMessage[T]) {
	fn := func() { gameMessageRequest(traceId, gameMessage) }
	async.AsyncRunCoroutine(fn)
}

/**
 * gameMessageRequest
 * 发送游戏消息
 *
 * @param traceId string - traceId 用于日志跟踪
 * @param gameMessage GameMessage - 所发送的游戏消息
 * @return int - 请求返回码
 */

func gameMessageRequest[T any](traceId string, gameMessage GameMessage[T]) {
	url := fmt.Sprintf("%v/feign/message/game/send", conf.GetPlatformInfoUrl())
	messageBuffer, _ := json.Marshal(gameMessage)
	message := string(messageBuffer)
	msg := fmt.Sprintf("gameMessageRequest traceId=%v, gameMessage=%+v, message=%v, url=%v", traceId, gameMessage, message, url)
	_ = runHttpPost(traceId, msg, url, gameMessage, nil)
}

/*
	UserMessage用户专属消息
*/

type UserMessage struct {
	GameRoomId     int64      `json:"gameRoomId"`     //游戏房间id
	UserId         int64      `json:"userId"`         //客户端id,如用户id,此集合不为空则指定推送
	MessageCommand MessageCmd `json:"messageCommand"` //消息指令
	Date           time.Time  `json:"date"`           //时间
	Body           string     `json:"body"`           //消息
}

/**
 * UserMessageRequest
 * 发送玩家专属消息
 *
 * @param traceId string - traceId 用于日志跟踪
 * @param gameMessage GameMessage - 所发送的玩家专属消息
 * @return int - 请求返回码
 */

func UserMessageRequest(traceId string, userMessage UserMessage) int {
	url := fmt.Sprintf("%v/feign/message/user/send", conf.GetPlatformInfoUrl())
	msg := fmt.Sprintf("UserMessageRequest traceId=%v, userMessage=%v, url=%v", traceId, userMessage, url)

	ret := runHttpPost(traceId, msg, url, userMessage, nil)
	return ret
}

/**
 * AsyncSendRoundMessage
 * 异步发送局消息到ws集群
 *
 * @traceId string - 跟踪id
 * @gameRoomId string - 房间id
 * @cmd MessageCmd - 局状态
 * @srcData interfaces{} - 消息体
 * @return  -
 */

func AsyncSendRoundMessage[T any](traceId, gameRoomId, gameRoundId string, cmd string, srcData T) {
	trace.Info("AsyncSendRoundMessage traceId=%v, gameRoomId=%v, gameRoundId=%v  cmd:%v",
		traceId, gameRoomId, gameRoundId, cmd)
	//异步事件通知取消
	message := GameMessage[T]{
		GameRoomId:     gameRoomId,
		GameRoundId:    gameRoundId,
		MessageCommand: cmd,
		Date:           time.Now().Unix(),
		Body:           srcData,
	}
	AsyncGameMessageRequest(traceId, message)
}
