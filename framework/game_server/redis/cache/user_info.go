package cache

import (
	"encoding/json"
	"fmt"
	"sl.framework.com/game_server/game/dao/redisdb"
	"sl.framework.com/game_server/game/service/type/dto"
	"sl.framework.com/game_server/redis/rediskey"
	"sl.framework.com/trace"
	"time"
)

const userInfoFileKeyPrefix = "UserInfoFile"

const (
	userInfoKeyPrefix     = "UserInfo"
	userInfoLockKeyPrefix = "UserInfoLock"
)

// ToTrialOn 根据用户信息的Type转化为TrialOn信息的'Y'或者'N'
func (u UserInfo) ToTrialOn() string {
	trialOn := "N"
	switch u.Type {
	case "Normal":
		trialOn = "Y"
	}

	return trialOn
}

// UserInfo redis缓存用户信息
type (
	UserInfo struct {
		Id           int64  `json:"id"`           //数据库字段:id id
		SiteUsername string `json:"siteUsername"` //数据库字段:site_username 站点用户名
		Username     string `json:"username"`     //数据库字段:username 用户名,需要根据站点用户名组合生成全局唯一
		Nickname     string `json:"nickname"`     //数据库字段:nickname 昵称
		Type         string `json:"type"`         //数据库字段:type 类型:游客 Visitor，试玩 Trial，正式 Normal ,测试 Test,带单 Capper,内部 Inner
		Status       string `json:"status"`       //数据库字段:status 状态:正常 Enable，登录锁定 Login_Lock ，游戏锁定 Game_Locked,重提锁定 Recharge_Locked
		BlackListOn  string `json:"blackListOn"`  //数据库字段:black_list_on 黑名单开启:是 Y,否 N
	}

	UserInfoCache struct {
		*BaseCache //基础缓存结构 封装解析与反解析函数

		Data *dto.UserDto //缓存中实际存储的结构

		TraceId string //用于日志跟踪

		/*操作缓存所需要的信息*/
		RoomId string
		UserId string
	}
)

var _ ICache = (*UserInfoCache)(nil)

/**
 * Get
 * 获取缓存信息 并将缓存信息解析到对应缓存结构中
 */

func (user *UserInfoCache) Get() (success bool) {
	var (
		val string
		err error
	)
	msgHeader := fmt.Sprintf("UserInfoCache traceId=%v, gameRoomId=%v, userId=%v",
		user.TraceId, user.RoomId, user.UserId)

	//获取缓存信息
	redisKey := user.redisKey()
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
	userInfo := new(dto.UserDto)
	if err = json.Unmarshal([]byte(val), userInfo); nil != err {
		trace.Error("%v, json unmarshal failed, table=%v, field=%v, val=%v, err=%v", msgHeader, redisKey.Key,
			redisKey.Field, val, err.Error())
		return
	}
	user.Data = userInfo
	trace.Debug("%v, userInfo=%+v", msgHeader, userInfo)

	return true
}

/**
 * Set
 * 设置缓存信息 将缓存结构序列化后以字符串的形式放入redis缓存中
 */

func (user *UserInfoCache) Set() (success bool) {
	var (
		val  int64
		err  error
		data []byte
	)

	msgHeader := fmt.Sprintf("SetUserInfo traceId=%v, gameRoomId=%v, userId=%v userinf=%+v",
		user.TraceId, user.RoomId, user.UserId, user.Data)
	if data, err = json.Marshal(user.Data); err != nil {
		trace.Error("%v, json marshal failed, error=%v", msgHeader, err.Error())
		return
	}
	redisKey := user.redisKey()
	if val, err = redisdb.HSet(redisKey.Key, redisKey.Field, string(data), redisKey.Expire); nil != err {
		trace.Error("%v, redis HSet failed, table=%v, field=%v, error=%v", msgHeader, redisKey.Key,
			redisKey.Field, err.Error())
		return
	}

	trace.Info("%v, val=%v, userInfo=%+v", msgHeader, val, user.Data)
	return true
}

/**
 * redisKey
 * 构建redisKey
 * @return *redisInfo - *redisInfo
 */

func (user *UserInfoCache) redisKey() *rediskey.RedisInfo {
	return rediskey.BuildRedisInfo(
		time.Duration(24)*time.Hour, //该hash保存房间内所有用户数据 应该一直存在
		user.UserId,                 //hash field字段
		userInfoFileKeyPrefix,
		userInfoKeyPrefix,
		user.RoomId,
	)
}

/**
 * redisLock
 * 构造锁信息
 *
 * @return *redisLockInfo - 锁信息
 */

func (user *UserInfoCache) redisLock() *rediskey.RedisLockInfo {
	return rediskey.BuildRedisLockInfo(
		time.Duration(10)*time.Second,
		userInfoFileKeyPrefix,
		userInfoLockKeyPrefix,
		user.RoomId,
		user.UserId,
	)
}
