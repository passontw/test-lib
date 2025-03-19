package rediskey

import (
	"sl.framework.com/game_server/redis/redis_tool"
	"sl.framework.com/game_server/redis/types"
	"strconv"
	"time"
)

const (
	gameNextGameRoundIdPrefix     = "NextGameRoundId"
	gameNextGameRoundIdLockPrefix = "NextGameRoundIdLock"
)

// GetNextRoundIdRedisInfo 下局局号信息
func GetNextRoundIdRedisInfo(gameRoundId int64) *types.RedisInfo {
	return redistool.BuildRedisInfo(
		redistool.GetRedisExpireDuration(),
		gameFileKeyPrefix,
		gameNextGameRoundIdPrefix,
		strconv.FormatInt(gameRoundId, 10),
	)
}

// GetNextRoundIdLockRedisInfo 下局局号信息key信息锁
func GetNextRoundIdLockRedisInfo(gameRoundId int64, requestId string) *types.RedisLockInfo {
	return redistool.BuildRedisLockInfo(
		types.RedisLockExpireDuration,
		gameFileKeyPrefix,
		gameNextGameRoundIdLockPrefix,
		strconv.FormatInt(gameRoundId, 10),
		requestId,
	)
}

const (
	gameRoundStatusPrefix     = "GameRoundStatus"
	gameRoundStatusLockPrefix = "GameRoundStatusLock"
)

// GetGameRoundDetailRedisInfo 游戏局状态redis信息
func GetGameRoundDetailRedisInfo(gameRoomId, gameRoundId string) *types.RedisInfo {
	return redistool.BuildRedisInfo(
		redistool.GetRedisExpireDuration(),
		gameFileKeyPrefix,
		gameRoundStatusPrefix,
		gameRoomId,
		gameRoundId,
	)
}

// GetGameRoundStatusLockRedisInfo 游戏局状态redis锁
func GetGameRoundStatusLockRedisInfo(gameRoomId, gameRoundId string) *types.RedisLockInfo {
	return redistool.BuildRedisLockInfo(
		time.Duration(10)*time.Second,
		gameFileKeyPrefix,
		gameRoundStatusLockPrefix,
		gameRoomId,
		gameRoundId,
	)
}
