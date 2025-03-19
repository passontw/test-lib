package cache

import (
	"encoding/json"
	"fmt"
	"sl.framework.com/game_server/game/dao/redisdb"
	types "sl.framework.com/game_server/game/service/type"
	"sl.framework.com/game_server/redis/rediskey"
	"sl.framework.com/trace"
	"strconv"
	"time"
)

const dynamicOddsFileKeyPrefix = "DynamicOddsFile"

const (
	dynamicOddKeyPrefix     = "DynamicOdds"
	dynamicOddLockKeyPrefix = "DynamicOddsLock"
)

type DynamicOddsCache struct {
	*BaseCache //基础缓存结构 封装解析与反解析函数

	Data *types.DynamicOddsInfo //缓存动态赔率信息

	TraceId     string //用于日志跟踪
	GameId      int64
	RoomId      int64
	GameRoundId int64
	WagerId     int64
}

/*
Get 获取缓存信息 并将缓存信息解析到对应缓存结构中
执行成功返回true 失败返回false
*/
func (c *DynamicOddsCache) Get() (success bool) {
	var (
		val string
		err error
	)
	msgHeader := fmt.Sprintf("[动态赔率缓存] 获取当个玩法的赔率信息 DynamicOddsCache traceId=%v, roomId=%v",
		c.TraceId, c.RoomId)

	//获取缓存信息
	redisKey := c.redisKey()
	if val, err = redisdb.HGet(redisKey.Key, redisKey.Field); nil != err {
		trace.Error("%v,  redis HGet failed, table=%v, field=%v, error=%v", msgHeader, redisKey.Key,
			redisKey.Field, err.Error())
		return
	}
	if len(val) <= 0 {
		trace.Notice("%v, redis HGet no data, table=%v, field=%v", msgHeader, redisKey.Key, redisKey.Field)
		return
	}

	//解析缓存信息
	oddsData := new(types.DynamicOddsInfo)
	if err = json.Unmarshal([]byte(val), oddsData); nil != err {
		trace.Error("%v, json unmarshal failed, table=%v, field=%v, val=%v, err=%v", msgHeader, redisKey.Key,
			redisKey.Field, val, err.Error())
		return
	}
	c.Data = oddsData
	trace.Debug("%v, data=%+v", msgHeader, oddsData)

	return true
}

/*
GetAll 获取当前游戏当前房间当前局所有玩法的赔率信息
执行成功返回true 失败返回false
*/
func (c *DynamicOddsCache) GetAll() (oddsList []*types.DynamicOddsInfo, success bool) {
	var (
		val map[string]string
		err error
	)
	msgHeader := fmt.Sprintf("[动态赔率缓存] 获取全部玩法赔率信息 DynamicOddsCache traceId=%v, roomId=%v",
		c.TraceId, c.RoomId)

	//获取缓存信息
	redisKey := c.redisKey()
	if val, err = redisdb.HGetAll(redisKey.Key); nil != err {
		trace.Error("%v,  redis HGet failed, table=%v, field=%v, error=%v", msgHeader, redisKey.Key,
			redisKey.Field, err.Error())
		return
	}
	if len(val) <= 0 {
		trace.Notice("%v, redis HGet no data, table=%v, field=%v", msgHeader, redisKey.Key, redisKey.Field)
		return
	}

	for _, s2 := range val {
		dynamicOddsInfo := new(types.DynamicOddsInfo)
		if err = json.Unmarshal([]byte(s2), dynamicOddsInfo); nil != err {
			trace.Error("%v, json unmarshal failed, table=%v, field=%v, val=%v, err=%v", msgHeader, redisKey.Key,
				redisKey.Field, val, err.Error())
			return make([]*types.DynamicOddsInfo, 0), false
		}
		oddsList = append(oddsList, dynamicOddsInfo)
	}

	if len(oddsList) <= 0 {
		trace.Notice("%v, 获取动态赔率列表为空, table=%v, field=%v", msgHeader, redisKey.Key,
			redisKey.Field)
		return make([]*types.DynamicOddsInfo, 0), false
	}
	trace.Info("%v, oddsList=%v", msgHeader, oddsList)

	return oddsList, true
}

/*
Set 设置缓存信息 将缓存结构序列化后以字符串的形式放入redis缓存中
执行成功返回true 失败返回false
*/
func (c *DynamicOddsCache) Set(srcDto *types.DynamicOddsInfo) (success bool) {
	var (
		val    int64
		err    error
		data   []byte
		curDto *types.DynamicOddsInfo
	)
	if srcDto == nil {
		curDto = c.Data
	} else {
		curDto = srcDto
	}
	msgHeader := fmt.Sprintf("[动态赔率缓存] 设置 DynamicOddsCache traceId=%v, roomId=%v",
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

	trace.Info("%v, val=%v, data=%+v", msgHeader, val, srcDto)
	return true
}

/**
 * redisKey
 * 构建redisKey
 * @return *redisInfo - *redisInfo
 */

func (c *DynamicOddsCache) redisKey() *rediskey.RedisInfo {
	return rediskey.BuildRedisInfo(
		time.Duration(24)*time.Hour,      //该hash保存房间内所有用户数据 应该一直存在
		strconv.FormatInt(c.WagerId, 10), //hash field字段
		dynamicOddsFileKeyPrefix,
		dynamicOddKeyPrefix,
		strconv.FormatInt(c.GameId, 10),
		strconv.FormatInt(c.RoomId, 10),
		strconv.FormatInt(c.GameRoundId, 10),
	)
}

/**
 * redisLock
 * 构造锁信息
 *
 * @return *redisLockInfo - 锁信息
 */

func (c *DynamicOddsCache) redisLock() *rediskey.RedisLockInfo {
	return rediskey.BuildRedisLockInfo(
		time.Duration(10)*time.Second,
		dynamicOddsFileKeyPrefix,
		dynamicOddLockKeyPrefix,
		strconv.FormatInt(c.GameId, 10),
		strconv.FormatInt(c.RoomId, 10),
		strconv.FormatInt(c.GameRoundId, 10),
	)
}
