package mq

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"sl.framework.com/async"
	"sl.framework.com/game_server/conf"
	"sl.framework.com/trace"
)

type (
	// RocketMessage 发送消息结构体
	RocketMessage struct {
		body      string
		topic     string
		tag       string //不同服务器根据tag过滤是否是自己关心的消息
		traceId   string //用于日志跟踪
		timestamp string //时间戳
	}

	// Producer 生产者对象封装
	Producer struct {
		producer     rocketmq.Producer   // rocketmq 生产者接口
		topic        string              // 该生产者对应的主题
		group        string              // 生产者组
		messageQueue chan *RocketMessage // 消息队列
		running      bool                //是否一起启动
	}
)

/**
 * newProducer
 * 创建一个新的生产者并启动
 *
 * @param topic - 生产者订阅的主题
 * @param group - 生产者
 * @param retries - 消息处理失败 重试的次数
 * @return *Producer - 返回新创建的生产者
 * @return bool - 创建生产者是否成功
 */

func newProducer(addrs []string, topic Topic, group string, retries int) (*Producer, bool) {
	//创建新的消费者包装对象
	p := &Producer{
		topic:        string(topic),
		group:        group,
		messageQueue: make(chan *RocketMessage, conf.GetRocketMQQueueMaxLen()),
	}

	//创建生产者
	var err error
	if p.producer, err = rocketmq.NewProducer(
		producer.WithNameServer(addrs),
		producer.WithRetry(retries),
		producer.WithGroupName(group),
	); nil != err {
		trace.Error("NewRocketProducer create producers failed,  addr=%v, groupName=%v, retries=%v, error=%v",
			addrs, group, retries, err.Error())
		return p, false
	}
	trace.Info("NewRocketProducer create producers success, addr=%v, groupName=%v, retries=%v",
		addrs, group, retries)

	//启动生产者
	p.start()

	return p, true
}

/**
 * start
 * 启动生产者
 *
 * @param
 * @return
 */

func (p *Producer) start() {
	if p.producer == nil {
		trace.Error("Producer start producers is nil, group name=%v, topic=%v", p.group, p.topic)
		return
	}

	//启动生产者
	fn := func() {
		if err := p.producer.Start(); nil != err {
			trace.Error("Producer start failed, group name=%v, topic=%v, error=%v", p.group, p.topic, err.Error())
			return
		}
		p.running = true

		trace.Info("Producer start success, group name=%v, topic=%v", p.group, p.topic)
		for m := range p.messageQueue {
			msg := &primitive.Message{
				Topic: m.topic,
				Body:  []byte(m.body),
			}
			msg.WithTag(m.tag)
			msg.WithProperty(string(propertyTraceId), m.traceId)
			msg.WithProperty(string(propertyTimestamp), m.timestamp)
			if result, err := p.producer.SendSync(context.Background(), msg); nil != err {
				trace.Error("Producer start send message failed, group name=%v, topic=%v, tag=%v, traceId=%v, "+
					"timestamp=%v, status=%v, msg id=%v, error=%v", p.group, p.topic, m.tag, m.traceId,
					m.timestamp, result.Status, result.MsgID, err.Error())
			}
			trace.Info("Producer send message done, group name=%v, topic=%v, tag=%v, traceId=%v, timestamp=%v, message=%v",
				p.group, p.topic, m.tag, m.traceId, m.timestamp, m.body)
		}
	}
	async.AsyncRunCoroutine(fn)
}

/**
 * shutdown
 * 关闭生产者
 *
 * @param
 * @return
 */

func (p *Producer) shutdown() {
	if p.producer == nil || !p.running {
		trace.Error("Producer shutdown producers is nil or running=%v, group name=%v, topic=%v",
			p.group, p.topic, p.running)
		return
	}

	//关闭生产者
	if err := p.producer.Shutdown(); nil != err {
		trace.Error("Producer shutdown failed, group name=%v, topic=%v,error=%v", p.group, p.topic, err.Error())
	} else {
		p.running = false
		trace.Info("Producer shutdown success, group name=%v, topic=%v", p.group, p.topic)
	}
}

/**
 * sendMsg
 * 发送消息rocketmq
 *
 * @param topic string - 消息主题
 * @param tag - messageq tag
 * @param traceId - 用于日志跟踪
 * @param msg string - 要发送的消息
 * @return
 */

func (p *Producer) sendMsg(topic Topic, tag, traceId, timestamp, msg string) {
	if len(p.messageQueue) >= conf.GetRocketMQQueueMaxLen() {
		trace.Error("sendMsg the queue is full, skip topic=%v, msg=%v", topic, msg)
		return
	}

	p.messageQueue <- &RocketMessage{
		topic:   string(topic),
		body:    msg,
		tag:     tag,
		traceId: traceId,
	}
}
