package types

type DynamicOddsInfo struct {
	WagerId int64   `json:"wager_id"` //玩法
	Enable  bool    `json:"enable"`   //是否启用动态倍率
	Odds    float32 `json:"odds"`     //动态倍率 为0时，取玩法默认赔率
}
