package handler

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sl.framework.com/async"
	"sl.framework.com/game_server/currency/impl"
	errcode "sl.framework.com/game_server/error_code"
	"sl.framework.com/game_server/game/service"
	types "sl.framework.com/game_server/game/service/type"
	"sl.framework.com/game_server/game/service/type/const_type"
	"sl.framework.com/game_server/game/service/type/dto"
	"sl.framework.com/game_server/redis/cache"
	rpcreq "sl.framework.com/game_server/rpc_client"
	"sl.framework.com/tool"
	"sl.framework.com/trace"
	"strconv"
	"time"
)

/**
 * OnBetConfirmHandler
 * 提交注单消息处理函数
 *
 *
 * @param traceId string - 跟踪Id
 * @param msgBody []byte 消息体
 * @return int - 返回码
 */

func OnGameDrawHandler(traceId string, msgBody []byte) int {
	var err error
	timeStart := time.Now()
	msgHeader := fmt.Sprintf("OnGameDrawHandler traceId=%v", traceId)
	trace.Info("MQ消息 消费者回调 %v", msgHeader)
	//解析MQ消息
	msgDrawGameDataDTO := new(types.GameDrawDataDTO)
	if err = json.Unmarshal(msgBody, &msgDrawGameDataDTO); nil != err {
		trace.Error("%v, json unmarshal failed, error=%v, data=%v", msgHeader, err.Error(), string(msgBody))
		return errcode.HttpErrorOrderParse
	}
	if !service.ValidateGameId(types.GameId(msgDrawGameDataDTO.GameId)) {
		trace.Error("%v, validateGameId failed, msg=%v", msgHeader, string(msgBody))
		return errcode.ErrorOk
	}

	trace.Info("%v, time elapse=%v, msg header=%+v", msgHeader, time.Since(timeStart), msgDrawGameDataDTO)
	ret := processDrawGame(traceId, msgDrawGameDataDTO)
	trace.Info("%v process game draw done, tme elapse=%v, gameRoundId=%v, gameRoundNo=%v, ret=%v",
		msgHeader, time.Since(timeStart), msgDrawGameDataDTO.GameRoundId, msgDrawGameDataDTO.GameRoundNo, ret)

	return ret
}

// processDrawGame 结算来自rocketmq的注单消息
func processDrawGame(traceId string, msgDrawGameDataDTO *types.GameDrawDataDTO) int {
	var (
		ret           = errcode.ErrorOk
		SettleDTOList []*types.SettleDTO
		BetOrdersList []*dto.BetDTO
	)
	msgHeader := fmt.Sprintf("processDrawGame traceId=%v, gameRoomId=%v, gameRoundId=%v, gameRoundNo=%v",
		traceId, msgDrawGameDataDTO.GameRoomId, msgDrawGameDataDTO.GameRoundId, msgDrawGameDataDTO.GameRoundNo)
	trace.Info("MQ消息 未派彩注单派彩 %v", msgHeader)
	//OrderNoList=0则此局没有注单 只推送游戏结果
	if len(msgDrawGameDataDTO.OrderList) == 0 {
		trace.Notice("%v, ParseGameResult no order, gameResult=%+v", msgHeader, msgDrawGameDataDTO.GameRoundResultDTO)
		////更新order_plan状态为Doing status:Ready,Done,Doing,Create,Failed
		//orderPlanId := int64(-1) //当平台中心的orderPlanId=0时表示没有注单结算 将orderPlanId重置为-1 利于平台中心判断
		//rpcreq.AsyncUpdateOrderPlanIdRequest(orderPlanId, msgDrawGameDataDTO.GameRoomId,
		//	msgDrawGameDataDTO.GameRoundId, traceId, "Done")
		return errcode.ErrorOk
	}

	pDog := tool.NewWatcher("批量结算")
	// 获取结算对象
	drawer := service.GetDrawer(traceId, types.GameId(msgDrawGameDataDTO.GameId))
	if reflect.ValueOf(drawer).IsNil() || drawer == nil {
		trace.Error("%v, no game draw.handler, invalid gameId=%v", msgHeader, msgDrawGameDataDTO.GameId)
		return errcode.ErrorOk //游戏类型错误 返回mq broker成功 不再重发该结算消息
	}
	defer service.PutDrawer(types.GameId(msgDrawGameDataDTO.GameId), drawer)
	//获取结算注单
	SettleDTOList = drawer.SettleOrder(traceId, &msgDrawGameDataDTO.GameRoundResultDTO, msgDrawGameDataDTO, &msgDrawGameDataDTO.OrderList)
	trace.Info("MQ消息 未派彩注单派彩  获取结算注单:%+v OrderList:%+v", msgHeader, SettleDTOList)
	pDog.Stop()

	pDog.Start("CheckSettlement")
	SettleDTOMap := make(map[int64]*types.SettleDTO)
	//结算订单结算完之后获取结算信息
	for _, v := range SettleDTOList {
		SettleDTOMap[v.OrderNo] = v
		//SettleDTOList = append(SettleDTOList, v)
	}

	if len(SettleDTOMap) == 0 {
		trace.Error("processDrawGame traceId:%v settle DTO list empty", traceId)
		return errcode.ErrorSettleEmpty
	}
	pDog.Stop()
	//缓存开奖小票
	trace.Info("MQ消息 未派彩注单派彩  缓存开奖小票:%+v", msgHeader)
	pDog.Start("Cache Receipt")
	orderAllList := cache.GetOrders(traceId, strconv.FormatInt(msgDrawGameDataDTO.GameRoomId, 10), strconv.FormatInt(msgDrawGameDataDTO.GameRoundId, 10))
	if orderAllList != nil {
		for _, order := range orderAllList {
			settleDTO := SettleDTOMap[order.OrderNo]
			if settleDTO != nil {
				//更新注单信息
				order.WinLostStatus = settleDTO.WinLostStatus
				strResult, _ := json.Marshal(settleDTO.GameRoundResult.Payload.Result)
				order.DrawOdds = settleDTO.DrawOdds
				order.PostStatus = string(const_type.PostStatusPaid)
				order.ClientStatus = string(const_type.ClientStatusSettled)
				order.GameResult = string(strResult)
				order.SettleTime = settleDTO.SettleTime //记录结算开始时间
				order.PostTime = time.Now()             //记录派彩完成时间
				order.SettleStatus = "Success"          //走到这里说明结算成功了
				order.UpdateTime = time.Now()
				order.WinAmount = impl.NewCurrency(order.Currency, settleDTO.WinAmount).CurrencyValue()
				order.AvailableBetAmount = impl.NewCurrency(order.Currency, settleDTO.AvailableBetAmount).CurrencyValue()
				BetOrdersList = append(BetOrdersList, order)
			}

		}
	}

	PutReceipts(traceId, msgDrawGameDataDTO.GameRoomId, msgDrawGameDataDTO.GameRoundId, BetOrdersList)
	pDog.Stop()

	//更新用户注单缓存
	cache.SetOrders(traceId, strconv.FormatInt(msgDrawGameDataDTO.GameRoomId, 10), strconv.FormatInt(msgDrawGameDataDTO.GameRoundId, 10), orderAllList)
	//更新注单入库
	pDog.Start("更新注单入库")
	dbGet := service.NewGameDBSaver(traceId, types.GameId(msgDrawGameDataDTO.GameId))
	if dbGet == nil {
		//*e.RetHandleEvent = errcode.HttpErrorInvalidParam
		trace.Error("MQ消息 未派彩注单派彩 NewGameDBSaver %v, new game order saver interfaces failed", traceId)
		return errcode.ErrorUnknown
	}
	dbGet.UpdateOrders(traceId, msgDrawGameDataDTO.GameRoomId, msgDrawGameDataDTO.GameRoundId, &BetOrdersList)
	pDog.Stop()

	//发送注单,更新平台中心注单状态
	//先获取orderlist
	trace.Info("MQ消息  发送注单,更新平台中心注单状态:%+v", msgHeader)
	rpcreq.Settle(traceId, strconv.FormatInt(msgDrawGameDataDTO.GameRoomId, 10), strconv.FormatInt(msgDrawGameDataDTO.GameRoundId, 10), SettleDTOList)

	//通知客户端开小票 需要等上一步执行完
	async.AsyncRunCoroutine(func() {
		sendReceiptToUser(traceId, msgDrawGameDataDTO, BetOrdersList)
	})
	//进行完成结算之后的逻辑处理
	drawer.AfterCompletion(traceId, msgDrawGameDataDTO.GameRoomId, msgDrawGameDataDTO.GameRoundId, SettleDTOList)
	trace.Info("MQ消息 未派彩注单派彩 完成:%+v", msgHeader)
	return ret
}

/**
 * PutReceipts
 * 缓存结算小票
 *
 * @param traceId string - 跟踪id
 * @param gameRoomId int64 - 房间id
 * @param gameRoundId int64 - 局id
 * @param betVOList []*types.BetOrderV2 - 用户下注订单
 * @return
 */

func PutReceipts(traceId string, gameRoomId, gameRoundId int64, betVOList []*dto.BetDTO) {
	trace.Info("PutReceipts traceId:%v, gameRoomId:%v, gameRoundId:%v betList size=%v",
		traceId, gameRoomId, gameRoundId, len(betVOList))
	var (
		receipts map[int64][]*types.BetReceipt
	)
	receipts = make(map[int64][]*types.BetReceipt)
	//遍历订单列表
	for _, val := range betVOList {
		var receipt = &types.BetReceipt{
			OrderNo:       val.OrderNo,
			GameWagerId:   val.GameWagerId,
			WinAmount:     val.WinAmount,
			DrawOdds:      val.DrawOdds,
			WinLostStatus: val.WinLostStatus,
		}
		if receipts[val.UserId] == nil {

			var betReceiptList []*types.BetReceipt
			betReceiptList = append(betReceiptList, receipt)
			receipts[val.UserId] = betReceiptList
		} else {
			receipts[val.UserId] = append(receipts[val.UserId], receipt)
		}
	}
	cache.PutReceipts(traceId, gameRoomId, gameRoundId, receipts)
}

/**
 * SendReceiptToUser
 * 给用户发送小票
 *
 * @param tracdId string - 跟踪id
 * @param gameDrawDataDTO *types.GameDrawDataDTO - 游戏结算数据对象
 */

func sendReceiptToUser(tracdId string, gameDrawDataDTO *types.GameDrawDataDTO, betOrdersList []*dto.BetDTO) {
	trace.Info("sendReceiptToUser traceId:%v,gameRoomId:%v,gameRoundId:%v,betOrderList:%v", tracdId, gameDrawDataDTO.GameRoomId, gameDrawDataDTO.GameRoundId, betOrdersList)
	var (
		UserIdSet []int64
		UserIdMap map[int64]*dto.BetDTO
	)
	UserIdSet = make([]int64, 0)
	UserIdMap = make(map[int64]*dto.BetDTO)
	for _, v := range betOrdersList {
		if UserIdMap[v.UserId] == nil {
			UserIdSet = append(UserIdSet, v.UserId)
			UserIdMap[v.UserId] = v
		} else {
			UserIdMap[v.UserId].WinAmount += v.WinAmount
		}
	}

	if len(UserIdSet) == 0 {
		trace.Error("sendReceiptToUser traceId:%v,gameRoomId:%v,gameRoundId:%v UserIdSet is empty.", tracdId, gameDrawDataDTO.GameRoomId, gameDrawDataDTO.GameRoundId)
		return
	}

	pWatch := tool.NewWatcher("sendReceiptToUser")
	for i := 0; i < len(UserIdSet); i++ {
		//查询数据库
		receipts := cache.GetReceipts(tracdId, gameDrawDataDTO.GameRoomId, gameDrawDataDTO.GameRoundId, UserIdSet[i])
		gameDrawResultVO := &types.GameDrawResultVO{
			GameRoundId: gameDrawDataDTO.GameRoundId,
			Receipts:    receipts,
		}
		message := types.UserMessageDTO{
			GameRoomId: gameDrawDataDTO.GameRoomId,
			UserId:     strconv.FormatInt(UserIdSet[i], 10),
			Command:    string(types.GameEventCommandBetReceipt),
			Time:       strconv.FormatInt(tool.Current(), 10),
			Body:       gameDrawResultVO,
		}
		async.AsyncRunCoroutine(func() {
			trace.Info("sendReceiptToUser traceId:%v,gameRoomId:%v,gameRoundId:%v,betOrderList:%v SendReceipts", tracdId, gameDrawDataDTO.GameRoomId, gameDrawDataDTO.GameRoundId, betOrdersList)
			rpcreq.SendReceipts(tracdId, gameDrawDataDTO.GameRoundNo, gameDrawDataDTO.GameRoundId, gameDrawDataDTO.GameRoomId, message)
		})

	}
	pWatch.Stop()

}
