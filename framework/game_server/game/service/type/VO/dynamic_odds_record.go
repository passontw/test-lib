package VO

type DynamicOddsRecord struct {
	Id        int64   `json:"id"`         //id 序号
	WagerId   int64   `json:"wager_id"`   //玩法
	MinWeight int     `json:"min_weight"` //最小权重
	MaxWeight int     `json:"max_weight"` //最大权重
	Status    bool    `json:"enable"`     //是否命中
	Odds      float32 `json:"odds"`       //赔率
}
