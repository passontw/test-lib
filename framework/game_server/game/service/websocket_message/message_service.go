package websocket_message

import (
	"sl.framework.com/game_server/conf"
	"sl.framework.com/game_server/game/service"
	types "sl.framework.com/game_server/game/service/type"
	"sl.framework.com/game_server/game/service/type/dto"
	"sl.framework.com/trace"
)

/**
 * HandleMessage
 * 消息处理函数
 *
 * @param traceId - 跟踪id
 * @param msg *dto.WsClientMessageDTO -消息体
 * @return RETURN -
 */

func HandleMessage(traceId string, msg *dto.WsClientMessageDTO) {
	gameId := conf.GetGameId()

	message := service.GetWebsocketMessage(traceId, types.GameId(gameId))
	if message == nil {
		trace.Error("WebSocketMessage 找不到消息处理handler")
	} else {
		message.OnMessage(traceId, msg)
	}

}
