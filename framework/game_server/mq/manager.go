package mq

import (
	"github.com/apache/rocketmq-client-go/v2/rlog"
	"sl.framework.com/game_server/conf"
	"sl.framework.com/game_server/mq/handler"

	//"sl.framework.com/game_server/game/service/listenner"
	"sl.framework.com/trace"
	"strings"
	"sync"
)

// Topic racket messageq 主题类型
type Topic string

const (
	TopicGameServerEvent  Topic = "game-server-event"  //游戏事件主题
	TopicGameDraw         Topic = "game-draw"          //游戏开奖主题
	TopicBetConfirm       Topic = "bet-confirm"        //游戏提交注单
	TopicJoinMessageRoom  Topic = "join-message-room"  //玩家进入房间
	TopicLeaveMessageRoom Topic = "leave-message-room" //玩家离开你房间

	TopicGameDrawOut0 Topic = "game-draw-0" //游戏开奖主题0
	TopicGameDrawOut1 Topic = "game-draw-1" //游戏开奖主题1
	TopicGameDrawOut2 Topic = "game-draw-2" //游戏开奖主题2
	TopicGameDrawOut3 Topic = "game-draw-3" //游戏开奖主题3
	TopicGameDrawOut4 Topic = "game-draw-4" //游戏开奖主题4
	TopicGameDrawOut5 Topic = "game-draw-5" //游戏开奖主题5
	TopicGameDrawOut6 Topic = "game-draw-6" //游戏开奖主题6
	TopicGameDrawOut7 Topic = "game-draw-7" //游戏开奖主题7
	TopicGameDrawOut8 Topic = "game-draw-8" //游戏开奖主题8
	TopicGameDrawOut9 Topic = "game-draw-9" //游戏开奖主题9

	TopicBetConfirmOut0 Topic = "bet-confirm-0" //注单提交主题0
	TopicBetConfirmOut1 Topic = "bet-confirm-1" //注单提交主题1
)

// Property racket mq属性变量
type Property string

const (
	propertyTraceId   Property = "traceId"
	propertyTimestamp Property = "timestamp"
)

var (
	consumerInitOnce sync.Once
	consumerManager  *ConsumerManager // 管理者对象
)

// ConsumerManager racket消费者的管理对象
type ConsumerManager struct {
	consumers []*Consumer //消费者列表
}

/**
 * stopConsumerManager
 * 关闭所有消费者
 *
 * @param
 * @return
 */

func stopConsumerManager() {
	for _, c := range consumerManager.consumers {
		c.shutdown()
	}
}

/**
 * initConsumerManger
 * Racket管理者初始化 供外部调用
 *
 * @param
 * @return
 */

func initConsumerManger() (ok bool) {
	trace.Info("initConsumerManger...")
	consumerInitOnce.Do(func() {
		ok = consumerManagerInitOnce()
	})

	return
}

/**
 * consumerManagerInitOnce
 * mq消息管理者初始化
 *
 * @param
 * @return
 */

func consumerManagerInitOnce() (ok bool) {
	nameServer := conf.GetRocketMQNameServer()
	//gameServerEventGroup := conf.GetGameServerEventGroup()
	//gameDrawGroup := conf.GetGameDrawGroup()
	joinMessageRoom := conf.GetJoinMessageRoomGroup()
	leaveMessageRoom := conf.GetLeaveMessageRoomGroup()
	retries := conf.GetRocketMQRetries()
	trace.Info("consumerManagerInitOnce nameServer=%v", nameServer)
	//workaround:域名需要增加http://头否则解析失败 与java保持一致 java用法没有这个http://
	if len(nameServer) > 0 && !strings.HasPrefix(nameServer[0], "http") {
		nameServer[0] = "http://" + nameServer[0]
	}
	if len(nameServer) <= 0 || retries == 0 {
		trace.Error("rocketMQInit failed, nameServer=%v,  joinMessageRoom=%v, "+
			"leaveMessageRoom=%v, retries=%v", nameServer, joinMessageRoom,
			leaveMessageRoom, retries)
		return false
	}

	//设置日志打印级别 打印>=warn
	//日志级别:trace<debug<info<warning<error<fatal<panic
	rlog.SetLogLevel("warn")

	for _, topicDTO := range conf.ServerConf.Rocketmq.GameDrawTopicsIn {
		ok = RegistMQConsumer(Topic(topicDTO.TopicName), topicDTO.TopicGroup)
		if !ok {
			return false
		}
	}

	for _, topicDTO := range conf.ServerConf.Rocketmq.BetConfirmTopicsIn {
		ok = RegistMQConsumer(Topic(topicDTO.TopicName), topicDTO.TopicGroup)
		if !ok {
			return false
		}
	}
	return ok
}

// 根据topic创建不同给的消费者
func RegistMQConsumer(topic Topic, topGroup string) (ok bool) {
	nameServer := conf.GetRocketMQNameServer()
	if len(nameServer) > 0 && !strings.HasPrefix(nameServer[0], "http://") {
		nameServer[0] = "http://" + nameServer[0]
	}
	trace.Info("RegistMQConsumer nameServer=%v, topic=%v,topGroup=%v", nameServer, topic, topGroup)
	// 创建消费者管理对象
	var c *Consumer
	consumerManager = &ConsumerManager{
		consumers: make([]*Consumer, 0, 4),
	}

	switch topic {
	case TopicGameDrawOut0, TopicGameDrawOut1, TopicGameDrawOut2,
		TopicGameDrawOut3, TopicGameDrawOut4, TopicGameDrawOut5,
		TopicGameDrawOut6, TopicGameDrawOut7, TopicGameDrawOut8,
		TopicGameDrawOut9:
		// game-draw 结算事件消费者
		if c, ok = newConsumer(nameServer, topic, topGroup, cluster,
			handler.OnGameDrawHandler, conf.GetRocketMQRetries()); !ok {
			return ok
		}
	case TopicBetConfirmOut0, TopicBetConfirmOut1:
		if c, ok = newConsumer(nameServer, topic, topGroup, cluster,
			handler.OnBetConfirmHandler, conf.GetRocketMQRetries()); !ok {
			return ok
		}
	default:
		trace.Info("RegistMQConsumer nameServer=%v, topic=%v,topGroup=%v default ", nameServer, topic, topGroup)
	}

	consumerManager.consumers = append(consumerManager.consumers, c)
	return ok
}

var (
	producerInitOnce sync.Once
	producerManager  *ProducerManager // 生产者管理对象
)

// ProducerManager rocket生产者的管理对象
type ProducerManager struct {
	producers map[Topic]*Producer //生产者映射 map[string]*Producer
}

/**
 * stopProducerManager
 * 关闭所有生产者
 *
 * @param
 * @return
 */

func stopProducerManager() {
	for _, c := range producerManager.producers {
		c.shutdown()
	}
}

/**
 * initConsumerManger
 * Racket管理者初始化 供外部调用
 *
 * @param
 * @return
 */

func initProducerManger() (ok bool) {
	trace.Info("initProducerManger...")
	producerInitOnce.Do(func() {
		ok = producerManagerInitOnce()
	})

	return
}

/**
 * producerManagerInitOnce
 * mq生产者管理者初始化
 *
 * @param
 * @return
 */

func producerManagerInitOnce() (ok bool) {
	nameServer := conf.GetRocketMQNameServer()
	retries := conf.GetRocketMQRetries()
	trace.Info("producerManagerInitOnce nameServer=%v", nameServer)
	//workaround:域名需要增加http://头否则解析失败 与java保持一致 java用法没有这个http://
	if len(nameServer) > 0 && !strings.HasPrefix(nameServer[0], "http://") {
		nameServer[0] = "http://" + nameServer[0]
	}
	if len(nameServer) <= 0 || retries == 0 {
		trace.Error("rocketMQInit failed, nameServer=%v, retries=%v", nameServer, retries)
		return
	}

	//设置日志打印级别 打印>=warn
	//日志级别:trace<debug<info<warning<error<fatal<panic
	rlog.SetLogLevel("warn")

	// 创建生产者管理对象
	var p *Producer
	producerManager = &ProducerManager{
		producers: make(map[Topic]*Producer),
	}

	// 创建生产者
	//gameDrawGroup := conf.GetGameDrawGroup()
	for _, item := range conf.ServerConf.Rocketmq.GameDrawTopicsOut {
		if p, ok = newProducer(nameServer, TopicGameDraw, item.TopicName, conf.GetRocketMQRetries()); !ok {
			trace.Error("producerManagerInitOnce failed, nameServer=%v, retries=%v topicname:%v", nameServer, retries, item.TopicName)
			return
		}
		topicName := Topic(item.TopicName)
		producerManager.producers[topicName] = p
	}
	for _, item := range conf.ServerConf.Rocketmq.BetConfirmTopicsOut {
		if p, ok = newProducer(nameServer, TopicBetConfirm, item.TopicName, conf.GetRocketMQRetries()); !ok {
			trace.Error("producerManagerInitOnce failed, nameServer=%v, retries=%v topicname:%v", nameServer, retries, item.TopicName)
			return
		}
		topicName := Topic(item.TopicName)
		producerManager.producers[topicName] = p
	}
	//listenner.SetMessageSender(SendMessage)

	return true
}

/**
 * SendMessage
 * 通过生产者发送消息
 *
 * @param tag - messageq tag
 * @param traceId - 用于日志跟踪
 * @param topic - 发送消息的主题
 * @param message - 要发送的消息
 * @return
 */

func SendMessage(topic, tag, traceId, timestamp string, message string) {
	t := Topic(topic)
	producer, ok := producerManager.producers[t]
	if !ok {
		trace.Error("SendMessage no producer for topic=%v, message=%v", topic, message)
		return
	}

	producer.sendMsg(t, tag, traceId, timestamp, message)
}

/**
 * InitRocketManager
 * 初始化rocket mq包括rocket mq的消费者和生产者
 *
 * @param
 * @return bool - 是否初始化成功
 */

func InitRocketManager() bool {
	okConsumer := initConsumerManger()
	okProducer := initProducerManger()

	//关闭启动成功的消费者或者生产者
	if !(okConsumer && okProducer) {
		stopConsumerManager()
		stopProducerManager()
	}

	return okConsumer && okProducer
}

/**
 * StopRocketManager
 * 关闭rocket mq包括rocket mq的消费者和生产者
 *
 * @param
 * @return
 */

func StopRocketManager() {
	stopConsumerManager()
	stopProducerManager()
}
