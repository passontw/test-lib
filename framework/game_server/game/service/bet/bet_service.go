package bet

import (
	"fmt"
	"sl.framework.com/async"
	"sl.framework.com/game_server/conf"
	snowflaker "sl.framework.com/game_server/conf/snow_flake_id"
	"sl.framework.com/game_server/error_code"
	"sl.framework.com/game_server/game/service"
	gamelogic "sl.framework.com/game_server/game/service/game"
	"sl.framework.com/game_server/game/service/interface/bet"
	types "sl.framework.com/game_server/game/service/type"
	"sl.framework.com/game_server/game/service/type/const_type"
	"sl.framework.com/game_server/game/service/type/dto"
	"sl.framework.com/game_server/redis/cache"
	rpcreq "sl.framework.com/game_server/rpc_client"
	"sl.framework.com/trace"
	"strconv"
	"sync"
	"time"
)

/**
 * validUserBalance
 * 玩家余额校验
 * 玩家下注金额和缓存注单的中的额度之和要小于玩家余额 否则为额度不足
 *
 * @param traceId string - traceId用于日志跟踪
 * @param gameRoomId int64 - 房间号
 * @param gameRoundId int64 - 局号
 * @param userId int64 - 用户Id
 * @param currency string - 货币类型
 * @param betAmount float64 - 下注金额
 * @return int - 校验返回码
 */

func validUserBalance(traceId, gameRoomId, gameRoundId, userId, currency string, betAmount float64) int {
	msgHeader := fmt.Sprintf("validUserBalance gameRoomId=%v, gameRoundId=%v, userId=%v, currency=%v, "+
		"betAmount=%v", gameRoomId, gameRoundId, userId, currency, betAmount)

	//查询余额
	balance, code := rpcreq.BalanceRequest(traceId, currency, userId)
	if code != errcode.ErrorOk {
		trace.Error("%v, balance request failed, code=%v", msgHeader, code)
		return code
	}

	//从缓存hash中获取该用户的订单数据 如果没有数据或者读数据出错那么返回nil
	orderList := cache.GetUserOrder(traceId, gameRoomId, gameRoundId, userId)
	//额度校验
	sumBetAmount := float64(0)
	for _, order := range orderList {
		trace.Debug("%v, order=%+v", msgHeader, order)
		sumBetAmount += order.BetAmount
	}
	if balance.Balance-sumBetAmount < 0 {
		trace.Error("%v, balance not enough. sumBetAmount=%v, balance=%+v", msgHeader, sumBetAmount, balance)
		return errcode.GameErrorBalanceNotEnough
	}

	trace.Info("%v, balance validate success. subBetAmount=%v, balance=%+v", msgHeader, sumBetAmount, balance)
	return errcode.ErrorOk
}

/**
 * initOrder
 * 初始化订单 根据缓存信息构造订单
 *
 * @param traceId string - traceId用于日志跟踪
 * @param userInfo *cache.UserInfo - 用户信息
 * @param gameRoundDetail *types.GameRoundDTO - 局详情信息
 * @param bets *types.BetVO - 玩家下注信息
 * @return []*types.BetOrderV2 - 构造好的注单列表
 */

func initOrder(traceId string, userInfo *dto.UserDto, gameRoundDetail *types.GameRoundDTO,
	bets *types.BetVO) []*dto.BetDTO {
	orderList := make([]*dto.BetDTO, 0)
	if userInfo == nil || gameRoundDetail == nil || bets == nil {
		return orderList
	}
	for i, betInfo := range bets.Bets {
		gameRoomId, _ := strconv.ParseInt(bets.GameRoomId, 10, 64)
		gameRoundId, _ := strconv.ParseInt(bets.GameRoundId, 10, 64)
		gameCategroyId, _ := strconv.ParseInt(gameRoundDetail.GameCategoryId, 10, 64)
		gameId, _ := strconv.ParseInt(gameRoundDetail.GameId, 10, 64)
		userId, _ := strconv.ParseInt(userInfo.Id, 10, 64)
		order := &dto.BetDTO{
			Id:              snowflaker.GetSnowFlakeInstance().GetUniqueId(),
			OrderNo:         snowflaker.GetSnowFlakeInstance().GetUniqueId(),
			UserId:          userId,
			SiteId:          conf.GetServerId(),
			Username:        userInfo.UserName,
			SiteUsername:    userInfo.SiteUserName,
			Nickname:        userInfo.NickName,
			GameRoomId:      gameRoomId,
			GameRoundId:     gameRoundId,
			GameCategoryId:  gameCategroyId,
			GameId:          gameId,
			GameRoundNo:     gameRoundDetail.RoundNo,
			GameWagerId:     int64(betInfo.GameWagerId),
			BetAmount:       betInfo.Chip,
			Currency:        bets.Currency,
			Num:             1,
			TrialOn:         userInfo.ToTrialOn(),
			ManualOn:        "Y", //写死为手动投注
			Type:            string(const_type.Bet),
			ClientType:      bets.Device.ClientType,
			BetStatus:       string(const_type.BetStatusUnpaid),
			WinLostStatus:   string(const_type.WinLostStatusCreate),
			PostStatus:      string(const_type.PostStatusCreate),
			ClientStatus:    string(const_type.ClientStatusSettled),
			AvailableStatus: string(const_type.AvailableStatusAvailable),
			CreateTime:      time.Now(),
			Sort:            i,
		}

		orderList = append(orderList, order)
		trace.Debug("initOrder traceId=%v, order=%+v", traceId, order)
	}

	return orderList
}

/**
 * validateOrders
 * 批量校验注单
 *
 * @param traceId string - traceId用于日志跟踪
 * @param bet.handler bet.IGameBettor - bettor下注对象
 * @param orders []*types.BetOrderV2 - 需要校验的订单信息
 * @return int - 注单校验结果
 */

func validateOrders(traceId string, bettor bet.IGameBettor, orders []*dto.BetDTO) int {
	for _, order := range orders {
		if ret := validate(traceId, bettor, order); ret != errcode.ErrorOk {
			return ret
		}
	}
	return errcode.ErrorOk
}

/**
 * validate
 * 单条注单校验
 *
 * @param traceId string - traceId用于日志跟踪
 * @param bet.handler bet.IGameBettor - bettor下注对象
 * @param order *types.BetOrderV2 - 需要校验的订单信息
 * @return int - 注单校验结果
 */

func validate(traceId string, bettor bet.IGameBettor, order *dto.BetDTO) int {
	var (
		retOdd, retUserLimit      int
		retRoomLimit, retPlayType int
		retExtraRule              int
		retCode                   = errcode.ErrorOk
		odds                      float32
		wg                        = new(sync.WaitGroup)
	)
	msgHeader := fmt.Sprintf(fmt.Sprintf("PlayerBet traceId=%v, gameRoomId=%v, gameRoundId=%v, gameRoundNo=%v, "+
		"userId=%v, wagerId=%v", traceId, order.GameRoomId, order.GameRoundId, order.GameRoundNo, order.UserId, order.GameWagerId))
	if conf.GetUserLimitSwitch() {
		// 1.玩家限红校验
		wg.Add(1)
		async.AsyncRunCoroutine(func() {
			defer wg.Done()
			retUserLimit = bettor.ValidateUserLimit(order)
		})
	}

	// 2.房间限红校验
	wg.Add(1)
	async.AsyncRunCoroutine(func() {
		defer wg.Done()
		retRoomLimit = bettor.ValidateRoomLimit(order)
	})

	//3.获取赔率信息
	wg.Add(1)
	async.AsyncRunCoroutine(func() {
		defer wg.Done()
		retOdd, odds = gamelogic.GetRoomOdd(traceId, order)
	})

	//4.主玩法与旁注玩法校验
	wg.Add(1)
	async.AsyncRunCoroutine(func() {
		defer wg.Done()
		retPlayType = bettor.ValidatePlayType(order)
	})

	//5.其他规则校验
	wg.Add(1)
	async.AsyncRunCoroutine(func() {
		defer wg.Done()
		retExtraRule = bettor.ValidateExtraRule(order)
	})
	wg.Wait()

	retCode = errcode.ErrorOk
	if retUserLimit != errcode.ErrorOk {
		retCode = retUserLimit
	} else if retRoomLimit != errcode.ErrorOk {
		retCode = retRoomLimit
	} else if retOdd != errcode.ErrorOk {
		retCode = retOdd
	} else if retPlayType != errcode.ErrorOk {
		retCode = retPlayType
	} else if retExtraRule != errcode.ErrorOk {
		retCode = retExtraRule
	}
	if retCode == errcode.ErrorOk {
		//缓存注单
		curOrderList := cache.GetUserOrder(traceId, strconv.FormatInt(order.GameRoomId, 10),
			strconv.FormatInt(order.GameRoundId, 10), strconv.FormatInt(order.UserId, 10))
		if len(curOrderList) == 0 {

			curOrderList = make([]*dto.BetDTO, 0)
		}
		curOrderList = append(curOrderList, order)
		cache.SetUserOrder(traceId, strconv.FormatInt(order.GameRoomId, 10), strconv.FormatInt(order.GameRoundId, 10),
			strconv.FormatInt(order.UserId, 10), curOrderList)

		trace.Info("%v, retUserLimit=%v, retRoomLimit=%v, retOdd=%v, retPlayType=%v, retExtraRule=%v, odds=%v, reCode=%v",
			msgHeader, retUserLimit, retRoomLimit, retOdd, retPlayType, retExtraRule, odds, retCode)
	}
	return retCode
}

/**
 * ServiceBet
 * 玩家下注业务层处理
 *
 * @param traceId string - traceId用于日志跟踪
 * @param userId int64 - 用户ID
 * @param bet *types.BetVO - 投注相关信息
 * @return []types.BetResult - 投注结果信息
 */

func ServiceBet(traceId, userId string, betParam *types.BetVO) (int, []types.BetResult) {
	msgHeader := fmt.Sprintf("ServiceBet traceId=%v, gameRoomId=%v, gameRoundId=%v",
		traceId, betParam.GameRoomId, betParam.GameRoundId)
	trace.Debug("%v, betParam=%+v", msgHeader, *betParam)

	//从缓存中拿局信息
	lUserId, _ := strconv.ParseInt(userId, 10, 64)
	gameRoundId, _ := strconv.ParseInt(betParam.GameRoundId, 10, 64)
	roomId, _ := strconv.ParseInt(betParam.GameRoomId, 10, 64)

	//获取下注对象
	bettor := service.GetBettor(traceId, types.GameId(conf.GetGameId()))
	if bettor == nil {
		trace.Error("%v, no game bet.handler, invalid gameId=%v", msgHeader, conf.GetGameId())
		return errcode.GameErrorBettorNotExist, make([]types.BetResult, 0)
	}
	defer service.PutBettor(types.GameId(conf.GetGameId()), bettor)

	roundCache := cache.GameRoundCache{TraceId: traceId, RoomId: roomId, GameRoundId: betParam.GameRoundId}
	roundCache.Get()
	gameRoundDetail := roundCache.Data
	if gameRoundDetail == nil || gameRoundDetail.Id != betParam.GameRoundId {
		trace.Error("%v, gameRoundId not exist, gameRoundDetail=%+v", msgHeader, gameRoundDetail)
		return errcode.GameErrorGameRoundIdNotExist, make([]types.BetResult, 0)
	}
	trace.Debug("%v, gameRoundInfo=%+v", msgHeader, gameRoundDetail)

	//查询用户信息
	userCache := cache.UserInfoCache{TraceId: traceId, RoomId: betParam.GameRoomId, UserId: userId}
	if !userCache.Get() {
		trace.Error("%v, user cache not exist, userId=%v", msgHeader, userId)
		return errcode.GameErrorUserIdNotExist, make([]types.BetResult, 0)
	}
	userInfo := userCache.Data
	trace.Debug("%v, userInfo=%+v", msgHeader, userInfo)

	//用户余额校验
	if code := validUserBalance(traceId, betParam.GameRoomId, betParam.GameRoundId, userId, betParam.Currency,
		betParam.BetAmount); code != errcode.ErrorOk {
		trace.Error("%v, balance not enough", msgHeader)
		return errcode.GameErrorBalanceNotEnough, make([]types.BetResult, 0)
	}

	//组合下注订单后并行校验
	orderList := initOrder(traceId, userInfo, gameRoundDetail, betParam)
	if orderList == nil || len(orderList) == 0 {
		trace.Error("%v, betParam param illegal", msgHeader)
		return errcode.GameErrorBetParamIllegal, make([]types.BetResult, 0)
	}

	//调用游戏服接口校验投注信息
	if code := validateOrders(traceId, bettor, orderList); code != errcode.ErrorOk {
		trace.Error("%v, validate failed, code=%v", msgHeader, code)
		return code, make([]types.BetResult, 0)
	}

	//通知中台投注信息
	sendBetGameMessage(traceId, betParam.GameRoomId, betParam.GameRoundId, orderList)
	//rpcreq.GameMessage[[]*types.BetOrderV2]{}(traceId, betParam.Currency, betParam.GameRoomId, betParam.GameRoundId, userId, )

	//返回投注结果
	results := make([]types.BetResult, 0)
	for _, val := range orderList {
		result := types.BetResult{
			OrderNo:     strconv.FormatInt(val.OrderNo, 10),
			GameWagerId: strconv.FormatInt(val.GameWagerId, 10),
			Currency:    val.Currency,
			BetAmount:   val.BetAmount,
		}
		results = append(results, result)
	}
	trace.Info("%v betParam redis_tool handle done, results=%+v", msgHeader, results)

	//下注完成回调
	bettor.AfterBetComplete(roomId, gameRoundId, lUserId, orderList)

	return errcode.ErrorOk, results
}

func sendBetGameMessage(traceId, gameRoomId, gameRoundId string, orders []*dto.BetDTO) {
	itemList := make([]*dto.BetSimpleDTO, 0)
	for _, order := range orders {

		item := new(dto.BetSimpleDTO)
		item.GameId = strconv.FormatInt(order.GameId, 10)
		item.UserId = strconv.FormatInt(order.UserId, 10)
		item.BetAmount = order.BetAmount
		item.GameRoundId = strconv.FormatInt(order.GameRoundId, 10)
		item.GameWagerId = strconv.FormatInt(order.GameWagerId, 10)
		item.Currency = order.Currency
		itemList = append(itemList, item)
	}
	message := rpcreq.GameMessage[[]*dto.BetSimpleDTO]{
		GameRoomId:     gameRoomId,
		GameRoundId:    gameRoundId,
		MessageCommand: string(types.GameEventCommandBet),
		Date:           time.Now().Unix(),
		Body:           itemList,
	}
	rpcreq.AsyncGameMessageRequest(traceId, message)
}
