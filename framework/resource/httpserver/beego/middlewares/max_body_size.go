package middlewares

/**
 * @Author: M
 * @Date: 2024/7/31 9:42
 * @Desc:
 */

import (
	"github.com/beego/beego/v2/server/web/context"
)

func MaxBodySizeMiddleware(maxSize int64) func(ctx *context.Context) {
	return func(ctx *context.Context) {
		if len(ctx.Input.RequestBody) > int(maxSize) {
			ctx.Output.SetStatus(413)
			ctx.Output.Body([]byte("Request Entity Too Large"))
			ctx.Abort(413, "Request Entity Too Large")
		}
	}
}
