package client

import (
	"fmt"
	errcode "sl.framework.com/game_server/error_code"
	"sl.framework.com/game_server/game/controller/base_controller"
	"sl.framework.com/game_server/game/service/bet"
	"sl.framework.com/tool"
	"sl.framework.com/trace"
)

type BetRecordController struct {
	base_controller.BaseController
}

/**
 * BetRecord
 * 游戏投注订单临时记录
 * 用于投注 注码的展示恢复使用
 */

func (p *BetRecordController) BetRecord() {
	controllerParserDTO := p.ParserFromClient(nil)
	if controllerParserDTO.Code != errcode.ErrorOk {
		trace.Error("BetController BetRecord parser error, code=%v", controllerParserDTO.Code)
		return
	}
	userId := p.Ctx.Input.Header(string(base_controller.TagUserId))
	gameRoomId := p.Ctx.Input.Param(":gameRoomId")
	gameRoundId := p.Ctx.Input.Param(":gameRoundId")

	msgHeader := fmt.Sprintf("玩家投注记录 BetController BetRecord traceId=%v, gameRoomId=%v, gameRoundId=%v, userId=%v",
		controllerParserDTO.TraceId, gameRoomId, gameRoundId, userId)
	pWatcher := tool.NewWatcher(msgHeader)
	trace.Info("%v", msgHeader)
	if len(gameRoundId) == 0 || len(gameRoomId) == 0 {
		trace.Error("%v, invalid param", msgHeader)
		p.ClientResponse(errcode.HttpErrorInvalidParam, controllerParserDTO.TraceId, nil)
		return
	}
	//业务层处理
	betRecords, code := bet.ServiceBetRecord(controllerParserDTO.TraceId, gameRoomId, gameRoundId, userId)
	pWatcher.Stop()
	p.ClientResponse(code, controllerParserDTO.TraceId, betRecords)
}
