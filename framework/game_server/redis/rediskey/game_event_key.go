package rediskey

import (
	"sl.framework.com/game_server/redis/redis_tool"
	"sl.framework.com/game_server/redis/types"
	"strconv"
	"time"
)

const (
	gameEventPrefix     = "GameEvent"
	gameEventLockPrefix = "GameEventLock"
)

// GetGameEventRedisInfo 游戏事件redis信息
func GetGameEventRedisInfo(gameRoundId int64) *types.RedisInfo {
	return redistool.BuildRedisInfo(
		time.Duration(10)*time.Minute,
		gameEventPrefix,
		gameCommandPrefix,
		strconv.FormatInt(gameRoundId, 10),
	)
}

// GetGameEventLockRedisInfo 游戏事件redis锁
func GetGameEventLockRedisInfo(requestId, gameRoundNo, command string) *types.RedisLockInfo {
	return redistool.BuildRedisLockInfo(
		time.Duration(10)*time.Second,
		gameEventPrefix,
		gameEventLockPrefix,
		requestId,
		gameRoundNo,
		command,
	)
}
