package middlewares

/**
 * @Author: M
 * @Date: 2024/7/31 10:26
 * @Desc:
 */

import (
	"github.com/beego/beego/v2/server/web/context"
	"sl.framework.com/trace"
)

// 记录这些参数
func RecordParamsMiddleware(ctx *context.Context) {
	// 获取路径参数
	pathParams := ctx.Input.Params()
	// 获取表单参数
	formParams := ctx.Input.RequestBody
	if len(pathParams) > 0 {
		trace.Info("request url path: %v, url  params: %v", ctx.Input.URL(), pathParams)
	}
	if len(formParams) > 0 {
		trace.Info("request url path: %v, form params: %s", ctx.Input.URL(), formParams)
	}
}
