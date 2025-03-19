package middlewares

/**
 * @Author: M
 * @Date: 2024/7/31 9:43
 * @Desc:
 */

import (
	"github.com/beego/beego/v2/server/web/context"
	"net/http"
	"sync"
	"time"
)

var (
	requestCounts = make(map[string]int)
	requestTimes  = make(map[string]time.Time)
	mu            sync.Mutex
)

func RateLimitingMiddleware(limit int, duration time.Duration) func(ctx *context.Context) {
	return func(ctx *context.Context) {
		// Skip rate limiting for /healthcheck endpoint
		if ctx.Input.URL() == "/healthcheck" {
			return
		}
		ip := ctx.Input.IP()
		mu.Lock()
		defer mu.Unlock()

		if time.Since(requestTimes[ip]) > duration {
			requestCounts[ip] = 0
		}

		if requestCounts[ip] >= limit {
			ctx.Output.SetStatus(http.StatusTooManyRequests)
			ctx.Output.Body([]byte("Too Many Requests"))
			ctx.Abort(http.StatusTooManyRequests, "Too Many Requests")
		} else {
			requestCounts[ip]++
			requestTimes[ip] = time.Now()
		}
	}
}
