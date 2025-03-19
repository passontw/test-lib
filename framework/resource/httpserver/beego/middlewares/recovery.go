package middlewares

/**
 * @Author: M
 * @Date: 2024/7/31 9:42
 * @Desc:
 */

import (
	"github.com/beego/beego/v2/server/web/context"
	"log"
)

func RecoveryMiddleware(ctx *context.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic: %v", r)
			ctx.Output.SetStatus(500)
			ctx.Output.Body([]byte("Internal Server Error"))
		}
	}()
}
