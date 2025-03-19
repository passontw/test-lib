package gameevent

import (
	"fmt"
	"sl.framework.com/async"
	"sl.framework.com/game_server/conf"
	errcode "sl.framework.com/game_server/error_code"
	gamelogic "sl.framework.com/game_server/game/service/game"
	types "sl.framework.com/game_server/game/service/type"
	"sl.framework.com/game_server/game/service/type/VO"
	"sl.framework.com/game_server/game/service/type/dto"
	"sl.framework.com/game_server/redis/cache"
	"sl.framework.com/tool"
	"sl.framework.com/trace"
	"sort"
	"strconv"
)

type GameStartEvent struct {
	types.EventBase
}

/**
 * NewGameStart
 * 创建游戏开始实例
 *
 * @param roomId int64 - 房间id
 * @param roundId int64 - 局id
 * @param nextRoundId int64 - 下一局id
 * @param event types.GameEventVO - 来自数据源的事件信息
 * @param roundDto *types.GameRoundDTO - 的游戏局信息
 * @return RETURN - 返回游戏事件实例
 */

func NewGameStart(event types.GameEventVO, roundDto *types.GameRoundDTO, gameEventInitVo *VO.GameEventInitVO) *GameStartEvent {
	return &GameStartEvent{
		EventBase: types.EventBase{
			Dto: &types.EventDTO{
				GameRoomId:      gameEventInitVo.RoomId,
				GameRoundId:     gameEventInitVo.RoundId,
				GameId:          conf.GetGameId(),
				NextGameRoundId: gameEventInitVo.NextRoundId,
				GameRoundNo:     event.GameRoundNo,
				Command:         string(event.Command),
				Time:            event.Time,
				ReceiveTime:     event.ReceiveTime,
				Payload:         event.Payload,
			},
			RoundDTO:       roundDto,
			TraceId:        gameEventInitVo.TraceId,
			RequestId:      gameEventInitVo.RequestId,
			RetHandleEvent: gameEventInitVo.Code,
			MsgHeader: fmt.Sprintf("command=%s  traceId=%v,RequestId=%v, roomId=%v, gameRoundId=%v, "+
				"nextGameRoundId=%v", event.Command, gameEventInitVo.TraceId, gameEventInitVo.RequestId, gameEventInitVo.RoomId, gameEventInitVo.RoundId, gameEventInitVo.NextRoundId),
		},
	}
}

/**
 * HandleEvent
 * 处理游戏事件函数
 *
 * @param traceId string - 跟踪id
 * @return RETURN
 */

func (e *GameStartEvent) HandleRondEvent() {
	//直接设置为成功
	*e.RetHandleEvent = errcode.ErrorOk
	//基础事务
	gameRound := types.GameRound{RoundId: strconv.FormatInt(e.Dto.GameRoundId, 10), RoundNo: e.Dto.GameRoundNo}
	trace.Info("[游戏开始] GameStart %v 局信息%+v", e.MsgHeader, gameRound)
	e.Dto.Payload = gameRound
	EventCommonSet(&e.EventBase, string(types.GameEventCommandBetStart), string(types.GameEventCommandBetStart))
	//延伸事务
	var players *types.PlayerInRoom
	if e.Dto.NextGameRoundId == 0 {
		//局号为0则说明局号不正确直接返回 且不当做错误处理
		trace.Notice("%v", e.MsgHeader)
		return
	}
	trace.Info("[游戏开始] %v start", e.MsgHeader)

	fnRoom := func() {
		//设置房间赔率 房间限红缓存
		trace.Info("[游戏开始] 设置房间赔率 房间限红缓存 %v 局信息%+v", e.MsgHeader, gameRound)
		gamelogic.NewEventRoomDetailInfo(e.TraceId, e.Dto.NextGameRoundId, e.Dto.GameRoomId).HandleEvent()
	}
	async.AsyncRunCoroutine(fnRoom)

	//获取redis中在线玩家信息
	trace.Info("[游戏开始] 获取redis中在线玩家信息 %v 局信息%+v", e.MsgHeader, gameRound)
	if players = gamelogic.GetUserIdsInRoom(e.Dto.GameRoomId); players == nil {
		*e.RetHandleEvent = errcode.RedisErrorDataIsEmpty
		trace.Notice("%v, no player in redis", e.MsgHeader)
		return
	}

	//并发设置个人限红
	trace.Info("[游戏开始] 设置个人限红 %v 局信息%+v", e.MsgHeader, gameRound)
	for _, player := range players.PlayerInfoSet {
		fn := func(usrId int64, currency string) {
			gamelogic.NewEventUserLimit(e.TraceId, currency, e.Dto.NextGameRoundId, usrId).HandleEvent()
		}
		async.AsyncRunWithAnyMulti[int64, string](fn, player.UserId, player.Currency)
	}

	//设置redis 设置下一局信息已经缓存
	trace.Info("[游戏开始] 异步设置下局信息到redis %v 局信息%+v", e.MsgHeader, gameRound)
	fn := func() { SetNextGameRoundId(e.Dto.GameRoundId, e.Dto.NextGameRoundId, e.RequestId) }
	async.AsyncRunCoroutine(fn)

	//初始化动态赔率
	trace.Info("[游戏开始] 根据动态赔率配置初始化动态赔率 %v", e.MsgHeader)
	e.InitDynamicOdds()
	return
}

/**
 * InitDynamicOdds
 * 初始化动态赔率
 * @return RETURN -
 */

func (e *GameStartEvent) InitDynamicOdds() {
	//获取是否开启动态赔率
	if !conf.ServerConf.GameConfig.DynamicOddsEnable {
		trace.Notice("[初始化动态赔率] 未开启！traceId=%v DynamicOddsEnable=%v", e.TraceId, conf.ServerConf.GameConfig.DynamicOddsEnable)
		return
	}
	//从redis中读取房间赔率
	roomOddInfoCache := &cache.WagerCache{TraceId: e.TraceId, GameId: e.Dto.GameId, GameRoomId: e.Dto.GameRoomId}
	roomOddInfoCache.Get()
	roomOddInfo := roomOddInfoCache.Data
	dynamicInfoList := make([]*types.DynamicOddsInfo, 0)
	//遍历出开启动态赔率的玩法
	for _, oddsInfo := range roomOddInfo {
		wagerId, _ := strconv.ParseInt(oddsInfo.Id, 10, 64)
		if oddsInfo.Type == "Random" {
			infoItem := e.GetDynamicInfo(wagerId, oddsInfo.Odds, oddsInfo.Rate, oddsInfo.OddsList)
			//缓存动态赔率
			dynamicOddsCache := cache.DynamicOddsCache{TraceId: e.TraceId, GameId: e.Dto.GameId, WagerId: wagerId, RoomId: e.Dto.GameRoomId, GameRoundId: e.Dto.GameRoundId}
			dynamicOddsCache.Set(infoItem)
			dynamicInfoList = append(dynamicInfoList, infoItem)
		}
	}
	trace.Info("[初始化动态赔率]  InitDynamicOdds traceId:%v dynamicInfoList:%+v size:%v", e.TraceId, dynamicInfoList, len(dynamicInfoList))
}

/**
 * FunctionName
 * 函数功能描述
 *
 * @param PARAM - 参数说明
 * @return RETURN - 返回值说明
 */

func (e *GameStartEvent) GetDynamicInfo(wagerId int64, oriOdds, rate float32, oddsList []dto.GameWagerOddsDTO) *types.DynamicOddsInfo {
	trace.Info("[获取动态赔率信息] wagerId=%v,oriOdds=%v,rate=%v,oddsList=%v", wagerId, oriOdds, rate, oddsList)
	//先随机看是否落在启用动态赔率权重内
	weight := int(rate * 100) //权重0~weight (weight 最大值10000=100*100)
	randomWeight := tool.GenerateRandomRange(0, 10000)
	dynamicOddsInfo := &types.DynamicOddsInfo{}
	dynamicOddsInfo.WagerId = wagerId
	dynamicOddsInfo.Odds = oriOdds
	if randomWeight <= weight {
		dynamicOddsInfo.Enable = true
	} else {
		trace.Notice("[获取动态赔率信息] 没有随机触发动态赔率 wagerId=%v,weight=%v,curRandom=%v", wagerId, weight, randomWeight)
		return dynamicOddsInfo
	}
	if len(oddsList) == 0 {
		trace.Notice("[获取动态赔率信息] 获取动态赔率列表为空 GetDynamicInfo oddsList is empty")
		return dynamicOddsInfo
	}
	//启用动态赔率，接下来根据各个赔率的权重，随机一个赔率
	//获取动态赔率权重
	dynamicOddsRecordList := make([]*VO.DynamicOddsRecord, 0)
	//先按照id进行升序排序
	sort.Slice(oddsList, func(i, j int) bool { return oddsList[i].Id < oddsList[j].Id })
	var (
		curMin              int = 0
		maxRandom           int = 0
		randomDynamicWeight int
	)

	for _, oddsDTO := range oddsList {
		nId, _ := strconv.ParseInt(oddsDTO.Id, 10, 64)
		nWagerId, _ := strconv.ParseInt(oddsDTO.GameWagerId, 10, 64)
		if oddsDTO.Status == "Enable" {
			item := &VO.DynamicOddsRecord{}
			item.Id = nId
			item.WagerId = nWagerId
			item.Status = true
			item.MinWeight = curMin
			item.MaxWeight = int(oddsDTO.Rate * 100)
			item.Odds = oddsDTO.Odds
			curMin = item.MaxWeight
			dynamicOddsRecordList = append(dynamicOddsRecordList, item)
			trace.Debug("[获取动态赔率信息] 玩法：%v,权重：%v~%v 赔率：%v", nWagerId, item.MinWeight, item.MaxWeight, item.Odds)
		}
	}
	maxRandom = curMin
	if maxRandom < 10000 {
		maxRandom = 10000
	}
	randomDynamicWeight = tool.GenerateRandomRange(0, maxRandom)
	//根据(record.MinWeight,record.MaxWeight] 前开后闭区间取赔率
	for _, record := range dynamicOddsRecordList {
		if randomDynamicWeight > record.MinWeight && randomWeight <= record.MaxWeight {
			dynamicOddsInfo.Odds = record.Odds
			break
		}
	}
	trace.Info("[获取动态赔率信息] wagerId=%v,oriOdds=%v,rate=%v,oddsList=%v,结果:%v", wagerId, oriOdds, rate, oddsList, dynamicOddsInfo)
	return dynamicOddsInfo
}
