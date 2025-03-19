package network

import (
	"fmt"
	"net"
	"sl.framework.com/async"
	"sl.framework.com/base"
	"strings"
)

type InboundConnection struct {
	*tcpsocket
}

func NewInboundConnection(conn *net.TCPConn, headerSize uint32) *InboundConnection {
	return &InboundConnection{
		tcpsocket: newTcpSocket(conn, headerSize),
	}
}

func (c *InboundConnection) Run() {
	async.AsyncRunCoroutine(func() {
		c.doWork()
	})
}

func (c *InboundConnection) SetIp(str string) {
	arr := strings.Split(c.remoteAddr, ":")
	if len(arr) == 2 {
		c.remoteAddr = fmt.Sprintf("%s:%s", str, arr[1])
	} else {
		c.remoteAddr = str
	}
}

func (c *InboundConnection) GetIp() uint32 {
	arr := strings.Split(c.remoteAddr, ":")
	if len(arr) > 0 {
		return uint32(base.IpToInt(arr[0]))
	}
	return 0
}

func (c *InboundConnection) Endpoint() string {
	return c.remoteAddr
}
