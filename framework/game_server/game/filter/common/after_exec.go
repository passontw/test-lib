package common

import (
	"github.com/beego/beego/v2/server/web/context"
	"sl.framework.com/trace"
)

/**
 * BeforeRouterCommonHandler
 * 路由前的公共处理
 *
 * @param ctx  *context.Context- 上下文
 * @return RETURN - 返回值说明
 */

func AfterExecCommonHandler(ctx *context.Context) {
	trace.Debug("控制器执行后的公共处理  请求头:%+v 状态码：%v", ctx.ResponseWriter.Status)
	responseBody := ctx.Output.Body(nil)
	trace.Debug("控制器执行后的公共处理  请求体:%+v", responseBody)
}
