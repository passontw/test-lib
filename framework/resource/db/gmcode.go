package db

import (
	"errors"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	snowflaker "sl.framework.com/resource/snow_flake_id"
	"time"
)

type GmcodeMapping struct {
	Id        int64     `orm:"auto;pk"` // 自增主键，唯一标识每一条记录
	OldGmcode string    `orm:"size(16)"`
	NewGmcode string    `orm:"size(16)"`
	CreatedAt time.Time `orm:"auto_now_add;type(datetime)"`
}

// WriteNewGMCode 方法将 oldgmcode 和 newgmcode 存储到数据库中
func WriteNewGMCode(oldgmcode string, newgmcode string) error {

	o := orm.NewOrm()
	gmcode := GmcodeMapping{
		Id:        snowflaker.UniqueId(),
		OldGmcode: oldgmcode,
		NewGmcode: newgmcode,
	}

	// 插入数据到数据库
	_, err := o.Insert(&gmcode)
	if err != nil {
		return fmt.Errorf("error storing newgmcode: %v", err)
	}
	return nil
}

// GetNewGMCode 方法根据 oldgmcode 获取最新的 newgmcode
func GetNewGMCode(oldgmcode string) (string, error) {

	o := orm.NewOrm()
	var gmcode GmcodeMapping

	// 根据 key_prefix 查找最新的 newgmcode
	err := o.QueryTable("gmcode_mapping").
		Filter("old_gmcode", oldgmcode).
		Limit(1).
		One(&gmcode)
	if errors.Is(err, orm.ErrNoRows) {
		return "", fmt.Errorf("no newgmcode found for this oldgmcode")
	} else if err != nil {
		return "", fmt.Errorf("error retrieving newgmcode: %v", err)
	}
	return gmcode.NewGmcode, nil
}
