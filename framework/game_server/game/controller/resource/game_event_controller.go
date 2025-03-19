package resource

import (
	"sl.framework.com/game_server/error_code"
	"sl.framework.com/game_server/game/controller/base_controller"
	"sl.framework.com/game_server/game/service/listenner"
	"time"

	//"sl.framework.com/game_server/game/service/listenner"
	types "sl.framework.com/game_server/game/service/type"
	"sl.framework.com/trace"
)

/*
	GameEventController 极速电子百家乐MQ游戏事件处理控制器
	游戏相关事件 包括 Bet_Start Bet_Stop Game_Data Change_Deck Game_Pause
	Bet_Start: 游戏开始 对下局个人限红 房间限红等做预缓存
	Bet_Stop: 游戏结束 如果Bet_Start没有缓存成功则在此消息中对个人限红 房间限红做预缓存
	Game_Data: 荷官发牌 游戏服务器接收到该协议暂不做任何处理
	Change_Deck: 换牌 游戏服务器接收到该协议暂不做任何处理
	Game_Pause: 游戏换靴 游戏服务器接收到该协议后，房间内牌的数量重置为0
*/

type GameEventController struct {
	base_controller.BaseController
}

/*
	GameEventController 极速电子百家乐MQ游戏事件处理控制器
	游戏相关事件 包括 Bet_Start Bet_Stop Game_Data Change_Deck Game_Pause
	Bet_Start: 游戏开始 对下局个人限红 房间限红等做预缓存
	Bet_Stop: 游戏结束 如果Bet_Start没有缓存成功则在此消息中对个人限红 房间限红做预缓存
	GameDraw: 游戏结算事件
	Game_Data: 荷官发牌 游戏服务器接收到该协议暂不做任何处理
	Change_Deck: 换牌 游戏服务器接收到该协议暂不做任何处理
	Game_Pause: 游戏换靴 游戏服务器接收到该协议后，房间内牌的数量重置为0
*/

/**
 * GameEvent
 * 从数据源接收事件数据
 *
 * @param
 * @return
 */

func (c *GameEventController) GameEvent() {
	var gameEvent types.GameEventMessageTmp
	controllerParserDTO := c.ParserFromDataSource(&gameEvent)
	if controllerParserDTO.Code != errcode.ErrorOk {
		trace.Error("GameEventController parser error, code=%v", controllerParserDTO.Code)
		return
	}
	trace.Info("[数据源游戏事件]，事件类型为:%v,GameEventController GameEvent trace id=%v, receive gameEvent=%+v",
		gameEvent.MessageCommand, controllerParserDTO.TraceId, gameEvent)

	command := types.GameEventCommand(gameEvent.MessageCommand)

	v2 := types.GameEventVO{
		GameRoomId:      gameEvent.GameRoomId,
		GameRoundNo:     gameEvent.RoundNo,
		NextGameRoundNo: gameEvent.NextRoundNo,
		Time:            gameEvent.Time,
		Command:         command,
		ReceiveTime:     time.Now().UnixMilli(),
		Payload:         gameEvent.Payload,
	}

	code := errcode.ErrorOk
	listenner.DispatchGameEventV2(controllerParserDTO, v2, &code)
	c.DataSourceResponse(code, controllerParserDTO.TraceId, nil)
}
