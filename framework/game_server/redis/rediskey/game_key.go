package rediskey

import (
	"strconv"
	"time"
)

const (
	gameInfoFileKeyPrefix     = "GameInfoWager"
	gameInfoInfoKeyPrefix     = "GameInfoInfo"
	gameInfoInfoLockKeyPrefix = "HGameInfoLock"
)

type GameInfoKey struct {
	Data *RedisInfo

	GameId int64
}

/**
 * Key
 * 获取游戏详情key
 *
 * @param PARAM - 参数说明
 * @return RETURN - 返回值说明
 */

func (c *GameInfoKey) Key() *RedisInfo {
	c.Data = BuildRedisInfo(
		GetRedisExpireDuration(),
		"", //key-value缓存 不设置field字段
		gameInfoFileKeyPrefix,
		gameInfoInfoKeyPrefix,
		strconv.FormatInt(c.GameId, 10),
	)
	return c.Data
}

/**
 * Key
 * 获取游戏详情key
 *
 * @param PARAM - 参数说明
 * @return RETURN - 返回值说明
 */

func (c *GameInfoKey) LockKey() *RedisLockInfo {
	return BuildRedisLockInfo(
		time.Duration(10)*time.Second,
		gameInfoFileKeyPrefix,
		gameInfoInfoLockKeyPrefix,
		strconv.FormatInt(c.GameId, 10),
	)
}
