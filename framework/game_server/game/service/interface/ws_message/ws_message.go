package ws_message

import "sl.framework.com/game_server/game/service/type/dto"

type IWsMessageHandler interface {
	/*
	 * OnMessage 消息处理函数
	 *@ traceId string 跟踪id
	 *@ message *dto.WsClientMessageDTO 消息结构体
	 *@ return -
	 */
	OnMessage(traceId string, message *dto.WsClientMessageDTO)
}
