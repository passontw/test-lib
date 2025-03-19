package controllers

import "github.com/beego/beego/v2/server/web"

type HealthController struct {
	web.Controller
}

func (p *HealthController) Healthcheck() {
	p.Data["json"] = map[string]string{"message": "ok"}
	_ = p.ServeJSON()
}

func (p *HealthController) PrometheusMetrics() {
	p.Data["json"] = map[string]string{"status": "UP"}
	_ = p.ServeJSON()
}
