package gamelogic

import (
	"encoding/json"
	"fmt"
	"sl.framework.com/game_server/error_code"
	"sl.framework.com/game_server/game/dao/redisdb"
	types "sl.framework.com/game_server/game/service/type"
	"sl.framework.com/game_server/redis/rediskey"
	"sl.framework.com/game_server/rpc_client"
	"sl.framework.com/tool"
	"sl.framework.com/trace"
	"strconv"
)

/*
	Redis个人限红信息
	表名:	HUserLimitInfo{GM054247311DA}{userId}{currency}
	field:	{userId}{gameId}{gameWagerId}
	value:	{currency:"",minAmount:"",maxAmount:""}
*/

// getRedisUserLimitInfo 从redis获取用户限红信息
func getRedisUserLimitInfo(traceId, currency string, gameRoundId, userId, gameId, gameWagerId int64) *types.LimitInfo {
	var (
		val       string
		err       error
		userLimit = new(types.LimitInfo)
	)
	msgHeader := fmt.Sprintf("getRedisUserLimitInfo traceId=%v, gameRoundId=%v, userId=%v, currency=%v, gameId=%v, gameWagerId=%v",
		traceId, gameRoundId, userId, currency, gameId, gameWagerId)
	trace.Info("%v start", msgHeader)

	//查询数据
	redisInfo := rediskey.GetUserLimitHRedisInfoEx(currency, gameRoundId, userId, gameId, gameWagerId)
	if val, err = redisdb.HGet(redisInfo.HTable, redisInfo.Filed); nil != err {
		trace.Error("%v, redis HGet failed, table=%v, field=%v, error=%v", msgHeader, redisInfo.HTable,
			redisInfo.Filed, err.Error())
		return nil
	}

	if len(val) <= 0 {
		trace.Notice("%v, redis HGet no data, table=%v, field=%v", msgHeader, redisInfo.HTable, redisInfo.Filed)
		return nil
	}
	if err = json.Unmarshal([]byte(val), userLimit); nil != err {
		trace.Error("%v, json unmarshal failed, table=%v, field=%v, val=%v", msgHeader, redisInfo.HTable,
			redisInfo.Filed, val)
		return nil
	}

	return userLimit
}

// GetUserLimitInfo 获取个人限红,redis中不存在则从能力中心读取
func GetUserLimitInfo(traceId, currency string, gameRoundId, userId, gameId, gameWagerId int64) *types.LimitInfo {
	var userLimit *types.LimitInfo
	msgHeader := fmt.Sprintf("GetUserLimitInfo traceId=%v, gameRoundId=%v, userId=%v, currency=%v, "+
		"gameId=%v, gameWagerId=%v", traceId, gameRoundId, userId, currency, gameId, gameWagerId)

	trace.Info("%v start", msgHeader)
	if userLimit = getRedisUserLimitInfo(traceId, currency, gameRoundId, userId, gameId, gameWagerId); nil != userLimit {
		return userLimit
	}

	//redis中不存在则重新从平台中心读取并存储到redis
	NewEventUserLimit(traceId, currency, gameRoundId, userId).HandleEvent()

	//再次从redis中读取
	if userLimit = getRedisUserLimitInfo(traceId, currency, gameRoundId, userId, gameId, gameWagerId); nil != userLimit {
		return userLimit
	}

	return userLimit
}

// NewEventUserLimit 创建一个个人限红信息事件对象
func NewEventUserLimit(traceId, currency string, gameRoundId, userId int64) *EventUserLimitInfo {
	return &EventUserLimitInfo{
		RedisEvent: RedisEvent{
			traceId: traceId,
		},
		gameRoundId: gameRoundId,
		userId:      userId,
		currency:    currency,
		msgHeader: fmt.Sprintf("NewEventUserLimit HandleEvent traceId=%v, gameRoundId=%v, userId=%v, "+
			"currency=%v", traceId, gameRoundId, userId, currency),
	}
}

var _ IEvent = (*EventUserLimitInfo)(nil)

// EventUserLimitInfo 更新个人限红信息 包括内存信息和redis信息
type EventUserLimitInfo struct {
	RedisEvent
	currency    string
	gameRoundId int64  //游戏局号 用于创建Redis Hash表
	userId      int64  //个人ID 用于
	msgHeader   string //打印信息头
}

func (u *EventUserLimitInfo) HandleEvent() {
	var (
		userLimitInfo *types.UserBetLimitInfo
		ret           int
	)
	pDog := tool.NewWatcher(u.msgHeader)
	if userLimitInfo, ret = rpcreq.GetUserLimitRequest(u.traceId, u.currency, u.userId, u.gameRoundId); errcode.ErrorOk != ret {
		trace.Error("%v, failed error code=%v", u.msgHeader, ret)
		return
	}
	pDog.Stop()

	NewUserLimitInfoStoreEvent(u.traceId, u.currency, u.gameRoundId, u.userId, userLimitInfo).HandleEvent()
	pDog.Stop()
}

// NewUserLimitInfoStoreEvent 创建一个更新个人信息事件对象
func NewUserLimitInfoStoreEvent(traceId, currency string, gameRoundId, userId int64,
	ptrUserLimitInfo *types.UserBetLimitInfo) *UserLimitInfoStoreEvent {
	return &UserLimitInfoStoreEvent{
		RedisEvent: RedisEvent{
			traceId: traceId,
		},
		gameRoundId:         gameRoundId,
		userId:              userId,
		currency:            currency,
		ptrUserBetLimitInfo: ptrUserLimitInfo,
		msgHeader: fmt.Sprintf("UserLimitInfoStoreEvent HandleEvent traceId=%v, gameRoundId=%v,"+
			"currency=%v", traceId, gameRoundId, currency),
	}
}

var _ IEvent = (*UserLimitInfoStoreEvent)(nil)

// UserLimitInfoStoreEvent 存储个人限红信息事件
type UserLimitInfoStoreEvent struct {
	RedisEvent
	currency            string                  //货币类型
	gameRoundId         int64                   //游戏局号 用于创建Redis Hash表
	userId              int64                   //个人ID 用于
	ptrUserBetLimitInfo *types.UserBetLimitInfo //用户个人限红信息
	msgHeader           string                  //打印信息头
}

func (r *UserLimitInfoStoreEvent) HandleEvent() {
	if r.ptrUserBetLimitInfo == nil {
		trace.Error("%v, ptrUserBetLimitInfo is nil", r.msgHeader)
		return
	}
	if len(r.ptrUserBetLimitInfo.BetLimitRuleList) == 0 {
		trace.Notice("%v, ptrUserBetLimitInfo limit rule is empty.", r.msgHeader)
		return
	}

	var (
		err          error
		val          int64
		jsonData     []byte
		limitInfoMap = make(map[string]string, 64)
	)
	for _, item := range r.ptrUserBetLimitInfo.BetLimitRuleList {
		limit := types.LimitInfo{
			Currency:  item.Currency,
			MinAmount: item.MinAmount,
			MaxAmount: item.MaxAmount,
		}
		if jsonData, err = json.Marshal(limit); nil != err {
			trace.Error("%v, json marshal failed, error=%v, limit=%+v", r.msgHeader, err.Error(), limit)
			continue
		}
		//workaround:从能力中心获取int64时候传过来的是string 所以在此转化下
		gameWagerId, _ := strconv.ParseInt(item.GameWagerId, 10, 64)
		gameId, _ := strconv.ParseInt(item.GameId, 10, 64)
		redisInfo := rediskey.GetUserLimitHRedisInfoEx(r.currency, r.gameRoundId, r.userId, gameId, gameWagerId)
		limitInfoMap[redisInfo.Filed] = string(jsonData)
	}
	if len(limitInfoMap) <= 0 {
		trace.Notice("%v, no data in map", r.msgHeader)
		return
	}
	redisLockInfo := rediskey.GetUserLimitLockHRedisInfo(r.gameRoundId, r.userId)
	if !redisdb.Lock(redisLockInfo) {
		trace.Error("%v, redis lock failed", r.msgHeader)
		return
	}
	defer redisdb.Unlock(redisLockInfo)

	redisInfo := rediskey.GetUserLimitHRedisInfo(r.gameRoundId, r.userId)
	if val, err = redisdb.HSetBatch(redisInfo.HTable, limitInfoMap, redisInfo.Expire); nil != err {
		trace.Error("%v, HSetBatch failed, val=%v, error=%v", r.msgHeader, val, err.Error())
		return
	}
	trace.Info("%v, HSetBatch success, key=%v, result=%v", r.msgHeader, redisInfo.HTable, val)
}
