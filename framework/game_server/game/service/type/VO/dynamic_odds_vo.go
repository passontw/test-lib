package VO

//gameServer => client ws 消息 通知赔率变化

type DynamicOddsVO struct {
	GameWagerId string  `json:"game_wager_id"` //玩法Id
	Odds        float64 `json:"odds"`          //赔率
	Type        string  `json:"type"`          //类型 Dynamic动态 Limit限红
}
