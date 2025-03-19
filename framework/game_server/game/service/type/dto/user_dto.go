package dto

type UserDto struct {
	Id           string `json:"id"`           //用户id
	SiteUserName string `json:"siteUsername"` //site_username 站点用户名
	UserName     string `json:"username"`     //username 用户名,需要根据站点用户名组合生成全局唯一
	NickName     string `json:"nickname"`     //nickname 昵称
	Type         string `json:"type"`         //type 类型:游客 Visitor，试玩 Trial，正式 Normal ,测试 Test,带单 Capper,内部 Inner
	Status       string `json:"status"`       //status 状态:正常 Enable，登录锁定 Login_Lock ，游戏锁定 Game_Locked,重提锁定 Recharge_Locked
	BlackListOn  string `json:"blackListOn"`  //black_list_on 黑名单开启:是 Y,否 N
}

func (u UserDto) ToTrialOn() string {
	trialOn := "N"
	switch u.Type {
	case "Normal":
		trialOn = "Y"
	}

	return trialOn
}
