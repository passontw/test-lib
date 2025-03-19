package gameevent

import (
	"fmt"
	"sl.framework.com/game_server/conf"
	types "sl.framework.com/game_server/game/service/type"
	"sl.framework.com/game_server/game/service/type/VO"
	"sl.framework.com/trace"
)

type GameChangeCardEvent struct {
	types.EventBase
}

/**
 * NewChangeCard
 * 创建改牌实例
 *
 * @param roomId int64 - 房间id
 * @param roundId int64 - 局id
 * @param nextRoundId int64 - 下一局id
 * @param event types.GameEventVO - 来自数据源的事件信息
 * @param roundDto *types.GameRoundDTO - 的游戏局信息
 * @return RETURN - 返回游戏事件实例
 */

func NewChangeCard(event types.GameEventVO, roundDto *types.GameRoundDTO, gameEventInitVo *VO.GameEventInitVO) *GameChangeCardEvent {
	return &GameChangeCardEvent{
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
 * 游戏换靴 游戏服务器接收到该协议后，房间内牌的数量重置为0
 *
 * @param traceId string - 跟踪id
 * @return RETURN
 */

func (e *GameChangeCardEvent) HandleRondEvent() {
	//直接设置为成功
	trace.Info("[游戏换牌] ChangeCard %v 局信息%+v", e.MsgHeader, e.Dto.Payload)
	EventCommonSet(&e.EventBase, string(types.GameEventCommandChangeDeck), string(types.GameEventCommandChangeDeck))

	return
}
