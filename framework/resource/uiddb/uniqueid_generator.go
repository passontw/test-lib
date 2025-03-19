/*
	uiddb 包用于管理 UID 数据库的连接和唯一 ID 生成器。
	通过此包提供的 ORM 实例连接 MySQL 数据库，并且实现了唯一 ID 生成器，用于分布式环境下生成唯一的服务器 ID。
	该包自动初始化数据库连接和唯一 ID 生成器实例。
*/

package uiddb

import (
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
	"sl.framework.com/async"
	"sl.framework.com/resource/conf"
	"sl.framework.com/tool"
	"sl.framework.com/trace"
	"time"
)

const section = "uid-database"

var (
	dbOrm       orm.Ormer          // 全局 ORM 实例，用于与 UID 数据库交互
	uniqueIdGen *UniqueIdGenerator // 唯一 ID 生成器的单例实例
)

func init() {
	// 初始化 ORM 和数据库连接
	username := conf.Section(section, "user")
	password := conf.Section(section, "password")
	host := conf.Section(section, "host")
	dbName := conf.Section(section, "dbname")
	alias := conf.Section(section, "alias")
	dataSource := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4", username, password, host, dbName)

	trace.Info("正在初始化 UID 数据库连接: alias=%s, dsn=%s", alias, dataSource)

	// 注册 MySQL 驱动
	if err := orm.RegisterDriver("mysql", orm.DRMySQL); err != nil {
		trace.Error("注册 MySQL 驱动失败：%v", err)
		return
	}

	// 注册数据库
	if err := orm.RegisterDataBase(alias, "mysql", dataSource); err != nil {
		trace.Error("注册数据库失败：%v", err)
		return
	}

	// 注册模型
	orm.RegisterModel(new(GoWorkerNode))
	orm.SetMaxOpenConns(alias, 10)
	orm.SetMaxIdleConns(alias, 5)
	dbOrm = orm.NewOrmUsingDB(alias)

	// 初始化唯一 ID 生成器实例
	uniqueIdGen = &UniqueIdGenerator{maxCounts: 1000}
}

// GetUidDBOrm 返回 ORM 实例，用于与 UID 数据库交互
// 该方法提供一个与 UID 数据库交互的全局 ORM 实例
func GetUidDBOrm() orm.Ormer {
	return dbOrm
}

// UniqueIdGenerator 负责生成唯一服务器 ID
type UniqueIdGenerator struct {
	maxCounts int64 // 最大尝试次数，用于重试逻辑
	serId     int64 // 服务器唯一 ID
}

// GetIdGenerator 获取唯一 ID 生成器的单例实例
// 如果实例不存在，将自动初始化。通过此方法调用的实例保证唯一
func GetIdGenerator() *UniqueIdGenerator {
	return uniqueIdGen
}

// GetSvrUniqueId 获取唯一 ID，带有重试逻辑
// 该方法用于在分布式环境中生成唯一 ID，最大尝试次数超过限制时返回错误
func GetSvrUniqueId() (int64, error) {
	var countLoop int64
	serverID, err := obtainUniqueId(countLoop)
	if err != nil {
		trace.Error("尝试获取唯一 ID 失败，次数=%d, 错误=%v", countLoop, err)
	}
	return serverID, err
}

// obtainUniqueId 尝试获取唯一服务器 ID，带有最大重试次数限制
// 此方法用于递增获取服务器 ID，若失败则自动重试
func obtainUniqueId(countLoop int64) (int64, error) {
	for countLoop <= uniqueIdGen.maxCounts {
		serverID, err := queryAndIncrementId()
		if err == nil {
			return serverID, nil
		}
		countLoop++
		trace.Warn("获取唯一 ID 重试中，第 %d 次尝试: %v", countLoop, err)
		time.Sleep(10 * time.Millisecond) // 间隔 10 毫秒进行重试
	}
	return 0, fmt.Errorf("在 %d 次尝试后获取唯一 ID 失败", uniqueIdGen.maxCounts)
}

// queryAndIncrementId 查询当前 ID，递增后保存新的 ID
// 此方法负责在 UID 数据库中查询当前最大 ID，并在生成唯一 ID 后递增保存
func queryAndIncrementId() (int64, error) {
	goWorkerNode := &GoWorkerNode{HostName: tool.GetLocalIp()}
	if err := goWorkerNode.QueryMaxId(); err != nil {
		return 0, err
	}

	goWorkerNode.Id++
	goWorkerNode.CreateTime = time.Now()
	if err := goWorkerNode.Insert(); err != nil {
		return 0, err
	}

	trace.Notice("成功生成唯一 ID: %d", goWorkerNode.Id)
	return goWorkerNode.Id, nil
}

// AsyncSetUniqueId 异步设置服务器 ID，在启动时调用
// 该方法异步执行生成唯一 ID 的过程，并设置服务器 ID
func AsyncSetUniqueId() {
	async.AsyncRunCoroutine(func() {
		if id, err := GetSvrUniqueId(); err == nil {
			uniqueIdGen.serId = id
		} else {
			trace.Error("异步设置服务器 ID 失败：%v", err)
		}
	})
}

// InitUniqueId 初始化服务器 ID，在启动时调用
// 该方法在同步生成唯一 ID 并设置服务器 ID，适合阻塞式调用
func InitUniqueId() int64 {
	var err error
	uniqueIdGen.serId, err = GetSvrUniqueId()
	if err != nil {
		trace.Error("设置服务器 ID 失败：%v", err)
		uniqueIdGen.serId = time.Now().UnixNano()
		trace.Notice("设置服务器 ID: %d", uniqueIdGen.serId)
	}
	return uniqueIdGen.serId
}
