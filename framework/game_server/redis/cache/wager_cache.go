package cache

import (
	"encoding/json"
	"fmt"
	"sl.framework.com/game_server/game/dao/redisdb"
	"sl.framework.com/game_server/game/service/type/dto"
	"sl.framework.com/game_server/redis/rediskey"
	"sl.framework.com/trace"
)

type WagerCache struct {
	*BaseCache //基础缓存结构 封装解析与反解析函数

	Data []*dto.GameWagerDTO //缓存中实际存储的结构

	TraceId    string //用于日志跟踪
	GameId     int64
	GameRoomId int64
}

/*
Get 获取缓存信息 并将缓存信息解析到对应缓存结构中
执行成功返回true 失败返回false
*/
func (c *WagerCache) Get() (success bool) {
	var (
		val string
		err error
	)
	msgHeader := fmt.Sprintf("WagerCache traceId=%v, gameId=%v, GameRoomId=%v",
		c.TraceId, c.GameId, c.GameRoomId)

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
	gameWagerDto := new([]*dto.GameWagerDTO)
	if err = json.Unmarshal([]byte(val), gameWagerDto); nil != err {
		trace.Error("%v, json unmarshal failed, table=%v, field=%v, val=%v, err=%v", msgHeader, redisKey.Key,
			redisKey.Field, val, err.Error())
		return
	}
	c.Data = *gameWagerDto
	trace.Debug("获取房间赔率信息 %v, gameWagerDto=%+v", msgHeader, gameWagerDto)

	return true
}

/*
Set 设置缓存信息 将缓存结构序列化后以字符串的形式放入redis缓存中
执行成功返回true 失败返回false
*/
func (c *WagerCache) Set(srcDto []*dto.GameWagerDTO) (success bool) {
	var (
		val    int64
		err    error
		data   []byte
		curDto []*dto.GameWagerDTO
	)
	if srcDto == nil {
		curDto = c.Data
	} else {
		curDto = srcDto
	}
	msgHeader := fmt.Sprintf("WagerCache traceId=%v, gameId=%v, GameRoomId=%v",
		c.TraceId, c.GameId, c.GameRoomId)
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

	trace.Info("存储房间赔率信息 %v, val=%v, userInfo=%+v", msgHeader, val, c.Data)
	return true
}

/*
redisKey 缓存的key信息 只缓存结构内部使用 不对外可用
*/
func (c *WagerCache) redisKey() *rediskey.RedisInfo {
	cKey := rediskey.GameWagerKey{GameId: c.GameId}
	return cKey.Key()
}

/*
redisLock 缓存的锁信息 只缓存结构内部使用 不对外可用
*/
func (c *WagerCache) redisLock() *rediskey.RedisLockInfo {
	return &rediskey.RedisLockInfo{}
}
