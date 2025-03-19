package rediskey

import (
	"sl.framework.com/game_server/redis/redis_tool"
	"sl.framework.com/game_server/redis/types"
	"strconv"
	"time"
)

const (
	gameDrawResultPrefix     = "GameDrawResult"
	gameDrawResultLockPrefix = "GameDrawResultLock"
)

// GetGameDrawResultRedisInfo 游戏结果redis信息
func GetGameDrawResultRedisInfo(gameRoomId, gameRoundId int64) *types.RedisInfo {
	return redistool.BuildRedisInfo(
		time.Duration(1)*time.Hour,
		gameFileKeyPrefix,
		gameDrawResultPrefix,
		strconv.FormatInt(gameRoomId, 10),
		strconv.FormatInt(gameRoundId, 10),
	)
}

// GetGameDrawResultLockRedisInfo 游戏结果redis锁
func GetGameDrawResultLockRedisInfo(gameRoomId, gameRoundId int64) *types.RedisLockInfo {
	return redistool.BuildRedisLockInfo(
		time.Duration(10)*time.Second,
		gameFileKeyPrefix,
		gameDrawResultLockPrefix,
		strconv.FormatInt(gameRoomId, 10),
		strconv.FormatInt(gameRoundId, 10),
	)
}
