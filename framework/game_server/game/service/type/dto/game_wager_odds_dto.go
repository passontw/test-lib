package dto

// 游戏玩法赔率数据对象
type GameWagerOddsDTO struct {
	Id          string  `json:"id"`          //id
	GameWagerId string  `json:"gameWagerId"` //game_wager_id 游戏玩法id
	Rate        float64 `json:"rate"`        //rate 随机赔率概率
	Odds        float32 `json:"odds"`        //odds 赔率
	Status      string  `json:"status"`      //status 状态:启用 Enable，停用 Disable，维护 Maintain
	CreateTime  string  `json:"createTime"`  //create_time 创建时间
	UpdateTime  string  `json:"updateTime"`  //update_time 更新时间
	Md5         string  `json:"md5"`         //md5 数据指纹
}
