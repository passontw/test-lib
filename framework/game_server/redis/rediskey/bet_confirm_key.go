package rediskey

import (
	redistool "sl.framework.com/game_server/redis/redis_tool"
	"sl.framework.com/game_server/redis/types"
	"strconv"
	"time"
)

const (
	betConfirmedPrefix     = "BetConfirmed"     //投注确认
	betConfirmedLockPrefix = "BetConfirmedLock" //投注确认锁
)

/**
 * GetBetConfirmedRedisInfo
 * 投注确认redis信息
 *
 * @param gameRoomId int64 - 房间Id
 * @param gameRoundId int64 - 局号Id
 * @return *RedisInfo - 投注确认redis信息
 */

func GetBetConfirmedRedisInfo(gameRoomId, gameRoundId int64) *types.RedisInfo {
	return redistool.BuildRedisInfo(
		time.Duration(60)*time.Second,
		betFileKeyPrefix,
		betConfirmedPrefix,
		strconv.FormatInt(gameRoomId, 10),
		strconv.FormatInt(gameRoundId, 10),
	)
}

/**
 * GetBetConfirmedLockRedisInfo
 * 投注确认redis锁信息
 *
 * @param gameRoomId int64 - 房间Id
 * @param gameRoundId int64 - 局号Id
 * @return *RedisLockInfo - redis锁信息
 */

func GetBetConfirmedLockRedisInfo(gameRoomId, gameRoundId, userId string) *types.RedisLockInfo {
	return redistool.BuildRedisLockInfo(
		time.Duration(10)*time.Second,
		betFileKeyPrefix,
		betConfirmedLockPrefix,
		gameRoomId,
		gameRoundId,
		userId,
	)
}

/**
 * GetBetConfirmedGameDataLockRedisInfo
 * 投注确认在第一次gameData消息的redis锁信息
 *
 * @param gameRoomId int64 - 房间Id
 * @param gameRoundId int64 - 局号Id
 * @return *RedisLockInfo - redis锁信息
 */

func GetBetConfirmedGameDataLockRedisInfo(gameRoomId, gameRoundId string) *types.RedisLockInfo {
	return redistool.BuildRedisLockInfo(
		time.Duration(10)*time.Second,
		betFileKeyPrefix,
		betConfirmedLockPrefix,
		gameRoomId,
		gameRoundId,
	)
}
