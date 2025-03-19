package snowflaker

import (
	"github.com/bwmarrin/snowflake"
	"sl.framework.com/async"
	"sl.framework.com/trace"
	"sync"
	"time"
)

var (
	once     sync.Once
	instance *SnowFlake
)

// UniqueId 生成唯一id
func UniqueId() int64 {
	id, ok := <-instance.chId
	if !ok {
		trace.Warning("UniqueId get from chan failed.")
		id = time.Now().Unix()
	}
	trace.Info("UniqueId: %d size: %d", id, len(instance.chId))
	return id
}

// Init 雪花算法实例单例
func Init(svrId int64) {
	once.Do(func() {
		// 设置节点ID,可以从0到1023之间选择一个唯一的ID
		nodeID := svrId % 1024
		node, err := snowflake.NewNode(nodeID)
		if err != nil {
			trace.Error("雪花算法初始化失败: nodeId(%d), error: %s", nodeID, err.Error())
			return
		}
		instance = &SnowFlake{
			node:   node,
			chId:   make(chan int64, 5000), //用有缓冲channel 提前生产Id放入其中 供消息者消费
			chQuit: make(chan struct{}),
		}
		instance.runningLoop()
		trace.Notice("雪花算法初始化成功: nodeId(%d)", nodeID)
	})
}

// Note:实测结果使用channel方式并不能更快,这种场景下反而慢了点

// Get 雪花算法实例单例
func Get(svrId int64) *SnowFlake {
	once.Do(func() {
		// 设置节点ID,可以从0到1023之间选择一个唯一的ID
		nodeID := svrId % 1024
		node, err := snowflake.NewNode(nodeID)
		if err != nil {
			trace.Error("Get failed, nodeId=%v, error=%v", nodeID, err.Error())
			return
		}
		instance = &SnowFlake{
			node:   node,
			chId:   make(chan int64, 64), //用有缓冲channel 提前生产Id放入其中 供消息者消费
			chQuit: make(chan struct{}),
		}
		instance.runningLoop()
		trace.Info("Get success, nodeId=%v", nodeID)
	})

	return instance
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
		trace.Warning("getSnowId node is nil, return timestamp")
		return time.Now().Unix()
	}

	return int64(s.node.Generate())
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
