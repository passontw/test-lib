package rediskey

import (
	redistool "sl.framework.com/game_server/redis/redis_tool"
	"sl.framework.com/game_server/redis/types"
	"strconv"
	"time"
)

const (
	betCancelPrefix         = "BetCancel"     //游戏局禁止投注
	betCancelLockLockPrefix = "BetCancelLock" //游戏局禁止投注锁
)

/**
 * GetBetCancelRedisInfo
 * 游戏局取消投注redis信息
 *
 * @param gameRoomId int64 - 房间Id
 * @param gameRoundId int64 - 局号Id
 * @param userId int64 - 用户Id
 * @return *RedisInfo - redis信息
 */

func GetBetCancelRedisInfo(gameRoomId, gameRoundId, userId int64) *types.RedisInfo {
	return redistool.BuildRedisInfo(
		redistool.GetRedisExpireDuration(),
		betFileKeyPrefix,
		betCancelPrefix,
		strconv.FormatInt(gameRoomId, 10),
		strconv.FormatInt(gameRoundId, 10),
		strconv.FormatInt(userId, 10),
	)
}

/**
 * GetBetCancelLockRedisInfo
 * 游戏局取消投注redis锁信息
 *
 * @param gameRoomId int64 - 房间Id
 * @param gameRoomId int64 - 房间Id
 * @param userId int64 - 用户Id
 * @return *RedisLockInfo - redis锁信息
 */

func GetBetCancelLockRedisInfo(gameRoomId, gameRoundId, userId string) *types.RedisLockInfo {
	return redistool.BuildRedisLockInfo(
		time.Duration(10)*time.Second,
		betFileKeyPrefix,
		betCancelLockLockPrefix,
		gameRoomId,
		gameRoundId,
		userId,
	)
}
