package rediskey

import (
	"sl.framework.com/game_server/conf"
	redistool "sl.framework.com/game_server/redis/redis_tool"
	"strings"
	"time"
)

type (
	// RedisInfo redis信息
	RedisInfo struct {
		Key    string        //key值 可以是hash表名 list名等
		Field  string        //当key为hash表名的时候 该值为hash下的一个field 非hash表则填空值
		Expire time.Duration //过期时间
	}

	// RedisInfo redis锁信息
	RedisLockInfo struct {
		Key    string        //key值
		Expire time.Duration //过期时间
		Owner  int64         //key持有者
	}
)

/**
 * getRedisExpireDuration
 * 查询缓存时间
 * 该缓存时间根据游戏的局时长设置 设置局时长的10倍
 * @param
 * @return time.Duration - 缓存时间
 */

func GetRedisExpireDuration() time.Duration {
	return time.Duration(10) * time.Minute
}

/**
 * BuildRedisInfo
 * 构造RedisInfo
 *
 * @param expiration time.Duration - 过期时间
 * @param field string - redis info中的field字段
 * @param keys - 可变长参数string类型 用于构造redis info中的key
 * @return *redisInfo - redis信息
 */

func BuildRedisInfo(expiration time.Duration, field string, keys ...string) *RedisInfo {
	var keyList []string
	keyList = append(keyList, redistool.GetKeyPrefix())
	for _, key := range keys {
		keyList = append(keyList, key)
	}

	hashRedisInfo := &RedisInfo{
		Key:    strings.Join(keyList, ":"),
		Field:  field,
		Expire: expiration,
	}
	return hashRedisInfo
}

/**
 * BuildRedisLockInfo
 * 构造RedisLockInfo
 *
 * @param expiration time.Duration - 过期时间
 * @param keys - 可变长参数string类型 用于构造redis lock info中的key
 * @return *redisLockInfo - redis锁信息
 */

func BuildRedisLockInfo(expiration time.Duration, keys ...string) *RedisLockInfo {
	var keyList []string
	keyList = append(keyList, redistool.GetKeyPrefix())
	for _, key := range keys {
		keyList = append(keyList, key)
	}

	hashRedisInfo := &RedisLockInfo{
		Key:    strings.Join(keyList, ":"),
		Expire: expiration,
		Owner:  conf.GetServerId(),
	}
	return hashRedisInfo
}
