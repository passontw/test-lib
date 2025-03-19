package handler

import (
	"encoding/json"
	"sl.framework.com/async"
	errcode "sl.framework.com/game_server/error_code"
	"sl.framework.com/game_server/game/dao/redisdb"
	"sl.framework.com/game_server/game/service"
	"sl.framework.com/game_server/game/service/bet"
	types "sl.framework.com/game_server/game/service/type"
	"sl.framework.com/game_server/game/service/type/rocket_mq"
	"sl.framework.com/game_server/redis/rediskey"
	"sl.framework.com/tool"
	"sl.framework.com/trace"
	"strconv"
)

/**
 * OnBetConfirmHandler
 * 提交注单消息处理函数
 *
 * @param traceId string - 跟踪Id
 * @param msgBody []byte 消息体
 * @return int - 返回码
 */

func OnBetConfirmHandler(traceId string, msgBody []byte) int {
	trace.Info("[提交注单消息处理函数] traceId=%v 包体长度=%v", traceId, string(msgBody))

	//解析MQ消息
	betConfirmMessagePayload := new(rocket_mq.BetConfirmMessagePayload)
	if err := json.Unmarshal(msgBody, &betConfirmMessagePayload); nil != err {
		trace.Error("[提交注单消息处理函数] traceId=%v json unmarshal failed, error=%v, data=%v", traceId, err.Error(), string(msgBody))
		return errcode.HttpErrorOrderParse
	}
	if !service.ValidateGameId(types.GameId(betConfirmMessagePayload.GameId)) {
		trace.Error("[提交注单消息处理函数] traceId=%v validate GameId failed, msg=%v", traceId, string(msgBody))
		return errcode.ErrorOk
	}
	pWatcher := tool.NewWatcher("注单提交处理")
	ret := processBetConfirm(traceId, betConfirmMessagePayload)
	trace.Info("[提交注单消息处理函数] traceId=%v 处理完成, GameRoomId=%v, GameRoundId=%v, ret=%v",
		traceId, betConfirmMessagePayload.GameRoomId, betConfirmMessagePayload.GameRoundId, ret)
	pWatcher.Stop()
	return errcode.ErrorOk
}

/**
 * processBetConfirm
 * 处理注单提交函数
 *
 * @param traceId string- 跟踪id
 * @param payload *rocket_mq.BetConfirmMessagePayload -提交的注单用户信息
 * @return int - 返回码
 */

func processBetConfirm(traceId string, payload *rocket_mq.BetConfirmMessagePayload) int {
	trace.Info("[处理提交注单] traceId=%v payload=%+v", traceId, payload)

	gameRoomId := payload.GameRoomId
	gameRoundId := payload.GameRoundId
	for _, userInfo := range payload.UserInfo {
		trace.Debug("[处理提交注单] traceId=%v userInfo=%+v", traceId, userInfo)
		async.AsyncRunCoroutine(func() {
			//分布式锁防止重复确认，无需主动释放，让它自然过期
			redisLockInfo := rediskey.GetBetConfirmedLockRedisInfo(strconv.FormatInt(gameRoomId, 10),
				strconv.FormatInt(gameRoundId, 10), userInfo.UserId)
			if !redisdb.TryLock(redisLockInfo) {
				trace.Notice("[处理提交注单] traceId=%v, redis lock failed, lock info=%+v", traceId, redisLockInfo)
				return
			}
			bet.ServiceBetConfirm(traceId, strconv.FormatInt(gameRoomId, 10),
				strconv.FormatInt(gameRoundId, 10), userInfo.UserId, userInfo.Currency)
		})
	}

	return errcode.ErrorOk
}
