/**
 * @Author: M
 * @Date: 2024/7/30
 * @Description: 该文件实现了从数据库中查询 `api_info` 表数据的功能。主要功能包括：
 *               1. 通过 Beego ORM 进行数据库查询，从 `api_info` 表中获取 API 信息。
 *               2. 使用 `TimerProfiler` 记录查询的执行时间，便于性能监控和优化。
 *               3. 在查询出现错误时记录错误日志，帮助追踪数据库查询失败的原因。
 *
 * @Dependencies:
 *               - `github.com/beego/beego/v2/client/orm`: Beego ORM，用于处理数据库查询。
 *               - `sl.framework.com/trace`: 用于日志记录，记录错误和信息日志。
 *               - `sl.framework.com/tool`: 用于性能监控，使用 `TimerProfiler` 记录数据库查询耗时。
 *
 * @Usage:
 *               1. 使用 `GetApiInfos` 函数从数据库中查询所有 API 信息。
 *               2. 查询结果会以 `[]ApiInfo` 形式返回，如果查询失败，错误会记录到日志中。
 */

package db

import (
	"github.com/beego/beego/v2/client/orm"
	"sl.framework.com/tool"
	"sl.framework.com/trace"
	"time"
)

// ApiInfo 数据库模型定义
type ApiInfo struct {
	Id            int64  `orm:"pk"`                       // 自增主键，唯一标识每一条记录
	RequestUrl    string `orm:"size(255);default('')"`    // 请求的 URL，将与订阅者表中的 endpoint 拼接成最终请求连接
	RequestMethod string `orm:"size(10);default('POST')"` // 请求方式，例如 POST、GET 等
	Description   string `orm:"size(255);null"`           // API 的附加说明或备注
}

// GetApiInfos 查询 `api_info` 表中的所有 API 信息。
// 返回：返回一个包含 `ApiInfo` 结构体的切片，表示从数据库中查询到的所有 API 信息。
//
// 功能：
//   - 使用 Beego ORM 查询数据库中的 `api_info` 表，获取所有 API 信息。
//   - 使用 `TimerProfiler` 记录查询操作的执行时间，并在查询执行完成时停止计时。
//   - 如果查询过程中出现错误，记录错误日志。
func GetApiInfos() []ApiInfo {
	timeProfiler := tool.NewTimerProfiler("get api info", 500*time.Millisecond)
	defer timeProfiler.Stop(true)
	o := orm.NewOrm()
	var subscribers []ApiInfo
	_, err := o.QueryTable("api_info").All(&subscribers)
	if err != nil {
		trace.Error("Failed to load api info from database: %v", err)
	}
	return subscribers
}
