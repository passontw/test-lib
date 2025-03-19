package cache

import (
	"encoding/json"
	"fmt"
	"sl.framework.com/game_server/game/dao/redisdb"
	"sl.framework.com/game_server/game/service/type/dto"
	"sl.framework.com/game_server/redis/rediskey"
	"sl.framework.com/trace"
	"strconv"
)

/**
 * GetUserOrder
 * 从缓存hash中获取该用户的订单数据
 *
 * @param traceId string - traceId用于日志跟踪
 * @param gameRoomId int64 - 房间Id
 * @param gameRoundId int64 - 局Id
 * @param userId int64 - 用户ID
 * @return *[]types.BetOrderV2 - 用户的投注订单信息
 */

func GetUserOrder(traceId, gameRoomId, gameRoundId, userId string) []*dto.BetDTO {
	msgHeader := fmt.Sprintf("GetUserOrder traceId=%v, gameRoomId=%v, gameRoundId=%v, userId=%v",
		traceId, gameRoomId, gameRoundId, userId)

	var (
		val string
		err error
	)
	redisInfo := rediskey.GetBetHRedisInfo(gameRoomId, gameRoundId, userId)
	if val, err = redisdb.HGet(redisInfo.HTable, redisInfo.Filed); nil != err {
		trace.Error("%v, redis HGet failed, table=%v, field=%v, error=%v", msgHeader, redisInfo.HTable,
			redisInfo.Filed, err.Error())
		return nil
	}
	if len(val) <= 0 {
		trace.Notice("%v, redis HGet no data, table=%v, field=%v,val=%v",
			msgHeader, redisInfo.HTable, redisInfo.Filed, val)
		return nil
	}

	orderList := make([]*dto.BetDTO, 0)
	if err = json.Unmarshal([]byte(val), &orderList); nil != err {
		trace.Error("%v, json unmarshal failed, table=%v, field=%v, val=%v, err=%v", msgHeader, redisInfo.HTable,
			redisInfo.Filed, val, err.Error())
		return nil
	}
	trace.Info("%v, key=%v ,field=%v,order list=%+v", msgHeader, redisInfo.HTable, redisInfo.Filed, orderList)

	return orderList
}

/**
 * GetOrders
 * 从缓存hash中获取该局订单数据
 *
 * @param traceId string - traceId用于日志跟踪
 * @param gameRoomId int64 - 房间Id
 * @param gameRoundId int64 - 局Id
 * @return *[]types.BetOrderV2 - 用户的投注订单信息
 */

func GetOrders(traceId, gameRoomId, gameRoundId string) []*dto.BetDTO {
	msgHeader := fmt.Sprintf("获取局全量注单 GetUserOrder traceId=%v, gameRoomId=%v, gameRoundId=%v",
		traceId, gameRoomId, gameRoundId)

	var (
		val       map[string]string
		err       error
		orderList []*dto.BetDTO
	)
	redisInfo := rediskey.GetBetHRedisInfo(gameRoomId, gameRoundId, "")
	if val, err = redisdb.HGetAll(redisInfo.HTable); nil != err {
		trace.Error("%v, redis GetOrders failed, table=%v, field=%v, error=%v", msgHeader, redisInfo.HTable,
			redisInfo.Filed, err.Error())
		return nil
	}
	if len(val) <= 0 {
		trace.Notice("%v, redis GetOrders no data, table=%v, field=%v", msgHeader, redisInfo.HTable, redisInfo.Filed)
		return nil
	}

	for _, s2 := range val {
		orderArray := make([]dto.BetDTO, 0)
		if err = json.Unmarshal([]byte(s2), &orderArray); nil != err {
			trace.Error("%v, json unmarshal failed, table=%v, field=%v, val=%v, err=%v", msgHeader, redisInfo.HTable,
				redisInfo.Filed, val, err.Error())
			return nil
		}
		for _, v2 := range orderArray {
			orderList = append(orderList, &v2)
		}

	}

	if len(orderList) <= 0 {
		trace.Error("%v, GetOrders empty, table=%v, field=%v, val=%v, err=%v", msgHeader, redisInfo.HTable,
			redisInfo.Filed, val, err.Error())
		return nil
	}
	trace.Info("%v, order list=%v", msgHeader, orderList)

	return orderList
}

/**
 * SetUserOrder
 * 设置hash中的玩家订单信息
 *
 * @param traceId string - traceId用于日志跟踪
 * @param gameRoomId int64 - 房间Id
 * @param gameRoundId int64 - 局Id
 * @param userId int64 - 用户ID
 * @param *[]types.BetOrderV2 - 用户的投注订单信息
 * @return
 */

func SetUserOrder(traceId, gameRoomId, gameRoundId, userId string, orderList []*dto.BetDTO) {
	msgHeader := fmt.Sprintf("更新用户注单信息 SetUserOrder traceId=%v, gameRoomId=%v, gameRoundId=%v, userId=%v",
		traceId, gameRoomId, gameRoundId, userId)

	var (
		val  int64
		err  error
		data []byte
	)
	if data, err = json.Marshal(orderList); err != nil {
		trace.Error("%v, json marshal failed, error=%v", msgHeader, err.Error())
		return
	}
	for _, betDTO := range orderList {
		trace.Info("%v orderInfo=%+v", msgHeader, *betDTO)
	}
	redisInfo := rediskey.GetBetHRedisInfo(gameRoomId, gameRoundId, userId)
	if val, err = redisdb.HSet(redisInfo.HTable, redisInfo.Filed, string(data), redisInfo.Expire); nil != err {
		trace.Error("%v, redis HGet failed, table=%v, field=%v, error=%v", msgHeader, redisInfo.HTable,
			redisInfo.Filed, err.Error())
		return
	}

	trace.Info("%v, val=%v, order list=%v", msgHeader, val, orderList)
}

/**
 * SetOrders
 * 设置所有当局注单信息
 *
 * @param traceId string - traceId用于日志跟踪
 * @param gameRoomId int64 - 房间Id
 * @param gameRoundId int64 - 局Id
 * @param *[]types.BetOrderV2 - 用户的投注订单信息
 * @return
 */

func SetOrders(traceId, gameRoomId, gameRoundId string, orderList []*dto.BetDTO) {
	msgHeader := fmt.Sprintf("更新全量用户注单信息 SetUserOrder traceId=%v, gameRoomId=%v, gameRoundId=%v order size=%v",
		traceId, gameRoomId, gameRoundId, len(orderList))

	var orderBetMap = make(map[string][]*dto.BetDTO, 0)

	for _, betDTO := range orderList {
		usrId := strconv.FormatInt(betDTO.UserId, 10)
		if len(orderBetMap[usrId]) == 0 {
			orderBetMap[usrId] = []*dto.BetDTO{}
			orderBetMap[usrId] = append(orderBetMap[usrId], betDTO)
		} else {
			orderBetMap[usrId] = append(orderBetMap[usrId], betDTO)
		}
	}

	for s, dtos := range orderBetMap {
		SetUserOrder(traceId, gameRoomId, gameRoundId, s, dtos)
	}
	trace.Info("%v 更新成功", msgHeader)
}
