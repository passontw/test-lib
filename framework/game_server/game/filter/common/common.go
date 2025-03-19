package common

import (
	"bytes"
	"github.com/beego/beego/v2/server/web/context"
	"net/http"
	"runtime/debug"
	"sl.framework.com/trace"
)

/**
 * BeforeRouterCommonHandler
 * 路由前的公共处理
 *
 * @param ctx  *context.Context- 上下文
 * @return RETURN - 返回值说明
 */

func BeforeRouterCommonHandler(ctx *context.Context) {
	trace.Debug("路由前的公共处理  请求路径:%+v", ctx.Request.URL)
}

/**
 * ErrorHandler
 * 错误公共处理
 *
 * @param ctx  *context.Context- 上下文
 * @return RETURN - 返回值说明
 */

type ErrorResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

func ErrorHandler(ctx *context.Context) {
	trace.Info("捕获到状态码:%v", ctx.ResponseWriter.Status)
	if ctx.ResponseWriter.Status != http.StatusOK {
		trace.Info("捕获到错误，状态码:%v", ctx.ResponseWriter.Status)

		// 返回 JSON 格式错误信息
		ctx.Output.SetStatus(ctx.ResponseWriter.Status)
		err := ctx.Output.JSON(ErrorResponse{
			Code:    ctx.ResponseWriter.Status,
			Message: http.StatusText(ctx.ResponseWriter.Status),
			Data:    "",
		}, false, false)
		if err != nil {
			return
		}
	}
	if err := recover(); err != nil {
		trace.Error("捕获到 panic:%v", err) // 记录日志

		// 记录 `panic` 堆栈信息
		var buf bytes.Buffer
		buf.Write(debug.Stack())
		trace.Info("堆栈信息:%v\n", buf.String())
		// 设置 500 错误码
		ctx.Output.SetStatus(http.StatusInternalServerError)

		// 返回 JSON 格式的错误信息
		err := ctx.Output.JSON(ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "server internal error",
			Data:    "",
		}, false, false)
		if err != nil {
			return
		}
	}
}
