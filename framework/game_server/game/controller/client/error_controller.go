package client

import (
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	errcode "sl.framework.com/game_server/error_code"
	"sl.framework.com/game_server/game/controller/base_controller"
)

// 自定义错误控制器
type ErrorController struct {
	beego.Controller
}

// 处理 400 错误
func (this *ErrorController) Error400() {
	res := base_controller.HttpResponse{
		Code: fmt.Sprintf("%04d", 400), Msg: errcode.GetErrMsg(400), Data: nil,
	}
	this.Data["json"] = res
	this.ServeJSON()
}

// 处理 "userNotFound" 错误
func (this *ErrorController) ErrorUserNotFound() {
	this.Data["json"] = map[string]string{"error": "User not found"}
	this.ServeJSON()
}
