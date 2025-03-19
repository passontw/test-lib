package rediskey

import (
	"sl.framework.com/game_server/redis/redis_tool"
	"sl.framework.com/game_server/redis/types"
	"strconv"
	"time"
)

const (
	gameResultPrefix     = "GameResult"
	gameResultLockPrefix = "GameResultLock"
)

// GetGameResultRedisInfo 游戏结果redis信息
func GetGameResultRedisInfo(gameRoomId, gameRoundId int64) *types.RedisInfo {
	return redistool.BuildRedisInfo(
		time.Duration(1)*time.Hour,
		gameFileKeyPrefix,
		gameResultPrefix,
		strconv.FormatInt(gameRoomId, 10),
		strconv.FormatInt(gameRoundId, 10),
	)
}

// GetGameResultLockRedisInfo 游戏结果redis锁
func GetGameResultLockRedisInfo(gameRoomId, gameRoundId int64) *types.RedisLockInfo {
	return redistool.BuildRedisLockInfo(
		time.Duration(10)*time.Second,
		gameFileKeyPrefix,
		gameResultLockPrefix,
		strconv.FormatInt(gameRoomId, 10),
		strconv.FormatInt(gameRoundId, 10),
	)
}

//
/**
 * GetGameRoundResultCacheKey
 * 获取 游戏局开奖 key
 *
 * @param gameRoundId - 游戏局id
 * @return *types.RedisInfo - redis 信息
 */

func GetGameRoundResultCacheKey(gameRoundId int64) *types.RedisInfo {
	return redistool.BuildRedisInfo(
		time.Duration(1)*time.Hour,
		gameFileKeyPrefix,
		gameResultPrefix,
		strconv.FormatInt(gameRoundId, 10),
	)
}
