package rediskey

import (
	"sl.framework.com/game_server/redis/redis_tool"
	"sl.framework.com/game_server/redis/types"
	"strconv"
	"time"
)

// 开始游戏 结束游戏相关redis key
const (
	// gameFileKeyPrefix 避免key重复 每个文件中要使用一个key与其他文件区别
	gameFileKeyPrefix = "Game"
)
const (
	gameCommandPrefix     = "GameCommand"
	gameCommandLockPrefix = "GameCommandLock"
)

// GetGameCommandRedisInfo 游戏命令redis信息
func GetGameCommandRedisInfo(gameRoomId, gameRoundId int64) *types.RedisInfo {
	return redistool.BuildRedisInfo(
		time.Duration(60)*time.Second,
		gameFileKeyPrefix,
		gameCommandPrefix,
		strconv.FormatInt(gameRoomId, 10),
		strconv.FormatInt(gameRoundId, 10),
	)
}

// GetGameCommandLockRedisInfo 游戏命令redis锁
func GetGameCommandLockRedisInfo(gameRoomId, gameRoundId int64) *types.RedisLockInfo {
	return redistool.BuildRedisLockInfo(
		time.Duration(10)*time.Second,
		gameFileKeyPrefix,
		gameCommandLockPrefix,
		strconv.FormatInt(gameRoomId, 10),
		strconv.FormatInt(gameRoundId, 10),
	)
}
