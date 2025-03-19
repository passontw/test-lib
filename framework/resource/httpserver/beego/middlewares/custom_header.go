package middlewares

/**
 * @Author: M
 * @Date: 2024/7/31 10:25
 * @Desc:
 */

import (
	"github.com/beego/beego/v2/server/web/context"
)

func CustomHeaderMiddleware(ctx *context.Context) {
	// 添加自定义的HTTP头
	ctx.Output.Header("X-Custom-Header", "Value")
}
