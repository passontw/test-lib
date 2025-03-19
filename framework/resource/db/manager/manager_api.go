package mgr

import (
	"fmt"
	"sl.framework.com/resource/db"
	"sl.framework.com/tool"
	"sl.framework.com/trace"
	"sync"
	"time"
)

var apiData sync.Map // 存储 API 信息

func Init() {
	go initAndUpdateAPIInfo()
	go initAndUpdatePostInfo()
	// 启动队列处理器
	go dispatcherQueueData(incomingQueue, processingQueue)
	go dispatcherQueueData(processingQueue, retryingQueue)
	go dispatcherQueueData(retryingQueue, nil)
	// no such host 异常单独处理 避免占用重试队列资源
	go dispatcherQueueData(dnsNotFoundQueue, nil)
}

// 初始化并定时更新 API 信息
func initAndUpdateAPIInfo() {
	syncAPIInfos()                            // 同步API信息
	go initAndAutoUpdateSubscribers()         // 自动更新订阅者信息
	ticker := time.NewTicker(5 * time.Minute) // 每5分钟更新一次
	for range ticker.C {
		syncAPIInfos() // 同步API信息
	}
}

// 同步 API 信息
func syncAPIInfos() {
	timeProfiler := tool.NewTimerProfiler("sync api info", 500*time.Millisecond)
	defer timeProfiler.Stop(true)
	subscribers := db.GetApiInfos()
	for _, sb := range subscribers {
		x := sb
		apiData.Store(sb.Id, &x)
		trace.Info("Updated api info: %d requestUrl: %s requestMethod: %s description: %s",
			x.Id, x.RequestUrl, x.RequestMethod, x.Description)
	}
}

// 返回指定 API ID 对应的请求 URL 和方法
func EndpointAndMethod(id int64, end string) (string, string) {
	if value, ok := apiData.Load(id); ok {
		info := value.(*db.ApiInfo)
		return fmt.Sprintf("%s%s", end, info.RequestUrl), info.RequestMethod
	}
	return "", ""
}

// BetTime 根据 vid（视频或房间标识符）从映射中查找对应房间的下注时间。
// 如果未找到或时间为 0，则返回默认值 def。
func BetTime(vid string, def int64) int64 {
	// 从 vidGameRoomIdMap 映射中查找 vid 对应的房间 ID
	if value, ok := vidGameRoomIdMap.Load(vid); ok {
		// 根据房间 ID 从 betTimeMap 映射中查找对应的下注时间
		ivalue := value.(int64)
		if v, ok := betTimeMap.Load(ivalue); ok {
			// 尝试将加载到的值断言为 int64 类型，并确保它不为 0
			info, ok := v.(int64)
			if ok && info != 0 {
				return info // 如果找到有效的下注时间，直接返回
			}
		}
	}
	// 如果未找到房间 ID 或下注时间，则返回默认值 def
	return def
}
