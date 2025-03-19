package client

import (
	"fmt"
	errcode "sl.framework.com/game_server/error_code"
	"sl.framework.com/game_server/game/controller/base_controller"
	"sl.framework.com/game_server/game/service/type/VO"
	"sl.framework.com/game_server/redis/cache"
	"sl.framework.com/trace"
	"strconv"
)

// DrawResultController 开奖结果查询接口
type DrawResultController struct {
	base_controller.BaseController
}

/* 临时结构 后续删除 begin berlyn 2025.01.21 */

type (
	Card struct {
		Banker []string `json:"dragon"` //格式如: "dragon":["11:1"]
		Player []string `json:"tiger"`  //格式如: "tiger":["11:1"]
	}

	GameResultBac struct {
		Raw         Card `json:"raw"`         //原始的游戏结果
		DragonPoint int8 `json:"dragonPoint"` //龙的点数
		TigerPoint  int8 `json:"tigerPoint"`  //虎的点数
		CardNum     int8 `json:"cardNum"`     //牌的数量
	}

	Payload struct {
		Data   any `json:"data"`
		Result any `json:"result"`
	}

	DrawResult struct {
		ID          string  `json:"id"`
		GameRoomId  string  `json:"gameRoomId"`
		GameRoundId string  `json:"gameRoundId"`
		Payload     Payload `json:"payload"`
	}
)

/*DrawResult 临时结构 后续删除 end berlyn 2025.01.21 */

/**
 * GetList
 * 获取游戏结果列表
 *
 * @param
 * @return
 */

func (p *DrawResultController) GetList() {
	controllerParserDTO := p.ParserFromClient(nil)
	if controllerParserDTO.Code != errcode.ErrorOk {
		trace.Error("BetController Bet parser error, code=%v", controllerParserDTO.Code)
		return
	}

	gameRoomId, _ := strconv.ParseInt(p.Ctx.Input.Param(":gameRoomId"), 10, 64)
	msgHeader := fmt.Sprintf("开奖结果查询 DrawResultController traceId=%v, gameRoomId=%v", controllerParserDTO.TraceId, gameRoomId)
	trace.Info("%v", msgHeader)
	if gameRoomId == 0 {
		trace.Error("%v, invalid param", msgHeader)
		p.ClientResponse(errcode.HttpErrorInvalidParam, controllerParserDTO.TraceId, nil)
		return
	}

	settleCache := cache.SettleCache{TraceId: controllerParserDTO.TraceId, RoomId: gameRoomId}
	settleCache.Get()
	results := settleCache.Data

	resVOList := make([]*VO.GameResultVO, 0)
	for _, result := range results {
		item := new(VO.GameResultVO)
		item.GameRoundId = result.GameRoundId
		item.Timestamp = result.Timestamp
		item.Payload = result.Payload
		resVOList = append(resVOList, item)
	}

	trace.Info("%v, drawResult len=%v", msgHeader, len(resVOList))
	p.ClientResponse(controllerParserDTO.Code, controllerParserDTO.TraceId, resVOList)
}
