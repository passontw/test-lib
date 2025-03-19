package cache

import (
	"encoding/json"
	"fmt"
	"sl.framework.com/game_server/game/dao/redisdb"
	"sl.framework.com/game_server/game/service/type/dto"
	"sl.framework.com/game_server/redis/rediskey"
	"sl.framework.com/trace"
	"strconv"
	"time"
)

const orderCheckBalanceFileKeyPrefix = "OrderCheckBalanceFile"

const (
	orderCheckBalanceKeyPrefix     = "OrderCheckBalance"
	orderCheckBalanceLockKeyPrefix = "OrderCheckBalanceLock"
)

// 此缓存至用户校验用户余额，不做状态记录
type OrderCheckBalanceCache struct {
	Data        []*dto.BetDTO
	TraceId     string
	GameId      int
	GameRoomId  int64
	GameRoundId int64
	UserId      int64
}

/**
 * Set
 * 存入数据
 *
 * @param srcDto []*dto.BetDTO - 原始数据
 * @return success bool - 返回是否设置成功
 */

func (c *OrderCheckBalanceCache) Set(srcDto []*dto.BetDTO) (success bool) {
	var (
		val        int64
		err        error
		data       []byte
		curDtoList []*dto.BetDTO
		dataJson   string
		userDto    []*dto.BetDTO
	)
	if srcDto == nil {
		curDtoList = c.Data
	} else {
		curDtoList = srcDto
	}
	msgHeader := fmt.Sprintf("[校验额度注单缓存] OrderCheckBalanceCache 设置 traceId=%v, roomId=%v,roundId=%v",
		c.TraceId, c.GameRoomId, c.GameRoundId)

	redisKey := c.redisUserIdKey()
	//userid为key进行先取再存
	//userid为key roomId:roundId:userId 为field 进行先取再存
	if dataJson, err = redisdb.HGet(redisKey.Key, redisKey.Field); err != nil {
		trace.Error("%v, redis HGet failed, table=%v, field=%v, error=%v", msgHeader, redisKey.Key,
			"", err.Error())
		return
	}
	if err = json.Unmarshal([]byte(dataJson), &userDto); err != nil {
		trace.Error("%v, json unmarshal failed, error=%v", msgHeader, err.Error())
		return
	}
	curDtoList = append(curDtoList, userDto...)
	if data, err = json.Marshal(curDtoList); err != nil {
		trace.Error("%v, json 序列化 userId表 roomId:roundId:userId 字段 失败, error=%v", msgHeader, err.Error())
		return
	}
	if val, err = redisdb.HSet(redisKey.Key, redisKey.Field, string(data), redisKey.Expire); nil != err {
		trace.Error("%v, redis 设置 userId表 roomId:roundId:userId 失败, table=%v, field=%v, error=%v", msgHeader, redisKey.Key,
			redisKey.Field, err.Error())
		return
	}
	trace.Info("%v, val=%v, data=%+v", msgHeader, val, c.Data)
	return true
}

/**
 * GetFromUserId
 * 根据用户名查询所有注单
 *
 * @param PARAM - 参数说明
 * @return RETURN - 返回值说明
 */

func (c *OrderCheckBalanceCache) GetFromUserId() []*dto.BetDTO {
	var (
		err        error
		curDtoList []*dto.BetDTO
		dataJson   string
	)

	msgHeader := fmt.Sprintf("[校验额度注单缓存] OrderCheckBalanceCache 根据用户名获取 traceId=%v, roomId=%v,roundId=%v",
		c.TraceId, c.GameRoomId, c.GameRoundId)

	//需要存两次
	redisKey := c.redisUserIdKey()
	//userid为key进行先取再存
	if dataJson, err = redisdb.HGet(redisKey.Key, ""); err != nil {
		trace.Error("%v, redis HGet failed, table=%v, field=%v, error=%v", msgHeader, redisKey.Key,
			"", err.Error())
		return make([]*dto.BetDTO, 0)
	}
	if err = json.Unmarshal([]byte(dataJson), &curDtoList); err != nil {
		trace.Error("%v, json 反序列化 失败, error=%v", msgHeader, err.Error())
		return make([]*dto.BetDTO, 0)
	}
	return curDtoList
}

/**
 * RemoveFromField
 * 删除roomId:roundId:userId 为field的记录
 *
 * @param PARAM - 参数说明
 * @return RETURN - 返回值说明
 */

func (c *OrderCheckBalanceCache) RemoveFromField() (success bool) {
	var (
		err   error
		count int64
	)
	msgHeader := fmt.Sprintf("[校验额度注单缓存] OrderCheckBalanceCache 删除roomId:roundId:userId 为field的记录 traceId=%v, roomId=%v,roundId=%v,userId=%v",
		c.TraceId, c.GameRoomId, c.GameRoundId, c.UserId)
	redisKey := c.redisUserIdKey()
	if count, err = redisdb.HDel(redisKey.Key, redisKey.Field); err != nil {
		trace.Error("%v, json 删除失败, error=%v", msgHeader, err.Error())
		return
	}
	trace.Debug("%v 删除数量=%v", count)
	return true
}

/**
 * redisUserIdKey
 * 构建以UserId为key的数据 roomId:roundId:userId为field的数据
 * @return *redisInfo - *redisInfo
 */

func (c *OrderCheckBalanceCache) redisUserIdKey() *rediskey.RedisInfo {
	return rediskey.BuildRedisInfo(
		time.Duration(24)*time.Hour,                                    //该hash保存redis内所有用户数据 应该一直存在
		fmt.Sprintf("%v:%v:%v", c.GameRoomId, c.GameRoundId, c.UserId), //hash field字段                          //hash field字段
		orderCheckBalanceFileKeyPrefix,
		orderCheckBalanceKeyPrefix,
		strconv.FormatInt(c.UserId, 10),
	)
}

/**
 * redisLock
 * 构造锁信息
 *
 * @return *redisLockInfo - 锁信息
 */

func (c *OrderCheckBalanceCache) redisLock() *rediskey.RedisLockInfo {
	return rediskey.BuildRedisLockInfo(
		time.Duration(10)*time.Second,
		orderCheckBalanceFileKeyPrefix,
		orderCheckBalanceLockKeyPrefix,
		strconv.FormatInt(c.GameRoomId, 10),
		strconv.FormatInt(c.GameRoundId, 10),
		strconv.FormatInt(c.UserId, 10),
	)
}
