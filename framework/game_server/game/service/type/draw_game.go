package types

import (
	"time"
)

const GameDrawResultCacheSize = 100 //游戏开奖结果缓存长度，默认保留最新的 100条
type (
	/*
		GameDraw
		能力平台发送该结构给游戏服
		游戏服组装GameDrawInput结构发送到mq 通过mq分发到其他的游戏服
	*/
	GameDraw struct {
		GameRoomId      int64       `json:"gameRoomId"`      //数据库字段:game_room_id 游戏房间id
		GameRoundId     int64       `json:"gameRoundId"`     //数据库字段:game_round_id 游戏局id
		GameRoundNo     string      `json:"gameRoundNo"`     //数据库字段:game_round_no 游戏局id
		GameCategoryId  int64       `json:"gameCategoryId"`  //数据库字段:game_category_id 游戏分类id
		GameId          int64       `json:"gameId"`          //数据库字段:game_id 游戏id
		Payload         interface{} `json:"payload"`         //游戏结果result 此处使用interface{}接收 在不同游戏中分别解析使用
		OrderPlanIdList []int64     `json:"orderPlanIdList"` //订单执行计划id列表

		/*以下字段会发送 但是当前不使用*/
		Id     int64  `json:"id"`     //数据库字段:id id
		Status string `json:"status"` //数据库字段:status 状态:创建 Create,就绪 Ready,处理中 Doing,失败 Failed,完成 Done
		Md5    string `json:"md5"`    //数据库字段:md5 数据指纹
	}

	/*
		GameDrawInput
		游戏结算接收结构体 由MQ通过game-draw消息传递给游戏服
		不同游戏Payload可能不同 通过interface{}传递具体游戏做解析使用
	*/
	GameDrawInput struct {
		GameRoomId     int64       //数据库字段:game_room_id 游戏房间id
		GameRoundId    int64       //数据库字段:game_round_id 游戏局id
		GameRoundNo    string      //数据库字段:game_round_no 游戏局id
		GameCategoryId int64       //数据库字段:game_category_id 游戏分类id
		GameId         int64       //数据库字段:game_id 游戏id
		Payload        interface{} //游戏结果result 此处使用interface{}接收 在不同游戏中分别解析使用
		OrderNoList    []int64     //要结算的注单的orderList
	}

	/*
		DrawOrder 玩家下注注单结构和结算注单结构类型有些不同 该结构为结算注单结构
		该注单数据从能力中心获取 int64类型以string类型传过来
	*/
	DrawOrder struct {
		Id                 string  //数据库字段:id id
		MerchantId         string  //数据库字段:merchant_id 商户id
		SiteId             string  //数据库字段:site_id 站点id
		UserId             string  //数据库字段:user_id 用户id
		Username           string  //数据库字段:username 用户名
		SiteUsername       string  //数据库字段:site_username 站点用户名
		Nickname           string  //数据库字段:nickname 用户昵称
		GroupId            string  //数据库字段:group_id 组id
		OrderNo            string  //数据库字段:order_no 订单号
		GameRoomId         string  //数据库字段:game_room_id 游戏房间id
		GameRoomNo         string  //数据库字段:game_room_no 房号
		WorkerId           string  //数据库字段:worker_id 现场员工id
		WorkerUsername     string  //数据库字段:worker_username 主播名称
		GameRoundId        string  //数据库字段:game_round_id 游戏局id
		GameRoundNo        string  //数据库字段:game_round_no 局号
		GameCategoryId     string  //数据库字段:game_category_id 游戏分类id
		GameCategory       string  //数据库字段:game_category 游戏分类名称
		GameId             string  //数据库字段:game_id 游戏id
		GameCode           string  //数据库字段:game_code 游戏编码
		Game               string  //数据库字段:game 游戏名称
		GameWagerId        string  //数据库字段:game_wager_id 玩法id
		GameWager          string  //数据库字段:game_wager 玩法名
		GameWagerCode      string  //数据库字段:game_wager_code 玩法编码
		Currency           string  //数据库字段:currency 币种
		Price              float64 //数据库字段:price 单价
		Num                int     //数据库字段:num 数量
		BetOdds            float32 //数据库字段:bet_odds 投注时赔率
		BetMultiple        int     //数据库字段:bet_multiple 投注倍数
		BetAmount          float64 //数据库字段:bet_amount 投注金额
		AvailableBetAmount float64 //数据库字段:available_bet_amount 有效投注金额
		GameResult         string  //数据库字段:game_result 游戏结果
		WinAmount          float64 //数据库字段:win_amount 如果赢金额,投注完成时计算好,派奖使用
		DrawOdds           float32 //数据库字段:draw_odds 开奖时赔率
		AvailableStatus    string  //数据库字段:available_status 有效状态: 有效 Available，取消 Cancel，重新结算 Resettle
		ClientStatus       string  //数据库字段:client_status 显示状态：已支付 Paid，已结算 Settled，取消 Cancel，结算失败 Settled_Failed
		BetStatus          string  //数据库字段:bet_status 投注状态:未支付 Unpaid，已支付 Paid，作废 Invalid,超时未支付 Timeout,支付失败 Failed，支付中 Paying，异常 Exception
		WinLostStatus      string  //数据库字段:win_lost_status 输赢状态：创建 Create,输 Lose，赢 Win,和 Tie
		PostStatus         string  //数据库字段:post_status 派奖状态：创建 Create，待派奖 Ready，派彩中 Doing，作废 Invalid，已派奖 Paid ，已退款 Refund , 失败 Failed , 重新结算 Resettle
		AdvanceRate        float64 //数据库字段:advance_rate 大奖池垫资比例
		ManualOn           string  //数据库字段:manual_on 手动投注投注:是 Y,否 N
		TrialOn            string  //数据库字段:trial_on 是否试玩: 是 Y,否 N
		ClientType         string  //数据库字段:client_type 客户端类型:安卓 Android，IOS IOS,电脑端H5 PC_H5，手机端H5 Mobile_H5
		BetDoneTime        string  //数据库字段:bet_done_time 投注完成时间 string(date-time)
		SettleTime         string  //数据库字段:settle_time 输赢结算时间 string(date-time)
		PostTime           string  //数据库字段:post_time 派彩完成时间 string(date-time)
		Sort               int     //数据库字段:sort 排序，同一组注单内排序
		OrderPlanId        string  //数据库字段:order_plan_id 订单计划id
		OperatorId         string  //数据库字段:operator_id 操作人id
		Operator           string  //数据库字段:operator 操作人
		CreateTime         string  //数据库字段:create_time 创建时间string(date-time)
		UpdateTime         string  //数据库字段:update_time 更新时间string(date-time)
		Summary            string  //数据库字段:summary 说明
		Md5                string  //数据库字段:md5 数据指纹
		SettleStatus       string  `json:"-"` //结算状态,取值为"Success" "Failed" 内部使用不发送出去
	}

	// OrderPlanIdItem 根据OrderPlanId从平台中心获取的数据对应的结构
	OrderPlanIdItem struct {
		SiteId    string       //数据库字段:site_id 站点id
		OrderList []*DrawOrder //订单集合
	}

	OrderPlanIdDrawResItem struct {
		OrderPlanId        string  `json:"orderPlanId"`        //数据库字段:order_plan_id 订单计划id
		OrderNo            string  `json:"orderNo"`            //数据库字段:order_no 订单号
		WinAmount          float64 `json:"winAmount"`          //数据库字段:win_amount 如果赢金额,投注完成时计算好,派奖使用
		AvailableBetAmount float64 `json:"availableBetAmount"` //数据库字段:available_bet_amount 有效投注金额
		DrawOdds           float64 `json:"drawOdds"`           //数据库字段:draw_odds 开奖时赔率
		WinLostStatus      string  `json:"winLostStatus"`      //数据库字段:win_lost_status 输赢状态：创建 Create,输 Lose，赢 Win,和 Tie
	}

	// OrderPlanIdSentStatus 更新平台中心注单状态信息
	OrderPlanIdSentStatus struct {
		OrderDrawResHead
		Code int
	}

	OrderDrawResHead struct {
		Token string `json:"Token"`
	}

	//OrderPlanIdDrawRes 结算完成以后推送到platform信息的结构
	OrderPlanIdDrawRes struct {
		OrderDrawResHead
		SettleDTOs []*SettleDTO `json:"SettleDTOs"`
	}

	// DrawResKey 用于将结算注单统计为发送的模式
	DrawResKey struct {
		SiteId  string
		RoomId  string
		RoundId string
	}
	DrawResMap map[DrawResKey][]*OrderPlanIdDrawResItem
)

// v2版本框架新增数据结构
type (
	/*
		   GameRoundDTO 每局游戏的详情信息
			用于redis缓存

		   Status有以下几种状态

		   	Create   //创建
		   	Start    //开始投注
		   	Stop     //停止投注
		   	Doing    //开奖中
		   	Pause    //暂停
		   	Done     //已开奖
		   	Cancel   //取消
		   	Resettle //重算
		   	Resettle //异常
	*/
	GameRoundDTO struct {
		Id             string    `json:"id"`             //数据库字段:id id
		GameCategoryId string    `json:"gameCategoryId"` //数据库字段:game_category_id 游戏分类id int64<->string
		GameId         string    `json:"gameId"`         //数据库字段:game_id 游戏id
		GameRoomId     string    `json:"gameRoomId"`     //数据库字段:game_room_id 游戏房间id
		RoundNo        string    `json:"roundNo"`        //数据库字段:round_no 局号
		Currency       string    `json:"currency"`       //数据库字段:currency 币种
		UserLostAmount float64   `json:"userLostAmount"` //数据库字段:user_lost_amount 用户输金额
		UserWinAmount  float64   `json:"userWinAmount"`  //数据库字段:user_win_amount 用户赢金额
		BetTotal       float64   `json:"betTotal"`       //数据库字段:bet_total 总投注金额
		Status         string    `json:"status"`         //数据库字段:status 状态:创建 Create,开始投注 Start，停止投注 Stop，开奖中 Doing， 已开奖 Done，取消 Cancel，重算 Resettle，异常 Exception
		StartTime      time.Time `json:"startTime"`      //数据库字段:start_time 开始时间，默认当下时间
		EndTime        time.Time `json:"endTime"`        //数据库字段:end_time 结束时间，无限期则设置99年后时间
		CreateTime     time.Time `json:"createTime"`     //数据库字段:create_time 创建时间
		Summary        string    `json:"summary"`        //数据库字段:summary 说明
	}
	//游戏局结果数据对象
	GameRoundResultDTO struct {
		GameRoundId string  `json:"gameRoundId"`
		Headers     *Heads  `json:"headers"`
		Payload     Payload `json:"payload"`
		Timestamp   int64   `json:"timestamp"`
	}
	//游戏开奖数据对象
	GameDrawDataDTO struct {
		GameRoomId         int64              `json:"gameRoomId"`
		GameRoundId        int64              `json:"gameRoundId"`
		GameId             int64              `json:"gameId"`
		GameRoundNo        string             `json:"gameRoundNo"`
		GameRoundResultDTO GameRoundResultDTO `json:"gameRoundResultDTOs"`
		OrderList          []int64            `json:"orderList"`
	}
	//游戏结算数据对象
	SettleDTO struct {
		OrderNo            int64              `json:"orderNo"`            //order_no 订单号
		WinAmount          float64            `json:"winAmount"`          //win_amount 如果赢金额,投注完成时计算好,派奖使用
		AvailableBetAmount float64            `json:"availableBetAmount"` //available_bet_amount 有效投注金额
		DrawOdds           float32            `json:"drawOdds"`           //draw_odds 开奖时赔率
		WinLostStatus      string             `json:"winLostStatus"`      //win_lost_status 输赢状态：创建 Create,输 Lose，赢 Win,和 Tie
		GameRoundResult    GameRoundResultDTO `json:"gameRoundResult"`    //游戏局开奖通用结果
		SettleTime         time.Time          `json:"-"`                  //结算时间
		Md5                string             `json:"md5"`                //数据签名,游戏记录中的数据签名
	}
	//游戏结算小票
	BetReceipt struct {
		OrderNo       int64   `json:"orderNo"`
		GameWagerId   int64   `json:"gameWagerId"`
		WinAmount     float64 `json:"winAmount"`
		DrawOdds      float32 `json:"drawOdds"`
		WinLostStatus string  `json:"winLostStatus"`
	}
	//结算结果对象
	GameDrawResultVO struct {
		GameRoundId int64         `json:"gameRoundId"`
		Receipts    []*BetReceipt `json:"receipts"`
	}
	//用户消息数据对象
	UserMessageDTO struct {
		GameRoomId int64  `json:"gameRoomId"`
		UserId     string `json:"userId"`
		Command    string `json:"command"`
		Time       string `json:"time"`
		Body       any    `json:"body"`
	}
)
