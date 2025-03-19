package redisdb

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"sl.framework.com/game_server/conf"
	snowflaker "sl.framework.com/game_server/conf/snow_flake_id"
	"sl.framework.com/game_server/error_code"
	redistool "sl.framework.com/game_server/redis/redis_tool"

	"sl.framework.com/game_server/redis/types"
	"sl.framework.com/trace"
	"strconv"
	"sync"
	"time"
)

var (
	redisInitOnce  sync.Once
	redisUniversal redis.UniversalClient //单节点和集群redis对象
)

/**
 * RedisClientClose
 * 关闭redis连接
 *
 * @param
 * @return
 */

func RedisClientClose() {
	if redisUniversal == nil {
		trace.Notice("RedisClientClose redisUniversal is nil")
		return
	}

	if err := redisUniversal.Close(); err != nil {
		trace.Error("RedisClientClose redisUniversal close failed, err=%v", err.Error())
		return
	}
	trace.Info("RedisClientClose redisUniversal close success")
}

func RedisClientInitOnce() (ok bool) {
	redisInitOnce.Do(func() {
		if redisUniversal == nil {
			ok = redisClientUniversalInit()
		} else {
			trace.Error("RedisClientInitOnce redisClient init already")
		}
	})

	return
}

// 初始化redis对象 适配集群和单节点模式的redis
// 单节点模式:[]string{"127.0.0.1:6379"}
// 集群模式:[]string{"127.0.0.1:7000", "127.0.0.1:7001", "127.0.0.1:7002"}
// 分片模式:[]string{"127.0.0.1:7000", "127.0.0.1:7001"}
func redisClientUniversalInit() (isSuccess bool) {
	addrs, _ := conf.GetRedisAddr()
	if len(addrs) == 0 {
		return
	}

	_, password := conf.GetRedisUserInfo()
	switch conf.GetRedisMode() {
	case conf.RedisModeSingle:
		redisUniversal = redis.NewClient(&redis.Options{
			Addr:     addrs[0],
			DB:       conf.GetRedisDb(),
			Password: password,
		})
	case conf.RedisModeCluster:
		redisUniversal = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    addrs,
			Password: password,
		})
	default:
		trace.Error("redisClientUniversalInit wrong redis mode =%v", conf.GetRedisMode())
		return
	}

	// 检测是否建立连接(需要传递上下文)
	timeoutCtx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()
	// 检测
	_, redisErr := redisUniversal.Ping(timeoutCtx).Result()
	if redisErr != nil {
		trace.Error("redisClientUniversalInit ping error=%s", redisErr.Error())
		return
	}

	// 设置redis key的公用前缀
	prefix := conf.GetKeyPrefix()
	if len(prefix) <= 0 {
		uniqueId := snowflaker.GetSnowFlakeInstance().GetUniqueId()
		prefix = strconv.FormatInt(uniqueId, 10)
		trace.Notice("GetKeyPrefix prefix is empty, generate one prefix=%v", prefix)
	}
	redistool.SetKeyPrefix(prefix)

	trace.Info("redisClientUniversalInit ping success")
	return true
}

/**
 * Lock
 * 为给定的信息加锁 加锁失败则尝试5次每次间隔20ms
 *
 * @param lock - 锁信息
 * @return bool - 返回是否加锁成功
 */

func Lock(lock *types.RedisLockInfo) bool {
	ctx := context.Background()
	bIsLocked := false
	//加锁失败则尝试加锁5次,每次间隔20ms
	tryInterval := 20
	for loop := 0; loop < 5; loop++ {
		bIsExist, err := redisUniversal.SetNX(ctx, lock.Key, strconv.FormatInt(lock.Owner, 10), lock.Expire).Result()
		if nil != err {
			trace.Error("redisdb setnx lock failed,key=%v, server id=%v, error=%v", lock.Key, lock.Owner, err.Error())
			time.Sleep(time.Millisecond * time.Duration(tryInterval))
			continue
		}
		//锁被占用
		if !bIsExist {
			trace.Error("redisdb setnx lock already locked, key=%v, server id=%v", lock.Key, lock.Owner)
			time.Sleep(time.Millisecond * time.Duration(tryInterval))
			continue
		}
		bIsLocked = true
		break
	}

	return bIsLocked
}

/**
 * TryLock
 * 为给定的信息加锁 加锁失败则直接返回
 *
 * @param lock - 锁信息
 * @return bool - 返回是否加锁成功
 */

func TryLock(lock *types.RedisLockInfo) (bIsLocked bool) {
	ctx := context.Background()

	bIsExist, err := redisUniversal.SetNX(ctx, lock.Key, strconv.FormatInt(lock.Owner, 10), lock.Expire).Result()
	if nil != err {
		trace.Error("redisdb setnx try lock failed, key=%v, server id=%v, error=%v", lock.Key, lock.Owner, err.Error())
		return
	}
	//锁被占用
	if !bIsExist {
		trace.Warning("redisdb setnx try lock already locked, key=%v, server id=%v", lock.Key, lock.Owner)
		return
	}
	bIsLocked = true

	return bIsLocked
}

// Unlock 解锁
func Unlock(lock *types.RedisLockInfo) {
	var (
		err      error
		serverId string
	)

	if serverId, err = Get(lock.Key); nil != err {
		trace.Error("Unlock key=%v get failed, error=%v", lock.Key, err.Error())
		return
	}
	//判断是否是锁的持有者
	if serverId != strconv.FormatInt(lock.Owner, 10) {
		trace.Notice("Unlock server id=%v try to unlock own server id=%v", lock.Owner, serverId)
		return
	}
	if _, err = Delete(lock.Key); nil != err {
		trace.Error("Unlock server id=%v, key=%v failed, error=%v", serverId, lock.Key, err.Error())
		return
	}
	return
}

// Set 增
func Set(key, value string, expired time.Duration) (string, error) {
	ctx := context.Background()
	val, err := redisUniversal.Set(ctx, key, value, expired).Result()
	if err != nil {
		trace.Error("Set key=%v, value=%v, err=%v", key, val, err.Error())
		return "", err
	}

	return val, nil
}

func Delete(key string) (int64, error) {
	ctx := context.Background()
	return redisUniversal.Del(ctx, key).Result()
}

// Get 查
func Get(key string) (string, error) {
	ctx := context.Background()
	result, err := redisUniversal.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil //key 不存在则认为是非错误,返回空字符串
		}
		trace.Error("Get key=%v, value=%v, err=%v", key, result, err.Error())
		return "", err
	}

	return result, nil
}

// GetKeys 使用keys命令查所有key Note:该命令在集群模式下不可靠
func GetKeys(key string) ([]string, error) {
	ctx := context.Background()
	result, err := redisUniversal.Keys(ctx, key).Result()
	if err != nil {
		trace.Error("Get key=%v, value=%v, err=%v", key, result, err.Error())
		return nil, err
	}

	return result, nil
}

// Scan  使用scan命令查所有key Note:该命令在集群模式下不可靠
func Scan(pattern string) (keys []string, err error) {
	var (
		allKeys []string
		cursor  = uint64(0)
		ctx     = context.Background()
	)

	for {
		allKeys, cursor, err = redisUniversal.Scan(ctx, cursor, pattern, 11).Result()
		if err != nil {
			trace.Error("redis scan error=%v", err.Error())
			break
		}
		keys = append(keys, allKeys...)
		// 如果游标为0，表示迭代结束
		if cursor == 0 {
			break
		}
	}
	trace.Info("redis scan, keys=%v", keys)
	return
}

// TTL 获取key的TTL
func TTL(key string) (time.Duration, error) {
	ctx := context.Background()

	ttl, err := redisUniversal.TTL(ctx, key).Result()
	if err != nil {
		trace.Error("TTL key=%v failed, error=%v", err.Error())
		return errcode.RedisErrorTTLNotValid, err
	}

	return ttl, nil
}

/**
 * HSetBatch
 * 批量设置redis的hash表数据 value以map的形式传入
 * 该方法设置的hash field是map的key field对应的值是map的value
 *
 * @param hashTable string - 参数说明
 * @param value map[string]string - 参数说明
 * @param expiration time.Duration - 参数说明
 * @return int64 - 设置成功的条数
 * @return error - 错误信息 如果设置没报错则返回nil
 */

func HSetBatch(hashTable string, value map[string]string, expiration time.Duration) (int64, error) {
	ctx := context.Background()
	val, err := redisUniversal.HSet(ctx, hashTable, value).Result()
	if err != nil {
		trace.Error("HSetBatch hashTable=%v, value=%+v, err=%v", hashTable, val, err.Error())
		return val, err
	}

	// 设置哈希表的过期时间
	err = redisUniversal.Expire(ctx, hashTable, expiration).Err()
	if err != nil {
		trace.Error("HSetBatch hashTable=%v, value=%+v, err=%v", err)
		return 0, err
	}
	return val, nil
}

/**
 * HSet
 * 设置redis的hash表field的值为value
 *
 * @param hashTable string - hash表明
 * @param field string - 表中的key
 * @param value string - value值
 * @param expiration time.Duration - 过期时间
 * @return int64 - 设置成功的条数
 * @return error - 错误信息 如果设置没报错则返回nil
 */

func HSet(hashTable, field, value string, expiration time.Duration) (int64, error) {
	ctx := context.Background()
	val, err := redisUniversal.HSet(ctx, hashTable, field, value).Result()
	if err != nil {
		trace.Error("HSet hashTable=%v, field=%v, value=%v, err=%v", hashTable, field, value, err.Error())
		return val, err
	}

	// 设置哈希表的过期时间
	err = redisUniversal.Expire(ctx, hashTable, expiration).Err()
	if err != nil {
		trace.Error("HSet hashTable=%v, field=%v, value=%v, err=%v", hashTable, field, value, err.Error())
		return 0, err
	}
	return val, nil
}

// HGet 读取redis的hash数据
func HGet(hashTable string, field string) (string, error) {
	ctx := context.Background()
	val, err := redisUniversal.HGet(ctx, hashTable, field).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil //key 不存在则认为是非错误,返回空字符串
		}
		trace.Error("HGet hashTable=%v, field=%v failed, err=%v", hashTable, field, err.Error())
		return "", err
	}

	return val, nil
}

// HGet 读取redis的hash数据
func HDel(hashTable string, field string) (int64, error) {
	ctx := context.Background()
	deletedCount, err := redisUniversal.HDel(ctx, hashTable, field).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, nil //key 不存在则认为是非错误,返回空字符串
		}
		trace.Error("HDel hashTable=%v, field=%v failed, err=%v", hashTable, field, err.Error())
		return 0, err
	}

	return deletedCount, nil
}

// HGetAll 读取hash表的所有字段和值
func HGetAll(hashTable string) (map[string]string, error) {
	dstMap := make(map[string]string)
	ctx := context.Background()
	val, err := redisUniversal.HGetAll(ctx, hashTable).Result()
	if err != nil {
		trace.Error("HGetAll hashTable=%v, failed, err=%v", hashTable, err.Error())
		return dstMap, err
	}
	for k, v := range val {
		dstMap[k] = v
	}
	return dstMap, nil
}

// LAppend 向list追加数据
func LAppend(list, val string, expiration time.Duration) error {
	ctx := context.Background()
	_, err := redisUniversal.RPush(ctx, list, val).Result()
	if nil != err {
		trace.Error("LAppend list name=%v, bet data=%v failed, error=%v", list, val, err.Error())
		return err
	}

	//设置list的过期时间
	if err = redisUniversal.Expire(ctx, list, expiration).Err(); nil != err {
		trace.Error("LAppend expire list name=%v, bet data=%v, expiration=%v failed, error=%v",
			list, val, expiration, err.Error())
		return err
	}

	return nil
}

/**
 * LLen
 * 获取List的长度
 *
 * @param key string - redis key
 * @return int64,error -
 */

func LLen(key string) (int64, error) {
	var (
		valuelen int64
		err      error
	)
	ctx := context.Background()
	valuelen, err = redisUniversal.LLen(ctx, key).Result()
	if nil != err {
		trace.Error("LLen list key=%v, failed, error=%v", key, err.Error())
		return valuelen, err
	}
	return valuelen, err
}

/**
 * LRemove
 * 删除list中的数据数据
 *
 * @param list string - list 名字
 * @param val string - 要删除元素
 * @return error - 错误信息 如果没错误则返回 nil
 */

func LRemove(list, val string) error {
	ctx := context.Background()
	_, err := redisUniversal.LRem(ctx, list, 0, val).Result()
	if nil != err {
		trace.Error("LRemove list name=%v, val=%v failed, error=%v", list, val, err.Error())
		return err
	}

	return nil
}

/**
 * LPop
 * 从列表的头部取一个元素
 *
 * @param key string - key
 * @return string, err -
 */

func LPop(key string) (string, error) {
	// 从列表左侧弹出元素
	ctx := context.Background()
	value, err := redisUniversal.LPop(ctx, key).Result()
	if err != nil {
		trace.Error("LPop key:%v Error:", key, err)
		return value, err
	}
	trace.Info("LPop key:%v success value:", key, value)
	return value, err
}

/**
 * LPick
 * 从列表pos取一个元素
 *
 * @param key string - key
 * @param pos int64 - 位置
 * @return string, err -
 */

func LPick(key string, pos int64) (string, error) {
	// 从列表左侧弹出元素
	ctx := context.Background()
	value, err := redisUniversal.LRange(ctx, key, pos, pos).Result()
	if err != nil || len(value) == 0 {
		trace.Error("LPick key:%v pos:%v Error:%v", key, pos, err)
		return "", err
	}
	trace.Info("LPick key:%v pos:%v success value:%v", key, pos, value)
	return value[0], err
}

// LAllMember 获取list中所有的数据
func LAllMember(list string) ([]string, bool) {
	ctx := context.Background()
	var (
		err error
		val []string
	)
	if val, err = redisUniversal.LRange(ctx, list, 0, -1).Result(); nil != err {
		trace.Error("LAllMember get all members list name=%v failed, error=%v", list, err.Error())
		return val, false
	}
	return val, true
}

/**
 * SetList
 * 把数据放到redis list内
 *
 * @param rediskey string - redis 存储用的key
 * @param data T - redis 存储用的数据

 * @return bool - 返回值说明
 */

func SetList[T any](rediskey string, duration time.Duration, data T) bool {
	//把数据转为string
	var (
		databyte []byte
		err      error
	)

	if databyte, err = json.Marshal(data); err != nil {
		trace.Error("SetList %v, json marshal failed, error=%v", data, err.Error())
		return false
	}

	if err = LAppend(rediskey, string(databyte), duration); nil != err {
		trace.Error("SetList %v, redis set failed, table=%v,  error=%v", data, rediskey,
			err.Error())
		return false
	}
	trace.Info("SetList key %v, val=%v", rediskey, data)

	return true
}

/**
 * GetListSize
 * 获取List的长度
 *
 * @param key string - redis的key值
 * @return int64 - List的长度
 */

func GetListSize(key string) int64 {
	var (
		vlen int64
		err  error
	)
	if vlen, err = LLen(key); err != nil {
		trace.Error("GetListSize key=%v, error=%v", key, err.Error())
		return vlen
	}
	return vlen
}

/**
 * RemoveFromList
 * 从redis list移除最前面元素
 *
 * @param redisKey string - redis key
 * @return err error - 错误
 */

func RemoveTopFromList(redisKey string) error {
	//把数据转为string
	var (
		val string
		err error
	)

	// 从列表左侧弹出元素
	val, err = LPop(redisKey)
	if err != nil {
		trace.Error("RemoveTopFromList key:%v Error:%v", redisKey, err)
		return err
	}
	trace.Info("RemoveTopFromList key %v, success value=%v", redisKey, val)
	return err
}
