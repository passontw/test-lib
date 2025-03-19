package middlewares

/**
 * @Author: M
 * @Date: 2024/7/31 10:25
 * @Desc:
 */

import (
	"github.com/beego/beego/v2/server/web/context"
	"github.com/google/uuid"
)

func RequestIDMiddleware(ctx *context.Context) {
	// 生成唯一的请求ID
	requestID := uuid.New().String()
	// 将请求ID添加到响应头中
	ctx.Output.Header("X-Request-ID", requestID)
}
