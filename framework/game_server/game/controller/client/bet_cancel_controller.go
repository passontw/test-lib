package client

import (
	"fmt"
	errcode "sl.framework.com/game_server/error_code"
	"sl.framework.com/game_server/game/controller/base_controller"
	"sl.framework.com/game_server/game/dao/redisdb"
	"sl.framework.com/game_server/game/service/bet"
	types "sl.framework.com/game_server/game/service/type"
	"sl.framework.com/game_server/redis/cache"
	"sl.framework.com/game_server/redis/rediskey"
	"sl.framework.com/tool"
	"sl.framework.com/trace"
	"strconv"
)

/**
 * BetCancel
 * 处理玩家投注取消操作
 * 控制器层处理 解析参数以及进行参数校验
 */

func (p *BetController) BetCancel() {
	param := types.BetCancelParam{}
	controllerParserDTO := p.ParserFromClient(&param)
	if controllerParserDTO.Code != errcode.ErrorOk {
		trace.Error("BetController BetCancel parser error, code=%v", controllerParserDTO.Code)
		return
	}
	userId := p.Ctx.Input.Header(string(base_controller.TagUserId))
	//requestId := p.Ctx.Input.Header(string(tagRequestId))

	//参数判断
	msgHeader := fmt.Sprintf("游戏投注取消 BetController handle traceId=%v, gameRoomId=%v, gameRoundId=%v, userId=%v, "+
		"order list=%+v", controllerParserDTO.TraceId, param.GameRoomId, param.GameRoundId, userId, param.OrderNoList)
	pWatcher := tool.NewWatcher(msgHeader)
	trace.Info("%v, param=%+v", msgHeader, param)
	if len(param.GameRoundId) == 0 || len(param.GameRoomId) == 0 || len(param.OrderNoList) == 0 {
		trace.Error("%v, invalid param", msgHeader)
		p.ClientResponse(errcode.HttpErrorInvalidParam, controllerParserDTO.TraceId, nil)
		return
	}

	//检查局状态
	roomId, _ := strconv.ParseInt(param.GameRoomId, 10, 64)
	roundCache := cache.GameRoundCache{TraceId: controllerParserDTO.TraceId, RoomId: roomId, GameRoundId: param.GameRoundId}
	roundCache.Get()
	gameDetail := roundCache.Data
	if gameDetail == nil || gameDetail.Status != "Bet_Start" {
		trace.Error("%v, current game detailed is %+v not Start, can't cancel", msgHeader, gameDetail)
		p.ClientResponse(errcode.GameErrorWrongGameRoundStatus, controllerParserDTO.TraceId, nil)
		return
	}

	//查询游戏事件缓存信息 检查游戏状态
	gameEventCache := cache.GameEventCache{TraceId: controllerParserDTO.TraceId, GameRoomId: param.GameRoomId, GameRoundId: param.GameRoundId}
	if !gameEventCache.Get() && cache.ConvertToEventCommandType(gameEventCache.Data.Command) != types.MessageCommandTypeBetStart {
		trace.Error("%v, current game detailed is %+v not Start, can't cancel", msgHeader, gameDetail)
		p.ClientResponse(errcode.GameErrorWrongGameRoundStatus, controllerParserDTO.TraceId, nil)
		return
	}

	//分布式锁防止重复取消
	redisLockInfo := rediskey.GetBetCancelLockRedisInfo(param.GameRoomId, param.GameRoundId, userId)
	if !redisdb.TryLock(redisLockInfo) {
		p.ClientResponse(errcode.GameErrorBetCancelTooFast, controllerParserDTO.TraceId, nil)
		trace.Error("%v, redis lock failed, lock info=%+v", msgHeader, redisLockInfo)
		return
	}
	defer redisdb.Unlock(redisLockInfo)

	//业务层处理
	controllerParserDTO.Code = bet.ServiceBetCancel(controllerParserDTO.TraceId, userId, &param)

	pWatcher.Stop()
	p.ClientResponse(controllerParserDTO.Code, controllerParserDTO.TraceId, "true")
}
