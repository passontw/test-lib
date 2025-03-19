package bet

import (
	errcode "sl.framework.com/game_server/error_code"
	"sl.framework.com/game_server/game/service/type/VO"
	"sl.framework.com/game_server/redis/cache"
	"sl.framework.com/trace"
)

/**
 * BetRecord
 * 游戏投注订单临时记录
 * 用于投注 注码的展示恢复使用
 *
 * @param traceId string - traceId用于日志跟踪
 * @param gameRoomId int64 - 房间Id
 * @param gameRoundId int64 - 局Id
 * @return int - 投注取消操作返回值
 */

func ServiceBetRecord(traceId, gameRoomId, gameRoundId, userId string) ([]*VO.BetRecordVO, int) {
	betRecordVOs := make([]*VO.BetRecordVO, 0)
	curOrderList := cache.GetUserOrder(traceId, gameRoomId, gameRoundId, userId)
	if len(curOrderList) == 0 {
		trace.Notice("ServiceBetRecord get user orders from cache empty")
		return make([]*VO.BetRecordVO, 0), errcode.ErrorOk
	}
	for _, v := range curOrderList {
		voItem := &VO.BetRecordVO{
			UserId:         v.UserId,
			GroupId:        v.GroupId,
			OrderNo:        v.OrderNo,
			GameRoomId:     v.GameRoomId,
			GameRoundId:    v.GameRoundId,
			GameCategoryId: v.GameCategoryId,
			GameId:         v.GameId,
			GameWagerId:    v.GameWagerId,
			Currency:       v.Currency,
			Price:          v.BetAmount,
			Num:            v.Num,
			BetOdds:        float64(v.BetOdds),
			BetMultiple:    1,
			BetAmount:      v.BetAmount,
		}
		betRecordVOs = append(betRecordVOs, voItem)
	}
	trace.Info("ServiceBetRecord success BetRecordVO len:%v", len(betRecordVOs))
	return betRecordVOs, errcode.ErrorOk
}
