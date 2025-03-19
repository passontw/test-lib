package base

import (
	"encoding/json"
	"sl.framework.com/async"
	"sl.framework.com/game_server/conf"
	"sl.framework.com/game_server/game/dao/redisdb"
	types "sl.framework.com/game_server/game/service/type"
	"sl.framework.com/game_server/redis/rediskey"
	"sl.framework.com/trace"
	"strconv"
	"sync"
	"time"
)

var (
	heartbeatInstanceOnce       sync.Once
	heartbeatTickerTaskInstance *HeartbeatLoopTask
)

// GetHeartbeatInstance 获取心跳单例
func GetHeartbeatInstance() *HeartbeatLoopTask {
	heartbeatInstanceOnce.Do(func() {
		heartbeatTickerTaskInstance = &HeartbeatLoopTask{
			CLoopTask: &CLoopTask{
				createTime: time.Now(),
				interval:   time.Duration(conf.GetHeartbeatInterval()) * time.Second,
			},
			clusterNum: 3, //默认值设置为三个
			serverId:   strconv.FormatInt(conf.GetServerId(), 10),
		}
	})

	return heartbeatTickerTaskInstance
}

var _ ILoopTask = (*HeartbeatLoopTask)(nil)

// HeartbeatLoopTask 心跳包循环任务
type HeartbeatLoopTask struct {
	*CLoopTask

	clusterNumMutex sync.Mutex
	clusterNum      int

	serverId string
}

// RunLoopTask 启动loop循环执行逻辑
func (h *HeartbeatLoopTask) RunLoopTask() {
	fn := func() {
		trace.Info("HeartbeatLoopTask RunLoopTask start running, heartbeat interval=%v(s)", h.interval)
		h.ticker = time.NewTicker(h.interval)
		defer h.ticker.Stop()
		for {
			select {
			case <-h.ticker.C:
				h.heartbeatSend()
			}
		}
	}

	h.once.Do(func() {
		async.AsyncRunCoroutine(fn)
	})
}

// 向redis发送心跳包
func (h *HeartbeatLoopTask) heartbeatSend() {
	var (
		err            error
		val            string
		jsonData       []byte
		redisInfo      = rediskey.GetHeartbeatStatusRedisInfo()
		redisLockInfo  = rediskey.GetHeartbeatStatusLockRedisInfo()
		serverStatus   = make(map[string]*types.HeartbeatStatus, 8)
		serverToDelete = make(map[string]*types.HeartbeatStatus, 8)
	)

	if !redisdb.Lock(redisLockInfo) {
		trace.Error("HeartbeatLoopTask heartbeatSend lock failed, server id=%v", h.serverId)
		return
	}
	defer redisdb.Unlock(redisLockInfo)

	//获取redis信息并根据时间做过滤
	if val, err = redisdb.Get(redisInfo.Key); nil != err {
		trace.Error("HeartbeatLoopTask heartbeatSend redis get failed, server id=%v, key=%v, err=%v",
			h.serverId, redisInfo.Key)
		return
	}
	if len(val) > 0 {
		if err = json.Unmarshal([]byte(val), &serverStatus); nil != err {
			trace.Error("HeartbeatLoopTask heartbeatSend json unmarshal failed, server id=%v, key=%v, val=%v, err=%v",
				h.serverId, redisInfo.Key, val, err.Error())
			return
		}
		for _, v := range serverStatus {
			if h.serverId == v.ServerId {
				v.UpdateTime = time.Now()
			}
			if time.Since(v.UpdateTime) > time.Duration(conf.GetHeartbeatExpired()) {
				serverToDelete[v.ServerId] = v
			}
		}
		for k, _ := range serverToDelete {
			delete(serverStatus, k)
		}
	} else {
		heartbeatStatus := &types.HeartbeatStatus{
			ServerId:   h.serverId,
			UpdateTime: time.Now(),
		}
		serverStatus[h.serverId] = heartbeatStatus
	}

	h.clusterNumMutex.Lock()
	h.clusterNum = len(serverStatus)
	h.clusterNumMutex.Unlock()

	//更新redis缓存
	if jsonData, err = json.Marshal(serverStatus); nil != err {
		trace.Error("HeartbeatLoopTask heartbeatSend json marshal failed, server id=%v, status=%+v, err=%v",
			h.serverId, serverStatus, err.Error())
		return
	}
	if _, err = redisdb.Set(redisInfo.Key, string(jsonData), redisInfo.Expire); nil != err {
		trace.Error("HeartbeatLoopTask heartbeatSend failed, key=%v, "+
			"server id=%v, error=%v", redisInfo.Key, h.serverId, err.Error())
		return
	}

	trace.Info("HeartbeatLoopTask heartbeatSend success key=%v, server id=%v, val=%v",
		redisInfo.Key, h.serverId, string(jsonData))
}

// getClusterOnline 外部接口,获取当前集群个数
func (h *HeartbeatLoopTask) getClusterOnline() int {
	h.clusterNumMutex.Lock()
	defer h.clusterNumMutex.Unlock()
	return h.clusterNum
}

// GetClusterOnlineWithDefault 外部接口,获取当前集群个数,默认值为3
func (h *HeartbeatLoopTask) GetClusterOnlineWithDefault() (clusterOnline int) {
	clusterOnline = h.getClusterOnline()
	if clusterOnline <= 0 {
		clusterOnline = 3
	}

	return
}
