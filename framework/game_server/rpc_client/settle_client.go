package rpcreq

import (
	"fmt"
	"sl.framework.com/game_server/conf"
	types "sl.framework.com/game_server/game/service/type"
	"sl.framework.com/game_server/rpc_client/config"
	"sl.framework.com/trace"
)

/*
	orderReckonResultPut 批量提交订单开奖结果
	调用接口:/v1/settle/feign/orderDraw/result/{siteId}/{gameRoomId}/{gameRoundId}
*/

func Settle(traceId string, gameRoomId, gameRoundId string, orderSettleList []*types.SettleDTO) int {
	url := fmt.Sprintf(config.SettleURL, conf.GetPlatformInfoUrl(),
		gameRoomId, gameRoundId)
	msg := fmt.Sprintf("orderReckonResultPut traceId=%v, gameRoomId=%v, gameRoundId=%v, url=%v orderSettleList=%v", traceId,
		gameRoomId, gameRoundId, url, orderSettleList)
	trace.Info("向中台发送订单开奖结果 msg=%v", msg)
	return runHttpPost(traceId, msg, url, orderSettleList, nil)
}

/*
	DrawResultPost 推送游戏结果
	调用接口:/v1/settle/feign/game/drawResult/{gameRoomId}/{gameRoundId}
*/

func DrawResultPost(traceId string, res *types.GameRoundResultDTO) int {
	url := fmt.Sprintf(config.DrawResultPostURL, conf.GetPlatformInfoUrl(),
		res.Headers.GameRoomId, res.GameRoundId)
	msg := fmt.Sprintf("DrawResultDataPut traceId=%v, gameRoundId=%v, url=%v res:=%+v", traceId, res.Headers.GameRoundId, url, res)
	return runHttpPut(traceId, msg, url, res)
}

/*
	OrderSettled 已经结算的注单
*/

type OrderSettled struct {
	OrderNo            int64   //数据库字段:order_no 订单号
	WinAmount          float64 //数据库字段:win_amount 如果赢金额,投注完成时计算好,派奖使用
	AvailableBetAmount float64 //数据库字段:available_bet_amount 有效投注金额
	DrawOdds           float64 //数据库字段:draw_odds 开奖时赔率
	OrderWinLostStatus string  //数据库字段:win_lost_status 输赢状态：创建 Create,输 Lose，赢 Win,和 Tie"
}

/**
 * SendReceipts
 * 结算之后发送小票到ws集群
 *
 * @param traceId string - traceId 用于日志跟踪
 * @param gameRoundNo string - gameRoundNo 局号
 * @param gameRoundId int64 - gameRoundId 局Id
 * @return int - 请求返回码
 */

func SendReceipts(traceId, gameRoundNo string, gameRoundId, gameRoomId int64, message types.UserMessageDTO) int {
	url := fmt.Sprintf("%v/feign/message/user/send", conf.GetPlatformInfoUrl())
	msg := fmt.Sprintf("SendReceipts traceId=%v, gameRoundId=%v, gameRoundNo=%v, gameRoomId=%v"+
		"url=%v,message=%v", traceId, gameRoundId, gameRoundNo, gameRoomId, url, message)
	trace.Info("结算之后发送小票到ws集群 msg=%v", msg)
	ret := runHttpPost(traceId, msg, url, message, nil)
	return ret
}
