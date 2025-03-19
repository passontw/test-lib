package gameevent

import (
	"fmt"
	"sl.framework.com/game_server/conf"
	errcode "sl.framework.com/game_server/error_code"
	gamelogic "sl.framework.com/game_server/game/service/game"
	types "sl.framework.com/game_server/game/service/type"
	"sl.framework.com/game_server/game/service/type/VO"
	"sl.framework.com/trace"
)

type GamePauseEvent struct {
	types.EventBase
}

/**
 * NewGamePause
 * 创建游戏暂停实例
 *
 * @param roomId int64 - 房间id
 * @param roundId int64 - 局id
 * @param nextRoundId int64 - 下一局id
 * @param event types.GameEventVO - 来自数据源的事件信息
 * @param roundDto *types.GameRoundDTO - 的游戏局信息
 * @return RETURN - 返回游戏事件实例
 */

func NewGamePause(event types.GameEventVO, roundDto *types.GameRoundDTO, gameEventInitVo *VO.GameEventInitVO) *GamePauseEvent {
	return &GamePauseEvent{
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

func (e *GamePauseEvent) HandleRondEvent() {
	//直接设置为成功
	*e.RetHandleEvent = errcode.ErrorOk
	trace.Info("[游戏暂停] GamePause %v 局信息%+v", e.MsgHeader, e.Dto.Payload)
	EventCommonSet(&e.EventBase, string(types.GameEventCommandGamePause), string(types.GameEventCommandGamePause))
	gamelogic.NewEventUpdateRoomCardNum(e.TraceId, e.Dto.GameRoomId, e.Dto.GameRoundNo, 0, types.RoomCardNumReset).HandleEvent()
	return
}
