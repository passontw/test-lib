package rediskey

import "strconv"

const (
	gameWagerFileKeyPrefix     = "GameWager"
	gameWagerInfoKeyPrefix     = "GameWagerInfo"
	gameWagerInfoLockKeyPrefix = "HGameWagerInfoLock"
)

type GameWagerKey struct {
	Data *RedisInfo

	GameId int64
}

/**
 * FunctionName
 * 函数功能描述
 *
 * @param PARAM - 参数说明
 * @return RETURN - 返回值说明
 */

func (c *GameWagerKey) Key() *RedisInfo {
	c.Data = BuildRedisInfo(
		GetRedisExpireDuration(),
		"", //key-value缓存 不设置field字段
		gameWagerFileKeyPrefix,
		gameWagerInfoKeyPrefix,
		strconv.FormatInt(c.GameId, 10),
	)
	return c.Data
}
