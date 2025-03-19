package client

import (
	"fmt"
	"sl.framework.com/game_server/error_code"
	"sl.framework.com/game_server/game/controller/base_controller"
	"sl.framework.com/game_server/game/dao/redisdb"
	"sl.framework.com/game_server/game/service/bet"
	types "sl.framework.com/game_server/game/service/type"
	"sl.framework.com/game_server/redis/cache"
	rediskey "sl.framework.com/game_server/redis/rediskey"
	"sl.framework.com/tool"
	"sl.framework.com/trace"
)

// BetController 玩家下注订单处理
type BetController struct {
	base_controller.BaseController
}

/**
 * Bet
 * 处理玩家下注信息
 */

func (p *BetController) Bet() {
	param := types.BetVO{}
	controllerParserDTO := p.ParserFromClient(&param)
	if controllerParserDTO.Code != errcode.ErrorOk {
		trace.Error("BetController Bet parser error, code=%v", controllerParserDTO.Code)
		return
	}
	userId := p.Ctx.Input.Header(string(base_controller.TagUserId))
	userName := p.Ctx.Input.Header(string(base_controller.TagUserName))
	userType := p.Ctx.Input.Header(string(base_controller.TagUserType))
	language := p.Ctx.Input.Header(string(base_controller.TagLanguage))
	currency := p.Ctx.Input.Header(string(base_controller.TagCurrency))
	clientType := p.Ctx.Input.Header(string(base_controller.TagClientType))

	msgHeader := fmt.Sprintf("游戏投注 BetController Bet traceId=%v, gameRoomId=%v, gameRoundId=%v",
		controllerParserDTO.TraceId, param.GameRoomId, param.GameRoundId)
	pWatcher := tool.NewWatcher(msgHeader)
	trace.Info("%v, userId=%v, userName=%v, userType=%v, language=%v, currency=%v, clientType=%v",
		msgHeader, userId, userName, userType, language, currency, clientType)
	trace.Info("%v, param Bets=%+v", msgHeader, param.Bets)
	trace.Info("%v, param LimitRule=%+v", msgHeader, param.LimitRule)
	trace.Info("%v, param Device=%+v", msgHeader, param.Device)
	if len(param.GameRoundId) == 0 || len(param.GameRoomId) == 0 || len(param.Bets) == 0 || param.BetAmount <= 0 ||
		len(userId) <= 0 || len(param.Currency) <= 0 {
		trace.Error("%v, invalid param", msgHeader)
		p.ClientResponse(errcode.HttpErrorInvalidParam, controllerParserDTO.TraceId, nil)
		return
	}

	//查询游戏事件缓存信息 检查游戏状态
	gameEventCache := cache.GameEventCache{TraceId: controllerParserDTO.TraceId, GameRoomId: param.GameRoomId, GameRoundId: param.GameRoundId}
	if !gameEventCache.Get() || cache.ConvertToEventCommandType(gameEventCache.Data.Command) != types.MessageCommandTypeBetStart {
		trace.Error("%v, game event cache not exist", msgHeader)
		p.ClientResponse(errcode.GameErrorWrongGameRoundStatus, controllerParserDTO.TraceId, nil)
		return
	}

	//分布式锁防止重复下注
	redisLockInfo := rediskey.GetBetLockRedisInfo(param.GameRoomId, param.GameRoundId, userId)
	if !redisdb.TryLock(redisLockInfo) {
		trace.Error("%v, redis lock failed, lock info=%+v", msgHeader, redisLockInfo)
		p.ClientResponse(errcode.GameErrorBetTooFast, controllerParserDTO.TraceId, nil)
		return
	}
	defer redisdb.Unlock(redisLockInfo)

	//下注业务处理
	ret, betResult := bet.ServiceBet(controllerParserDTO.TraceId, userId, &param)

	pWatcher.Stop()
	p.ClientResponse(ret, controllerParserDTO.TraceId, betResult)
}
