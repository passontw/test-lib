package cache

import (
	"sl.framework.com/game_server/redis/rediskey"
)

/**
 * ICache 缓存接口
 * 自定义缓存结构 要提供读取、设置缓存、锁的能力
 */

type ICache interface {
	/*
		Get 获取缓存信息 并将缓存信息解析到对应缓存结构中
		执行成功返回true 失败返回false
	*/
	Get() bool

	/*
		Set 设置缓存信息 将缓存结构序列化后以字符串的形式放入redis缓存中
		执行成功返回true 失败返回false
	*/
	Set() bool

	/*
		redisKey 缓存的key信息 只缓存结构内部使用 不对外可用
	*/
	redisKey() *rediskey.RedisInfo

	/*
	   redisLock 缓存的锁信息 只缓存结构内部使用 不对外可用
	*/
	redisLock() *rediskey.RedisLockInfo
}
