package rediskey

import (
	"sl.framework.com/game_server/conf"
	"sl.framework.com/game_server/redis/redis_tool"
	"sl.framework.com/game_server/redis/types"
	"time"
)

const (
	// heartbeatFileKeyPrefix 避免key重复 每个文件中要使用一个key与其他文件区别
	heartbeatFileKeyPrefix = "Heartbeat"
)

const (
	/*
		HeartbeatStatus保存服务心跳状态的key
	*/
	heartbeatStatusPrefix     = "HeartbeatStatus"
	heartbeatStatusLockPrefix = "HeartbeatStatusLock"
)

// GetHeartbeatStatusRedisInfo 服务器心跳状态的key信息
// 如:G32FastBacGameServer:Heartbeat:HeartbeatStatus:65
func GetHeartbeatStatusRedisInfo() *types.RedisInfo {
	return redistool.BuildRedisInfo(
		time.Duration(conf.GetHeartbeatExpired())*time.Second,
		heartbeatFileKeyPrefix,
		heartbeatStatusPrefix,
	)
}

// GetHeartbeatStatusLockRedisInfo 服务器心跳状态的key信息锁
func GetHeartbeatStatusLockRedisInfo() *types.RedisLockInfo {
	return redistool.BuildRedisLockInfo(
		time.Duration(conf.GetHeartbeatExpired())*time.Second,
		heartbeatStatusLockPrefix,
		gameNextGameRoundIdLockPrefix)
}
