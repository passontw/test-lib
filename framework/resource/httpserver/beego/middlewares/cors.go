package middlewares

/**
 * @Author: M
 * @Date: 2024/7/31 9:41
 * @Desc:
 */

import (
	"github.com/beego/beego/v2/server/web/context"
)

func CORSMiddleware(ctx *context.Context) {
	ctx.Output.Header("Access-Control-Allow-Origin", "*")
	ctx.Output.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	ctx.Output.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
	if ctx.Input.Method() == "OPTIONS" {
		ctx.ResponseWriter.WriteHeader(200)
	}
}
