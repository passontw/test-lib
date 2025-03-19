package rediskey

import "strconv"

const (
	userSessionFileKeyPrefix = "UserSession"
	userSessionKeyPrefix     = "UserSessionInfo"
	userSessionLockKeyPrefix = "HUserSessionLock"
)

type UserSessionKey struct {
	Data *RedisInfo

	GameRoomId int64
	UserId     int64
}

/**
 * Key
 * 获取游戏详情key
 *
 * @param PARAM - 参数说明
 * @return RETURN - 返回值说明
 */

func (c *UserSessionKey) Key() *RedisInfo {
	c.Data = BuildRedisInfo(
		GetRedisExpireDuration(),
		strconv.FormatInt(c.UserId, 10), //key-value缓存 不设置field字段
		userSessionFileKeyPrefix,
		userSessionKeyPrefix,
		strconv.FormatInt(c.GameRoomId, 10),
	)
	return c.Data
}
