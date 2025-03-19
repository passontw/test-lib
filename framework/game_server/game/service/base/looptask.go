package base

import (
	"sync"
	"time"
)

// ILoopTask 定时接口
type ILoopTask interface {
	RunLoopTask()
}

// CLoopTask 定时任务基类
type CLoopTask struct {
	ticker     *time.Ticker
	interval   time.Duration //定时任务时间间隔,单位ms
	createTime time.Time     //task创建时间
	once       sync.Once
}
