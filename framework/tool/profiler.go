package tool

import (
	"sl.framework.com/async"
	"sl.framework.com/trace"
	"time"
)

const (
	profilerDefaultWarningDoor = time.Duration(100) * time.Millisecond //执行时间超过默认100ms则以告警形式打印信息
)

type TimeProfiler struct {
	startTime  time.Time
	noticeSpan time.Duration
	strLog     string
}

// NewTimerProfiler 创建一个新的 TimeProfiler 实例
func NewTimerProfiler(log string, noticeSpan time.Duration) *TimeProfiler {
	stack := async.GetCallStack(1, 1)
	partPath := ExtractPathFromIndex(stack[0][1], -3)
	trace.Info("%s <= (%s:%s %s)", log, partPath, stack[0][2], stack[0][0])
	return &TimeProfiler{
		startTime:  time.Now(),
		noticeSpan: noticeSpan,
		strLog:     log,
	}
}

// Stop 停止计时并记录日志
func (p *TimeProfiler) Stop(notice bool) {
	span := time.Since(p.startTime)
	if notice && span >= p.noticeSpan {
		trace.Notice("%s, span=%v, performance warning", p.strLog, span)
	} else {
		trace.Debug("%s, span=%v", p.strLog, span)
	}
}

// TimeoutWarning 用于对整个函数计时,使用方法 defer TimeoutWarning("{msg}")
func TimeoutWarning(fnName string, warningDoor time.Duration) func() {
	begin := time.Now()
	return func() {
		span := time.Since(begin)
		iWarningDoor := profilerDefaultWarningDoor
		if warningDoor > profilerDefaultWarningDoor {
			iWarningDoor = warningDoor
		}
		if span > iWarningDoor {
			trace.Notice("%v cost span=%v", fnName, span)
		}
	}
}

/**
 * NewWatcher
 * 创建一个代码测量对象
 *
 * @param msg string - 对象内Id 用于打印
 * @return *Watcher - 代码测量对象指针
 */

func NewWatcher(id string) *Watcher {
	return &Watcher{
		msg:         id,
		warningDoor: profilerDefaultWarningDoor,
		start:       time.Now(),
		statistics:  make(map[string]time.Duration),
	}
}

/*
	Watcher 代码计时工具,可以像秒表一样使用
	也可以在最后一起打印出每个阶段的执行耗时
*/

type Watcher struct {
	msg         string                   //打印耗时信息的时候打印信息,每次New对象都要传入
	warningDoor time.Duration            //打印告警信息阈值
	start       time.Time                //对象创建时候初始化 用作计算耗时的起点
	statistics  map[string]time.Duration //用
}

/**
 * FunctionName
 * 函数功能描述
 *
 * @param PARAM - 参数说明
 * @return RETURN - 返回值说明
 */

func (p *Watcher) Start(msg string) {
	p.msg = msg
	p.start = time.Now()
}

/**
 * Stop
 * 计时并打印耗时信息 不放入统计信息中
 * 调用的时候会打印代码执行耗时 以对象创建的时间点为起点
 *
 */

func (p *Watcher) Stop() {
	timeElapse := time.Since(p.start)
	if timeElapse >= p.warningDoor {
		trace.SetLogFuncCallDepth(4)
		trace.Notice("%v cost %v", p.msg, timeElapse)
		trace.SetLogFuncCallDepth(3)
	} else {
		trace.SetLogFuncCallDepth(4)
		trace.Info("%v cost %v", p.msg, timeElapse)
		trace.SetLogFuncCallDepth(3)
	}
}

/**
 * StopOnStat
 * 计时并打印耗时信息
 * 并将该地方耗时记录到统计信息中
 *
 * @param id string - 信息Id一个对象每次传入的值都不能一样
 * @return
 */

func (p *Watcher) StopOnStat(id string) {
	p.statistics[id] = time.Since(p.start)
	p.Stop()
}

/**
 * StopOnEnd
 * 计时并打印耗时信息
 * 并将该地方耗时记录到统计信息中 并打印出所有的统计信息
 *
 * @param id string - 信息Id一个对象每次传入的值都不能一样
 * @return
 */

func (p *Watcher) StopOnEnd(id string) {
	p.statistics[id] = time.Since(p.start)

	trace.SetLogFuncCallDepth(4)
	trace.Info("---------------------------------------------------")
	for k, v := range p.statistics {
		trace.Info("%v cost %v", k, v)
	}
	trace.Info("---------------------------------------------------")
	trace.SetLogFuncCallDepth(3)
}
