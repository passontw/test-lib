package middle_platform

import (
	"fmt"
	errcode "sl.framework.com/game_server/error_code"
	"sl.framework.com/game_server/game/controller/base_controller"
	"sl.framework.com/game_server/game/service/type/dto"
	"sl.framework.com/game_server/game/service/websocket_message"
	"sl.framework.com/tool"
	"sl.framework.com/trace"
)

type WsClientMessageController struct {
	base_controller.BaseController
}

/**
 * WebSocketMessage
 * 从中台ws推送过来的消息
 *
 * @param PARAM - 参数说明
 * @return RETURN - 返回值说明
 */

func (p *WsClientMessageController) WebSocketMessage() {
	param := &dto.WsClientMessageDTO{}
	controllerParserDTO := p.ParserFromClient(&param)
	if controllerParserDTO.Code != errcode.ErrorOk {
		trace.Error("BetController BetRecord parser error, code=%v", controllerParserDTO.Code)
		return
	}
	userId := p.Ctx.Input.Header(string(base_controller.TagUserId))
	gameRoomId := p.Ctx.Input.Param(":gameRoomId")
	gameRoundId := p.Ctx.Input.Param(":gameRoundId")

	msgHeader := fmt.Sprintf("消息处理函数 WebSocketMessage traceId=%v, gameRoomId=%v, gameRoundId=%v, userId=%v",
		controllerParserDTO.TraceId, gameRoomId, gameRoundId, userId)
	pWatcher := tool.NewWatcher("消息处理")
	trace.Info("%v", msgHeader)
	if len(gameRoundId) == 0 || len(gameRoomId) == 0 {
		trace.Error("%v, invalid param", msgHeader)
		p.ClientResponse(errcode.HttpErrorInvalidParam, controllerParserDTO.TraceId, nil)
		return
	}
	websocket_message.HandleMessage(controllerParserDTO.TraceId, param)
	pWatcher.Stop()
	p.ClientResponse(controllerParserDTO.Code, controllerParserDTO.TraceId, nil)
}
