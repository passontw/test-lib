package db

import (
	"github.com/beego/beego/v2/client/orm"
	"sl.framework.com/tool"
	"sl.framework.com/trace"
	"time"
)

// SubscriberInfo 用于存储订阅者信息的结构
type SubscriberInfo struct {
	Id             int64  `orm:"pk"`                     // 自增主键，唯一标识每一条记录
	Name           string `orm:"size(30)"`               // 订阅者名称，例如 MainServer
	Token          string `orm:"size(32)"`               // 订阅者鉴权 Token，用于身份验证，用户自行设置
	SubscribedVids string `orm:"size(255)"`              // 订阅的 vid 例如 B001
	Endpoint       string `orm:"size(255)"`              // 订阅者的 Endpoint，包含完整的 URL 或 IP 地址及路径
	Status         string `orm:"size(32);default('启用')"` // 状态，0（禁用） 1（启用）
	IsOnline       string `orm:"size(32);default('在线')"` // 在线状态，0（离线） 1（在线）
	GameRoomId     int64  `orm:""`                       // 房间id 与能力平台房间id对应与订阅的vid是1对1关系
}

func GetSubscribers() []SubscriberInfo {
	timeProfiler := tool.NewTimerProfiler("get subscribers", 500*time.Millisecond)
	defer timeProfiler.Stop(true)
	o := orm.NewOrm()
	// 原生 SQL 查询
	sql := `SELECT a.* FROM subscriber_info a RIGHT JOIN video_info b ON a.subscribed_vids = b.vid WHERE b.gmtype = 'SHB'`
	var subscribers []SubscriberInfo
	_, err := o.Raw(sql).QueryRows(&subscribers)
	if err != nil {
		trace.Error("Failed to load subscribers info from database: %v", err)
	}
	return subscribers
}
