package gameevent

import (
	"fmt"
	"sl.framework.com/async"
	"sl.framework.com/game_server/conf"
	errcode "sl.framework.com/game_server/error_code"
	"sl.framework.com/game_server/game/service/game"
	types "sl.framework.com/game_server/game/service/type"
	"sl.framework.com/game_server/game/service/type/VO"
	"sl.framework.com/game_server/game/service/type/const_type"
	"sl.framework.com/game_server/redis/cache"
	rpcreq "sl.framework.com/game_server/rpc_client"
	"sl.framework.com/trace"
	"strconv"
)

type GameStopEvent struct {
	types.EventBase
}

/**
 * NewGameStop
 * 创建停止下注实例
 *
 * @param roomId int64 - 房间id
 * @param roundId int64 - 局id
 * @param nextRoundId int64 - 下一局id
 * @param event types.GameEventVO - 来自数据源的事件信息
 * @param roundDto *types.GameRoundDTO - 的游戏局信息
 * @return RETURN - 返回游戏事件实例
 */

func NewGameStop(event types.GameEventVO, roundDto *types.GameRoundDTO, gameEventInitVo *VO.GameEventInitVO) *GameStopEvent {
	return &GameStopEvent{
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
			MsgHeader: fmt.Sprintf("command=%s  traceId=%v,requestId=%v, roomId=%v, gameRoundId=%v, "+
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

func (e *GameStopEvent) HandleRondEvent() {
	//直接设置为成功
	*e.RetHandleEvent = errcode.ErrorOk
	gameRound := types.GameRound{RoundId: strconv.FormatInt(e.Dto.GameRoundId, 10), RoundNo: e.Dto.GameRoundNo}
	trace.Info("[游戏停止] GameStop %v 局信息%+v", e.MsgHeader, gameRound)
	e.Dto.Payload = gameRound
	EventCommonSet(&e.EventBase, string(types.GameEventCommandBetStop), string(types.GameEventCommandBetStop))
	var players *types.PlayerInRoom
	if e.Dto.NextGameRoundId == 0 {
		//局号为0则说明局号不正确直接返回 且不当做错误处理
		trace.Notice("[游戏停止] %v, invalid nextGameRoundId", e.MsgHeader)
		return
	}
	//停止投注之后 向中台ws推送动态赔率

	if conf.ServerConf.GameConfig.DynamicOddsEnable {
		trace.Info("[游戏停止] 动态赔率开启，推送动态赔率信息。 traceId=%v,gameRoomId=%v,gameRoundId=%v", e.TraceId, e.Dto.GameRoomId, e.Dto.GameRoundId)
		//从缓存中获取动态赔率信息
		//从缓存中获取动态赔率，设置结算赔率
		dynamicOddsCache := cache.DynamicOddsCache{TraceId: e.TraceId, GameId: e.Dto.GameId, RoomId: e.Dto.GameRoomId, GameRoundId: e.Dto.GameRoundId}
		if OddList, err := dynamicOddsCache.GetAll(); !err {
			trace.Notice("[游戏停止] 推送动态赔率到ws %v, 获取动态赔率失败", e.MsgHeader)
			return
		} else {
			//把动态赔率信息发送给ws
			dynamicOddsVOList := make([]VO.DynamicOddsVO, 0)
			for _, info := range OddList {
				if info.Enable {
					vo := VO.DynamicOddsVO{
						GameWagerId: strconv.FormatInt(info.WagerId, 10),
						Odds:        float64(info.Odds),
						Type:        string(const_type.DynamicTypeDynamic),
					}
					dynamicOddsVOList = append(dynamicOddsVOList, vo)
				}
			}
			if len(dynamicOddsVOList) > 0 {
				trace.Info("[游戏停止] 推送动态赔率到ws %v, 获取动态赔率:%+v", e.MsgHeader, dynamicOddsCache)
				rpcreq.AsyncSendRoundMessage(e.TraceId, strconv.FormatInt(e.Dto.GameRoomId, 10), e.RoundDTO.Id, string(types.WSMessageCommandDynamicOdds), dynamicOddsVOList)
			} else {
				trace.Info("[游戏停止] 推送动态赔率到ws %v, 动态赔率列表为空", e.MsgHeader)
			}

		}

	}
	//判断是否已经缓存上一局信息
	nextGameRoundId := int64(0)
	if nextGameRoundId = GetNextGameRoundId(e.Dto.GameRoundId); nextGameRoundId == e.Dto.NextGameRoundId {
		trace.Info("[游戏停止] %v, already in redis. skip it. ", e.MsgHeader)
		return
	}

	trace.Info("[游戏停止] %v, nextGameRoundId=%v not in redis", e.MsgHeader, nextGameRoundId)
	fnRoom := func() {
		//设置房间赔率 房间限红缓存
		gamelogic.NewEventRoomDetailInfo(e.TraceId, e.Dto.NextGameRoundId, e.Dto.GameRoomId).HandleEvent()
	}
	async.AsyncRunCoroutine(fnRoom)

	//获取redis中在线玩家信息
	if players = gamelogic.GetUserIdsInRoom(e.Dto.GameRoomId); players == nil {
		*e.RetHandleEvent = errcode.RedisErrorDataIsEmpty
		trace.Notice("[游戏停止] %v, getUserIdsInRoom no player", e.MsgHeader)
		return
	}

	//并发设置个人限红
	for _, player := range players.PlayerInfoSet {
		fn := func(usrId int64, currency string) {
			gamelogic.NewEventUserLimit(e.TraceId, currency, e.Dto.NextGameRoundId, usrId).HandleEvent()
		}
		async.AsyncRunWithAnyMulti[int64, string](fn, player.UserId, player.Currency)
	}

	return
}
