package rediskey

import (
	redistool "sl.framework.com/game_server/redis/redis_tool"
	"sl.framework.com/game_server/redis/types"
	"strconv"
	"time"
)

const (
	betReceiptPrefix         = "BetReceipt"     //游戏下注小票
	betReceiptLockLockPrefix = "BetReceiptLock" //游戏下注小票注锁
)

/**
 * GetBetReceiptKey
 * 获取下注小票key
 *
 * @param gameRoomId int64 - 房间Id
 * @param gameRoundId int64 - 局号Id
 * @return *RedisInfo - redis信息
 */

func GetBetReceiptKey(gameRoomId, gameRoundId int64) *types.HashRedisInfo {
	return redistool.BuildHashRedisInfo(
		redistool.GetRedisExpireDuration(),
		betFileKeyPrefix,
		betReceiptPrefix,
		strconv.FormatInt(gameRoomId, 10),
		strconv.FormatInt(gameRoundId, 10),
	)
}

/**
 * GetBetReceiptLockKey
 * 获取下注小票锁
 *
 * @param gameRoomId int64 - 房间Id
 * @param gameRoomId int64 - 房间Id
 * @return *RedisLockInfo - redis锁信息
 */

func GetBetReceiptLockKey(gameRoomId, gameRoundId, userId string) *types.RedisLockInfo {
	return redistool.BuildRedisLockInfo(
		time.Duration(10)*time.Second,
		userId,
		betFileKeyPrefix,
		betReceiptLockLockPrefix,
		gameRoomId,
		gameRoundId,
	)
}
