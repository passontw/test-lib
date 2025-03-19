package cache

import (
	"encoding/json"
	"fmt"
	"sl.framework.com/game_server/game/dao/redisdb"
	types "sl.framework.com/game_server/game/service/type"
	"sl.framework.com/game_server/game/service/type/dto"
	"sl.framework.com/game_server/redis/rediskey"
	rpcreq "sl.framework.com/game_server/rpc_client"
	err "sl.framework.com/resource/error"
	"sl.framework.com/trace"
	"strconv"
	"time"
)

const gameEventFileKeyPrefix = "GameEvent"

const (
	gameEventKeyPrefix     = "GameEvent"
	gameEventLockKeyPrefix = "GameEventLock"
)

// GameEvent redis缓存游戏事件信息
type (
	GameEventCache struct {
		*BaseCache //基础缓存结构 封装解析与反解析函数

		Data    *types.GameEventVO //缓存中实际存储的结构
		TraceId string             //用于日志跟踪

		/*操作缓存所需要的信息*/
		GameRoomId  string
		GameRoundId string
	}
)

/**
 * Get
 * 获取缓存信息 并将缓存信息解析到对应缓存结构中
 */

func (event *GameEventCache) Get() (success bool) {
	var (
		val string
		err error
	)
	msgHeader := fmt.Sprintf("GameEventCache get traceId=%v, gameRoomId=%v, gameRoundId=%v",
		event.TraceId, event.GameRoomId, event.GameRoundId)

	//获取缓存信息
	redisKey := event.redisKey()
	if val, err = redisdb.Get(redisKey.Key); nil != err {
		trace.Error("%v, redis Get failed, key=%v,error=%v", msgHeader, redisKey.Key, err.Error())
		return
	}
	if len(val) <= 0 {
		trace.Notice("%v, redis HGet no data, key=%v", msgHeader, redisKey.Key)
		return
	}

	//解析缓存信息
	gameEvent := new(types.GameEventVO)
	if err = json.Unmarshal([]byte(val), gameEvent); nil != err {
		trace.Error("%v, json unmarshal failed, key=%v, val=%v, err=%v", msgHeader, redisKey.Key, val, err.Error())
		return
	}
	event.Data = gameEvent
	trace.Debug("%v, userInfo=%+v", msgHeader, gameEvent)

	return true
}

/**
 * Set
 * 设置缓存信息 将缓存结构序列化后以字符串的形式放入redis缓存中
 */

func (event *GameEventCache) Set(srcDto *types.GameEventVO) (success bool) {
	var (
		val    string
		err    error
		data   []byte
		curDto *types.GameEventVO
	)
	if srcDto == nil {
		curDto = event.Data
	} else {
		curDto = srcDto
	}
	msgHeader := fmt.Sprintf("GameEventCache set traceId=%v, gameRoomId=%v, gameRoundId=%v",
		event.TraceId, event.GameRoomId, event.GameRoundId)
	if data, err = json.Marshal(curDto); err != nil {
		trace.Error("%v, json marshal failed, error=%v", msgHeader, err.Error())
		return
	}
	redisKey := event.redisKey()
	if val, err = redisdb.Set(redisKey.Key, string(data), redisKey.Expire); nil != err {
		trace.Error("%v, redis set failed, table=%v, field=%v, error=%v", msgHeader, redisKey.Key,
			redisKey.Field, err.Error())
		return
	}

	trace.Info("%v, val=%v, gameEvent data=%+v", msgHeader, val, event.Data)
	return true
}

/**
 * SetList
 * 设置缓存信息 将缓存结构序列化后以字符串的形式放入redis list缓存中
 */

func (event *GameEventCache) SetList() (success bool) {
	var (
		val  string
		err  error
		data []byte
	)

	msgHeader := fmt.Sprintf("GameEventCache set traceId=%v, gameRoomId=%v, gameRoundId=%v",
		event.TraceId, event.GameRoomId, event.GameRoundId)
	if data, err = json.Marshal(event.Data); err != nil {
		trace.Error("%v, json marshal failed, error=%v", msgHeader, err.Error())
		return
	}
	redisKey := event.redisKey()
	val = string(data)
	if err = redisdb.LAppend(redisKey.Key, string(data), redisKey.Expire); nil != err {
		trace.Error("%v, redis set failed, table=%v, field=%v, error=%v", msgHeader, redisKey.Key,
			redisKey.Field, err.Error())
		return
	}

	trace.Info("%v, val=%v, gameEvent data=%+v", msgHeader, val, event.Data)
	return true
}

/**
 * Notify
 * 缓存游戏事件之后通知客户端
 */

func (event *GameEventCache) Notify() bool {
	trace.Info(" 缓存游戏事件之后通知客户端 GameEventCache Notify traceId=%v, gameRoomId=%v, gameRoundId=%v event data=%+v",
		event.TraceId, event.GameRoomId, event.GameRoundId, event.Data)
	//缓存游戏事件
	if !event.Set(event.Data) {
		return false
	}
	gameRoomId, _ := strconv.ParseInt(event.GameRoomId, 10, 64)
	gameRoundId, _ := strconv.ParseInt(event.GameRoundId, 10, 64)
	//通知客户端
	msgData := &dto.GameCommandDTO{
		GameRoomId:  gameRoomId,
		GameRoundId: gameRoundId,
		Command:     string(event.Data.Command),
		Payload:     event.Data.Payload,
		CreateTime:  event.Data.Time,
	}
	if err.ERR_OK != rpcreq.BuildGameEventRequest(event.TraceId, gameRoomId, gameRoundId, msgData) {
		trace.Error("GameEventCache Notify traceId=%v, gameRoomId=%v, gameRoundId=%v send BuildGameEventRequest failed.",
			event.TraceId, event.GameRoomId, event.GameRoundId)
		return false
	}
	return true
}

/**
 * redisKey
 * 构建redisKey
 * @return *redisInfo - *redisInfo
 */

func (event *GameEventCache) redisKey() *rediskey.RedisInfo {
	return rediskey.BuildRedisInfo(
		rediskey.GetRedisExpireDuration(),
		"", //key-value缓存 不设置field字段
		gameEventFileKeyPrefix,
		gameEventKeyPrefix,
		event.GameRoomId,
		event.GameRoundId,
	)
}

/**
 * redisLock
 * 构造锁信息
 *
 * @return *redisLockInfo - 锁信息
 */

func (event *GameEventCache) redisLock() *rediskey.RedisLockInfo {
	return rediskey.BuildRedisLockInfo(
		time.Duration(10)*time.Second,
		gameEventFileKeyPrefix,
		gameEventLockKeyPrefix,
		event.GameRoomId,
		event.GameRoundId,
	)
}

/**
 * ConvertToEventCommandType
 * 生成下一局信息
 *
 * @param srcCmd types.GameEventCommand - 用于日志跟踪
 * @return cache.MessageCommandType
 */

func ConvertToEventCommandType(gameEventCommand types.GameEventCommand) types.MessageCommandType {
	command := types.MessageCommandTypeInvalid
	switch types.GameEventCommand(gameEventCommand) {
	case types.GameEventCommandBetStart:
		command = types.MessageCommandTypeBetStart
	case types.GameEventCommandBetStop:
		command = types.MessageCommandTypeBetStop
	case types.GameEventCommandGameDraw:
		command = types.MessageCommandTypeGameDraw
	case types.GameEventCommandGameData:
		command = types.MessageCommandTypeGameData
	case types.GameEventCommandGamePause:
		command = types.MessageCommandTypeGamePause
	case types.GameEventCommandGameEnd:
		command = types.MessageCommandTypeGameEnd
	default:
		trace.Error("command between game server and data source not exist. gameEvent=%+v", command)
		return command
	}
	return command
}
