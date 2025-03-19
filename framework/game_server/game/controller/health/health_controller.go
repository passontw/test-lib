package health

import (
	beego "github.com/beego/beego/v2/server/web"
)

// HealthStatus 健康状态

type HealthStatus struct {
	Status string `json:"status"`
}

/**
 * HealthCheck
 * 健康检查回包
 *
 * @return
 */

func (c *HealthController) HealthCheck() {
	h := HealthStatus{Status: "UP"}
	c.Data["json"] = h

	c.ServeJSON()
}

/**
 * HealthController
 * 健康检查控制器
 */

type HealthController struct {
	beego.Controller
}
