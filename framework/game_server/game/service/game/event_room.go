package gamelogic

import (
	"encoding/json"
	"fmt"
	"sl.framework.com/async"
	"sl.framework.com/game_server/conf"
	"sl.framework.com/game_server/error_code"
	"sl.framework.com/game_server/game/dao/redisdb"
	"sl.framework.com/game_server/game/service/base"
	"sl.framework.com/game_server/game/service/type/dto"
	"sl.framework.com/game_server/redis/cache"
	"sl.framework.com/game_server/redis/redis_tool"
	"sl.framework.com/game_server/redis/rediskey"

	types "sl.framework.com/game_server/game/service/type"
	rpcreq "sl.framework.com/game_server/rpc_client"
	"sl.framework.com/trace"
	"strconv"
	"time"
)

// getUserIdsInRoom 从redis获取房间内用户
func GetUserIdsInRoom(roomId int64) *types.PlayerInRoom {
	var (
		val           string
		err           error
		playersInRoom = new(types.PlayerInRoom)
	)

	redisInfo := rediskey.GetRoomActionRedisInfo(roomId)
	if val, err = redisdb.Get(redisInfo.Key); nil != err {
		trace.Error("getUserIdsInRoom redis get failed, roomId=%v, key=%v, error=%v", roomId, redisInfo.Key, err.Error())
		return nil
	}

	if len(val) <= 0 {
		trace.Notice("getUserIdsInRoom redis get no data roomId=%v, key=%v", roomId, redisInfo.Key)
		playersInRoom.PlayerInfoSet = make([]types.PlayerInfo, 0, 64)
		return playersInRoom
	}

	if err = json.Unmarshal([]byte(val), playersInRoom); nil != err {
		trace.Error("EventJoinRoom json unmarshal failed, roomId=%v, key=%v, val=%v", roomId, redisInfo.Key, val)
		return nil
	}

	return playersInRoom
}

// SetUserIdsInRoom 设置房间内用户数据
func SetUserIdsInRoom(roomId int64, playersInRoom *types.PlayerInRoom) int {
	var (
		err      error
		val      string
		jsonData []byte
	)

	redisInfo := rediskey.GetRoomActionRedisInfo(roomId)
	if jsonData, err = json.Marshal(playersInRoom); nil != err {
		trace.Error("setUserIdsInRoom json marshal failed, roomId=%v, players=%v, error=%v",
			roomId, playersInRoom, err.Error())
		return errcode.JsonErrorMarshal
	}

	if val, err = redisdb.Set(redisInfo.Key, string(jsonData), redisInfo.Expire); nil != err {
		trace.Error("setUserIdsInRoom redis dao set failed, roomId=%v, key=%v, data=%v",
			roomId, redisInfo.Key, jsonData)
		return errcode.RedisErrorSet
	}

	trace.Info("setUserIdsInRoom redis dao set roomId=%v, key=%v, val=%v", roomId, redisInfo.Key, val)
	return errcode.ErrorOk
}

/**
 * OnJoinRoom
 * 玩家进入房间消息 接收到此消息将玩家信息放入缓存 供对玩家个人限红进行预缓存的时候使用
 *
 * @param traceId string - 用于日志跟踪
 * @param roomAction types.JoinLeaveGameRoom - 进入房间事件信息
 * @return
 */

func OnJoinRoom(traceId string, roomAction types.JoinLeaveGameRoom) {
	ret := 0
	NewEventJoinRoom(traceId, roomAction.RoomId, roomAction.Currency, roomAction.UserId, &ret).HandleEvent()
}

/*
	OnPlayerJoinRoom 玩家进入房间离开房间事件回调
	玩家进入房间消息 接收到此消息将玩家信息放入缓存 供对玩家个人限红进行预缓存的时候使用
*/

func OnPlayerJoinRoom(traceId string, msgBody []byte) int {
	roomAction := types.JoinLeaveGameRoom{}
	if err := json.Unmarshal(msgBody, &roomAction); nil != err {
		trace.Error("OnPlayerJoinRoom unmarshal failed, traceId=%v, msg=%v", traceId, string(msgBody))
		return errcode.JsonErrorUnMarshal
	}
	trace.Info("OnPlayerJoinRoom traceId=%v, roomAction=%+v, msg=%v", traceId, roomAction, string(msgBody))

	ret := 0
	NewEventJoinRoom(traceId, roomAction.RoomId, roomAction.Currency, roomAction.UserId, &ret).HandleEvent()

	return ret
}

// CEventRoomAction 房间事件
type CEventRoomAction struct {
	traceId        string // traceId
	roomId         int64  // 房间ID
	userId         int64  // 用户ID
	currency       string //货币类型
	retHandleEvent *int   // 事件处理返回值
	msgHeader      string //打印信息头
}

// NewEventJoinRoom 创建进入房间事件
func NewEventJoinRoom(traceId, roomId, currency string, userId int64, ret *int) *EventJoinRoom {
	//workaround:获取int64时候传过来的是string 所以在此转化下
	lRoomId, _ := strconv.ParseInt(roomId, 10, 64)
	return &EventJoinRoom{
		CEventRoomAction: CEventRoomAction{
			traceId:        traceId,
			roomId:         lRoomId,
			userId:         userId,
			currency:       currency,
			retHandleEvent: ret,
			msgHeader: fmt.Sprintf("EventJoinRoom HandleEvent traceId=%v, roomId=%v, userId=%v, "+
				"currency=%v", traceId, roomId, userId, currency),
		},
	}
}

var _ IEvent = (*EventJoinRoom)(nil)

type EventJoinRoom struct {
	CEventRoomAction
}

// HandleEvent 处理事件
func (e *EventJoinRoom) HandleEvent() {
	var (
		ret           int
		playersInRoom *types.PlayerInRoom
	)
	trace.Info("%v start", e.msgHeader)

	//从redis获取房间内用户
	*e.retHandleEvent = errcode.ErrorOk
	if playersInRoom = GetUserIdsInRoom(e.roomId); playersInRoom == nil {
		trace.Error("%v, getUserIdsInRoom the return is nil", e.msgHeader)
		*e.retHandleEvent = errcode.RedisErrorGet
		return
	}
	//判断玩家是否在redis中 如果在则直接跳过不更新redis
	for _, player := range playersInRoom.PlayerInfoSet {
		if player.UserId == e.userId {
			trace.Info("%v, user already in redis, skip it.", e.msgHeader)
			return
		}
	}
	playerInfo := types.PlayerInfo{UserId: e.userId, Currency: e.currency}
	userInfo, ret := rpcreq.GetUserClientInfo(e.traceId, e.userId)
	trace.Info("%v, GetUserClientInfo from middle agent userinfo=%+v.", e.msgHeader, userInfo)
	if ret != errcode.ErrorOk {
		trace.Info("%v, can't get user info from middle agent.", e.msgHeader)
		return
	}
	playersInRoom.PlayerInfoSet = append(playersInRoom.PlayerInfoSet, playerInfo)
	//设置redis房间内用户
	userCache := cache.UserInfoCache{
		TraceId: e.traceId,
		RoomId:  strconv.FormatInt(e.roomId, 10),
		UserId:  strconv.FormatInt(e.userId, 10),
		Data:    userInfo,
	}
	userCache.Set()
	if ret = SetUserIdsInRoom(e.roomId, playersInRoom); errcode.ErrorOk != ret {
		trace.Error("%v, setUserIdsInRoom ret=%v", e.msgHeader, ret)
		*e.retHandleEvent = ret
		return
	}

	return
}

/**
 * OnLeaveRoom
 * 玩家离开房间消息 接收到此消息将玩家信息从缓存中去除
 *
 * @param traceId string - 用于日志跟踪
 * @param roomAction types.JoinLeaveGameRoom - 进入房间事件信息
 * @return
 */

func OnLeaveRoom(traceId string, roomAction types.JoinLeaveGameRoom) {
	ret := 0
	NewEventLeaveRoom(traceId, roomAction.RoomId, roomAction.Currency, roomAction.UserId, &ret).HandleEvent()
}

/*
	OnPlayerLeaveRoom 玩家进入房间离开房间事件回调
	玩家离开房间消息 接收到此消息将玩家信息从缓存中去除
*/

func OnPlayerLeaveRoom(traceId string, msgBody []byte) int {
	ret := 0
	//解析msg并将 玩家 房间信息传入函数
	roomAction := types.JoinLeaveGameRoom{}
	if err := json.Unmarshal(msgBody, &roomAction); nil != err {
		trace.Error("OnPlayerLeaveRoom unmarshal failed, traceId=%v, msg=%v", traceId, string(msgBody))
		return errcode.JsonErrorUnMarshal
	}

	trace.Info("OnPlayerLeaveRoom traceId=%v, roomAction=%+v, msg=%v", traceId, roomAction, string(msgBody))
	NewEventLeaveRoom(traceId, roomAction.RoomId, roomAction.Currency, roomAction.UserId, &ret).HandleEvent()

	return ret
}

// NewEventLeaveRoom 创建离开房间事件
func NewEventLeaveRoom(traceId, roomId, currency string, userId int64, ret *int) *EventLeaveRoom {
	//workaround:获取int64时候传过来的是string 所以在此转化下
	lRoomId, _ := strconv.ParseInt(roomId, 10, 64)
	return &EventLeaveRoom{
		CEventRoomAction: CEventRoomAction{
			traceId:        traceId,
			roomId:         lRoomId,
			userId:         userId,
			currency:       currency,
			retHandleEvent: ret,
			msgHeader: fmt.Sprintf("EventLeaveRoom HandleEvent traceId=%v, roomId=%v, userId=%v, currency=%v",
				traceId, roomId, userId, currency),
		},
	}
}

var _ IEvent = (*EventLeaveRoom)(nil)

type EventLeaveRoom struct {
	CEventRoomAction
}

// HandleEvent 处理事件
func (e *EventLeaveRoom) HandleEvent() {
	var (
		ret           int
		index         = errcode.ErrorInvalid
		playersInRoom *types.PlayerInRoom
	)
	trace.Info("%v start", e.msgHeader)

	//从redis获取房间内用户
	if playersInRoom = GetUserIdsInRoom(e.roomId); playersInRoom == nil {
		trace.Error("%v, getUserIdsInRoom the return is nil", e.msgHeader)
		*e.retHandleEvent = errcode.RedisErrorGet
		return
	}

	//删除掉离开房间的玩家
	for idx, v := range playersInRoom.PlayerInfoSet {
		if v.UserId == e.userId {
			index = idx
			break
		}
	}
	if errcode.ErrorInvalid == index {
		trace.Notice("%v, user not exist", e.msgHeader)
		//*e.retHandleEvent = errcode.OtherUserNotExist //不存在不返回错误只打印信息
		return
	}
	playersInRoom.PlayerInfoSet = append(playersInRoom.PlayerInfoSet[:index], playersInRoom.PlayerInfoSet[index+1:]...)

	//设置redis房间内用户
	if ret = SetUserIdsInRoom(e.roomId, playersInRoom); errcode.ErrorOk != ret {
		trace.Error("%v, ret=%v", e.msgHeader, ret)
		*e.retHandleEvent = ret
		return
	}

	*e.retHandleEvent = errcode.ErrorOk
	return
}

/*
赔率倍率
	表名:	HRoomOddsInfo{GM054247311DA}
	field:	{userId}{gameId}{gameWagerId}
	value:	{Odds:"24"}
房间限红信息
	表名:	HRoomLimitInfo{GM054247311DA}
	field:	{userId}{gameId}{gameWagerId}
	value:	{currency:"PHP",minAmount:"20",maxAmount:"2000"}
*/

// getRedisRoomLimitInfo 从redis获取房间限红信息
func getRedisRoomLimitInfo(traceId, currency string, gameRoundId, roomId, gameId, gameWagerId int64) *types.LimitInfo {
	var (
		val       string
		err       error
		limitInfo = new(types.LimitInfo)
	)
	msgHeader := fmt.Sprintf("getRedisRoomLimitInfo traceId=%v, gameRoundId=%v, roomId=%v, gameId=%v, "+
		"gameWagerId=%v, currency=%v", traceId, gameRoundId, roomId, gameId, gameWagerId, currency)

	redisInfo := rediskey.GetRoomLimitHRedisInfoEx(gameRoundId, roomId, gameId, gameWagerId, currency)
	//查询数据
	if val, err = redisdb.HGet(redisInfo.HTable, redisInfo.Filed); nil != err {
		trace.Error("%v, redis hget failed, table=%v, field=%v", msgHeader, redisInfo.HTable,
			redisInfo.Filed, err.Error())
		return nil
	}
	if len(val) <= 0 {
		trace.Error("%v, redis hget no data, table=%v, field=%v", msgHeader, redisInfo.HTable, redisInfo.Filed)
		return nil
	}

	if err = json.Unmarshal([]byte(val), &limitInfo); nil != err {
		trace.Error("%v, json unmarshal failed, table=%v, field=%v, err=%v", msgHeader, redisInfo.HTable,
			redisInfo.Filed, err.Error())
		return nil
	}
	trace.Info("%v, max=%0.4f, min=%0.4f", msgHeader, limitInfo.MaxAmount, limitInfo.MinAmount)

	return limitInfo

}

// GetRoomLimitInfo 获取房间限红,redis中不存在则从能力中心读取
func GetRoomLimitInfo(traceId, currency string, gameRoundId, roomId, gameId, gameWagerId int64) *types.LimitInfo {
	var limitInfo *types.LimitInfo
	trace.Info("GetRoomLimitInfo traceId=%v, gameRoundId=%v, roomId=%v, gameId=%v, gameWagerId=%v, currency=%v",
		traceId, gameRoundId, roomId, gameId, gameWagerId, currency)
	if limitInfo = getRedisRoomLimitInfo(traceId, currency, gameRoundId, roomId, gameId, gameWagerId); nil != limitInfo {
		return limitInfo
	}

	//redis中不存在则重新从平台中心读取并存储到redis
	NewEventRoomDetailInfo(traceId, gameRoundId, roomId).HandleEvent()

	//再次从redis中读取
	if limitInfo = getRedisRoomLimitInfo(traceId, currency, gameRoundId, roomId, gameId, gameWagerId); nil != limitInfo {
		return limitInfo
	}

	return limitInfo
}

// GetRedisRoomOddInfo 从redis获取房间赔率信息
func GetRedisRoomOddInfo(traceId string, gameRoundId, roomId, gameId, gameWagerId int64) *types.OddInfo {
	var (
		val     string
		err     error
		oddInfo = new(types.OddInfo)
	)

	//查询数据
	redisInfo := rediskey.GetRoomOddHRedisInfoEx(gameRoundId, roomId, gameId, gameWagerId)
	msgHeader := fmt.Sprintf("getRedisRoomOddInfo traceId=%v, gameRoundId=%v, roomId=%v, gameId=%v, "+
		"gameWagerId=%v, table=%v, field=%v", traceId, gameRoundId, roomId, gameId, gameWagerId,
		redisInfo.HTable, redisInfo.Filed)
	if val, err = redisdb.HGet(redisInfo.HTable, redisInfo.Filed); nil != err {
		trace.Error("%v, hget failed, error=%v", msgHeader, err.Error())
		return nil
	}
	if len(val) <= 0 {
		trace.Error("%v, redis hget no data", msgHeader)
		return nil
	}

	if err = json.Unmarshal([]byte(val), &oddInfo); nil != err {
		trace.Error("%v, json unmarshal failed error=%v", msgHeader, err.Error())
		return nil
	}
	trace.Info("获取房间赔率信息 %v oddInfo:%+v", msgHeader, oddInfo)
	return oddInfo
}

// GetRoomOddInfo4Draw 赔率信息 结算时候使用
func GetRoomOddInfo4Draw(traceId string, gameRoomId, gameRoundId, GameId, gameWagerId int64) (int, float32) {
	//workaround:从能力中心获取int64时候传过来的是string 所以在此转化下
	msgHeader := fmt.Sprintf("GetRoomOddInfo4Draw traceId=%v, gameRoundId=%v, roomId=%v, gameId=%v, "+
		"gameWagerId=%v", traceId, gameRoundId, gameRoomId, GameId, gameWagerId)

	oddInfo := GetRoomOddInfo(traceId, gameRoundId, gameRoomId, GameId, gameWagerId)
	if oddInfo == nil {
		trace.Error("%v, oddInfo is nil", msgHeader)
		return errcode.ValidateErrorUserLimitNotExist, errcode.ErrorInvalid
	}

	trace.Info("%v, odds=%v", msgHeader, oddInfo.Odds)
	return errcode.ErrorOk, oddInfo.Odds
}

/*
	getRoomOdd 获取赔率信息
	order *types.BetOrder 玩家下注订单
	返回值:1.赔率是否获取成功 2.赔率
	获取赔率信息 从redis缓存中赔率信息 如果不存在则返回错误码
	如果不需要赔率信息则错误码设置为errcode.ErrorOk
*/

func GetRoomOdd(traceId string, order *dto.BetDTO) (int, float32) {
	oddInfo := GetRoomOddInfo(traceId, order.GameRoundId, order.GameRoomId, order.GameId, order.GameWagerId)
	if oddInfo == nil {
		trace.Error("getRoomOdd failed, traceId=%v, gameRoundId=%v, roomId=%v, gameId=%v, gameWagerId=%v",
			traceId, order.GameRoundId, order.GameRoomId, order.GameId, order.GameWagerId)
		return errcode.ValidateErrorUserLimitNotExist, errcode.ErrorInvalid

	}
	//修改赔率信息
	order.BetOdds = oddInfo.Odds
	trace.Info("getRoomOdd traceId=%v, gameRoundId=%v, roomId=%v, gameId=%v, gameWagerId=%v, odds=%v",
		traceId, order.GameRoundId, order.GameRoomId, order.GameId, order.GameWagerId, oddInfo.Odds)
	return errcode.ErrorOk, oddInfo.Odds
}

// GetRoomOddInfo 获取房间赔率,redis中不存在则从能力中心读取
func GetRoomOddInfo(traceId string, gameRoundId, roomId, gameId, gameWagerId int64) *types.OddInfo {
	var oddInfo *types.OddInfo
	if oddInfo = GetRedisRoomOddInfo(traceId, gameRoundId, roomId, gameId, gameWagerId); nil != oddInfo {
		return oddInfo
	}

	//redis中不存在且已经在读取该局号的房间信息数据则不再读取 等待100ms后再读取redis缓存
	if base.GetCacheManager().CheckRoomDetailedInfoDoor(gameRoundId) {
		time.Sleep(time.Duration(100) * time.Millisecond)
	} else {
		//redis中不存在则重新从平台中心读取并存储到redis
		base.GetCacheManager().AddRoomDetailedInfoDoor(gameRoundId)
		NewEventRoomDetailInfo(traceId, gameRoundId, roomId).HandleEvent()
	}

	//再次从redis中读取
	if oddInfo = GetRedisRoomOddInfo(traceId, gameRoundId, roomId, gameId, gameWagerId); nil != oddInfo {
		return oddInfo
	}

	return oddInfo
}

// BetDataItem 下注信息 序列化后存入redis
type BetDataItem struct {
	OrderNo     int64   //订单号
	GameRoundId int64   //game round id
	UserId      int64   //用户Id
	WagerId     int64   //玩法
	Currency    string  //币种
	BetAmount   float64 //下注金额
}

// GetBetTotalAmount 获取玩家userId的currency币种下gameWagerId玩法的总下注数据
func GetBetTotalAmount(traceId, currency string, gameRoundId, userId, gameWagerId int64, betAmount float64) float64 {
	var totalBetAmount float64
	betDataList := GetBetData(traceId, currency, gameRoundId)
	for _, bet := range betDataList {
		if userId == bet.UserId && gameWagerId == bet.WagerId &&
			currency == bet.Currency {
			totalBetAmount += bet.BetAmount
		}
	}

	return totalBetAmount + betAmount
}

// GetBetData 获取gameRoundId下所有currency币种玩家的下注数据
func GetBetData(traceId, currency string, gameRoundId int64) []*BetDataItem {
	var (
		err         error
		betDataList = make([]*BetDataItem, 0, 128)
	)

	redisInfo := rediskey.GetRoomLimitLRedisInfo(gameRoundId, currency)
	betVal, _ := redisdb.LAllMember(redisInfo.Key)
	for _, bet := range betVal {
		item := new(BetDataItem)
		if err = json.Unmarshal([]byte(bet), item); nil != err {
			trace.Error("GetBetData jason unmarshal failed, traceId=%v, gameRoundId=%v, currency=%v, "+
				"data=%v, error=%v", traceId, gameRoundId, currency, bet, err.Error())
			continue
		}
		betDataList = append(betDataList, item)
	}

	return betDataList
}

/**
 * AsyncAppendBetDataList
 * 异步设置缓存列表
 *
 * @param traceId string - traceId 用于日志跟踪
 * @param currency string - currency 货币类型
 * @param gameRoundId int64 - gameRoundId 局Id
 * @param userId int64 - userId 用户Id
 * @param gameWagerId int64 - gameWagerId 玩法Id
 * @param orderNo int64 - orderNo 订单号
 * @return
 */

func AsyncAppendBetDataList(traceId, currency string, gameRoundId, userId int64, betList *[]*types.BetOrderV2) {
	fn := func() {
		for _, bet := range *betList {
			AppendBetData(traceId, currency, gameRoundId, userId, bet.GameWagerId, bet.OrderNo, bet.BetAmount)

		}
	}
	async.AsyncRunCoroutine(fn)
}

// AsyncAppendBetData 异步将下注数据追加到redis list中
func AsyncAppendBetData(traceId, currency string, gameRoundId, userId, gameWagerId, orderNo int64, betAmount float64) {
	fn := func() {
		AppendBetData(traceId, currency, gameRoundId, userId, gameWagerId, orderNo, betAmount)
	}
	async.AsyncRunCoroutine(fn)
}

// AppendBetData 将玩家下注数据追加到redis list中 用于限红判断
func AppendBetData(traceId, currency string, gameRoundId, userId, gameWagerId, orderNo int64, betAmount float64) {
	bet := &BetDataItem{
		OrderNo:     orderNo,
		GameRoundId: gameRoundId,
		UserId:      userId,
		WagerId:     gameWagerId,
		Currency:    currency,
		BetAmount:   betAmount,
	}

	var (
		betData []byte
		err     error
	)
	msgHeader := fmt.Sprintf("appendBetData traceId=%v, gameRoundId=%v, userId=%v, currency=%v, gameWagerId=%v, "+
		"betAmount=%v", traceId, gameRoundId, userId, currency, gameWagerId, betAmount)
	if betData, err = json.Marshal(bet); nil != err {
		trace.Error("%v json marsh failed, error=%v", msgHeader, err.Error())
		return
	}

	redisInfo := rediskey.GetRoomLimitLRedisInfo(gameRoundId, currency)
	if err = redisdb.LAppend(redisInfo.Key, string(betData), redisInfo.Expire); nil != err {
		trace.Error("%v, redisdb append failed, list=%v, data=%v, error=%v", msgHeader, redisInfo.Key,
			string(betData), err.Error())
		return
	}
	trace.Info("%v, redisdb append success, list=%v, data=%v", msgHeader, redisInfo.Key, string(betData))
}

// NewEventRoomDetailInfo 创建个人限红信息事件对象
func NewEventRoomDetailInfo(traceId string, gameRoundId, roomId int64) *EventRoomDetailInfo {
	return &EventRoomDetailInfo{
		RedisEvent: RedisEvent{
			traceId: traceId,
		},
		gameRoundId: gameRoundId,
		roomId:      roomId,
	}
}

var _ IEvent = (*EventRoomDetailInfo)(nil)

// EventRoomDetailInfo 房间详情信息事件
type EventRoomDetailInfo struct {
	RedisEvent
	gameRoundId int64 //游戏局号 用于创建Redis Hash表
	roomId      int64 //房间ID
}

func (r *EventRoomDetailInfo) HandleEvent() {
	var (
		roomDetailedInfo *types.RoomDetailedInfo
		ret              int
	)

	if roomDetailedInfo, ret = rpcreq.GetRoomInfoRequest(r.traceId, r.roomId, r.gameRoundId); errcode.ErrorOk != ret {
		trace.Error("EventRoomDetailInfo HandleEvent traceId=%v, roomId=%v, failed error code=%v", r.traceId, r.roomId, ret)
		return
	}

	//动态设置redis缓存时长
	redistool.SetRedisExpire(roomDetailedInfo.Game.Duration)
	//缓存房间限红数据
	NewRoomLimitInfoStoreEvent(r.traceId, r.gameRoundId, r.roomId, &roomDetailedInfo.BetLimitRuleList).HandleEvent()
	//缓存玩法数据 其中包括赔率信息
	NewRoomOddInfoStoreEvent(r.traceId, r.gameRoundId, r.roomId, &roomDetailedInfo.GameWagerList).HandleEvent()
}

// NewRoomLimitInfoStoreEvent 创建房间限红信息存储事件
func NewRoomLimitInfoStoreEvent(traceId string, gameRoundId, roomId int64, ptrUserLimitInfo *[]*types.BetLimitRule) *RoomLimitInfoStoreEvent {
	return &RoomLimitInfoStoreEvent{
		RedisEvent: RedisEvent{
			traceId: traceId,
		},
		gameRoundId:      gameRoundId,
		roomId:           roomId,
		ptrRoomLimitInfo: ptrUserLimitInfo,
		msgHeader: fmt.Sprintf("RoomLimitInfoStoreEvent HandleEvent traceId=%v, gameRoundId=%v, roomId=%v",
			traceId, gameRoundId, roomId),
	}
}

var _ IEvent = (*RoomLimitInfoStoreEvent)(nil)

// RoomLimitInfoStoreEvent 存储房间限红信息事件
type RoomLimitInfoStoreEvent struct {
	RedisEvent
	gameRoundId      int64                  //游戏局号 用于创建Redis Hash表
	roomId           int64                  //房间ID
	ptrRoomLimitInfo *[]*types.BetLimitRule //房间限红信息
	msgHeader        string                 //打印信息头
}

func (r *RoomLimitInfoStoreEvent) HandleEvent() {
	if r.ptrRoomLimitInfo == nil {
		trace.Error("%v, ptrRoomLimitInfo is nil", r.msgHeader)
		return
	}
	if len(*r.ptrRoomLimitInfo) == 0 {
		trace.Notice("%v, ptrRoomLimitInfo is empty", r.msgHeader)
		return
	}

	var (
		err          error
		val          int64
		jsonData     []byte
		limitInfoMap = make(map[string]string, 64)
	)
	//组装数据
	for _, item := range *r.ptrRoomLimitInfo {
		limit := types.LimitInfo{
			Currency:  item.Currency,
			MinAmount: item.MinAmount,
			MaxAmount: item.MaxAmount,
		}
		if jsonData, err = json.Marshal(limit); nil != err {
			trace.Error("%v, json marshal failed, err=%v, limit=%+v", r.msgHeader, limit)
			continue
		}
		//workaround:从能力中心获取int64时候传过来的是string 所以在此转化下
		gameWagerId, _ := strconv.ParseInt(item.GameWagerId, 10, 64)
		gameId, _ := strconv.ParseInt(item.GameId, 10, 64)
		redisInfo := rediskey.GetRoomLimitHRedisInfoEx(r.gameRoundId, r.roomId, gameId, gameWagerId, item.Currency)
		//trace.Debug("%v redisInfo.key:%v,filed:%v", r.msgHeader, redisInfo.HTable, redisInfo.Filed)
		limitInfoMap[redisInfo.Filed] = string(jsonData)
	}

	//写Redis数据
	redisLockInfo := rediskey.GetRoomLimitLockHRedisInfo(r.gameRoundId)
	if !redisdb.Lock(redisLockInfo) {
		trace.Error("%v, redis lock failed", r.msgHeader)
		return
	}
	defer redisdb.Unlock(redisLockInfo)

	redisInfo := rediskey.GetRoomLimitHRedisInfo(r.gameRoundId)
	if val, err = redisdb.HSetBatch(redisInfo.HTable, limitInfoMap, redisInfo.Expire); nil != err {
		trace.Error("%v, hset failed, limitInfoMap=%+v, val=%+v, error=%v", r.msgHeader, limitInfoMap, val, err.Error())
		return
	}

	trace.Info("%v, hset success, result=%v", r.msgHeader, val)
}

// NewRoomOddInfoStoreEvent 创建房间赔率信息存储事件
func NewRoomOddInfoStoreEvent(traceId string, gameRoundId, roomId int64, ptrRoomOddInfo *[]*dto.GameWagerDTO) *RoomOddInfoStoreEvent {
	return &RoomOddInfoStoreEvent{
		RedisEvent: RedisEvent{
			traceId: traceId,
		},
		gameRoundId:    gameRoundId,
		roomId:         roomId,
		ptrRoomOddInfo: ptrRoomOddInfo,
		msgHeader: fmt.Sprintf("RoomOddInfoStoreEvent HandleEvent traceId=%v, gameRoundId=%v, roomId=%v",
			traceId, gameRoundId, roomId),
	}
}

var _ IEvent = (*RoomOddInfoStoreEvent)(nil)

// RoomOddInfoStoreEvent 存储房间赔率信息事件
type RoomOddInfoStoreEvent struct {
	RedisEvent
	gameRoundId    int64                //游戏局号 用于创建Redis Hash表
	roomId         int64                //房间ID
	ptrRoomOddInfo *[]*dto.GameWagerDTO //房间赔率信息
	msgHeader      string               //打印信息头
}

func (r *RoomOddInfoStoreEvent) HandleEvent() {
	if r.ptrRoomOddInfo == nil {
		trace.Error("%v, ptrRoomLimitInfo is nil", r.msgHeader)
		return
	}

	var (
		err        error
		val        int64
		jsonData   []byte
		oddInfoMap = make(map[string]string, 32)
	)
	//存储房间赔率信息
	roomOddInfoCache := &cache.WagerCache{TraceId: r.traceId, GameId: conf.GetGameId(), GameRoomId: r.roomId}
	roomOddInfoCache.Set(*r.ptrRoomOddInfo)
	//设置玩法赔率
	for _, item := range *r.ptrRoomOddInfo {
		odd := types.OddInfo{
			Odds: item.Odds,
		}
		if jsonData, err = json.Marshal(odd); nil != err {
			trace.Error("%v, json marshal failed, limit=%+v, error=%v", r.msgHeader, odd, err.Error())
			continue
		}
		//workaround:从能力中心获取int64时候传过来的是string 所以在此转化下
		id, _ := strconv.ParseInt(item.Id, 10, 64)
		gameId, _ := strconv.ParseInt(item.GameId, 10, 64)
		redisInfo := rediskey.GetRoomOddHRedisInfoEx(r.gameRoundId, r.roomId, gameId, id) //item.Id指玩法Id
		//trace.Debug("%v redisInfo.key:%v,filed:%v", r.msgHeader, redisInfo.HTable, redisInfo.Filed)
		oddInfoMap[redisInfo.Filed] = string(jsonData)
	}
	redisLockInfo := rediskey.GetRoomOddLockHRedisInfo(r.gameRoundId)
	if !redisdb.Lock(redisLockInfo) {
		trace.Error("%v, redis lock failed", r.msgHeader)
		return
	}
	defer redisdb.Unlock(redisLockInfo)

	redisInfo := rediskey.GetRoomOddHRedisInfo(r.gameRoundId)
	if val, err = redisdb.HSetBatch(redisInfo.HTable, oddInfoMap, redisInfo.Expire); nil != err {
		trace.Error("%v, HSetBatch failed, oddInfoMap=%+v, val=%v, error=%v", r.msgHeader, oddInfoMap, val, err.Error())
		return
	}
	trace.Info("%v, HSetBatch success, result=%v", r.msgHeader, val)
}

// GetRoomCardStat 获取当前房间牌的数量
func GetRoomCardStat(traceId string, gameRoomId int64) (cardNum *types.RoomCardStat) {
	var (
		err error
		val string
	)
	msgHeader := fmt.Sprintf("GetRoomCardStat traceId=%v, gameRoomId=%v", traceId, gameRoomId)
	cardNum = new(types.RoomCardStat)

	redisInfo := rediskey.GetGameCardNumRedisInfo(gameRoomId)
	if val, err = redisdb.Get(redisInfo.Key); nil != err {
		trace.Error("%v, redis dao get failed, error=%v", msgHeader, err.Error())
		return
	}
	if len(val) <= 0 {
		trace.Notice("%v, redis get no data", msgHeader)
		return
	}

	if err = json.Unmarshal([]byte(val), cardNum); nil != err {
		trace.Error("%v, json unmarshal failed, error=%v", msgHeader, err.Error())
		return
	}
	trace.Info("%v, roomCardStat=%+v", msgHeader, cardNum)

	return
}

// NewEventUpdateRoomCardNum 创建房间牌的数量事件
func NewEventUpdateRoomCardNum(traceId string, gameRoomId int64, gameRoundNo string, cardNumDelta int8, op types.RoomCardStatOp) *RoomCardNumUpdateEvent {
	return &RoomCardNumUpdateEvent{
		RedisEvent: RedisEvent{
			traceId: traceId,
		},
		cardNumDelta: cardNumDelta,
		gameRoomId:   gameRoomId,
		gameRoundNo:  gameRoundNo,
		operation:    op,
		msgHeader: fmt.Sprintf("RoomCardNumUpdateEvent HandleEvent traceId=%v, cardNumDelta=%v, gameRoomId=%v, "+
			"gameRoundNo=%v, operation=%v", traceId, cardNumDelta, gameRoomId, gameRoundNo, op),
	}
}

var _ IEvent = (*RoomCardNumUpdateEvent)(nil)

// RoomCardNumUpdateEvent 存储房间赔率信息事件
type RoomCardNumUpdateEvent struct {
	RedisEvent
	gameRoomId   int64  //房间ID
	gameRoundNo  string //局号
	cardNumDelta int8   //房间内牌数量的增量
	operation    types.RoomCardStatOp
	msgHeader    string //打印信息
}

func (r *RoomCardNumUpdateEvent) HandleEvent() {
	var (
		err      error
		jsonData []byte
		cardNum  = &types.RoomCardStat{
			GameRoomId:  r.gameRoomId,  //房间号
			GameRoundNo: r.gameRoundNo, //当前局号 可能存在多节点更新牌的数量所以需要判断局号
		}
	)

	fn := func() {
		//设置房间牌的数量
		formerStat := GetRoomCardStat(r.traceId, r.gameRoomId)
		if formerStat.GameRoundNo == r.gameRoundNo {
			trace.Notice("%v, update already", r.msgHeader)
			return
		}

		if r.operation == types.RoomCardNumUpdate {
			//设置房间牌的数量
			cardNum.CardNum = int(r.cardNumDelta) + formerStat.CardNum
		} else {
			cardNum.CardNum = 0
		}
		if jsonData, err = json.Marshal(cardNum); nil != err {
			trace.Error("%v, json marshal failed, error=%v", r.msgHeader, err.Error())
			return
		}

		redisInfo := rediskey.GetGameCardNumRedisInfo(r.gameRoomId)
		val, err := redisdb.Set(redisInfo.Key, string(jsonData), redisInfo.Expire)
		if nil != err {
			trace.Error("%v, redis dao set failed, error=%v", r.msgHeader, err.Error())
			return
		}

		trace.Info("%v, redis dao set key=%v, jsonData=%v, val=%v", r.msgHeader, redisInfo.Key, string(jsonData), val)
	}
	async.AsyncRunCoroutine(fn)
}
