package resource

import (
	"sl.framework.com/game_server/conf"
	errcode "sl.framework.com/game_server/error_code"
	"sl.framework.com/game_server/game/controller/base_controller"
	gameevent "sl.framework.com/game_server/game/service/game_event"
	"sl.framework.com/game_server/game/service/listenner"
	types "sl.framework.com/game_server/game/service/type"
	"sl.framework.com/game_server/redis/cache"
	"sl.framework.com/trace"
	"strconv"
)

type GameRoundController struct {
	base_controller.BaseController
}

/**
 * GameEventRoundNo
 * 函数功能描述
 *
 * @param PARAM - 参数说明
 * @return RETURN - 返回值说明
 */

func (c *GameRoundController) GameEventRoundNo() {
	var (
		ret  int
		data *types.GameEventVO
	)
	ret = errcode.ErrorOk
	controllerParserDTO := c.ParserFromDataSource(nil)
	trace.Info(" 数据源调用 game/round接口,GameRoundController GameEventRoundNo traceId=%v,data=%v", controllerParserDTO.TraceId, string(c.Ctx.Input.RequestBody))
	if controllerParserDTO.Code != errcode.ErrorOk {
		trace.Error("BetController GameEventRoundNo parser error, code=%v", controllerParserDTO.Code)
		return
	}
	gameRoomId := c.Ctx.Input.Param(":gameRoomId")
	gameRoundNo := c.Ctx.Input.Param(":gameRoundNo")
	nGameRoomId, _ := strconv.ParseInt(gameRoomId, 10, 64)

	roundDTO := listenner.GetRoundInfo(controllerParserDTO.TraceId, gameRoundNo, nGameRoomId)

	eventCache := cache.GameEventCache{TraceId: controllerParserDTO.TraceId, GameRoomId: gameRoomId, GameRoundId: roundDTO.Id}
	ok := eventCache.Get()
	if !ok {
		data = &types.GameEventVO{}
	} else {
		data = eventCache.Data
	}
	//trace.Info(" GameRoundController GameEventRoundNo traceId=%v,get event cache ok:%v", traceId, ok)
	//获取游戏信息
	gameCache := cache.GameCache{TraceId: controllerParserDTO.TraceId, GameId: conf.GetGameId()}
	ok = gameCache.Get()
	//trace.Info(" GameRoundController GameEventRoundNo traceId=%v,get GameCache ok:%v", traceId, ok)
	if !ok {
		//如果没有则获取局信息
		nGameRoomId, _ := strconv.ParseInt(gameRoomId, 10, 64)
		nGameRoundId, _ := strconv.ParseInt(roundDTO.Id, 10, 64)
		ok := gameevent.GameRoomCache(controllerParserDTO.TraceId, nGameRoomId, nGameRoundId, conf.GetGameId())
		if ok {
			gameCache.Get()
		} else {
			gameCache.Data = new(types.Game)
			gameCache.Data.Countdown = 15
			gameCache.Data.Countdown = 30
		}
		//trace.Info(" GameRoundController GameEventRoundNo traceId=%v,get GameRoomCache data:%v", gameCache.Data)
	}

	dstData := &types.GameEventResultVO{
		GameRoomId:      data.GameRoomId,
		GameRoundNo:     data.GameRoundNo,
		NextGameRoundNo: data.NextGameRoundNo,
		CountDown:       int64(gameCache.Data.Countdown),
		Duration:        int64(gameCache.Data.Duration),
		Command:         data.Command,
		Time:            data.Time,
		Payload:         data.Payload,
	}
	c.DataSourceResponse(ret, controllerParserDTO.TraceId, dstData)
}
