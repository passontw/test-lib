package redistool

import (
	"sl.framework.com/game_server/conf"
	"sl.framework.com/game_server/redis/types"
	"sl.framework.com/trace"
	"strings"
	"time"
)

/**
 * BuildRedisInfo
 * 构造Redis信息
 *
 * @param expiration - 过期时间
 * @param keys - 可变长参数string类型
 * @return *RedisInfo - 锁信息
 */

func BuildRedisInfo(expiration time.Duration, keys ...string) *types.RedisInfo {
	var keyList []string
	keyList = append(keyList, types.ServerRedisKeyPrefix)
	for _, key := range keys {
		keyList = append(keyList, key)
	}

	redisInfo := &types.RedisInfo{
		Key:    strings.Join(keyList, ":"),
		Expire: expiration,
	}
	return redisInfo
}

/**
 * BuildHashRedisInfo
 * 构造Hash Redis信息
 *
 * @param expiration - 过期时间
 * @param keys - 可变长参数string类型
 * @return *RedisInfo - 锁信息
 */

func BuildHashRedisInfo(expiration time.Duration, field string, keys ...string) *types.HashRedisInfo {
	var keyList []string
	keyList = append(keyList, types.ServerRedisKeyPrefix)
	for _, key := range keys {
		keyList = append(keyList, key)
	}

	hashRedisInfo := &types.HashRedisInfo{
		HTable: strings.Join(keyList, ":"),
		Filed:  field,
		Expire: expiration,
	}
	return hashRedisInfo
}

// SetKeyPrefix 设置redis key的公用前缀
func SetKeyPrefix(prefix string) {
	types.ServerRedisKeyPrefix = prefix
	trace.Info("SetServerPrefix prefix=%v", prefix)
}

/**
 * GetKeyPrefix
 * 获取缓存前缀
 *
 * @param
 * @return string - 缓存前缀
 */

func GetKeyPrefix() string {
	return types.ServerRedisKeyPrefix
}

// SetRedisExpire 设置redis key过期时间
func SetRedisExpire(duration int32) {
	if time.Duration(duration)*time.Second < types.BacDuration {
		trace.Notice("SetRedisExpire too small duration the duration=%v", duration)
		return
	}

	//设置为3倍 结算的时候需要room odd数据 从设置odd到使用odd 有时候会超过60s
	types.RedisExpireDurationRWMutex.Lock()
	types.RedisExpireDuration = time.Duration(duration*3) * time.Second
	types.RedisExpireDurationRWMutex.Unlock()
	//trace.Info("SetRedisExpire duration=%v, redisExpireDuration=%v", duration, redisExpireDuration)
}

func GetRedisExpireDuration() time.Duration {
	types.RedisExpireDurationRWMutex.RLock()
	defer types.RedisExpireDurationRWMutex.RUnlock()
	return types.RedisExpireDuration
}

/**
 * BuildRedisLockInfo
 * 构造Redis锁信息
 *
 * @param expiration - 过期时间
 * @param keys - 可变长参数string类型
 * @return *RedisLockInfo - 锁信息
 */

func BuildRedisLockInfo(expiration time.Duration, keys ...string) *types.RedisLockInfo {
	redisLockInfo := &types.RedisLockInfo{
		RedisInfo: BuildRedisInfo(expiration, keys...),
		Owner:     conf.GetServerId(),
	}

	return redisLockInfo
}
