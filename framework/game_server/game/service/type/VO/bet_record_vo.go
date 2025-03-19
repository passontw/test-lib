package VO

type BetRecordVO struct {
	UserId         int64   `json:"userId"`         //数据库字段:user_id 用户id 与客户端交互int64->string
	GroupId        int64   `json:"groupId"`        //数据库字段:group_id 组id 与客户端交互int64->string
	OrderNo        int64   `json:"orderNo"`        //数据库字段:order_no 订单号 与客户端交互int64->string
	GameRoomId     int64   `json:"gameRoomId"`     //数据库字段:game_room_id 游戏房间id与客户端交互 int64->string
	GameRoundId    int64   `json:"gameRoundId"`    //数据库字段:game_round_id 游戏局id与客户端交互 int64->string
	GameCategoryId int64   `json:"gameCategoryId"` //数据库字段:game_category_id 游戏分类id与客户端交互 int64->string
	GameId         int64   `json:"gameId"`         //数据库字段:game_id 游戏id 与客户端交互int64->string
	GameWagerId    int64   `json:"gameWagerId"`    //数据库字段:game_wager_id 玩法id 与客户端交互int64->string
	Currency       string  `json:"currency"`       //数据库字段:currency 币种
	Price          float64 `json:"price"`          //数据库字段:price 单价
	Num            int     `json:"num"`            //数据库字段:num 数量
	BetOdds        float64 `json:"betOdds"`        //数据库字段:bet_odds 投注时赔率
	BetMultiple    int     `json:"betMultiple"`    //数据库字段:bet_multiple 投注倍数
	BetAmount      float64 `json:"betAmount"`      //数据库字段:bet_amount 投注金额
}
