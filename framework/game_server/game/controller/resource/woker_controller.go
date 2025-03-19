package resource

import (
	errcode "sl.framework.com/game_server/error_code"
	"sl.framework.com/game_server/game/controller/base_controller"
	"sl.framework.com/game_server/game/service/type/VO"
	"sl.framework.com/game_server/game/service/worker"
	"sl.framework.com/trace"
)

type WorkerController struct {
	base_controller.BaseController
}

/**
 * SignOn
 * 游戏主播登录
 *
 * @param -
 * @return -
 */

func (c *WorkerController) SignOn() {
	var workerSignVO VO.WorkerSignVO
	controllerParserDTO := c.ParserFromDataSource(&workerSignVO)
	if controllerParserDTO.Code != errcode.ErrorOk {
		trace.Error("GameEventController parser error, code=%v", controllerParserDTO.Code)
		return
	}
	trace.Info("【数据源TO游戏服】  事件类型为:主播登录,WorkerController trace id=%v, receive data=%+v",
		controllerParserDTO.TraceId, workerSignVO)

	ret := worker.SignOn(controllerParserDTO.TraceId, &workerSignVO)
	c.DataSourceResponse(ret, controllerParserDTO.TraceId, nil)
}
