package rediskey

import (
	redistool "sl.framework.com/game_server/redis/redis_tool"
	"sl.framework.com/game_server/redis/types"
	"strconv"
)

const (
	// roomLimitFileKeyPrefix 避免key重复 每个文件中要使用一个key与其他文件区别
	roomLimitFileKeyPrefix = "RoomLimit"
)

const (
	roomLimitInfoKeyPrefix     = "HRoomLimitInfo"
	roomLimitInfoLockKeyPrefix = "HRoomLimitInfoLock"
)

// GetRoomLimitHKey 房间限红信息 获取表名和过期时间
// 桌台限红以及桌台赔率信息是一次性设置进去所以HTable只使用gameRoundId
func GetRoomLimitHKey(gameId int64) *types.HashRedisInfo {
	return redistool.BuildHashRedisInfo(
		redistool.GetRedisExpireDuration(),
		"",
		roomLimitFileKeyPrefix,
		roomLimitInfoKeyPrefix,
		strconv.FormatInt(gameId, 10),
	)
}
