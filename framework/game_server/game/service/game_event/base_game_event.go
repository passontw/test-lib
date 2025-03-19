package gameevent

import (
	"encoding/json"
	"fmt"
	"sl.framework.com/game_server/conf"
	errcode "sl.framework.com/game_server/error_code"
	"sl.framework.com/game_server/game/dao/redisdb"
	"sl.framework.com/game_server/game/service/interface/events"
	types "sl.framework.com/game_server/game/service/type"
	"sl.framework.com/game_server/game/service/type/VO"
	"sl.framework.com/game_server/redis/cache"
	rediskey "sl.framework.com/game_server/redis/rediskey"
	rpcreq "sl.framework.com/game_server/rpc_client"
	"sl.framework.com/tool"
	"sl.framework.com/trace"
	"strconv"
	"time"
)

/**
 * FunctionName
 * 函数功能描述
 *
 * @param PARAM - 参数说明
 * @return RETURN - 返回值说明
 */

func CreateInstance(event types.GameEventVO, roundDto *types.GameRoundDTO, gameEventInitVo *VO.GameEventInitVO) events.IEventHandler {
	trace.Info("InitEvent traceId:%s,gametype:%d,roomid:%v,gameRoundId:%v,nextGameRoundId:%v",
		gameEventInitVo.TraceId, event.Command, gameEventInitVo.RoomId, gameEventInitVo.RoundId, gameEventInitVo.NextRoundId)
	switch event.Command {
	case types.GameEventCommandBetStart:
		return NewGameStart(event, roundDto, gameEventInitVo)
	case types.GameEventCommandBetStop:
		return NewGameStop(event, roundDto, gameEventInitVo)
	case types.GameEventCommandGamePause:
		return NewGamePause(event, roundDto, gameEventInitVo)
	case types.GameEventCommandGameDraw:
		return NewGameDraw(event, roundDto, gameEventInitVo)
	case types.GameEventCommandGameData:
		return NewGameData(event, roundDto, gameEventInitVo)
	case types.GameEventCommandGameEnd:
		return NewGameEnd(event, roundDto, gameEventInitVo)
	case types.GameEventCommandCancelRound:
		return NewCancelRound(event, roundDto, gameEventInitVo)
	case types.GameEventCommandChangeDeck:
		return NewChangeCard(event, roundDto, gameEventInitVo)
	default:
		return &types.EventBase{
			Dto: &types.EventDTO{
				GameRoomId:      gameEventInitVo.RoomId,
				GameRoundId:     gameEventInitVo.RoundId,
				GameId:          conf.GetGameId(),
				NextGameRoundId: gameEventInitVo.NextRoundId,
				GameRoundNo:     event.GameRoundNo,
				Command:         string(event.Command),
				Payload:         event.Payload,
			},
			RoundDTO:       roundDto,
			TraceId:        gameEventInitVo.TraceId,
			RequestId:      gameEventInitVo.RequestId,
			RetHandleEvent: gameEventInitVo.Code,
			MsgHeader: fmt.Sprintf("%s HandleEvent traceId=%v,requestId=%v, roomId=%v, gameRoundId=%v, "+
				"nextGameRoundId=%v", event.Command, gameEventInitVo.TraceId, gameEventInitVo.RequestId, gameEventInitVo.RoomId, gameEventInitVo.RoundId, gameEventInitVo.NextRoundId),
		}
	}
}

// 设置下局局信息的redis缓存 供Bet_Stop消息查询 如果没有下局信息则设置下局局缓存
func SetNextGameRoundId(gameRoundId, gameNextRoundId int64, requestId string) {
	msgHeader := fmt.Sprintf("setNextGameRoundId gameRoundId=%v, gameNextRoundId=%v", gameRoundId, gameNextRoundId)
	trace.Info("%v", msgHeader)

	redisInfo := rediskey.GetNextRoundIdRedisInfo(gameRoundId)
	redisLockInfo := rediskey.GetNextRoundIdLockRedisInfo(gameRoundId, requestId)
	if !redisdb.Lock(redisLockInfo) {
		trace.Error("%v, redis lock key=%v, lock failed", msgHeader, redisLockInfo.Key)
		return
	}
	defer redisdb.Unlock(redisLockInfo)

	var (
		val string
		err error
	)
	if val, err = redisdb.Set(redisInfo.Key, strconv.FormatInt(gameNextRoundId, 10), redisInfo.Expire); nil != err {
		trace.Error("%v, key=%v redis set err, err=%v", msgHeader, redisInfo.Key, err.Error())
		return
	}
	trace.Info("%v, key=%v, val=%v, redis set success", msgHeader, redisInfo.Key, val)
}

// 获取下局局信息的redis缓存 如果没有则缓存下局消息 有则跳过
func GetNextGameRoundId(gameRoundId int64) int64 {
	var (
		val string
		err error
	)

	redisInfo := rediskey.GetNextRoundIdRedisInfo(gameRoundId)
	if val, err = redisdb.Get(redisInfo.Key); nil != err {
		trace.Error("getNextGameRoundId gameRoundId=%v, key=%v, redis get err, err=%v",
			gameRoundId, redisInfo.Key, err.Error())
		return 0
	}
	trace.Info("getNextGameRoundId gameRoundId=%v, key=%v, val=%v, redis get success",
		gameRoundId, redisInfo.Key, val)
	nextGameRoundId, _ := strconv.ParseInt(val, 10, 64)

	return nextGameRoundId
}

/**
 * EventCommonSet
 * 事件公共设置
 *
 * @param PARAM - 参数说明
 * @return RETURN - 返回值说明
 */

func EventCommonSet(e *types.EventBase, status string, cmd string) {
	e.RoundDTO.Status = status
	//缓存局信息
	roundCache := cache.GameRoundCache{TraceId: e.TraceId, RoomId: e.Dto.GameRoomId, GameRoundId: e.RoundDTO.Id}
	roundCache.Set(e.RoundDTO)

	//发送到ws的时间
	sendTimestamp := time.Now().UnixMilli()
	sendTime := tool.FormatTime(sendTimestamp)
	nEventReceiveOffset := e.Dto.ReceiveTime - e.Dto.Time
	eventTime := tool.FormatTime(e.Dto.Time)
	receiveTime := tool.FormatTime(e.Dto.ReceiveTime)
	nSendReceiveOffset := sendTimestamp - e.Dto.ReceiveTime

	trace.Notice("[数据源时间转发ws] traceId:%v,事件类型：%v\r\n,事件接收时间：%v,时间发生时间：%v\r\n,事件发生到接收时间差：%v毫秒\r\n,发送时间：%v,接收和发送时差：%v毫秒\r\n",
		e.TraceId, e.Dto.Command, receiveTime, eventTime, nEventReceiveOffset, sendTime, nSendReceiveOffset)
	//发送到ws集群
	rpcreq.AsyncSendRoundMessage[interface{}](e.TraceId, strconv.FormatInt(e.Dto.GameRoomId, 10), e.RoundDTO.Id, cmd, e.Dto.Payload)

}

/**
 * GameRoomCache
 * 游戏玩法查询并缓存
 * 房间限红查询并缓存
 *
 * @return
 */

func GameRoomCache(traceId string, gameRoomId, gameRoundId, gameId int64) (success bool) {
	trace.Info("[游戏玩法查询并缓存] gameRoomCache traceId:%v,gameRoomId:%v,gameRoundId:%v", traceId, gameRoomId, gameRoundId)

	msgHeader := fmt.Sprintf("gameRoomCache traceId:%v,gameRoomId:%v,gameRoundId:%v", traceId, gameRoomId, gameRoundId)
	//从能力平台中获取房间详情
	trace.Info(" [游戏玩法查询并缓存] 从能力平台中获取房间详情 gameRoomCache traceId:%v,gameRoomId:%v,gameRoundId:%v", traceId, gameRoomId, gameRoundId)
	gameInfo, err := rpcreq.GetRoomInfoRequest(traceId, gameRoomId, gameRoundId)
	if err != errcode.ErrorOk {
		trace.Error("[游戏玩法查询并缓存] gameRoomCache traceId:%v,gameRoomId:%v,gameRoundId:%v,err:%v", traceId, gameRoomId, gameRoundId, err)
		return
	}
	//设置房间限红缓存
	trace.Info("[游戏玩法查询并缓存] 设置房间限红缓存 gameRoomCache traceId:%v,gameRoomId:%v,gameRoundId:%v", traceId, gameRoomId, gameRoundId)
	roomLimitKey := rediskey.GetRoomLimitHKey(gameId)
	if data, err := json.Marshal(gameInfo.BetLimitRuleList); err != nil {
		trace.Error("[游戏玩法查询并缓存] %v, roomLimitKey Marshal failed, table=%v, field=%v", msgHeader, roomLimitKey.HTable,
			roomLimitKey.Filed, err.Error())
		return
	} else {
		if _, err := redisdb.HSet(roomLimitKey.HTable, roomLimitKey.Filed, string(data), roomLimitKey.Expire); err != nil {
			trace.Error("[游戏玩法查询并缓存] %v, redis HSet failed, table=%v, field=%v", msgHeader, roomLimitKey.HTable,
				roomLimitKey.Filed, err.Error())
		}
	}
	//设置游戏玩法缓存
	trace.Info("[游戏玩法查询并缓存] 设置游戏玩法缓存 gameRoomCache traceId:%v,gameRoomId:%v,gameRoundId:%v", traceId, gameRoomId, gameRoundId)
	wagerCache := cache.WagerCache{TraceId: traceId, GameId: gameId, GameRoomId: gameRoomId}
	wagerCache.Set(gameInfo.GameWagerList)
	//设置游戏缓存
	trace.Info("[游戏玩法查询并缓存] 设置游戏缓存 gameRoomCache traceId:%v,gameRoomId:%v,gameRoundId:%v", traceId, gameRoomId, gameRoundId)
	gameCache := cache.GameCache{TraceId: traceId, GameId: gameId}
	gameCache.Set(gameInfo.Game)
	//设置游戏房间缓存
	trace.Info("[游戏玩法查询并缓存] 设置游戏房间缓存 gameRoomCache traceId:%v,gameRoomId:%v,gameRoundId:%v", traceId, gameRoomId, gameRoundId)
	gameRoomCache := cache.GameRoomCache{RoomId: gameRoomId}
	gameRoomCache.Set(gameInfo)
	return true
}
