package uiddb

import (
	"errors"
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
	"sl.framework.com/async"
	"sl.framework.com/game_server/conf"
	"sl.framework.com/game_server/error_code"
	"sl.framework.com/tool"
	"sl.framework.com/trace"
	"sync"
	"time"
)

/*
	beego orm 并不提供显式的关闭数据库连接的方法，通常依赖于数据库驱动的连接池管理
*/

var uidDbInitOnce sync.Once

func GetUidDBOrm() orm.Ormer {
	alias := conf.GetUidDbAliasName()

	return orm.NewOrmUsingDB(alias)
}

func OrmUidDbInit() (err error) {
	fn := func() {
		dsn := conf.GetMySqlDsnUidDb()
		alias := conf.GetUidDbAliasName()
		trace.Info("uid dao Initialize alias=%v, dsn=%v", alias, dsn)

		// 注册驱动 RegisterDataBase
		if err = orm.RegisterDriver("mysql", orm.DRMySQL); nil != err {
			trace.Error("uid dao initialize RegisterDriver, error=%v", err.Error())
			return
		}

		// 设置默认数据库
		if err = orm.RegisterDataBase(alias, "mysql", dsn); nil != err {
			trace.Error("uid dao initialize RegisterDataBase, error=%v", err.Error())
			return
		}

		// 注册定义的 uid model
		orm.RegisterModel(new(GoWorkerNode))

		//设置最大连接数
		orm.SetMaxOpenConns(alias, 5)
		//设置最大空闲连接数
		orm.SetMaxIdleConns(alias, 5)
	}
	uidDbInitOnce.Do(fn)

	return
}

/*
	UniqueIdGenerator连接g32_uid数据库，读表go_worker_node
*/

var (
	once              sync.Once
	uniqueIdGenerator *UniqueIdGenerator
)

type UniqueIdGenerator struct {
	maxCounts int64 //最大尝试次数
}

func GetUniqueIdGeneratorInstance() *UniqueIdGenerator {
	if uniqueIdGenerator == nil {
		once.Do(func() {
			uniqueIdGenerator = &UniqueIdGenerator{
				maxCounts: errcode.MaxLoopCount,
			}
		})
	}
	return uniqueIdGenerator
}

// GetUniqueId 供UniqueIdController使用，可以考虑将请求放入队列中逐个处理
func (u *UniqueIdGenerator) GetUniqueId() (err error, serverId int64) {
	var countLoop int64 = 0
	return u.mustGetUniqueId(countLoop)
}

// AsyncSetUniqueId 程序启动时候异步设置本服务Server ID
func (u *UniqueIdGenerator) AsyncSetUniqueId() {
	async.AsyncRunCoroutine(func() {
		var countLoop int64 = 0
		err, id := u.mustGetUniqueId(countLoop)
		if nil != err {
			trace.Error("AsyncSetUniqueId failed, count loop=%v", countLoop)
			return
		}

		async.AsyncRunCoroutine(func() {
			conf.SetServerId(id)
		})
	})
}

// SetUniqueId 程序启动时候异步设置本服务Server ID
func (u *UniqueIdGenerator) SetUniqueId() {
	var countLoop int64 = 0
	err, id := u.mustGetUniqueId(countLoop)
	if nil != err {
		trace.Error("AsyncSetUniqueId failed, count loop=%v", countLoop)
		return
	}
	conf.SetServerId(id)
}

// MustGetServerId 获取ServerId并将Id+1写入数据库
func (u *UniqueIdGenerator) mustGetUniqueId(countLoop int64) (err error, serverId int64) {
	retQuery, retInsert := -1, -1

	goWorkerNode := &GoWorkerNode{HostName: tool.GetLocalIp()}
	retQuery = goWorkerNode.QueryMaxId()
	if errcode.DBErrorOK == retQuery {
		trace.Info("mustGetUniqueId uniqueIdQuery ret=%v, id=%d, count=%v", retQuery, goWorkerNode.Id, countLoop)
		countLoop++
		if countLoop > u.maxCounts {
			err = errors.New("mustGetUniqueId uniqueIdQuery failed")
			return
		}
		goWorkerNode.Id = goWorkerNode.Id + 1
		goWorkerNode.CreateTime = time.Now()
		if retInsert = goWorkerNode.Insert(); errcode.DBErrorOK != retInsert {
			time.Sleep(time.Duration(10) * time.Millisecond)
			u.mustGetUniqueId(countLoop)
		}
		if errcode.DBErrorOK == retInsert {
			serverId = goWorkerNode.Id
		} else {
			trace.Error("mustGetUniqueId uniqueIdInsert failed, retInsert=%v, id=%v, count=%v",
				retInsert, goWorkerNode.Id, countLoop)
		}
	} else {
		trace.Error("mustGetUniqueId uniqueIdQuery failed, retQuery=%v, id=%v", retQuery, goWorkerNode.Id)
		countLoop++
		if countLoop > u.maxCounts {
			err = errors.New("mustGetUniqueId uniqueIdQuery failed")
			return
		}
		time.Sleep(time.Duration(10) * time.Millisecond)
		u.mustGetUniqueId(countLoop)
	}

	return
}
