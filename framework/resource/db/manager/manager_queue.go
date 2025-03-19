package mgr

import (
	"math"
	"math/rand"
	"sl.framework.com/base"
	fdb "sl.framework.com/resource/db"
	"sl.framework.com/tool"
	"sl.framework.com/trace"
	"strings"
	"time"
)

var (
	incomingQueue    = make(chan *RequestData, 10000) // 输入队列
	processingQueue  = make(chan *RequestData, 10000) // 处理队列
	retryingQueue    = make(chan *RequestData, 10000) // 重试队列
	dnsNotFoundQueue = make(chan *RequestData, 10000) // 重试队列

	retryingQueueEnableSillTimes = 7000 // 开启sill延迟入队的次数 重试队列len大于这个次数后 新的重试将会延迟写入重试队列
	maxTriedTimes                = 1000 // 最大重试次数 这个次数后会认为失败
	notSuchHost                  = "no such host"
)

func Push2Queue(data *RequestData) {
	retryingQueue <- data
}

func initAndUpdatePostInfo() {
	syncPost()
	ticker := time.NewTicker(3 * time.Minute) // 每3分钟更新一次
	for range ticker.C {
		syncPost() // 同步发布信息
	}
}

// 同步发布信息并将其放入重试队列
func syncPost() {
	timeProfiler := tool.NewTimerProfiler("sync post info", 500*time.Millisecond)
	defer timeProfiler.Stop(true)

	// 每次重新生成随机数生成器
	rand.NewSource(time.Now().UnixNano())

	subscribers := fdb.GetPosts(maxTriedTimes)
	trace.Info("[start] sync post info: %d , retryingQueue size: %d", len(subscribers), len(retryingQueue))
	for _, sb := range subscribers {
		if strings.HasSuffix(sb.ResponseBody, notSuchHost) {
			insertDNSNotFoundQueue(&sb)
			continue
		}
		if len(retryingQueue) < retryingQueueEnableSillTimes {
			insertRetryingQueue(&sb)
			continue
		}
		// 计算基于指数退避和抖动的延迟时间
		jitter := getExponentialBackoffWithJitter(sb.RetryCount, 1*time.Second, 2*time.Hour)
		trace.Warning("post info: %s , retryTimes: %d , retryingQueue size: %d , will try again in %d seconds", sb.TraceId, sb.RetryCount, jitter, len(retryingQueue))
		base.RunAfter(time.Duration(jitter)*time.Second, func() {
			insertRetryingQueue(&sb)
		})
	}
	trace.Info("[end  ] sync post info: %d , retryingQueue size: %d", len(subscribers), len(retryingQueue))
}

func newRequestData(sb *fdb.HttpPostRequests) *RequestData {
	d := &RequestData{
		// 构建请求数据并放入重试队列
		EndpointId:  sb.EndpointId,
		Endpoint:    sb.Endpoint,
		Method:      sb.Method,
		Body:        sb.RequestBody,
		TraceID:     sb.TraceId,
		GMCode:      sb.Gmcode,
		TryMaxTimes: int64(sb.RetryCount),
	}
	return d
}

func insertRetryingQueue(sb *fdb.HttpPostRequests) {
	d := newRequestData(sb)
	retryingQueue <- d
}

func insertDNSNotFoundQueue(sb *fdb.HttpPostRequests) {
	d := newRequestData(sb)
	d.isDNSNotFound = true
	dnsNotFoundQueue <- d
}

// getExponentialBackoffWithJitter 计算基于指数退避和抖动的延迟时间。
// 适用于防止雪崩效应，避免大量请求同时重试。
// 根据当前的重试次数，基础延迟时间以及随机抖动，生成一个下一次操作的延迟时间。
//
// 参数:
//
//	retryCount (int): 当前的重试次数，从1开始。
//	baseDelay (time.Duration): 基础延迟时间，通常为一个固定值（如1秒）。
//	maxDelay (time.Duration): 最打延迟事件，2小时。
//
// 返回:
//
//	int: 计算出的下一次操作的延迟时间（以整数秒为单位）。
//
// 算法说明:
//  1. 使用指数退避公式计算基础延迟：baseDelay * 2 ^ (retryCount - 1)。
//  2. 应用抖动，在基础延迟的 50% 到 100% 之间随机生成延迟时间。
//  3. 最终返回四舍五入后的整数秒数。
//
// 抖动的指数退避算法函数，增加最大延迟限制
func getExponentialBackoffWithJitter(retryCount int, baseDelay time.Duration, maxDelay time.Duration) int {
	// 指数退避计算基础延迟时间
	base := baseDelay.Seconds() * math.Pow(2, float64(retryCount-1))

	// 设置最大延迟，确保延迟不会超过 maxDelay
	if base > maxDelay.Seconds() {
		base = maxDelay.Seconds()
	}

	// 生成抖动：在 50% 到 100% 的范围内随机化延迟
	jitter := base/2 + rand.Float64()*base/2

	// 确保 jitter 是非负数
	if jitter < 0 {
		jitter = 0
	}

	// 返回四舍五入后的整数秒数
	return int(math.Round(jitter))
}
