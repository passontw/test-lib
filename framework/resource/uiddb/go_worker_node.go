package uiddb

import (
	"errors"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"sl.framework.com/trace"
	"time"
)

const UidTableName = "go_worker_node"

// UniqueIdData 请求UniqueId的返回数据结构
type UniqueIdData struct {
	UniqueId int64  `json:"uniqueId"`
	Code     int    `json:"code"` //0:正确,其他错误
	Msg      string `json:"msg"`  //code非0时的错误消息
}

type GoWorkerNode struct {
	Id         int64     `orm:"pk;description(id)"`
	HostName   string    `orm:"size(64);description(主机Ip)"`
	CreateTime time.Time `orm:"type(timestamp);description(创建时间)"`
}

// Insert 增
func (g *GoWorkerNode) Insert() error {
	o := GetUidDBOrm()

	effected, err := o.Insert(g)
	if nil != err {
		trace.Error("GoWorkerNode Insert failed, effected line=%v, error=%v, data=%+v", effected, err.Error(), g)
		return fmt.Errorf("GoWorkerNode Insert failed, effected line=%v, error=%+v", effected, err)
	}
	trace.Notice("GoWorkerNode Insert success, effected line=%v, data=%+v", effected, g)

	return nil
}

// QueryMaxId 查询Id最大的数据
func (g *GoWorkerNode) QueryMaxId() error {
	o := GetUidDBOrm()

	err := o.QueryTable(UidTableName).OrderBy("-id").Limit(1).One(g)
	if errors.Is(err, orm.ErrNoRows) {
		g.Id = 1 //表中没有数据则初始化为1
		trace.Notice("GoWorkerNode QueryMaxId no data in table, get default id=1")
	} else if err != nil {
		return fmt.Errorf("GoWorkerNode QueryMaxId query error: %v", err.Error())
	}

	trace.Notice("GoWorkerNode QueryMaxId query success, data=%+v", g)
	return nil
}

// Delete 删
func (g *GoWorkerNode) Delete() {
	o := GetUidDBOrm()

	effected, err := o.Delete(g)
	if nil != err {
		trace.Error("GoWorkerNode Delete failed, effected line=%v, data=%+v", effected, g)
	}
	trace.Info("GoWorkerNode Delete success, effected line=%v, data=%+v", effected, g)
}

// Update 更新cols字段 改
func (g *GoWorkerNode) Update(cols ...string) {
	o := GetUidDBOrm()

	effected, err := o.Update(g, cols...)
	if nil != err {
		trace.Error("GoWorkerNode Update failed, effected line=%v, data=%+v", effected, g)
	}
	trace.Info("GoWorkerNode Update success, effected line=%v, data=%+v", effected, g)
}

// Query 查询cols字段 查
func (g *GoWorkerNode) Query(cols ...string) {
	o := GetUidDBOrm()

	if err := o.Read(g, cols...); nil != err {
		trace.Error("GoWorkerNode Query failed, data=%+v", g)
		switch {
		case errors.Is(err, orm.ErrNoRows):
			trace.Error("Query GoWorkerNode error no rows")
		case errors.Is(err, orm.ErrMissPK):
			trace.Error("Query GoWorkerNode error no primary key")
		default:
			trace.Error("GoWorkerNode Query failed, unknown error=%v", err.Error())
		}
		return
	}
	trace.Info("GoWorkerNode Query, data=%+v", g)
}

// QueryRaw 使用SQL查询
func (g *GoWorkerNode) QueryRaw() {
	var nodes []GoWorkerNode
	o := GetUidDBOrm()

	rs := o.Raw("select * from go_worker_node where id > ? order by id desc", 10)
	num, err := rs.QueryRows(&nodes)
	if err != nil {
		trace.Error("QueryRaw failed, err=%v", err.Error())
		return
	}
	trace.Info("QueryRaw success, num=%v", num)
	for index, node := range nodes {
		trace.Info("QueryRaw success index=%v, node=%+v", index, node)
	}
}

// QueryConditon 使用SQL查询
func (g *GoWorkerNode) QueryConditon() {
	var (
		nodes []GoWorkerNode
		node  GoWorkerNode
	)
	o := GetUidDBOrm()

	qs := o.QueryTable(node).Filter("id__lt", 19) //.Limit(2, 6)
	num, err := qs.All(&nodes)
	if nil != err {
		trace.Error("GoWorkerNode QueryConditon failed, err=%v", err.Error())
		return
	}
	trace.Info("GoWorkerNode QueryConditon success, num=%v", num)
	for index, node := range nodes {
		trace.Info("GoWorkerNode QueryConditon success, index=%v, node=%v", index, node)
	}
}
