package dao

import (
	"sl.framework.com/game_server/game/dao/gamedb"
	"sl.framework.com/game_server/game/dao/uiddb"
	"sl.framework.com/trace"
)

// OrmInit 初始化UidDb GameDb Orm 并注册表
func OrmInit() error {
	//Database game order Orm初始化
	if err := gamedb.OrmGameDbInit(); nil != err {
		trace.Error("OrmInit OrmGameDbInit failed, error=%v", err.Error())
		return err
	}

	//Database uid order Orm初始化
	if err := uiddb.OrmUidDbInit(); nil != err {
		trace.Error("OrmInit OrmUidDbInit failed, error=%v", err.Error())
		return err
	}

	// 创建 table 如果存在则跳过 不使用orm创建表避免出现问题
	//if err := orm.RunSyncdb("default", false, true); nil != err {
	//	trace.Error("OrmInit initialize RunSync dao, error=%v", err.Error())
	//	return err
	//}

	return nil
}
