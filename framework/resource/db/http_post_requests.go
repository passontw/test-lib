package db

import (
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	snowflaker "sl.framework.com/resource/snow_flake_id"
	"sl.framework.com/tool"
	"sl.framework.com/trace"
	"strings"
	"time"
)

// HttpPostRequests 数据库模型定义
type HttpPostRequests struct {
	Id            int64     `orm:"pk"`                   // 自增主键，唯一标识每一条记录
	TraceId       string    `orm:"size(36)"`             // 请求的唯一标识符，用于追踪重试的请求
	Method        string    `orm:"size(10)"`             // 请求方法
	Gmcode        string    `orm:"size(16)"`             // 对局的唯一标识符，用于关联请求和具体的对局
	EndpointId    int64     `orm:""`                     // 请求的目标端点 URL ID
	Endpoint      string    `orm:"size(255)"`            // 请求的目标端点 URL
	RequestBody   string    `orm:"size(1024)"`           // 请求的参数内容，以 JSON 格式存储
	ResponseCode  int       `orm:"null"`                 // 请求的响应代码，例如 200 表示成功
	ResponseBody  string    `orm:"size(1024);null"`      // 响应的内容
	RequestTime   time.Time `orm:"type(timestamp)"`      // 请求的时间
	ResponseTime  time.Time `orm:"null;type(timestamp)"` // 响应的时间
	Status        string    `orm:"size(32)"`             // 请求的状态，枚举值包括 pending, success, failed, retrying
	RetryCount    int       `orm:"default(0)"`           // 重试次数，默认为 0
	LastRetryTime time.Time `orm:"null;type(timestamp)"` // 最后一次重试的时间
}

func GetPosts(max int) []HttpPostRequests {
	timeProfiler := tool.NewTimerProfiler("get post info", 500*time.Millisecond)
	defer timeProfiler.Stop(true)
	o := orm.NewOrm()
	var subscribers []HttpPostRequests
	_, err := o.Raw("select * from http_post_requests where status <> 'success' and retry_count <= ?").SetArgs(max).QueryRows(&subscribers)
	if err != nil {
		trace.Error("Failed to load post info from database: %v", err)
	}
	return subscribers
}

func (r *HttpPostRequests) Insert() error {
	timeProfiler := tool.NewTimerProfiler(fmt.Sprintf("notice sub gmcode=%s", r.Gmcode), 500*time.Millisecond)
	defer timeProfiler.Stop(true)
	var err error

	tx := orm.NewOrm()
	if err != nil {
		return fmt.Errorf("notice sub transactions start failed: %v", err)
	}

	// 执行插入操作，将新的游戏回合信息插入到 http_post_requests 表中
	res, err := tx.Raw("INSERT INTO http_post_requests (id, gmcode, endpoint_id, method, endpoint, request_body, request_time, trace_id, status) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		snowflaker.UniqueId(), r.Gmcode, r.EndpointId, r.Method, r.Endpoint, r.RequestBody, r.RequestTime, r.TraceId, r.Status).Exec()
	if err != nil { // 插入操作失败
		return fmt.Errorf("insert post failed: %v", err)
	}
	// 检查更新操作是否成功
	if i, err := res.RowsAffected(); err != nil || i == 0 {
		// 更新操作未影响任何行
		return fmt.Errorf("insert post failed, the number of rows affected is %d. err is: %v", i, err)
	}
	return nil
}

func (r *HttpPostRequests) Update() error {
	timeProfiler := tool.NewTimerProfiler(fmt.Sprintf("update post gmcode=%s", r.Gmcode), 500*time.Millisecond)
	defer timeProfiler.Stop(true)
	var err error

	tx := orm.NewOrm()
	// 执行插入操作，将新的游戏回合信息插入到 http_post_requests 表中
	res, err := tx.Raw("update http_post_requests set response_body=?, response_code=?, request_time=?, response_time=?, status=? where trace_id=?",
		r.ResponseBody, r.ResponseCode, r.RequestTime, r.ResponseTime, r.Status, r.TraceId).Exec()
	if err != nil { // 插入操作失败
		return fmt.Errorf("update post failed: %v", err)
	}
	// 检查更新操作是否成功
	if i, err := res.RowsAffected(); err != nil || i == 0 {
		// 更新操作未影响任何行
		return fmt.Errorf("update post failed, the number of rows affected is %d. err is: %v", i, err)
	}
	return nil
}

func (r *HttpPostRequests) RecordFailure() error {
	timeProfiler := tool.NewTimerProfiler(fmt.Sprintf("update post gmcode=%s", r.Gmcode), 500*time.Millisecond)
	defer timeProfiler.Stop(true)

	var err error
	var params []interface{}
	var setClauses []string

	// 动态构建 SQL 语句的 SET 子句
	if r.ResponseBody != "" {
		setClauses = append(setClauses, "response_body=?")
		params = append(params, r.ResponseBody)
	}

	if r.ResponseCode > 0 {
		setClauses = append(setClauses, "response_code=?")
		params = append(params, r.ResponseCode)
	}

	if !r.RequestTime.IsZero() {
		setClauses = append(setClauses, "request_time=?")
		params = append(params, r.RequestTime)
	}

	if !r.ResponseTime.IsZero() {
		setClauses = append(setClauses, "response_time=?")
		params = append(params, r.ResponseTime)
	}

	if r.Status != "" {
		setClauses = append(setClauses, "status=?")
		params = append(params, r.Status)
		if r.Status != "pending" {
			// retry_count 需要递增
			setClauses = append(setClauses, "retry_count=retry_count+1")
		}
	}

	// 构建最终的 SQL 语句
	sql := fmt.Sprintf("UPDATE http_post_requests SET %s WHERE trace_id=?", strings.Join(setClauses, ", "))
	params = append(params, r.TraceId)

	// 执行 SQL 语句
	tx := orm.NewOrm()
	res, err := tx.Raw(sql, params...).Exec()
	if err != nil {
		trace.Error("update post failed: %v", err)
		return fmt.Errorf("update post failed: %v", err)
	}

	// 检查更新操作是否成功
	if i, errr := res.RowsAffected(); errr != nil || i == 0 {
		trace.Error("update post failed, the number of rows affected is %d. err is: %v", i, errr)
		return fmt.Errorf("update post failed, the number of rows affected is %d. err is: %v", i, errr)
	}

	return nil
}
