package cache

import (
	"encoding/json"
	"fmt"
	"sl.framework.com/game_server/game/dao/redisdb"
	types "sl.framework.com/game_server/game/service/type"
	"sl.framework.com/game_server/redis/rediskey"
	"sl.framework.com/trace"
)

type GameRoomCache struct {
	*BaseCache //基础缓存结构 封装解析与反解析函数

	Data *types.RoomDetailedInfo //缓存中实际存储的结构

	TraceId string //用于日志跟踪
	RoomId  int64
}

/*
Get 获取缓存信息 并将缓存信息解析到对应缓存结构中
执行成功返回true 失败返回false
*/
func (c *GameRoomCache) Get() (success bool) {
	var (
		val string
		err error
	)
	msgHeader := fmt.Sprintf("GameRoomCache traceId=%v, roomId=%v",
		c.TraceId, c.RoomId)

	//获取缓存信息
	redisKey := c.redisKey()
	if val, err = redisdb.HGet(redisKey.Key, redisKey.Field); nil != err {
		trace.Error("%v, redis HGet failed, table=%v, field=%v, error=%v", msgHeader, redisKey.Key,
			redisKey.Field, err.Error())
		return
	}
	if len(val) <= 0 {
		trace.Notice("%v, redis HGet no data, table=%v, field=%v", msgHeader, redisKey.Key, redisKey.Field)
		return
	}

	//解析缓存信息
	roomData := new(types.RoomDetailedInfo)
	if err = json.Unmarshal([]byte(val), roomData); nil != err {
		trace.Error("%v, json unmarshal failed, table=%v, field=%v, val=%v, err=%v", msgHeader, redisKey.Key,
			redisKey.Field, val, err.Error())
		return
	}
	c.Data = roomData
	trace.Debug("%v, roomData=%+v", msgHeader, roomData)

	return true
}

/*
Set 设置缓存信息 将缓存结构序列化后以字符串的形式放入redis缓存中
执行成功返回true 失败返回false
*/
func (c *GameRoomCache) Set(srcDto *types.RoomDetailedInfo) (success bool) {
	var (
		val    int64
		err    error
		data   []byte
		curDto *types.RoomDetailedInfo
	)
	if srcDto == nil {
		curDto = c.Data
	} else {
		curDto = srcDto
	}
	msgHeader := fmt.Sprintf("GameRoomCache traceId=%v, roomId=%v",
		c.TraceId, c.RoomId)
	if data, err = json.Marshal(curDto); err != nil {
		trace.Error("%v, json marshal failed, error=%v", msgHeader, err.Error())
		return
	}
	redisKey := c.redisKey()
	if val, err = redisdb.HSet(redisKey.Key, redisKey.Field, string(data), redisKey.Expire); nil != err {
		trace.Error("%v, redis HSet failed, table=%v, field=%v, error=%v", msgHeader, redisKey.Key,
			redisKey.Field, err.Error())
		return
	}

	trace.Info("%v, val=%v, userInfo=%+v", msgHeader, val, c.Data)
	return true
}

/*
redisKey 缓存的key信息 只缓存结构内部使用 不对外可用
*/
func (c *GameRoomCache) redisKey() *rediskey.RedisInfo {
	keyInfo := rediskey.GetRoomDetailedInfoRedisInfo(c.RoomId)
	dst := &rediskey.RedisInfo{
		Key:    keyInfo.Key,
		Expire: keyInfo.Expire,
	}
	return dst
}

/*
redisLock 缓存的锁信息 只缓存结构内部使用 不对外可用
*/
func (c *GameRoomCache) redisLock() *rediskey.RedisLockInfo {
	return &rediskey.RedisLockInfo{}
}
