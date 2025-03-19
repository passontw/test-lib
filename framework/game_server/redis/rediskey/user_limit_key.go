package rediskey

import (
	"fmt"
	"sl.framework.com/game_server/redis/redis_tool"
	"sl.framework.com/game_server/redis/types"
	"strconv"
)

/*
个人限红信息
	表名:	HUserLimitInfo{GM054247311DA}
	field:	{UserId}{gameId}{gameWagerId}
	value:	{currency:"",minAmount:"",maxAmount:""}
*/

const (
	// userLimitFileKeyPrefix 避免key重复 每个文件中要使用一个key与其他文件区别
	userLimitFileKeyPrefix = "UserLimit"
)

const (
	userLimitInfoKeyPrefix     = "HUserLimitInfo"
	userLimitInfoLockKeyPrefix = "HUserLimitInfoLock"
)

// GetUserLimitHRedisInfo 获取用户限红的hash信息,只获取Hash表名和过期时间
// 如: hash表名 userLimitInfoKeyPrefix:UserLimit:HUserLimitInfo:GM054247311DA:{UserId}:{Currency}
func GetUserLimitHRedisInfo(gameRoundId, userId int64) *types.HashRedisInfo {
	return redistool.BuildHashRedisInfo(
		redistool.GetRedisExpireDuration(),
		"",
		userLimitFileKeyPrefix,
		userLimitInfoKeyPrefix,
		strconv.FormatInt(gameRoundId, 10),
		strconv.FormatInt(userId, 10),
	)
}

// GetUserLimitHRedisInfoEx 获取用户限红的hash信息,获取表名 Field字段 过期时间
// 如: hash表名 userLimitInfoKeyPrefix:HUserLimit:HUserLimitInfo:GM054247311DA:{UserId}
func GetUserLimitHRedisInfoEx(currency string, gameRoundId, userId, gameId, gameWagerId int64) *types.HashRedisInfo {
	return redistool.BuildHashRedisInfo(
		redistool.GetRedisExpireDuration(),
		fmt.Sprintf("%v:%v:%v:%v", userId, currency, gameId, gameWagerId),
		userLimitFileKeyPrefix,
		userLimitInfoKeyPrefix,
		strconv.FormatInt(gameRoundId, 10),
		strconv.FormatInt(userId, 10),
	)
}

// GetUserLimitHRedisInfoUserId 获取用户限红的hash信息,获取表名 Field字段 过期时间
// 如: hash表名 userLimitInfoKeyPrefix:HUserLimit:HUserLimitInfo:GM054247311DA:{UserId}
func GetUserLimitHRedisInfoUserId(userId int64) *types.HashRedisInfo {
	return redistool.BuildHashRedisInfo(
		redistool.GetRedisExpireDuration(),
		fmt.Sprintf("%v", userId),
		userLimitFileKeyPrefix,
		userLimitInfoKeyPrefix,
		strconv.FormatInt(userId, 10),
	)
}

// GetUserLimitLockHRedisInfo 用户限红信息分布式锁
// 玩家个人限红更新是逐个更新所以将userId作为Key
func GetUserLimitLockHRedisInfo(gameRoundId, userId int64) *types.RedisLockInfo {
	return redistool.BuildRedisLockInfo(
		types.RedisLockExpireDuration,
		userLimitFileKeyPrefix,
		userLimitInfoLockKeyPrefix,
		strconv.FormatInt(gameRoundId, 10),
		strconv.FormatInt(userId, 10),
	)
}
