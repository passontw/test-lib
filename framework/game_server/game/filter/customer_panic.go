package filter

import (
	"bytes"
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"runtime/debug"
	"sl.framework.com/trace"
	"strconv"
)

/**
 * RecoverPanic
 * 自定义panic函数
 *
 * @param ctx  *context.Context- 上下文
 * @return RETURN - 返回值说明
 */

func RecoverPanic(ctx *context.Context, config *beego.Config) {
	if err := recover(); err != nil {
		if err == beego.ErrAbort {
			return
		}

		if !beego.BConfig.RecoverPanic {
			panic(err)
		}
		if beego.BConfig.EnableErrorsShow {
			if _, ok := beego.ErrorMaps[fmt.Sprintf("%v", err)]; ok {
				nErrorCode, _ := strconv.ParseInt(fmt.Sprintf("%v", err), 10, 64)
				beego.Exception(uint64(nErrorCode), ctx)
			}
		}
		trace.Error("捕获到 panic:%v", err) // 记录日志

		// 记录 `panic` 堆栈信息
		var buf bytes.Buffer
		buf.Write(debug.Stack())
		trace.Error("堆栈信息:%v\n", buf.String())
	}
}
