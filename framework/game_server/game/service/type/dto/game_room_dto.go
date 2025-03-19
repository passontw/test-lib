package dto

import "time"

type GameRoomDTO struct {
	Id              string    `json:"id"`              //数据库字段:id id	integer(int64)
	GamePlatformId  string    `json:"gamePlatformId"`  //数据库字段:game_platform_id 游戏平台id	integer(int64)
	GameCategoryId  string    `json:"gameCategoryId"`  //数据库字段:game_category_id 游戏分类id	integer(int64)
	GameId          string    `json:"gameId"`          //数据库字段:game_id 游戏id	integer(int64)
	RoomNo          string    `json:"roomNo"`          //数据库字段:room_no 房号	string
	KickOutLimit    int32     `json:"kickOutLimit"`    //数据库字段:kick_out_limit 踢出局数,不投注踢出局数
	TableNo         string    `json:"tableNo"`         //数据库字段:table_no 桌台编码	string
	UserLimit       int32     `json:"userLimit"`       //数据库字段:user_limit 人数限制	integer(int32)
	UserTotal       int32     `json:"userTotal"`       //数据库字段:user_total 加入房间的总人数	integer(int32)
	OnlineUserTotal int32     `json:"onlineUserTotal"` //数据库字段:online_user_total 实时在线用户数	integer(int32)
	JackpotRate     float64   `json:"jackpotRate"`     //数据库字段:jackpot_rate 大奖垫资比例	number(double)
	UserLostAmount  float64   `json:"userLostAmount"`  //数据库字段:user_lost_amount 用户输金额	number(double)
	UserWinAmount   float64   `json:"userWinAmount"`   //数据库字段:user_win_amount 用户赢金额	number(double)
	BackGroundColor string    `json:"backgroundColor"` //数据库字段:background_color 背景色编码
	GameDataDelay   int32     `json:"gameDataDelay"`   //数据库字段:game_data_delay 延迟显示游戏结果,单位毫秒
	DrawDelay       int32     `json:"drawDelay"`       //数据库字段:draw_delay 延迟显示游戏开奖结果,单位毫秒
	Status          string    `json:"status"`          //数据库字段:status 状态:创建 Create,启用 Enable,停用 Disable,超时 Timeout string
	Type            string    `json:"type"`            //数据库字段:type 类型:普通 Normal，专属 Special	string
	OperatorId      int64     `json:"operatorId"`      //数据库字段:operator_id 操作人id
	Operator        string    `json:"operator"`        //数据库字段:operator 操作人
	CreateTime      time.Time `json:"createTime"`      //数据库字段:create_time 创建时间
	UpdateTime      time.Time `json:"updateTime"`      //数据库字段:update_time 更新时间
	Summary         string    `json:"summary"`         //数据库字段:summary 说明
}
