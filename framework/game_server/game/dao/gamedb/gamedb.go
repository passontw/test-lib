package gamedb

import (
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
	"sl.framework.com/game_server/conf"
	types "sl.framework.com/game_server/game/service/type"
	"sl.framework.com/trace"
	"sync"
)

var gameDbInitOnce sync.Once

func GetGameDBOrm() orm.Ormer {
	alias := conf.GetUidDbAliasName()
	trace.Info("GetGameDBOrm alias=%v", alias)
	return orm.NewOrmUsingDB(alias)
}
func GetGameGDBOrm() orm.Ormer {
	alias := conf.GetGameDbAliasName()
	trace.Info("GetGameGDBOrm alias=%v", alias)
	return orm.NewOrmUsingDB(alias)
}
func OrmGameDbInit() (err error) {
	fn := func() {
		err = nil
		aliasGdb := conf.GetGameDbAliasName()
		dsnGdb := conf.GetMySqlGameDb()
		trace.Info("game dao Initialize alias=%v, dsn=%v", aliasGdb, dsnGdb)

		// 注册驱动 RegisterDataBase
		if err = orm.RegisterDriver("mysql", orm.DRMySQL); nil != err {
			trace.Error("game dao initialize RegisterDriver, error=%v", err.Error())
			return
		}

		// 设置默认数据库
		if err = orm.RegisterDataBase(aliasGdb, "mysql", dsnGdb); nil != err {
			trace.Error("game dao initialize RegisterDataBase, error=%v", err.Error())
			return
		}
		// 注册定义的 game order model
		//orm.RegisterModel(new(types.FastBacOrder))
		// 注册定义的 游戏注单记录表
		orm.RegisterModel(new(types.BetOrderV2))

		//设置最大连接数
		orm.SetMaxOpenConns(aliasGdb, 5)
		//设置最大空闲连接数
		orm.SetMaxIdleConns(aliasGdb, 5)
	}
	gameDbInitOnce.Do(fn)

	return
}
