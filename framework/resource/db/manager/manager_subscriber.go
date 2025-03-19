package mgr

import (
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
	fdb "sl.framework.com/resource/db"
	"sl.framework.com/resource/protocol"
	"sl.framework.com/tool"
	"sl.framework.com/trace"
	"strings"
	"sync"
	"time"
)

var (
	subscriber       sync.Map // 存储订阅者信息
	subscriberStatus sync.Map // 存储订阅者在线状态
	vidMap           sync.Map // 存储 Vid 对应的订阅者
	vidGameRoomIdMap sync.Map // 存储 Vid 对应的GameRoomId
	betTimeMap       sync.Map // 存储 Vid 对应的下注时间
)

// 初始化并定时更新订阅者信息
func initAndAutoUpdateSubscribers() {
	syncSubscribers()           // 首次同步订阅者
	go updateSubscriberOnline() // 开启更新订阅者状态的协程

	ticker := time.NewTicker(1 * time.Minute) // 每1分钟更新一次
	for range ticker.C {
		processSubscribers(false) // 定期更新订阅者信息
	}
}

// 同步所有订阅者信息
func syncSubscribers() {
	timeProfiler := tool.NewTimerProfiler("sync subscribers", 500*time.Millisecond)
	defer timeProfiler.Stop(true)
	processSubscribers(true) // 初始化同步
}

// 处理订阅者信息（初始化或更新）
func processSubscribers(isInit bool) {
	subscribers := fdb.GetSubscribers()
	_vidMap := make(map[string][]*fdb.SubscriberInfo)

	for i := range subscribers {
		sb := &subscribers[i]
		processSubscriber(sb, _vidMap, isInit)
	}

	// 如果不是初始化，则更新 Vid map
	if !isInit {
		updateVidMap(_vidMap)
	}
}

// 处理单个订阅者逻辑
func processSubscriber(sb *fdb.SubscriberInfo, vidMap map[string][]*fdb.SubscriberInfo, isInit bool) {
	// 检查订阅者状态，状态无效则删除对应的检查URL
	if sb.Status == protocol.SUB_STATUS_CLOSED {
		url := c.getCheckUrlMap(sb.GameRoomId)
		if url != "" {
			c.delCheckUrlMap(sb.GameRoomId)
			trace.Warning("Subscriber game_room_id: %d vid: %s status: closed, clean it from check map.", sb.GameRoomId, sb.SubscribedVids)
		}
		return
	}

	// 处理 VID 相关的逻辑
	vid := strings.TrimSpace(sb.SubscribedVids)
	lastRound := fdb.GetLastRound(vid)

	// 存储对应关系
	vidGameRoomIdMap.Store(vid, sb.GameRoomId)

	if lastRound != "" {
		newLastRound, _ := fdb.GetNewGMCode(lastRound)
		if newLastRound != "" {
			trace.Notice("old round code: %s(len=%d) transfer to subscription round: %s(len=%d) success", lastRound, len(lastRound), newLastRound, len(newLastRound))
			lastRound = newLastRound
		} else {
			trace.Info("[%s] is Non-subscription round", lastRound)
		}
		updateCheckMapAndStatus(sb.GameRoomId, sb.Endpoint, lastRound)
	}

	// 存储订阅者信息
	subscriber.Store(sb.GameRoomId, sb)

	// 如果是初始化，添加订阅者到 Vid map，否则批量处理更新
	if isInit {
		addVidMap(vid, sb)
	} else {
		vidMap[vid] = append(vidMap[vid], sb)
	}
}

// 批量更新 Vid map
func updateVidMap(vidMap map[string][]*fdb.SubscriberInfo) {
	for vid, subs := range vidMap {
		swapVidMap(vid, subs)
	}
}

// swapVidMap 替换 Vid 对应的订阅者信息
func swapVidMap(vid string, sb []*fdb.SubscriberInfo) {
	vidMap.Swap(vid, sb)
}

// addVidMap 向 Vid 中添加新的订阅者信息
func addVidMap(vid string, sb *fdb.SubscriberInfo) {
	value, ok := vidMap.Load(vid)
	var infos []*fdb.SubscriberInfo
	if ok { // 如果已有条目，使用现有的切片，否则初始化新切片
		infos = value.([]*fdb.SubscriberInfo)
	}

	found := false // 查找是否已存在相同 ID 的 SubscriberInfo
	for i, info := range infos {
		if info.GameRoomId == sb.GameRoomId {
			infos[i] = sb // 存在则更新
			found = true
			break
		}
	}
	if !found { // 如果不存在相同 ID 的 SubscriberInfo，添加新的
		infos = append(infos, sb)
	}
	vidMap.Store(vid, infos) // 存储更新后的 SubscriberInfo 列表
}

// 定时更新订阅者的在线状态
func updateSubscriberOnline() {
	checkSubscriberOnline()                   // 首次检查订阅者在线状态
	ticker := time.NewTicker(1 * time.Minute) // 每1分钟更新一次
	for range ticker.C {
		checkSubscriberOnline()
	}
}

// GetSubscribersByVid 获取某个 Vid 对应的订阅者列表
func GetSubscribersByVid(vid string) []*fdb.SubscriberInfo {
	value, ok := vidMap.Load(vid)
	var infos []*fdb.SubscriberInfo
	if ok {
		infos = value.([]*fdb.SubscriberInfo)
		trace.Debug("vid=%s subscribers count=%d", vid, len(infos))
	}
	return infos
}

// isOnline 检查订阅者是否在线
func isOnline(id int64) bool {
	isOnlineAny, ok := subscriberStatus.Load(id)
	if ok {
		return isOnlineAny.(bool)
	}
	return false
}

// isWorking 检查网络是否正常 TODO进一步完善需要
// 参数：
//   - vid：视频id
func isWorking(vid string) {
	req, resp := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)
	// 尝试其他的请求是否也有异常
	data := new(RequestData)
	data.Endpoint = getCheckUrlByVid(vid) // 获取检查地址
	data.Method = `GET`
	data.TraceID = strings.ReplaceAll(uuid.New().String(), "-", "")
	setupRequest(req, data) // 配置请求

	if err := fasthttp.Do(req, resp); err != nil || resp.StatusCode() != fasthttp.StatusOK {
		trace.Error("Double check internet is not working now. url:[%s] code=%d, err: %v", data.Endpoint, resp.StatusCode(), err)
	} else {
		trace.Info("Double check internet is working now. code=%d, err: %v", resp.StatusCode(), err)
	}
}
