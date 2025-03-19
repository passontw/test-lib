package middle_platform

import (
	"encoding/json"
	beego "github.com/beego/beego/v2/server/web"
	"sl.framework.com/game_server/game/dao/uiddb"
	"sl.framework.com/trace"
)

// UniqueIDController 获取唯一ID控制类
type UniqueIDController struct {
	beego.Controller
}

// GetUniqueId 只提供一个Get方法供外部使用
func (u *UniqueIDController) GetUniqueId() {
	data := uiddb.UniqueIdData{}
	err, id := uiddb.GetUniqueIdGeneratorInstance().GetUniqueId()
	if nil != err {
		trace.Error("GetUniqueId failed, error=%v", err.Error())
		data.Code = 1
		data.Msg = err.Error()
	} else {
		data.Code = 0
		data.UniqueId = id
	}
	jsonData, _ := json.Marshal(data)
	trace.Info("GetUniqueId json data=%v", string(jsonData))

	u.Data["json"] = data
	if errServe := u.ServeJSON(); nil != errServe {
		trace.Error("GetUniqueId server json failed, error=%v", errServe.Error())
	}
}
