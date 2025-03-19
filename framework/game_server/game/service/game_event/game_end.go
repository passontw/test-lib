package gameevent

import (
	"fmt"
	"sl.framework.com/game_server/conf"
	errcode "sl.framework.com/game_server/error_code"
	types "sl.framework.com/game_server/game/service/type"
	"sl.framework.com/game_server/game/service/type/VO"
	"sl.framework.com/game_server/game/service/type/dto"
	"sl.framework.com/game_server/redis/cache"
	rpcreq "sl.framework.com/game_server/rpc_client"
	"sl.framework.com/trace"
	"strconv"
)

type GameEndEvent struct {
	types.EventBase
}

/**
 * NewGameEnd
 * 创建游戏结束实例
 *
 * @param roomId int64 - 房间id
 * @param roundId int64 - 局id
 * @param nextRoundId int64 - 下一局id
 * @param event types.GameEventVO - 来自数据源的事件信息
 * @param roundDto *types.GameRoundDTO - 的游戏局信息
 * @return RETURN - 返回游戏事件实例
 */

func NewGameEnd(event types.GameEventVO, roundDto *types.GameRoundDTO, gameEventInitVo *VO.GameEventInitVO) *GameEndEvent {
	return &GameEndEvent{
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
 * @return RETURN
 */

func (e *GameEndEvent) HandleRondEvent() {
	//直接设置为成功
	*e.RetHandleEvent = errcode.ErrorOk
	trace.Info("[游戏结束] GameEnd %v 局信息%+v", e.MsgHeader, e.Dto.Payload)
	EventCommonSet(&e.EventBase, string(types.GameEventCommandGameEnd), string(types.GameEventCommandGameEnd))

	//房间缓存
	GameRoomCache(e.TraceId, e.Dto.GameRoomId, e.Dto.GameRoundId, e.Dto.GameId)
	//用户限红缓存
	e.userLimitCache()
	return
}

/**
 * userLimitCache
 * 用户限红缓存
 *
 * @return
 */

func (e *GameEndEvent) userLimitCache() {
	msgHeader := fmt.Sprintf("userLimitCache traceId:%v,gameRoomId:%v,gameRoundId:%v", e.TraceId, e.Dto.GameRoomId, e.Dto.GameRoundId)
	trace.Info("%v", msgHeader)
	//查询用户缓存
	sessionCache := cache.UserSessionCache{GameRoomId: e.Dto.GameRoomId}
	userSessionDTOList, _ := sessionCache.GetALL()
	if len(userSessionDTOList) == 0 {
		trace.Notice("%v get session sign_dto failed", msgHeader)
		return
	}
	//过滤在线用户
	curSessionList := make([]*dto.UserSessionDTO, 0)
	for _, val := range userSessionDTOList {
		if val.Online == true {
			curSessionList = append(curSessionList, val)
		}
	}
	if len(curSessionList) == 0 {
		trace.Notice("%v no online user", msgHeader)
		return
	}
	var reqList types.UserBetLimitBatchRequest
	for _, val := range curSessionList {
		rsq := types.UserBetLimitRequest{
			UserId:   strconv.FormatInt(val.UserId, 10),
			Currency: val.Currency,
		}
		reqList.UserBetLimitList = append(reqList.UserBetLimitList, rsq)
	}
	userBetLimitList, ret := rpcreq.GetUserLimitBatchRequest(e.TraceId, e.Dto.GameRoundNo, e.Dto.GameRoundId, reqList)
	if ret != errcode.ErrorOk {
		trace.Notice("%v GetUserLimitBatchRequest failed", msgHeader)
		return
	}
	for _, val := range userBetLimitList {
		userId, _ := strconv.ParseInt(val.UserId, 10, 64)
		userLimitCache := cache.UserLimitCache{TraceId: e.TraceId, UserId: userId}
		userLimitCache.Set(val)
	}

}
