package rediskey

import (
	"sl.framework.com/game_server/redis/redis_tool"
	"sl.framework.com/game_server/redis/types"
	"strconv"
	"time"
)

const (
	gameOrderReceiptPrefix     = "GameOrderReceipt"
	gameOrderReceiptLockPrefix = "GameOrderReceiptLock"
)

// GetGameOrderReceiptRedisInfo 游戏结果redis信息
func GetGameOrderReceiptRedisInfo(gameRoomId, gameRoundId int64) *types.RedisInfo {
	return redistool.BuildRedisInfo(
		time.Duration(1)*time.Hour,
		gameFileKeyPrefix,
		gameResultPrefix,
		strconv.FormatInt(gameRoomId, 10),
		strconv.FormatInt(gameRoundId, 10),
	)
}

// GetGameOrderReceiptLockRedisInfo 游戏结果redis锁
func GetGameOrderReceiptLockRedisInfo(gameRoomId, gameRoundId int64) *types.RedisLockInfo {
	return redistool.BuildRedisLockInfo(
		time.Duration(10)*time.Second,
		gameFileKeyPrefix,
		gameResultLockPrefix,
		strconv.FormatInt(gameRoomId, 10),
		strconv.FormatInt(gameRoundId, 10),
	)
}
