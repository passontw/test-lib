package dto

type WorkerDTO struct {
	Id              string `json:"id"`              //数据库字段:id id
	UserName        string `json:"username"`        //数据库字段:username 用户名
	Password        string `json:"password"`        //数据库字段:password 密码
	Name            string `json:"name"`            //数据库字段:name 姓名
	Nickname        string `json:"nickname"`        //数据库字段:nickname 昵称
	Birthday        string `json:"birthday"`        //数据库字段:birthday 生日
	HomeTown        string `json:"homeTown"`        //数据库字段:home_town 籍贯
	Favorite        string `json:"favorite"`        //数据库字段:favorite 爱好,逗号分割
	Facebook        string `json:"facebook"`        //数据库字段:facebook 脸书账号
	Avatar          string `json:"avatar"`          //数据库字段:avatar 头像地址
	Type            string `json:"type"`            //数据库字段:type 类型:主播 Anchor，荷官 Dealer
	LiveStatus      string `json:"liveStatus"`      //数据库字段:live_status 直播状态:空闲 Free,直播 Live
	LiveCount       int    `json:"liveCount"`       //数据库字段:live_count 直播次数
	GiftAmountTotal string `json:"giftAmountTotal"` //数据库字段:gift_amount_total 礼物总收入
	Status          string `json:"status"`          //数据库字段:status 状态:创建 Create,启用 Enable，停用 Disable
	OperatorId      string `json:"operatorId"`      //数据库字段:operator_id 操作人id
	Operator        string `json:"operator"`        //数据库字段:operator 操作人
	CreateTime      string `json:"createTime"`      //数据库字段:create_time 创建时间
	UpdateTime      string `json:"updateTime"`      //数据库字段:update_time 更新时间
	Summary         string `json:"summary"`         //数据库字段:summary 说明
}
