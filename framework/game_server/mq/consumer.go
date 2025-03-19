package mq

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"sl.framework.com/game_server/error_code"
	"sl.framework.com/game_server/game/dao/redisdb"
	"sl.framework.com/game_server/game/service"
	"sl.framework.com/game_server/redis/redis_tool"
	"sl.framework.com/trace"
	"time"
)

type (
	// 消费者收到数据时候的回调
	fnOnMessage func(string, []byte) int

	/*ConsumerMode 消费者消费模式*/
	ConsumerMode string
)

const (
	cluster   ConsumerMode = "Cluster"   //集群模式
	broadcast ConsumerMode = "Broadcast" //广播模式
)

/*
	Consumer消费者对象封装
*/

type Consumer struct {
	consumer rocketmq.PushConsumer //消费者接口
	topic    string                //监听的topic
	mode     ConsumerMode          //集合模式 集群模式或者广播模式
	handler  fnOnMessage           //处理函数
	group    string                //组名
	retries  int                   //失败重试次数
	running  bool                  //是否正在运行
}

/**
 * newConsumer
 * 创建一个新的消费者并启动
 *
 * @param topic - 消费者订阅的主题
 * @param mode - 消费者的模式 集群或者广播
 * @param interfaces - 消费者处理消息的回调函数
 * @param group - 消费者组
 * @return *Consumer - 返回新创建的消费者
 * @return bool - 创建消费者是否成功
 */

func newConsumer(addrs []string, topic Topic, group string, mode ConsumerMode, handler fnOnMessage, retries int) (*Consumer, bool) {
	consumerMode := consumer.Clustering
	if mode == broadcast {
		consumerMode = consumer.BroadCasting
	}

	//创建新的消费者包装对象
	c := &Consumer{
		topic:   string(topic),
		mode:    cluster,
		handler: handler,
		group:   group,
		retries: retries,
	}

	//创建消费者
	var err error
	if c.consumer, err = rocketmq.NewPushConsumer(
		consumer.WithNameServer(addrs),
		consumer.WithConsumeFromWhere(consumer.ConsumeFromFirstOffset),
		consumer.WithGroupName(group),            // 消费组名称
		consumer.WithConsumerModel(consumerMode), //注册模式
		consumer.WithRetry(retries),
	); nil != err {
		trace.Error("newConsumer create failed addr=%+v, group name=%v, topic=%v, mode=%v, retries=%v, error=%v",
			addrs, c.group, c.topic, mode, retries, err.Error())
		return nil, false
	}
	trace.Info("newConsumer create success addr=%+v, group name=%v, topic=%v, mode=%v, retries=%v",
		addrs, c.group, c.topic, mode, c.retries)

	//启动消费者作业
	c.work()

	return c, true
}

/**
 * work
 * 启动消费者相关流程
 *
 * @param
 * @return
 */

func (c *Consumer) work() {
	c.subscribe()
	c.start()
}

/**
 * subscribe
 * 订阅主题
 *
 * @param
 * @return
 */

func (c *Consumer) subscribe() {
	if c.consumer == nil {
		trace.Error("subscribe topic=%v, group=%v, mode=%v failed", c.topic, c.group, c.mode)
		return
	}

	//rocket mq收到消息后的回调
	cb := func(ctx context.Context, messages ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		retSum := errcode.ErrorOk
		ret := consumer.ConsumeSuccess
		for _, msg := range messages {
			trace.Info("subscribe receive msg=%+v", msg)
			traceId := msg.GetProperty("traceId")

			//这里不再进行校验 因为有可能同一个游戏服节点手动同一个traceId下的多个orderPlanId处理消息
			/*
				timestamp := msg.GetProperty("timestamp")
				if processStatusCheck(traceId, timestamp) {
					trace.Notice("subscribe traceId=%v, timestamp=%v, is processing, skip it.",
						traceId, timestamp)
					continue
				}
			*/
			retSum += c.handler(traceId, msg.Body)
		}
		if errcode.ErrorOk != retSum {
			ret = consumer.ConsumeRetryLater
		}
		return ret, nil
	}

	//只接收tag形式的mq消息
	messageSelector := consumer.MessageSelector{
		Type:       consumer.TAG,
		Expression: service.BuildMQTags(),
	}
	if err := c.consumer.Subscribe(c.topic, messageSelector, cb); nil != err {
		trace.Error("subscribe failed topic=%v, group=%v, mode=%v failed, error=%v",
			c.topic, c.group, c.mode, err.Error())
		return
	}
	trace.Info("subscribe success topic=%v, group=%v, mode=%v",
		c.topic, c.group, c.mode)
}

/**
 * start
 * 启动消费者
 *
 * @param
 * @return
 */

func (c *Consumer) start() {
	if c.consumer == nil {
		trace.Error("Consumer start consumer is nil, group name=%v, topic=%v, mode=%v, retries=%v",
			c.group, c.topic, c.mode, c.retries)
		return
	}

	// 启动消费者
	if err := c.consumer.Start(); nil != err {
		trace.Error("Consumer start failed group name=%v, topic=%v, mode=%v, retries=%v, error=%v",
			c.group, c.topic, c.mode, c.retries, err.Error())
		return
	}
	c.running = true
	trace.Info("Consumer start success group name=%v, topic=%v, mode=%v, retries=%v",
		c.group, c.topic, c.mode, c.retries)
}

/**
 * shutdown
 * 关闭消费者
 *
 * @param
 * @return
 */

func (c *Consumer) shutdown() {
	if c.consumer == nil || !c.running {
		trace.Error("Consumer shutdown consumer is nil or running=%v, group name=%v, topic=%v, mode=%v, retries=%v",
			c.group, c.topic, c.mode, c.retries, c.running)
		return
	}

	//关闭集群消费
	if err := c.consumer.Shutdown(); nil != err {
		trace.Error("Consumer shutdown failed, group name=%v, topic=%v, mode=%v, retries=%v, error=%v",
			c.group, c.topic, c.mode, c.retries, err.Error())
	} else {
		c.running = false
		trace.Info("Consumer shutdown success, group name=%v, topic=%v, mode=%v, retries=%v",
			c.group, c.topic, c.mode, c.retries)
	}
}

/**
 * processStatusCheck
 * double check当前消息是否正在处理
 * 不主动释放锁 等待锁过期后自动释放
 *
 * @param traceId - traceId 全系统唯一
 * @param timestamp - 时间戳
 * @return bool - 是否正在处理 true:正在处理 false:没有在处理
 */

func processStatusCheck(traceId, timestamp string) (isProcessing bool) {
	redisLockInfo := redistool.BuildRedisLockInfo(time.Duration(15)*time.Second, "DoubleCheck", traceId, timestamp)

	lockStatus := redisdb.TryLock(redisLockInfo)
	isProcessing = !lockStatus //获取锁成功也即没有正在处理
	return
}
