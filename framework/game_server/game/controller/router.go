package controller

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"net/http"
	"sl.framework.com/game_server/game/controller/client"
	"sl.framework.com/game_server/game/controller/health"
	"sl.framework.com/game_server/game/controller/middle_platform"
	"sl.framework.com/game_server/game/controller/resource"
	"sl.framework.com/trace"
	"strings"
)

// 判断是否是浏览器
func isBrowser(userAgent string) bool {
	browsers := []string{"Mozilla", "Chrome", "Safari", "Opera", "Edge", "Trident"}
	for _, browser := range browsers {
		if strings.Contains(userAgent, browser) {
			return true
		}
	}
	return false
}

// 对来自浏览器的请求做过滤
func browserFilter(ctx *context.Context) {
	userAgent := ctx.Input.Header("User-Agent")
	if isBrowser(userAgent) {
		trace.Info("browserFilter the request is from browser, User-Agent=%v", userAgent)
		ctx.ResponseWriter.WriteHeader(http.StatusForbidden)
		ctx.ResponseWriter.Write([]byte("Forbidden"))
	}

	return
}

/**
 * RegisterHealthRouter
 * 注册健康检查端口路由
 * 以下四个路由必须实现 否则k8s健康检查不过会重启服务
 *
 * @param server *beego.HttpServer - http服务
 * @return
 */

func RegisterHealthRouter(server *beego.HttpServer) {
	server.Router("/actuator/health/readiness", &health.HealthController{}, "get:HealthCheck")
	server.Router("/actuator/health/liveness", &health.HealthController{}, "get:HealthCheck")
	server.Router("/actuator/prometheus", &health.HealthController{}, "get:HealthCheck")
	server.Router("/healthcheck", &health.HealthController{}, "get:HealthCheck")
}

/*
	RegisterRouter 注册路由
*/

func RegisterRouter() {
	//beego.InsertFilter("*", beego.BeforeRouter, browserFilter)

	/* 获取唯一节点Id */
	beego.Router("/uniqueId", &middle_platform.UniqueIDController{}, "get:GetUniqueId")

	/* 与客户端交互路由 包括 投注 投注取消 投注确认 投注记录查询 */
	beego.Router("/bet", &client.BetController{}, "post:Bet")
	beego.Router("/bet/cancel", &client.BetController{}, "put:BetCancel")
	beego.Router("/bet/confirmed", &client.BetController{}, "post:BetConfirm")
	beego.Router("/bet/records/:gameRoomId/:gameRoundId", &client.BetRecordController{}, "get:BetRecord")
	beego.Router("/settle/draw/list/:gameRoomId/", &client.DrawResultController{}, "get:GetList")
	/* 处理事件 包括游戏事件 玩家进入房间或者离开房间事件 */
	//beego.Router("/v1/gameEvent", &GameEventController{}, "post:GameEvent")
	beego.Router("/v1/joinOrLeave", &client.JoinOrLeaveController{}, "post:JoinOrLeaveRoom")
	beego.Router("/joinOrLeave", &client.JoinOrLeaveController{}, "post:JoinOrLeaveRoom")

	/*从数据源接收数据路由*/
	beego.Router("/game/event", &resource.GameEventController{}, "post:GameEvent")
	beego.Router("/game/round/:gameRoomId/:gameRoundNo", &resource.GameRoundController{}, "get:GameEventRoundNo")
	beego.Router("/game/round", &resource.WorkerController{}, "post:SignOn")
	/*从中台ws推送过来的消息*/
	beego.Router("/v1/websocket/message", &middle_platform.WsClientMessageController{}, "post:WebSocketMessage")

	/*错误处理*/
	//beego.ErrorController(&ErrorController{})
}
