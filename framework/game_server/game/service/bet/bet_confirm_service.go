package bet

import (
	"fmt"
	"sl.framework.com/async"
	"sl.framework.com/game_server/conf"
	errcode "sl.framework.com/game_server/error_code"
	"sl.framework.com/game_server/game/service"
	types "sl.framework.com/game_server/game/service/type"
	"sl.framework.com/game_server/game/service/type/const_type"
	"sl.framework.com/game_server/game/service/type/dto"
	"sl.framework.com/game_server/redis/cache"
	rpcreq "sl.framework.com/game_server/rpc_client"
	"sl.framework.com/trace"
	"strconv"
	"time"
)

/**
 * BetConfirm
 * 投注确认业务处理函数
 * todo:事务性函数 需要做事务性处理
 *
 * @param traceId string - traceId用于日志跟踪
 * @param gameRoomId int64 - 房间Id
 * @param gameRoundId int64 - 局号
 * @param gameId int64 - 游戏Id
 * @param userId int64 - 用户Id
 * @param currency string - 货币类型
 * @return int - 投注确认操作返回值
 */

func ServiceBetConfirm(traceId, gameRoomId, gameRoundId, userId string, currency string) int {
	msgHeader := fmt.Sprintf("ServiceBetConfirm traceId=%v, gameRoomId=%v, gameRoundId=%v, "+
		"userId=%v, currency=%v", traceId, gameRoomId, gameRoundId, currency, userId)
	trace.Info("[注单提交业务处理] %v", msgHeader)
	llGameRoomId, _ := strconv.ParseInt(gameRoomId, 10, 64)
	llGameRoundId, _ := strconv.ParseInt(gameRoundId, 10, 64)

	//从缓存hash中获取该用户的订单数据
	orderList := cache.GetUserOrder(traceId, gameRoomId, gameRoundId, userId)
	if orderList == nil {
		trace.Info("[注单提交业务处理] %v, redis order get user order failed", msgHeader)
		return errcode.RedisErrorGet
	}
	//批量签名 中台还未实现 展示注释
	//signor := impl.SignContextImpl{TraceId: traceId, GameId: types.GameId(conf.GetGameId())}
	//////获取验证通过的注单列表
	//signOrderList := signor.Sign(orderList)
	//通知中台投注注单，用于扣款
	if ret := rpcreq.BetRequest(traceId, currency, gameRoomId, gameRoundId, userId, &orderList); ret != errcode.ErrorOk {
		trace.Error("%v, bet http request failed. return code=%v", ret)
		return ret
	}

	//异步调用具体游戏服接口批量入库 避免具体游戏服数据库写入操作耗时太久而阻塞游戏框架流程
	dbSaver := service.NewGameDBSaver(traceId, types.GameId(conf.GetGameId()))
	if dbSaver == nil {
		trace.Error("[注单提交业务处理] %v, new game order saver interfaces failed", msgHeader)
		return errcode.GameErrorNoGameDBSaverRegistered
	}
	trace.Info("[注单提交业务处理] %v, SaveDBBatch orderlist:%+v", msgHeader, orderList)
	dstOrderList := make([]dto.BetDTO, 0)
	for _, v2 := range orderList {
		if conf.ServerConf.GameConfig.DynamicOddsEnable {
			//从缓存中获取动态赔率，设置结算赔率
			dynamicOddsCache := cache.DynamicOddsCache{TraceId: traceId, GameId: v2.GameId, WagerId: v2.GameWagerId, RoomId: v2.GameRoomId, GameRoundId: v2.GameRoundId}
			dynamicOddsCache.Get()
			dynamicOddsInfo := dynamicOddsCache.Data
			if dynamicOddsInfo != nil && dynamicOddsInfo.Enable {
				v2.DrawOdds = dynamicOddsInfo.Odds
			}
			if v2.DrawOdds == 0 {
				v2.DrawOdds = v2.BetOdds
			}
		} else {
			v2.DrawOdds = v2.BetOdds
		}

		v2.CreateTime = time.Now()
		v2.BetDoneTime = time.Now()
		v2.PostStatus = string(const_type.PostStatusReady)
		v2.BetStatus = string(const_type.BetStatusUnpaid)
		dstOrderList = append(dstOrderList, *v2)
	}
	fn := func() { dbSaver.SaveDBBatch(traceId, llGameRoomId, llGameRoundId, &dstOrderList) }
	async.AsyncRunCoroutine(fn)

	//更新缓存
	cache.SetUserOrder(traceId, gameRoomId, gameRoundId, userId, orderList)

	////更新注单入库
	//dbGet := service.NewGameDBSaver(traceId, types.GameId(conf.GetGameId()))
	//if dbGet == nil {
	//	//*e.RetHandleEvent = errcode.HttpErrorInvalidParam
	//	trace.Error("MQ消息 未派彩注单派彩 NewGameDBSaver %v, new game order saver interfaces failed", traceId)
	//	return errcode.ErrorUnknown
	//}
	//nGameRoomId, _ := strconv.ParseInt(gameRoomId, 10, 64)
	//nGameRoundId, _ := strconv.ParseInt(gameRoundId, 10, 64)
	//dbGet.UpdateOrders(traceId, nGameRoomId, nGameRoundId, &orderList)

	//下注确认完成回调
	//取消完成回调
	// 获取下注对象
	gameRoomIdInt, _ := strconv.ParseInt(gameRoomId, 10, 64)
	gameRoundIdInt, _ := strconv.ParseInt(gameRoundId, 10, 64)
	usrId, _ := strconv.ParseInt(userId, 10, 64)
	bettor := service.GetBettor(traceId, types.GameId(conf.GetGameId()))
	if bettor == nil {
		trace.Error("[注单提交业务处理] %v, no game bet.handler, invalid gameId:%v", msgHeader, conf.GetGameId())
		return errcode.ErrorUnknown
	}
	defer service.PutBettor(types.GameId(conf.GetGameId()), bettor)
	bettor.AfterConfirmedComplete(gameRoomIdInt, gameRoundIdInt, usrId, orderList)

	return errcode.ErrorOk
}
