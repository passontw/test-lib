package dto

import "time"

type BetDTO struct {
	Id                 int64     `json:"id"`                 //数据库字段:id id
	SiteId             int64     `json:"siteId"`             //数据库字段:site_id 站点id
	UserId             int64     `json:"userId"`             //数据库字段:user_id 用户id
	Username           string    `json:"username"`           //数据库字段:username 用户名
	SiteUsername       string    `json:"siteUsername"`       //数据库字段:site_username 站点用户名
	Nickname           string    `json:"nickname"`           //数据库字段:nickname 用户昵称
	GroupId            int64     `json:"groupId"`            //数据库字段:group_id 组id
	OrderNo            int64     `json:"orderNo"`            //数据库字段:order_no 订单号
	GameRoomId         int64     `json:"gameRoomId"`         //数据库字段:game_room_id 游戏房间id
	GameRoundId        int64     `json:"gameRoundId"`        //数据库字段:game_round_id 游戏局id
	GameRoundNo        string    `json:"gameRoundNo"`        //数据库字段:round_no 局号
	GameCategoryId     int64     `json:"gameCategoryId"`     //数据库字段:game_category_id 游戏分类id
	GameId             int64     `json:"gameId"`             //数据库字段:game_id 游戏id
	GameWagerId        int64     `json:"gameWagerId"`        //数据库字段:game_wager_id 玩法id
	Currency           string    `json:"currency"`           //数据库字段:currency 币种
	Num                int       `json:"num"`                //数据库字段:num 数量
	BetOdds            float32   `json:"betOdds"`            //数据库字段:bet_odds 投注时赔率
	DrawOdds           float32   `json:"drawOdds"`           //数据库字段:draw_odds 开奖时赔率
	Type               string    `json:"type"`               //数据库字段:type 类型:投注 Bet，比赛 Match，测试 Test
	BetAmount          float64   `json:"betAmount"`          //数据库字段:bet_amount 投注金额
	WinAmount          float64   `json:"winAmount"`          //数据库字段:win_amount 如果赢金额,投注完成时计算好,派奖使用
	AvailableStatus    string    `json:"availableStatus"`    //数据库字段:available_status 有效状态: 有效 Available，取消 Cancel，重新结算 Resettle
	AvailableBetAmount float64   `json:"availableBetAmount"` //数据库字段:available_bet_amount 有效投注金额
	ClientStatus       string    `json:"clientStatus"`       //数据库字段:client_status 显示状态：已支付 Paid，已结算 Settled，取消 Cancel，结算失败 Settled_Failed
	BetStatus          string    `json:"betStatus"`          //数据库字段:bet_status 投注状态:未支付 Unpaid，已支付 Paid，作废 Invalid,超时未支付 Timeout,支付失败 Failed，支付中  Paying，异常 Exception
	WinLostStatus      string    `json:"winLostStatus"`      //数据库字段:win_lost_status 输赢状态：创建 Create,输 Lose，赢 Win,和 Tie
	PostStatus         string    `json:"postStatus"`         //数据库字段:post_status 派奖状态：创建 Create，待派奖 Ready，派彩中 Doing，作废 Invalid，已派奖 Paid ，已退款 Refund , 失败 Failed , 重新结算 Resettle
	ClientType         string    `json:"clientType"`         //数据库字段:client_type 客户端类型:安卓 Android，IOS IOS,电脑端H5 PC_H5，手机端H5 Mobile_H5
	GameResult         string    `json:"gameResult"`         //数据库字段:game_result 游戏结果
	SettleTime         time.Time `json:"settleTime"`         //数据库字段:settle_time 输赢结算时间
	BetDoneTime        time.Time `json:"betDoneTime"`        //数据库字段:bet_done_time 投注完成时间
	PostTime           time.Time `json:"postTime"`           //数据库字段:post_time 派彩完成时间
	CreateTime         time.Time `json:"createTime"`         //数据库字段:create_time 创建时间
	UpdateTime         time.Time `json:"updateTime"`         //数据库字段:update_time 更新时间
	Summary            string    `json:"summary"`            //数据库字段:summary 说明
	ManualOn           string    `json:"manualOn"`           //数据库字段:manual_on 手动投注投注:是 Y,否 N
	TrialOn            string    `json:"trialOn"`            //数据库字段:trial_on 是否试玩: 是 Y,否 N
	Sort               int       `json:"sort"`               //数据库字段:sort 排序，同一组注单内排序
	Md5                string    `json:"md5"`                //数据库字段:md5 签名
	SettleStatus       string    `json:"-"`                  //结算状态,取值为"Success" "Failed" 内部使用不发送出去
}
