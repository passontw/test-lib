package client

import (
	"fmt"
	errcode "sl.framework.com/game_server/error_code"
	"sl.framework.com/game_server/game/controller/base_controller"
	"sl.framework.com/game_server/game/dao/redisdb"
	"sl.framework.com/game_server/game/service/bet"
	types "sl.framework.com/game_server/game/service/type"
	"sl.framework.com/game_server/redis/cache"
	rediskey "sl.framework.com/game_server/redis/rediskey"
	"sl.framework.com/tool"
	"sl.framework.com/trace"
)

/**
 * BetConfirm
 * 处理玩家投注确认操作
 */

func (p *BetController) BetConfirm() {
	param := types.BetConfirmParam{}
	controllerParserDTO := p.ParserFromClient(&param)
	if controllerParserDTO.Code != errcode.ErrorOk {
		trace.Error("BetController BetConfirm parser error, code=%v", controllerParserDTO.Code)
		return
	}
	userId := p.Ctx.Input.Header(string(base_controller.TagUserId))
	userName := p.Ctx.Input.Header(string(base_controller.TagUserName))
	userType := p.Ctx.Input.Header(string(base_controller.TagUserType))
	language := p.Ctx.Input.Header(string(base_controller.TagLanguage))
	currency := p.Ctx.Input.Header(string(base_controller.TagCurrency))
	clientType := p.Ctx.Input.Header(string(base_controller.TagClientType))

	msgHeader := fmt.Sprintf("游戏投注确认 BetController BetConfirm traceId=%v, gameRoomId=%v, gameRoundId=%v",
		controllerParserDTO.TraceId, param.GameRoomId, param.GameRoundId)
	pWatcher := tool.NewWatcher(msgHeader)
	trace.Info("%v, userId=%v, userName=%v, userType=%v, language=%v, currency=%v, clientType=%v, param=%+v",
		msgHeader, userId, userName, userType, language, currency, clientType, param)
	if len(param.GameRoundId) == 0 || len(param.GameRoomId) == 0 || len(param.Currency) == 0 {
		trace.Error("%v, invalid param", msgHeader)
		p.ClientResponse(errcode.HttpErrorInvalidParam, controllerParserDTO.TraceId, nil)
		return
	}

	//查询游戏事件缓存信息
	gameEventCache := cache.GameEventCache{TraceId: controllerParserDTO.TraceId, GameRoomId: param.GameRoomId, GameRoundId: param.GameRoundId}
	ok := gameEventCache.Get()
	if !ok || gameEventCache.Data == nil ||
		cache.ConvertToEventCommandType(gameEventCache.Data.Command) == types.MessageCommandTypeGameDraw {
		trace.Error("%v, game event cache not exist", msgHeader)
		p.ClientResponse(errcode.GameErrorBetConfirmLater, controllerParserDTO.TraceId, nil)
		return
	}

	//分布式锁防止重复确认，无需主动释放，让它自然过期
	redisLockInfo := rediskey.GetBetConfirmedLockRedisInfo(param.GameRoomId, param.GameRoundId, userId)
	if !redisdb.TryLock(redisLockInfo) {
		trace.Notice("%v, redis lock failed, lock info=%+v", msgHeader, redisLockInfo)
		p.ClientResponse(errcode.ErrorOk, controllerParserDTO.TraceId, nil)
		return
	}

	//业务层处理
	controllerParserDTO.Code = bet.ServiceBetConfirm(controllerParserDTO.TraceId, param.GameRoomId, param.GameRoundId, userId, currency)

	pWatcher.Stop()

	trace.Info("%v, userId=%v, userName=%v, bet confirm ClientResponse",
		msgHeader, userId, userName)
	p.ClientResponse(controllerParserDTO.Code, controllerParserDTO.TraceId, "true")
}
