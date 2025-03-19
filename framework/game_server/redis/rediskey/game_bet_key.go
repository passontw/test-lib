package rediskey

import (
	"sl.framework.com/game_server/redis/redis_tool"
	"sl.framework.com/game_server/redis/types"
	"strconv"
	"time"
)

// 下注redis key
const (
	// betFileKeyPrefix 避免key重复 每个文件中要使用一个key与其他文件区别
	betFileKeyPrefix = "Bet"
)

const (
	betPrefix     = "Bet"     //投注
	betLockPrefix = "BetLock" //投注锁
)

/**
 * GetBetHRedisInfo
 * 下注redis信息
 *
 * @param gameRoomId int64 - 房间Id 用于构建hash表名
 * @param gameRoundId int64 - 局号Id 用于构建hash表名
 * @param userId int64 - 用户Id 用于当做hash表的field字段
 * @return *RedisInfo - 下注redis信息
 */

func GetBetHRedisInfo(gameRoomId, gameRoundId, userId string) *types.HashRedisInfo {
	return redistool.BuildHashRedisInfo(
		redistool.GetRedisExpireDuration(),
		userId,
		betFileKeyPrefix,
		betPrefix,
		gameRoomId,
		gameRoundId,
	)
}

/**
 * GetBetInfo
 * 下注redis信息
 *
 * @param gameRoomId int64 - 房间Id 用于构建hash表名
 * @param gameRoundId int64 - 局号Id 用于构建hash表名
 * @return *RedisInfo - 下注redis信息
 */

func GetBetInfo(gameRoomId, gameRoundId string) *types.HashRedisInfo {
	return redistool.BuildHashRedisInfo(
		redistool.GetRedisExpireDuration(),
		betFileKeyPrefix,
		betPrefix,
		gameRoomId,
		gameRoundId,
	)
}

/**
 * GetBetLockRedisInfo
 * 下注redis锁信息
 *
 * @param gameRoomId int64 - 房间Id
 * @param gameRoundId int64 - 局号Id
 * @param userId int64 - 用户Id
 * @return *RedisLockInfo - redis锁信息
 */

func GetBetLockRedisInfo(gameRoomId, gameRoundId, userId string) *types.RedisLockInfo {
	return redistool.BuildRedisLockInfo(
		time.Duration(10)*time.Second,
		betFileKeyPrefix,
		betLockPrefix,
		gameRoomId,
		gameRoundId,
		userId,
	)
}

const (
	betDenyPrefix         = "BetDeny"     //游戏局禁止投注
	betDenyLockLockPrefix = "BetDenyLock" //游戏局禁止投注锁
)

/**
 * GetBetDenyRedisInfo
 * 游戏局禁止投注redis信息
 *
 * @param gameRoomId int64 - 房间Id
 * @param gameRoundId int64 - 局号Id
 * @return *RedisInfo - redis信息
 */

func GetBetDenyRedisInfo(gameRoomId, gameRoundId int64) *types.RedisInfo {
	return redistool.BuildRedisInfo(
		redistool.GetRedisExpireDuration(),
		betFileKeyPrefix,
		betDenyPrefix,
		strconv.FormatInt(gameRoomId, 10),
		strconv.FormatInt(gameRoundId, 10),
	)
}

/**
 * GetBetDenyLockRedisInfo
 * 游戏局禁止投注redis锁信息
 *
 * @param gameRoomId int64 - 房间Id
 * @param gameRoundId int64 - 局号Id
 * @return *RedisLockInfo - redis锁信息
 */

func GetBetDenyLockRedisInfo(gameRoomId, gameRoundId int64) *types.RedisLockInfo {
	return redistool.BuildRedisLockInfo(
		time.Duration(10)*time.Second,
		betFileKeyPrefix,
		betDenyLockLockPrefix,
		strconv.FormatInt(gameRoomId, 10),
		strconv.FormatInt(gameRoundId, 10),
	)
}
