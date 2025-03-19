package db

import (
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"sl.framework.com/resource/protocol"
	snowflaker "sl.framework.com/resource/snow_flake_id"
	"sl.framework.com/tool"
	"sl.framework.com/trace"
	"time"
)

// VideoInfo represents the video_info table.
type VideoInfo struct {
	Id            int64  `orm:"pk"`                                           // 自增主键，唯一标识每一条记录
	Vid           string `orm:"size(8)" json:"vid"`                           // 视频标识符，例如 B001
	Status        string `orm:"size(32);default('激活')" json:"status"`         // 状态，0（未激活），1（激活）
	Gmtype        string `orm:"size(8)" json:"gmtype"`                        // 游戏类型，例如 BAC、TBAC
	Bettime       int8   `orm:"default(0)" json:"bettime"`                    // 投注时间，单位为秒，默认值为 0
	Shoe          int    `orm:"null" json:"shoe"`                             // 换靴码（用于标识换靴的唯一性），例如 46701
	CurrentGmcode string `orm:"size(20);null" json:"currentGmcode"`           // 当前 GMCODE
	CurrentDealer string `orm:"size(50);null" json:"currentDealer"`           // 当前荷官
	GameStatus    string `orm:"size(64);default('CLOSED')" json:"gameStatus"` // 视频开启时的状态 'CLOSED', 'CAN_BET', 'GAME_DATA', 'NEW_SHOE'
}

// IsGameClose 局是否关闭
func (r *VideoInfo) IsGameClose() bool {
	return r.GameStatus == protocol.GetStatusString(protocol.GAME_STATUS_CLOSED)
}

// IsVideoClosed 视频是否关闭
func (r *VideoInfo) IsVideoClosed() bool {
	return r.Status == protocol.VIDEO_STATUS_CLOSED
}

func (r *VideoInfo) String() string {
	return fmt.Sprintf("vid=%s, status=%s, gmtype=%s, bettime=%d, shoe=%d, "+
		"currentgmcode=%s, currentDealer=%s, gamestatus=%s",
		r.Vid, r.Status, r.Gmtype, r.Bettime, r.Shoe, r.CurrentGmcode, r.CurrentDealer, r.GameStatus)
}

func (r *VideoInfo) Get() *VideoInfo {
	timeProfiler := tool.NewTimerProfiler(fmt.Sprintf("get video=%s", r.Vid), 500*time.Millisecond)
	defer timeProfiler.Stop(true)
	// 默认使用 default，你可以指定为其他数据库
	o := orm.NewOrm()
	err := o.Raw("SELECT * FROM video_info WHERE vid = ?", r.Vid).QueryRow(r)
	if err != nil {
		trace.Error("failed to get table: video_info data: %s", err.Error())
		return nil
	}
	trace.Info("%s", r)
	return r
}

func (r *VideoInfo) NewShoe() error {
	timeProfiler := tool.NewTimerProfiler(fmt.Sprintf("new shoe=%s", r.Vid), 500*time.Millisecond)
	defer timeProfiler.Stop(true)
	// 默认使用 default，你可以指定为其他数据库
	o := orm.NewOrm()
	x, err := o.Raw("UPDATE video_info SET shoe = ?, game_status = ? WHERE vid = ?", r.Shoe, r.GameStatus, r.Vid).Exec()
	if err != nil {
		return fmt.Errorf("failed to update table: video_info data, error: %v", err)
	}
	affected, err := x.RowsAffected()
	trace.Info("update video_info game_status success, the number of rows affected is %d. err is: %v", affected, err)
	return nil
}

func (r *VideoInfo) UpdateGameStatus() error {
	timeProfiler := tool.NewTimerProfiler(fmt.Sprintf("update video=%s state", r.Vid), 500*time.Millisecond)
	defer timeProfiler.Stop(true)
	// 默认使用 default，你可以指定为其他数据库
	o := orm.NewOrm()
	x, err := o.Raw("UPDATE video_info SET game_status = ? WHERE vid = ?", r.GameStatus, r.Vid).Exec()
	if err != nil {
		return fmt.Errorf("failed to update table: video_info data, error: %v", err)
	}
	affected, err := x.RowsAffected()
	trace.Info("update video_info game_status success, the number of rows affected is %d. err is: %v", affected, err)

	return nil
}

func ExternalGet(vid string) *VideoInfo {
	timeProfiler := tool.NewTimerProfiler(fmt.Sprintf("get video=%s", vid), 500*time.Millisecond)
	defer timeProfiler.Stop(true)
	// 默认使用 default，你可以指定为其他数据库
	o := orm.NewOrm()
	var err error
	r := new(VideoInfo)
	err = o.Raw("SELECT * FROM video_info WHERE vid = ?", vid).QueryRow(r)
	if err != nil {
		trace.Error("failed to get table: video_info data: %s", err.Error())
		return nil
	}
	trace.Info("%s", r)
	return r
}

func ExternalVideoInsert(vid, gmtype, gmcode, dealer string, shoe int) {
	timeProfiler := tool.NewTimerProfiler(fmt.Sprintf("insert video=%s", vid), 500*time.Millisecond)
	defer timeProfiler.Stop(true)
	// 默认使用 default，你可以指定为其他数据库
	db := orm.NewOrm()

	// 执行插入操作，将订阅的Vid信息插入到 video_info 表中
	res, err := db.Raw("INSERT INTO video_info (id, vid, status, gmtype, shoe, current_gmcode, current_dealer, game_status) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		snowflaker.UniqueId(), vid, protocol.VIDEO_STATUS_OPENED, gmtype, shoe, gmcode, dealer, protocol.GetStatusString(protocol.GAME_STATUS_CLOSED)).Exec()
	if err != nil { // 插入操作失败
		trace.Error("insert vid info failed: %v", err)
		return
	}
	// 检查更新操作是否成功
	if i, err := res.RowsAffected(); err != nil || i == 0 {
		// 更新操作未影响任何行
		trace.Error("insert vid info failed, the number of rows affected is 0. err is: %v", err)
	}
}

func ExternalNewShoe(vid, status string, shoe int) error {
	timeProfiler := tool.NewTimerProfiler(fmt.Sprintf("new vid=%s shoe=%d", vid, shoe), 500*time.Millisecond)
	defer timeProfiler.Stop(true)
	// 默认使用 default，你可以指定为其他数据库
	o := orm.NewOrm()

	x, err := o.Raw("UPDATE video_info SET shoe = ?, game_status = ? WHERE vid = ?", shoe, status, vid).Exec()
	if err != nil {
		return fmt.Errorf("failed to update table: video_info data, error: %v", err)
	}
	affected, err := x.RowsAffected()
	trace.Info("update video_info=%s shoe=%d success, the number of rows affected is %d. err is: %v", vid, shoe, affected, err)
	return nil
}

func ExternalUpdateGameStatus(vid, state string) error {
	timeProfiler := tool.NewTimerProfiler(fmt.Sprintf("update video=%s state=%s", vid, state), 500*time.Millisecond)
	defer timeProfiler.Stop(true)
	// 默认使用 default，你可以指定为其他数据库
	o := orm.NewOrm()
	x, err := o.Raw("UPDATE video_info SET game_status = ? WHERE vid = ?", state, vid).Exec()
	if err != nil {
		return fmt.Errorf("failed to update table: video_info data, error: %v", err)
	}
	affected, err := x.RowsAffected()
	trace.Info("update video_info=%s game_status=%s success, the number of rows affected is %d. err is: %v", vid, state, affected, err)
	return nil
}

func GetLastRound(vid string) string {
	timeProfiler := tool.NewTimerProfiler(fmt.Sprintf("get current round vid: %s", vid), 500*time.Millisecond)
	defer timeProfiler.Stop(true)
	o := orm.NewOrm()
	var err error
	var gmcode string
	// 找倒数第二条数据 第一条可能状态未就绪
	err = o.Raw("SELECT current_gmcode FROM video_info where vid = ?", vid).QueryRow(&gmcode)
	if err != nil {
		trace.Error("failed to get table: video_info data, error: %s", err.Error())
		return gmcode
	}
	return gmcode
}
