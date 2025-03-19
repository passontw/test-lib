package mgr

import (
	"fmt"
	"sl.framework.com/trace"
	"sync"
)

var checkOnlineStateId = 2

// cr 定义了一个含有读写锁的结构体，用于同步处理检查URL
type cr struct {
	sync.RWMutex
	checkMapUrl map[int64]string // 存储订阅者的检查URL
}

var c *cr

// 初始化 URL 管理结构
func init() {
	c = &cr{
		checkMapUrl: make(map[int64]string),
	}
}

// 获取地址
func (r *cr) getCheckUrlMap(gameRoomId int64) string {
	r.RLock()
	defer r.RUnlock()
	return r.checkMapUrl[gameRoomId]
}

// getCheckUrlByVid 根据 vid（视频或房间标识符）从映射中查找对应房间的CheckUrl。
// 如果未找到为 ""。
func getCheckUrlByVid(vid string) string {
	// 从 vidGameRoomIdMap 映射中查找 vid 对应的房间 ID
	if value, ok := vidGameRoomIdMap.Load(vid); ok {
		return c.checkMapUrl[value.(int64)]
	}
	return ""
}

// 添加游戏房间的检查URL
func (r *cr) addCheckUrlMap(gameRoomId int64, checkUrl string) {
	r.Lock()
	defer r.Unlock()
	r.checkMapUrl[gameRoomId] = checkUrl
}

// 删除游戏房间的检查URL
func (r *cr) delCheckUrlMap(gameRoomId int64) {
	r.Lock()
	defer r.Unlock()
	delete(r.checkMapUrl, gameRoomId)
}

// 更新检查URL和订阅者在线状态
func updateCheckMapAndStatus(gameRoomId int64, endpoint, gmcode string) {
	checkUrl, _ := EndpointAndMethod(int64(checkOnlineStateId), endpoint)
	checkUrl = fmt.Sprintf(checkUrl, gameRoomId, gmcode)
	c.addCheckUrlMap(gameRoomId, checkUrl)
	trace.Info("online state check info: game_room_id: %d, checkUrl: %s set to online", gameRoomId, checkUrl)
	subscriberStatus.Store(gameRoomId, true) // 更新订阅者在线状态为在线
}
