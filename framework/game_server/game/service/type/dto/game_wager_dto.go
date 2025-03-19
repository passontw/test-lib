package dto

// 游戏玩法数据对象

type GameWagerDTO struct {
	Id              string             //数据库字段:id id    integer(int64)
	GameCategoryId  string             //数据库字段:game_category_id 游戏分类id    integer(int64)
	GameId          string             //数据库字段:game_id 游戏id    integer(int64)
	Code            int32              //数据库字段:code 玩法编码    integer(int32)
	Name            string             //数据库字段:name 玩法名称    string
	Odds            float32            //数据库字段:odds 赔率    number(double)
	RefundRate      float32            //数据库字段:refund_rate 退款比例，存储比例除以100后结果
	Type            string             //数据库字段:type 类型:固定 Fix,随机 Random
	Rate            float32            //数据库字段:rate 当type是随机类型时，比例必须有默认值，可以是0不可以是负数
	SettleCount     int32              //数据库字段:settle_count 结算球数，张数等    integer(int32)
	BetImplement    string             //数据库字段:bet_implement 投注算法名称    string
	SettleImplement string             //数据库字段:settle_implement 结算算法实现名称    string
	Status          string             //数据库字段:status 状态:启用 Enable，停用 Disable，维护 Maintain    string
	OperatorId      string             //数据库字段:operator_id 操作人id    integer(int64)
	Operator        string             //数据库字段:operator 操作人    string
	CreateTime      string             //数据库字段:create_time 创建时间    string(date-time)
	UpdateTime      string             //数据库字段:update_time 更新时间    string(date-time)
	Summary         string             //数据库字段:summary 说明    string
	Md5             string             //数据库字段:md5 数据指纹    string
	OddsList        []GameWagerOddsDTO //随机玩法列表
}
