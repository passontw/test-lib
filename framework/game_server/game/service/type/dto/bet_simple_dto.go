package dto

type BetSimpleDTO struct {
	UserId      string  `json:"userId"`      //user_id 用户id
	GameRoundId string  `json:"gameRoundId"` //game_round_id 游戏局id
	GameId      string  `json:"gameId"`      //game_id 游戏id
	GameWagerId string  `json:"gameWagerId"` //game_wager_id 玩法id
	Currency    string  `json:"currency"`    //currency 币种
	BetAmount   float64 `json:"amount"`      //bet_amount 投注金额
}
