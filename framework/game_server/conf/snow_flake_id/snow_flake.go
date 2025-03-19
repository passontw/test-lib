package snowflaker

import (
	"github.com/bwmarrin/snowflake"
	"sl.framework.com/async"
	"sl.framework.com/game_server/conf"
	"sl.framework.com/trace"
	"sync"
	"time"
)

var (
	snowFlakeOnce     sync.Once
	snowFlakeInstance *SnowFlake
)

// Note:实测结果使用channel方式并不能更快,这种场景下反而慢了点

// GetSnowFlakeInstance 雪花算法实例单例
func GetSnowFlakeInstance() *SnowFlake {
	snowFlakeOnce.Do(func() {
		// 设置节点ID,可以从0到1023之间选择一个唯一的ID
		nodeID := conf.GetServerId() % 1024
		node, err := snowflake.NewNode(nodeID)
		if err != nil {
			trace.Error("GetSnowFlakeInstance failed, nodeId=%v, error=%v", nodeID, err.Error())
			return
		}
		snowFlakeInstance = &SnowFlake{
			node:   node,
			chId:   make(chan int64, 64), //用有缓冲channel 提前生产Id放入其中 供消息者消费
			chQuit: make(chan struct{}),
		}
		snowFlakeInstance.runningLoop()
		trace.Info("GetSnowFlakeInstance success, nodeId=%v", nodeID)
	})

	return snowFlakeInstance
}

// SnowFlake 雪花算法对象
type SnowFlake struct {
	node     *snowflake.Node
	chId     chan int64
	chQuit   chan struct{}
	quitOnce sync.Once
}

// runningLoop 启动循环生产Id
func (s *SnowFlake) runningLoop() {
	async.AsyncRunCoroutine(func() {
		trace.Info("SnowFlake runningLoop start")
	loop:
		for {
			select {
			case s.chId <- s.getSnowId():
			case _, ok := <-s.chQuit:
				if !ok {
					if s.chId != nil {
						close(s.chId)
					}
					s.node = nil
					trace.Notice("SnowFlake runningLoop stop")
					break loop
				}
			}
		}
	})

	return
}

// getSnowId 生成唯一id
func (s *SnowFlake) getSnowId() int64 {
	if s.node == nil {
		trace.Info("getSnowId node is nil, return timestamp")
		return time.Now().Unix()
	}

	return int64(s.node.Generate())
}

// GetUniqueId 生成唯一id
func (s *SnowFlake) GetUniqueId() int64 {
	id, ok := <-s.chId
	if !ok {
		trace.Notice("SnowFlake GetUniqueId get from chan failed.")
		return time.Now().Unix() //失败则返回当前时间戳
	}
	return id
}

// StopLoop 停止产生Id
func (s *SnowFlake) StopLoop() {
	s.quitOnce.Do(func() {
		trace.Info("SnowFlake StopLoop")
		if s.chQuit != nil {
			close(s.chQuit)
			s.chQuit = nil
		}
	})
}
