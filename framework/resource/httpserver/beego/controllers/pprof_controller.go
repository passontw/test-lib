package controllers

import (
	"github.com/beego/beego/v2/server/web"
	"net/http"
	"net/http/pprof"
)

type PprofController struct {
	web.Controller
}

func (p *PprofController) Get() {
	pp := p.Ctx.Input.Param(":pp")
	switch pp {
	default:
		pprof.Index(p.Ctx.ResponseWriter, p.Ctx.Request)
	case "":
		pprof.Index(p.Ctx.ResponseWriter, p.Ctx.Request)
	case "cmdline":
		pprof.Cmdline(p.Ctx.ResponseWriter, p.Ctx.Request)
	case "profile":
		pprof.Profile(p.Ctx.ResponseWriter, p.Ctx.Request)
	case "symbol":
		pprof.Symbol(p.Ctx.ResponseWriter, p.Ctx.Request)
	}
	p.Ctx.ResponseWriter.WriteHeader(http.StatusOK)
}
