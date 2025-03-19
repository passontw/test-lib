package bet

import (
	"fmt"
	errcode "sl.framework.com/game_server/error_code"
	"sl.framework.com/game_server/game/service"
	types "sl.framework.com/game_server/game/service/type"
	"sl.framework.com/game_server/game/service/type/dto"
	"sl.framework.com/game_server/redis/cache"
	rpcreq "sl.framework.com/game_server/rpc_client"
	"sl.framework.com/trace"
	"strconv"
)

/**
 * isIncluding
 * 判断list是否有对应的数据
 *
 * @param member int64 - 参数说明
 * @param list *[]int64 - 参数说明
 * @return bool - true:存在 false:不存在
 */

func isIncluding[T comparable](member T, list []T) (bExist bool) {
	for _, val := range list {
		if member == val {
			return true
		}
	}
	return
}

/**
 * ServiceBetCancel
 * 投注取消业务处理函数
 * 参数在上层已经做过校验 无需再次校验
 *
 * @param traceId string - traceId用于日志跟踪
 * @param userId int64 - 用户Id
 * @param gameId int64 - 游戏Id
 * @param betCancel *types.BetCancelParam - 投注取消相关信息
 * @return int - 投注取消操作返回值
 */

func ServiceBetCancel(traceId, userId string, betCancel *types.BetCancelParam) int {
	msgHeader := fmt.Sprintf("ServiceBetCancel traceId=%v, userId=%v, gameRoomId=%v, gameRound=%v",
		traceId, userId, betCancel.GameRoomId, betCancel.GameRoundId)

	//从缓存hash中获取该用户的订单数据
	orderList := cache.GetUserOrder(traceId, betCancel.GameRoomId, betCancel.GameRoundId, userId)
	if orderList == nil {
		trace.Info("%v, redis order get user order failed", msgHeader)
		return errcode.RedisErrorGet
	}
	trace.Info("%v, GetUserOrder orderList size:%v", msgHeader, len(orderList))
	//剔除掉取消的订单
	orderRemoved := make([]*dto.BetDTO, 0)
	orderKeep := make([]*dto.BetDTO, 0)
	var gameId int64
	for _, order := range orderList {
		if isIncluding[string](strconv.FormatInt(order.OrderNo, 10), betCancel.OrderNoList) {
			gameId = order.GameId
			orderRemoved = append(orderRemoved, order)
			trace.Info("%v, append to remove List no:%v,size:%v", msgHeader, order.OrderNo, len(orderRemoved))
			continue
		}

		orderKeep = append(orderKeep, order)
	}
	if len(orderRemoved) == 0 {
		trace.Info("%v, bet cancel failed", msgHeader)
		return errcode.GameErrorBetCancelFailed
	}
	//重新设置缓存
	trace.Info("%v,剩余注单列表%+v", msgHeader, orderKeep)
	cache.SetUserOrder(traceId, betCancel.GameRoomId, betCancel.GameRoundId, userId, orderKeep)

	//发送下注取消事件
	var betSimpleDTOList []*dto.BetSimpleDTO
	for _, val := range orderRemoved {
		item := new(dto.BetSimpleDTO)
		item.BetAmount = val.BetAmount
		item.Currency = val.Currency
		item.UserId = strconv.FormatInt(val.UserId, 10)
		item.GameRoundId = strconv.FormatInt(val.GameRoundId, 10)
		item.GameWagerId = strconv.FormatInt(val.GameWagerId, 10)
		item.GameId = strconv.FormatInt(val.GameId, 10)
		betSimpleDTOList = append(betSimpleDTOList, item)
	}

	//sendBetCancelGameMessage(traceId, betCancel.GameRoomId, betCancel.GameRoundId, orderRemoved)
	rpcreq.AsyncSendRoundMessage[[]*dto.BetSimpleDTO](traceId, betCancel.GameRoomId, betCancel.GameRoundId, string(types.GameEventCommandCancelBet), betSimpleDTOList)
	//取消完成回调
	// 获取下注对象
	gameRoomId, _ := strconv.ParseInt(betCancel.GameRoomId, 10, 64)
	gameRoundId, _ := strconv.ParseInt(betCancel.GameRoundId, 10, 64)
	usrId, _ := strconv.ParseInt(userId, 10, 64)
	bettor := service.GetBettor(traceId, types.GameId(gameId))
	if bettor == nil {
		trace.Error("%v, no game bet.handler, invalid gameId:%v", msgHeader, gameId)
		return errcode.ErrorUnknown
	}
	defer service.PutBettor(types.GameId(gameId), bettor)
	bettor.AfterCancelComplete(gameRoomId, gameRoundId, usrId, orderRemoved)
	trace.Info("%v,取消注单完成 finished", msgHeader)
	return errcode.ErrorOk
}

/**
 * sendBetCancelGameMessage
 * 函数功能描述
 *
 * @param PARAM - 参数说明
 * @return RETURN - 返回值说明
 */

func sendBetCancelGameMessage(traceId, gameRoomId, gameRoundId string, orders []*types.BetOrderV2) {
	//data, err := json.Marshal(orders)
	//if err != nil {
	//	trace.Error("sendBetCancelGameMessage traceId=%v, gameRoomId=%v, gameRoundId=%v json marshal failed, error=%v",
	//		traceId, gameRoomId, gameRoundId, err.Error())
	//	return
	//}

	//异步事件通知取消
	//llGameRoomId, _ := strconv.ParseInt(gameRoomId, 10, 64)
	//llGameRoundId, _ := strconv.ParseInt(gameRoundId, 10, 64)
	message := rpcreq.GameMessage[[]*types.BetOrderV2]{
		GameRoomId:     gameRoomId,
		GameRoundId:    gameRoundId,
		MessageCommand: string(types.GameEventCommandCancelBet),
		Body:           orders,
	}
	rpcreq.AsyncGameMessageRequest(traceId, message)
}
