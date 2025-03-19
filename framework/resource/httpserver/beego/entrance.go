package beego

import (
	"fmt"
	"github.com/beego/beego/v2/server/web"
	_ "net/http/pprof"
	"sl.framework.com/resource/conf"
	"sl.framework.com/resource/httpserver/beego/controllers"
	"sl.framework.com/resource/httpserver/beego/middlewares"
	"sl.framework.com/trace"

	"strconv"
	"strings"
)

func HttpEntrance() {
	// Register middleware
	//web.InsertFilter("*", web.BeforeRouter, middlewares.CORSMiddleware)
	//web.InsertFilter("*", web.BeforeRouter, middlewares.AuthMiddleware)
	//web.InsertFilter("*", web.BeforeRouter, middlewares.CustomHeaderMiddleware)
	//web.InsertFilter("*", web.BeforeRouter, middlewares.MaxBodySizeMiddleware(1024*1024)) // 1 MB limit
	//web.InsertFilter("*", web.BeforeRouter, middlewares.RequestIDMiddleware)
	//web.InsertFilter("*", web.BeforeRouter, middlewares.SecurityHeadersMiddleware)
	//web.InsertFilter("*", web.BeforeRouter, middlewares.RateLimitingMiddleware(50, 60)) // 5 requests per 60 seconds
	web.InsertFilter("*", web.BeforeRouter, middlewares.RecoveryMiddleware)
	web.InsertFilter("*", web.BeforeRouter, middlewares.RecordParamsMiddleware)

	// Register controllers
	c := &controllers.HealthController{}
	web.Router("/healthcheck", c, "*:Healthcheck")
	web.Router("/actuator/health/liveness", c, "get:PrometheusMetrics")
	web.Router("/actuator/health/readiness", c, "get:PrometheusMetrics")

	_pprof, _ := strconv.ParseBool(conf.Section("beego", "enablePprof"))
	if _pprof {
		pc := &controllers.PprofController{}
		trace.Notice("[pprof] performance monitoring is [enabled]")
		web.Router("/debug/pprof/", pc, "get:Get")
		web.Router(`/debug/pprof/:pp([\w]+)`, pc, "get:Get")
	}

	AllRouters()
	port := conf.Section("web", "port")
	web.Run(fmt.Sprintf("0.0.0.0:%s", port))
}

func AllRouters() {
	infos := web.BeeApp.Handlers.GetAllControllerInfo()
	processed := make(map[*web.ControllerInfo]struct{})
	for _, info := range infos {
		infoPointer := info
		if _, exists := processed[infoPointer]; exists {
			continue
		}
		processed[infoPointer] = struct{}{}
		trace.Notice("path: [%s] %s", info.GetPattern(), cleanMethod(fmt.Sprintf("%v", info.GetMethod())))
	}
}

func cleanMethod(method string) string {
	input := method
	trimmed := strings.TrimPrefix(input, "map[")
	trimmed = strings.TrimSuffix(trimmed, "]")
	parts := strings.Split(trimmed, ":")
	if len(parts) == 2 {
		_method := parts[0]
		endpoint := parts[1]
		return fmt.Sprintf("method: [%s] endpoint: [%s]", _method, endpoint)
	}
	return ""
}
