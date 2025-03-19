package base

import (
	"sl.framework.com/async"
	"sl.framework.com/trace"
	"sync"
	"time"
)

var (
	cacheManagerOnce     sync.Once
	cacheManagerInstance *CacheManagerLoopTask
)

// GetCacheManager 缓存管理
func GetCacheManager() *CacheManagerLoopTask {
	cacheManagerOnce.Do(func() {
		cacheManagerInstance = &CacheManagerLoopTask{
			CLoopTask: &CLoopTask{
				createTime: time.Now(),
				interval:   time.Duration(5) * time.Second,
			},
			roomDetailedInfoMap: make(map[int64]time.Time, 128),
		}
	})

	return cacheManagerInstance
}

var _ ILoopTask = (*CacheManagerLoopTask)(nil)

// CacheManagerLoopTask 内存管理任务
type CacheManagerLoopTask struct {
	*CLoopTask

	roomDetailedInfoMutex sync.Mutex
	roomDetailedInfoMap   map[int64]time.Time //key = GameRoundId
}

// RunLoopTask 启动缓存任务
func (c *CacheManagerLoopTask) RunLoopTask() {
	fn := func() {
		trace.Info("CacheManagerLoopTask RunLoopTask start running, heartbeat interval=%v(s)", c.interval)
		c.ticker = time.NewTicker(c.interval)
		defer c.ticker.Stop()
		for {
			select {
			case <-c.ticker.C:
				c.cacheCleanUp()
			}
		}
	}

	//启动循环任务
	c.once.Do(func() {
		async.AsyncRunCoroutine(fn)
	})
}

// 清理缓存任务
func (c *CacheManagerLoopTask) cacheCleanUp() {
	c.roomDetailedInfoMutex.Lock()
	defer c.roomDetailedInfoMutex.Unlock()

	toBeDeletedSet := make([]int64, 0, 8) //cap默认设置为8
	for k, v := range c.roomDetailedInfoMap {
		if time.Since(v) > time.Duration(1)*time.Minute {
			toBeDeletedSet = append(toBeDeletedSet, k)
		}
	}

	for _, v := range toBeDeletedSet {
		delete(c.roomDetailedInfoMap, v)
	}
}

// AddRoomDetailedInfoDoor 增加
func (c *CacheManagerLoopTask) AddRoomDetailedInfoDoor(gameRoundId int64) {
	c.roomDetailedInfoMutex.Lock()
	defer c.roomDetailedInfoMutex.Unlock()
	c.roomDetailedInfoMap[gameRoundId] = time.Now()
}

// CheckRoomDetailedInfoDoor 检查
func (c *CacheManagerLoopTask) CheckRoomDetailedInfoDoor(gameRoundId int64) bool {
	c.roomDetailedInfoMutex.Lock()
	defer c.roomDetailedInfoMutex.Unlock()

	bIsExist := false
	if _, ok := c.roomDetailedInfoMap[gameRoundId]; ok {
		bIsExist = true
		trace.Notice("CheckRoomDetailedInfoDoor already send room detailed info request to platform, gameRoundId=%v", gameRoundId)
	}
	return bIsExist
}
