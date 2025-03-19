package cache

import (
	"encoding/json"
	"fmt"
	"sl.framework.com/game_server/game/dao/redisdb"
	types "sl.framework.com/game_server/game/service/type"
	"sl.framework.com/game_server/game/service/type/dto"
	"sl.framework.com/game_server/redis/rediskey"
	"sl.framework.com/trace"
)

type UserLimitCache struct {
	*BaseCache //基础缓存结构 封装解析与反解析函数

	Data *types.UserBetLimitInfo //缓存中实际存储的结构

	TraceId string //用于日志跟踪
	UserId  int64
}

/*
Get 获取缓存信息 并将缓存信息解析到对应缓存结构中
执行成功返回true 失败返回false
*/
func (c *UserLimitCache) Get() (success bool) {
	var (
		val string
		err error
	)
	msgHeader := fmt.Sprintf("UserLimitCache traceId=%v, gameId=%v",
		c.TraceId, c.UserId)

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
	userSessionDto := new(types.UserBetLimitInfo)
	if err = json.Unmarshal([]byte(val), userSessionDto); nil != err {
		trace.Error("%v, json unmarshal failed, table=%v, field=%v, val=%v, err=%v", msgHeader, redisKey.Key,
			redisKey.Field, val, err.Error())
		return
	}
	c.Data = userSessionDto
	trace.Debug("%v, userInfo=%+v", msgHeader, userSessionDto)

	return true
}

/*
GetALL 获取缓存信息 并将缓存信息解析到对应缓存结构中
执行成功返回true 失败返回false
*/
func (c *UserLimitCache) GetALL() (sessionList []*dto.UserSessionDTO, success bool) {
	var (
		val string
		err error
	)
	msgHeader := fmt.Sprintf("UserLimitCache traceId=%v, gameId=%v",
		c.TraceId, c.UserId)

	//获取缓存信息
	redisKey := c.redisKey()
	if val, err = redisdb.HGet(redisKey.Key, ""); nil != err {
		trace.Error("%v, redis HGet failed, table=%v, field=%v, error=%v", msgHeader, redisKey.Key,
			redisKey.Field, err.Error())
		return
	}
	if len(val) <= 0 {
		trace.Notice("%v, redis HGet no data, table=%v, field=%v", msgHeader, redisKey.Key, redisKey.Field)
		return
	}

	//解析缓存信息
	userSessionDto := new([]*dto.UserSessionDTO)
	if err = json.Unmarshal([]byte(val), userSessionDto); nil != err {
		trace.Error("%v, json unmarshal failed, table=%v, field=%v, val=%v, err=%v", msgHeader, redisKey.Key,
			redisKey.Field, val, err.Error())
		return
	}
	trace.Debug("%v, userInfo=%+v", msgHeader, *userSessionDto)

	return *userSessionDto, true
}

/*
Set 设置缓存信息 将缓存结构序列化后以字符串的形式放入redis缓存中
执行成功返回true 失败返回false
*/
func (c *UserLimitCache) Set(srcDto *types.UserBetLimitInfo) (success bool) {
	var (
		val    int64
		err    error
		data   []byte
		curDto *types.UserBetLimitInfo
	)
	if srcDto == nil {
		curDto = c.Data
	} else {
		curDto = srcDto
	}
	msgHeader := fmt.Sprintf("UserLimitCache traceId=%v, gameId=%v",
		c.TraceId, c.UserId)
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
func (c *UserLimitCache) redisKey() *rediskey.RedisInfo {
	keyInfo := rediskey.GetUserLimitHRedisInfoUserId(c.UserId)
	dst := &rediskey.RedisInfo{
		Key:    keyInfo.HTable,
		Field:  keyInfo.Filed,
		Expire: keyInfo.Expire,
	}
	return dst
}

/*
redisLock 缓存的锁信息 只缓存结构内部使用 不对外可用
*/
func (c *UserLimitCache) redisLock() *rediskey.RedisLockInfo {
	return &rediskey.RedisLockInfo{}
}
