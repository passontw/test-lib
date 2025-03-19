package middlewares

/**
 * @Author: M
 * @Date: 2024/7/31 10:25
 * @Desc:
 */

import (
	"github.com/beego/beego/v2/server/web/context"
)

func SecurityHeadersMiddleware(ctx *context.Context) {
	// 添加常见的安全相关HTTP头
	ctx.Output.Header("X-Content-Type-Options", "nosniff")
	ctx.Output.Header("X-Frame-Options", "DENY")
	ctx.Output.Header("X-XSS-Protection", "1; mode=block")
}
