package client

import (
	"fmt"
	"sl.framework.com/game_server/conf"
	"sl.framework.com/game_server/error_code"
	"sl.framework.com/game_server/game/controller/base_controller"
	"sl.framework.com/game_server/game/service"
	gamelogic "sl.framework.com/game_server/game/service/game"
	types "sl.framework.com/game_server/game/service/type"
	"sl.framework.com/trace"
)

// JoinOrLeaveController 加入房间或者离开房间控制器
type JoinOrLeaveController struct {
	base_controller.BaseController
}

// JoinOrLeaveRoom 加入房间或者离开房间业务处理
func (c *JoinOrLeaveController) JoinOrLeaveRoom() {
	var joinLeaveRoom types.JoinLeaveGameRoom
	controllerParserDTO := c.ParserFromClient(&joinLeaveRoom)
	if controllerParserDTO.Code != errcode.ErrorOk {
		trace.Error("JoinOrLeaveController parser error, code=%v", controllerParserDTO.Code)
		return
	}
	msgHeader := fmt.Sprintf(" 客户端调用JoinOrLeave接口,JoinOrLeaveController JoinOrLeaveRoom trace id=%v", controllerParserDTO.TraceId)
	trace.Info("%v, receive data=%+v", msgHeader, joinLeaveRoom)
	if conf.GetGameId() != joinLeaveRoom.GameId {
		trace.Error("%v, wrong gameId=%v,  gameId should be=%v", joinLeaveRoom.GameId, conf.GetGameId())
		c.ClientResponse(errcode.GameErrorWrongGameId, controllerParserDTO.TraceId, nil)
		return
	}

	code := errcode.ErrorOk
	switch joinLeaveRoom.Type {
	case types.RoomActionJoin:
		gamelogic.OnJoinRoom(controllerParserDTO.TraceId, joinLeaveRoom)
	case types.RoomActionLeave:
		gamelogic.OnLeaveRoom(controllerParserDTO.TraceId, joinLeaveRoom)
	default:
		code = errcode.HttpErrorDataFailed
		trace.Info("%v, wrong type=%v", msgHeader, joinLeaveRoom.Type)
	}
	c.ClientResponse(code, controllerParserDTO.TraceId, nil)

	//异步通知到具体游戏服
	service.AsyncNotifyJoinLeaveGameRoomListener(controllerParserDTO.TraceId, joinLeaveRoom)
}
