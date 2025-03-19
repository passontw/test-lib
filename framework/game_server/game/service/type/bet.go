package types

import (
	"sl.framework.com/game_server/game/service/type/dto"
	"time"
)

// 不同游戏创建不同数据库 不同数据库中有记录订单信息的表 该表名字统一为game_record
const tableNameGameRecord = "game_record"

// TableName 结构对应的表名
func (b *BetOrderV2) TableName() string {
	return tableNameGameRecord // 返回结构体对应数据库中的表名
}

type (
	// BetOrderV2 玩家下注接口提
	BetOrderV2 struct {
		Id                 int64     `json:"id" orm:"column(id)"`                                     //数据库字段:id id
		SiteId             int64     `json:"site_id" orm:"column(site_id)"`                           //数据库字段:site_id 站点id
		UserId             int64     `json:"userId" orm:"column(user_id)"`                            //数据库字段:user_id 用户id
		Username           string    `json:"username" orm:"size(32);column(username)"`                //数据库字段:username 用户名
		SiteUsername       string    `json:"siteUsername" orm:"size(64);column(site_username)"`       //数据库字段:site_username 站点用户名
		Nickname           string    `json:"nickname" orm:"size(32);column(nickname)"`                //数据库字段:nickname 用户昵称
		GroupId            int64     `json:"groupId" orm:"column(group_id)"`                          //数据库字段:group_id 组id
		OrderNo            int64     `json:"orderNo" orm:"column(order_no)"`                          //数据库字段:order_no 订单号
		GameRoomId         int64     `json:"gameRoomId" orm:"column(game_room_id)"`                   //数据库字段:game_room_id 游戏房间id
		GameRoundId        int64     `json:"gameRoundId" orm:"column(game_round_id)"`                 //数据库字段:game_round_id 游戏局id
		GameRoundNo        string    `json:"gameRoundNo" orm:"size(32);column(game_round_no)"`        //数据库字段:round_no 局号
		GameCategoryId     int64     `json:"gameCategoryId" orm:"column(game_category_id)"`           //数据库字段:game_category_id 游戏分类id
		GameId             int64     `json:"gameId" orm:"column(game_id)"`                            //数据库字段:game_id 游戏id
		GameWagerId        int64     `json:"gameWagerId" orm:"column(game_wager_id)"`                 //数据库字段:game_wager_id 玩法id
		Currency           string    `json:"currency" orm:"size(16);column(currency)"`                //数据库字段:currency 币种
		Num                int       `json:"num" orm:"column(num)"`                                   //数据库字段:num 数量
		BetOdds            float32   `json:"betOdds" orm:"column(bet_odds)"`                          //数据库字段:bet_odds 投注时赔率
		DrawOdds           float32   `json:"drawOdds" orm:"column(draw_odds)"`                        //数据库字段:draw_odds 开奖时赔率
		Type               string    `json:"type" orm:"column(type)"`                                 //数据库字段:type 类型:投注 Bet，比赛 Match，测试 Test
		BetAmount          float64   `json:"betAmount" orm:"column(bet_amount)"`                      //数据库字段:bet_amount 投注金额
		WinAmount          float64   `json:"winAmount" orm:"column(win_amount)"`                      //数据库字段:win_amount 如果赢金额,投注完成时计算好,派奖使用
		AvailableStatus    string    `json:"availableStatus" orm:"size(16);column(available_status)"` //数据库字段:available_status 有效状态: 有效 Available，取消 Cancel，重新结算 Resettle
		AvailableBetAmount float64   `json:"availableBetAmount" orm:"column(available_bet_amount)"`   //数据库字段:available_bet_amount 有效投注金额
		ClientStatus       string    `json:"clientStatus" orm:"size(24);column(client_status)"`       //数据库字段:client_status 显示状态：已支付 Paid，已结算 Settled，取消 Cancel，结算失败 Settled_Failed
		BetStatus          string    `json:"betStatus" orm:"size(16);column(bet_status)"`             //数据库字段:bet_status 投注状态:未支付 Unpaid，已支付 Paid，作废 Invalid,超时未支付 Timeout,支付失败 Failed，支付中  Paying，异常 Exception
		WinLostStatus      string    `json:"winLostStatus" orm:"size(16);column(win_lost_status)"`    //数据库字段:win_lost_status 输赢状态：创建 Create,输 Lose，赢 Win,和 Tie
		PostStatus         string    `json:"postStatus" orm:"size(16);column(post_status)"`           //数据库字段:post_status 派奖状态：创建 Create，待派奖 Ready，派彩中 Doing，作废 Invalid，已派奖 Paid ，已退款 Refund , 失败 Failed , 重新结算 Resettle
		ClientType         string    `json:"clientType" orm:"size(16);column(client_type)"`           //数据库字段:client_type 客户端类型:安卓 Android，IOS IOS,电脑端H5 PC_H5，手机端H5 Mobile_H5
		GameResult         string    `json:"gameResult" orm:"size(16);column(game_result)"`           //数据库字段:game_result 游戏结果
		SettleTime         time.Time `json:"settleTime" orm:"column(settle_time)"`                    //数据库字段:settle_time 输赢结算时间
		BetDoneTime        time.Time `json:"betDoneTime" orm:"column(bet_done_time);type(datetime)"`  //数据库字段:bet_done_time 投注完成时间
		PostTime           time.Time `json:"postTime" orm:"column(post_time);type(datetime)"`         //数据库字段:post_time 派彩完成时间
		CreateTime         time.Time `json:"createTime" orm:"column(create_time);type(datetime)"`     //数据库字段:create_time 创建时间
		UpdateTime         time.Time `json:"updateTime" orm:"column(update_time);type(datetime)"`     //数据库字段:update_time 更新时间
		Summary            string    `json:"summary" orm:"size(255);column(summary)"`                 //数据库字段:summary 说明
		ManualOn           string    `json:"manualOn" orm:"size(16);column(manual_on)"`               //数据库字段:manual_on 手动投注投注:是 Y,否 N
		TrialOn            string    `json:"trialOn" orm:"size(8);column(trial_on)"`                  //数据库字段:trial_on 是否试玩: 是 Y,否 N
		Sort               int       `json:"sort" orm:"column(sort)"`                                 //数据库字段:sort 排序，同一组注单内排序
		Md5                string    `json:"md5" orm:"size(32);column(md5)"`                          //数据库字段:md5 签名
		SettleStatus       string    `json:"-" orm:"-"`                                               //结算状态,取值为"Success" "Failed" 内部使用不发送出去
	}

	/*
		BetOrder 玩家下注注单结构和结算注单结构类型有些不同 该结构为玩家下注注单结构
		该注单数据从能力中心获取 int64类型以string类型传过来
	*/
	BetOrder struct {
		Id                 int64     `json:"id"`                   //数据库字段:id id
		MerchantId         int64     `json:"merchantId"`           //数据库字段:merchant_id 商户id
		SiteId             int64     `json:"siteId"`               //数据库字段:site_id 站点id
		UserId             int64     `json:"userId"`               //数据库字段:user_id 用户id
		Username           string    `json:"username"`             //数据库字段:username 用户名
		SiteUsername       string    `json:"siteUsername"`         //数据库字段:site_username 站点用户名
		Nickname           string    `json:"nickname"`             //数据库字段:nickname 用户昵称
		GroupId            int64     `json:"groupId"`              //数据库字段:group_id 组id
		OrderNo            int64     `json:"orderNo"`              //数据库字段:order_no 订单号
		GameRoomId         int64     `json:"gameRoomId"`           //数据库字段:game_room_id 游戏房间id
		GameRoomNo         string    `json:"gameRoomNo"`           //数据库字段:game_room_no 房号
		WorkerId           int64     `json:"workerId"`             //据库字段:worker_id 现场员工id
		WorkerUsername     string    `json:"workerUsername"`       //数据库字段:worker_username 主播名称
		GameRoundId        int64     `json:"gameRoundId"`          //数据库字段:game_round_id 游戏局id
		GameRoundNo        string    `json:"gameRoundNo"`          //数据库字段:game_round_no 局号
		GameCategoryId     int64     `json:"gameCategoryId"`       //数据库字段:game_category_id 游戏分类id
		GameCategory       string    `json:"gameCategory"`         //数据库字段:game_category 游戏分类名称
		GameId             int64     `json:"gameId"`               //数据库字段:game_id 游戏id
		GameCode           int64     `json:"gameCode"`             //数据库字段:game_code 游戏编码
		Game               string    `json:"game"`                 //数据库字段:game 游戏名称
		GameWagerId        int64     `json:"gameWagerId"`          //数据库字段:game_wager_id 玩法id
		GameWager          string    `json:"gameWager"`            //数据库字段:game_wager 玩法名
		GameWagerCode      string    `json:"gameWagerCode"`        //数据库字段:game_wager_code 玩法编码
		Currency           string    `json:"currency"`             //数据库字段:currency 币种
		Price              float64   `json:"price"`                //数据库字段:price 单价
		Num                int32     `json:"num"`                  //数据库字段:num 数量
		BetOdds            float64   `json:"betOdds"`              //数据库字段:bet_odds 投注时赔率
		BetMultiple        int32     `json:"betMultiple"`          //数据库字段:bet_multiple 投注倍数
		BetAmount          float64   `json:"betAmount"`            //数据库字段:bet_amount 投注金额
		AvailableBetAmount float64   `json:"availableBetAmountid"` //数据库字段:available_bet_amount 有效投注金额
		GameResult         string    `json:"gameResult"`           //数据库字段:game_result 游戏结果
		WinAmount          float64   `json:"winAmount"`            //数据库字段:win_amount 如果赢金额,投注完成时计算好,派奖使用
		DrawOdds           float64   `json:"drawOdds"`             //数据库字段:draw_odds 开奖时赔率
		ClientStatus       string    `json:"clientStatus"`         //数据库字段:client_status 显示状态：已支付 Paid，已结算 Settled，取消 Cancel，结算失败 Settled_Failed
		BetStatus          string    `json:"betStatus"`            //数据库字段:bet_status 投注状态:未支付 Upaid，已支付 Paid，作废 Invalid,超时未支付 Timeout,支付失败 Failed，支付中  Paying，异常 Exception
		WnLostStatus       string    `json:"wnLostStatus"`         //数据库字段:win_lost_status 输赢状态：创建 Create,输 Lose，赢 Win,和 Tie
		PostStatus         string    `json:"postStatus"`           //数据库字段:post_status 派奖状态：创建 Create，待派奖 Ready，派彩中 Doing，作废 Invalid，已派奖 Paid ，已退款 Refund ,失败 Failed
		AdvanceRate        float64   `json:"advanceRate"`          //数据库字段:advance_rate 大奖池垫资比例
		ManualOn           string    `json:"manualOn"`             //数据库字段:manual_on 手动投注投注:是 Y,否 N
		TrialOn            string    `json:"trialOn"`              //数据库字段:trial_on 是否试玩: 是 Y,否 N
		ClientType         string    `json:"clientType"`           //数据库字段:client_type 客户端类型:安卓 Android，IOS IOS,电脑端H5 PC_H5，手机端H5 Mobile_H5
		BetDoneTime        time.Time `json:"betDoneTime"`          //数据库字段:bet_done_time 投注完成时间
		SettleTime         time.Time `json:"settleTime"`           //数据库字段:settle_time 输赢结算时间
		PostTime           time.Time `json:"postTime"`             //数据库字段:post_time 派彩完成时间
		OperatorId         int64     `json:"operatorId"`           //数据库字段:operator_id 操作人id
		Operator           string    `json:"operator"`             //数据库字段:operator 操作人
		CreateTime         time.Time `json:"createTime"`           //数据库字段:create_time 创建时间
		UpdateTime         time.Time `json:"updateTime"`           //数据库字段:update_time 更新时间
		Summary            string    `json:"summary"`              //数据库字段:summary 说明
		Md5                string    `json:"md5"`                  //数据库字段:md5 数据指纹
	}

	//BetOrderExtraResp 玩家下注回包需要增加的部分
	BetOrderExtraResp struct {
		BetOdds            float64 `json:"betOdds"`            //数据库字段:bet_odds 投注时赔率
		AvailableBetAmount float64 `json:"availableBetAmount"` //数据库字段:available_bet_amount 有效投注金额
		WinAmount          float64 `json:"winAmount"`          //数据库字段:win_amount 如果赢金额,投注完成时计算好,派奖使用
	}

	//BetOrderResp 玩家下注返回结构体
	BetOrderResp struct {
		*BetOrder
		*BetOrderExtraResp
	}

	//BetLimitRule 投注限红规则
	BetLimitRule struct {
		Id                    string  `json:"id"`                    //数据库字段:id int64
		Name                  string  `json:"name"`                  //数据库字段:name 名称	string
		BetLimitRuleGroupId   string  `json:"betLimitRuleGroupId"`   //数据库字段:bet_limit_rule_group_id 限红分组id	integer(int64)
		GameCategoryId        string  `json:"gameCategoryId"`        //数据库字段:game_category_id 游戏分类id integer(int64)
		GameId                string  `json:"gameId"`                //数据库字段:game_id 游戏id integer(int64)
		GameWagerId           string  `json:"gameWagerId"`           //数据库字段:game_wager_id 玩法id integer(int64)
		Currency              string  `json:"currency"`              //数据库字段:currency 币种	string
		MinAmount             float64 `json:"minAmount"`             //数据库字段:min_amount 最小金额 number(double)
		MaxAmount             float64 `json:"maxAmount"`             //数据库字段:max_amount 最大金额 number(double)
		Summary               string  `json:"summary"`               //数据库字段:summary 说明 string
		BetLimitRuleGroupName string  `json:"betLimitRuleGroupName"` //数据库字段:betLimitRuleGroupName 盘口名称
	}

	// UserBetLimitInfo 用户个人限红信息
	UserBetLimitInfo struct {
		Id                  string         //数据库字段:id id	int64
		UserId              string         //数据库字段:user_id 用户id	int64
		BetLimitRuleGroupId string         //数据库字段:bet_limit_rule_group_id int64
		Currency            string         //数据库字段:currency 币种编码 string
		BetLimitRuleList    []BetLimitRule //限红规则分组集合	array	BetLimitRuleDTO
	}

	// UserBetLimitRequest 个人限红请求
	UserBetLimitRequest struct {
		UserId   string //数据库字段:user_id 用户id	int64
		Currency string //数据库字段:bet_limit_rule_group_id int64
	}

	// UserBetLimitBatchRequest 批量个人限红请求
	UserBetLimitBatchRequest struct {
		UserBetLimitList []UserBetLimitRequest
	}

	//Game 游戏信息
	Game struct {
		Id             string //数据库字段:id id	integer(int64)
		GameCategoryId string //数据库字段:game_category_id 游戏分类id	integer(int64)
		Code           string //数据库字段:code 分类编码	integer(int64)
		Name           string //数据库字段:name 游戏名称	string
		Countdown      int32  //数据库字段:countdown 倒计时时长,秒	integer(int32)
		Duration       int32  //数据库字段:duration 游戏总时长,秒	integer(int32)
		PieceCount     int32  //数据库字段:piece_count 多少球，多少张牌等	integer(int32)
		Type           string //数据库字段:type 类型：游戏 Game，大奖游戏 Jackpot，礼物小游戏 Gift	string
		TableOn        string //数据库字段:table_on 是否桌台游戏:是 Y,否 N	string
		Status         string //数据库字段:status 状态:启用 Enable，停用 Disable，维护 Maintain	string
		Api            string //数据库字段:api 接口地址,通常为域名	string
		Token          string //数据库字段:token 授权访问token	string
		ImplementType  string //数据库字段:implement_type 算法实现类型：组件 Component，接口 API	string
		OperatorId     string //数据库字段:operator_id 操作人id	integer(int64)
		Operator       string //数据库字段:operator 操作人	string
		CreateTime     string //数据库字段:create_time 创建时间	string(date-time)
		UpdateTime     string //数据库字段:update_time 更新时间	string(date-time)
		Summary        string //数据库字段:summary 说明	string
	}

	//GameWager 游戏玩法信息
	GameWager struct {
		Id              string  //数据库字段:id id    integer(int64)
		GameCategoryId  string  //数据库字段:game_category_id 游戏分类id    integer(int64)
		GameId          string  //数据库字段:game_id 游戏id    integer(int64)
		Code            int32   //数据库字段:code 玩法编码    integer(int32)
		Name            string  //数据库字段:name 玩法名称    string
		Odds            float64 //数据库字段:odds 赔率    number(double)
		SettleCount     int32   //数据库字段:settle_count 结算球数，张数等    integer(int32)
		BetImplement    string  //数据库字段:bet_implement 投注算法名称    string
		SettleImplement string  //数据库字段:settle_implement 结算算法实现名称    string
		Status          string  //数据库字段:status 状态:启用 Enable，停用 Disable，维护 Maintain    string
		OperatorId      string  //数据库字段:operator_id 操作人id    integer(int64)
		Operator        string  //数据库字段:operator 操作人    string
		CreateTime      string  //数据库字段:create_time 创建时间    string(date-time)
		UpdateTime      string  //数据库字段:update_time 更新时间    string(date-time)
		Summary         string  //数据库字段:summary 说明    string
		Md5             string  //数据库字段:md5 数据指纹    string
	}

	//RoomDetailedInfo 房间内桌台限红和玩法信息与赔率
	RoomDetailedInfo struct {
		Id               string  `json:"id"`              //数据库字段:id id	integer(int64)
		GamePlatformId   string  `json:"gamePlatformId"`  //数据库字段:game_platform_id 游戏平台id	integer(int64)
		GameCategoryId   string  `json:"gameCategoryId"`  //数据库字段:game_category_id 游戏分类id	integer(int64)
		GameId           string  `json:"gameId"`          //数据库字段:game_id 游戏id	integer(int64)
		RoomNo           string  `json:"roomNo"`          //数据库字段:room_no 房号	string
		KickOutLimit     int32   `json:"kickOutLimit"`    //数据库字段:kick_out_limit 踢出局数,不投注踢出局数
		TableNo          string  `json:"tableNo"`         //数据库字段:table_no 桌台编码	string
		UserLimit        int32   `json:"userLimit"`       //数据库字段:user_limit 人数限制	integer(int32)
		UserTotal        int32   `json:"userTotal"`       //数据库字段:user_total 加入房间的总人数	integer(int32)
		OnlineUserTotal  int32   `json:"onlineUserTotal"` //数据库字段:online_user_total 实时在线用户数	integer(int32)
		JackpotRate      float64 `json:"jackpotRate"`     //数据库字段:jackpot_rate 大奖垫资比例	number(double)
		UserLostAmount   float64 `json:"userLostAmount"`  //数据库字段:user_lost_amount 用户输金额	number(double)
		UserWinAmount    float64 `json:"userWinAmount"`   //数据库字段:user_win_amount 用户赢金额	number(double)
		BackGroundColor  string  `json:"backgroundColor"` //数据库字段:background_color 背景色编码
		GameDataDelay    int32   `json:"gameDataDelay"`   //数据库字段:game_data_delay 延迟显示游戏结果,单位毫秒
		DrawDelay        int32   `json:"drawDelay"`       //数据库字段:draw_delay 延迟显示游戏开奖结果,单位毫秒
		Status           string  `json:"status"`          //数据库字段:status 状态:创建 Create,启用 Enable,停用 Disable,超时 Timeout string
		Type             string  `json:"type"`            //数据库字段:type 类型:普通 Normal，专属 Special	string
		Game             *Game
		GameWagerList    []*dto.GameWagerDTO
		BetLimitRuleList []*BetLimitRule //投注限红规则集合
	}

	Currency struct {
		Scale                int  `json:"scale"`                //小数位精度
		IsTransferFreeWallet bool `json:"isTransferFreeWallet"` //是否免转钱包 :是 true,否 false
	}

	/*
		BetWager 投注信息
	*/
	BetWager struct {
		GameWagerId int32   `json:"gameWagerId"` //玩法Id 与客户端交互int64->string
		Chip        float64 `json:"chip"`        //注码
	}
	/*
		LimitRule限红规则
	*/
	LimitRule struct {
		Currency  string  `json:"currency"`  //数据库字段:currency 币种
		MinAmount float64 `json:"minAmount"` //数据库字段:min_amount 最小金额
		MaxAmount float64 `json:"maxAmount"` //数据库字段:max_amount 最大金额
	}
	/*
		BetDeviceInfo 投注时设备信息 投注时使用
	*/
	BetDeviceInfo struct {
		ClientType string `json:"clientType"` //客户端类型
		ClientIp   string `json:"clientIp"`   //客户端ip
		UserAgent  string `json:"userAgent"`  //客户端信息
		Width      int    `json:"width"`      //客户端屏幕宽
		Height     int    `json:"height"`     //客户端屏幕高
		OS         string `json:"os"`         //操作系统信息
		Language   string `json:"language"`   //语言信息
		Browser    string `json:"browser"`    //浏览器信息
		Latency    int32  `json:"latency"`    //网络延迟
		Network    string `json:"network"`    //网络3g,4g,5g等
		ISP        string `json:"isp"`        //网络运营商
		AppVersion string `json:"appVersion"` //应用版本号
	}

	/*
		BetVO 投注参数接收对象
	*/
	BetVO struct {
		GameRoomId  string        `json:"gameRoomId"`  //游戏房间id 与客户端交互int64->string
		GameRoundId string        `json:"gameRoundId"` //游戏局id 与客户端交互int64->string
		Currency    string        `json:"currency"`    //币种
		BetAmount   float64       `json:"betAmount"`   //投注金额,计算bets的投注总计
		GameTime    string        `json:"gameTime"`    //游戏时间
		Bets        []BetWager    `json:"bets"`        //投注玩法信息集合
		LimitRule   LimitRule     `json:"limitRule"`   //限红规则
		Device      BetDeviceInfo `json:"device"`      //设备信息
	}

	/*
		BetResult 投注返回结果
	*/
	BetResult struct {
		OrderNo     string  `json:"orderNo"`     //数据库字段:order_no 订单号
		GameWagerId string  `json:"gameWagerId"` //数据库字段:game_wager_id 玩法id
		Currency    string  `json:"currency"`    //数据库字段:currency 币种
		BetAmount   float64 `json:"betAmount"`   //数据库字段:bet_amount 投注金额
	}

	/*
		BetCancelParam 投注取消参数接收对象
	*/
	BetCancelParam struct {
		GameRoomId  string        `json:"gameRoomId"`  //游戏房间id 与客户端交互int64->string
		GameRoundId string        `json:"gameRoundId"` //游戏局id 与客户端交互int64->string
		OrderNoList []string      `json:"orderNoList"` //订单号集合
		Device      BetDeviceInfo `json:"device"`      //设备信息
	}

	/*
		BetConfirmParam 投注确认参数接收对象
	*/
	BetConfirmParam struct {
		GameRoomId   string        `json:"gameRoomId"`   //游戏房间id 与客户端交互int64->string
		GameRoundId  string        `json:"gameRoundId"`  //游戏局id 与客户端交互int64->string
		Currency     string        `json:"currency"`     //币种
		BetAmount    float64       `json:"betAmount"`    //投注金额
		GameTime     string        `json:"gameTime"`     //游戏时间
		OriginalBets []BetWager    `json:"originalBets"` //原始投注,记录所有投注信息
		Bets         []BetWager    `json:"bets"`         //游投注玩法信息,最终投注信息
		LimitRule    LimitRule     `json:"limitRule"`    //游投注玩法信息,最终投注信息
		Device       BetDeviceInfo `json:"device"`       //游投注玩法信息,最终投注信息
	}

	/*
		BetRecord 投注记录对象 返回给前端
	*/
	BetRecord struct {
		UserId         string  `json:"userId"`         //数据库字段:user_id 用户id 与客户端交互int64->string
		GroupId        string  `json:"groupId"`        //数据库字段:group_id 组id 与客户端交互int64->string
		OrderNo        string  `json:"orderNo"`        //数据库字段:order_no 订单号 与客户端交互int64->string
		GameRoomId     string  `json:"gameRoomId"`     //数据库字段:game_room_id 游戏房间id与客户端交互 int64->string
		GameRoundId    string  `json:"gameRoundId"`    //数据库字段:game_round_id 游戏局id与客户端交互 int64->string
		GameCategoryId string  `json:"gameCategoryId"` //数据库字段:game_category_id 游戏分类id与客户端交互 int64->string
		GameId         string  `json:"gameId"`         //数据库字段:game_id 游戏id 与客户端交互int64->string
		GameWagerId    string  `json:"gameWagerId"`    //数据库字段:game_wager_id 玩法id 与客户端交互int64->string
		Currency       string  `json:"currency"`       //数据库字段:currency 币种
		Price          float64 `json:"price"`          //数据库字段:price 单价
		Num            int     `json:"num"`            //数据库字段:num 数量
		BetOdds        float64 `json:"betOdds"`        //数据库字段:bet_odds 投注时赔率
		BetMultiple    int     `json:"betMultiple"`    //数据库字段:bet_multiple 投注倍数
		BetAmount      float64 `json:"betAmount"`      //数据库字段:bet_amount 投注金额
	}
)
