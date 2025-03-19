package network

import (
	"errors"
	"net"
	"sl.framework.com/async"
	"sl.framework.com/trace"
)

type TcpConnectorNew struct {
	sock          *tcpsocket
	OnConnected   fnCallback
	OnClosed      fnClosedCallback
	OnCheckPacket fnCheckPacketCallback
	OnRecvPacket  fnRecvCallback
	IsLogin       bool
	HeaderSize    uint32
	Addr          string
}

func NewTcpConnectorNew(dest string, headerSize uint32) *TcpConnectorNew {
	return &TcpConnectorNew{Addr: dest, HeaderSize: headerSize}
}

func (c *TcpConnectorNew) Connect() error {
	trace.Notice("try to connect %s ...", c.Addr)
	addr, err := net.ResolveTCPAddr("tcp", c.Addr)
	if err != nil {
		trace.Notice("ResolveTCPAddr %s failed, err=%v", c.Addr, err)
	}
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		trace.Error("connect %s failed, %v", c.Addr, err)
		return err
	}
	trace.Notice("connected %s.", c.Addr)
	c.sock = newTcpSocket(conn, c.HeaderSize)
	if c.sock != nil {
		c.sock.SetAlive(true)
		c.sock.OnCheckPacket = c.OnCheckPacket
		c.sock.OnRecvPacket = c.OnRecvPacket
		c.sock.OnClosed = c.OnClosed
		c.connected()
		async.AsyncRunCoroutine(func() {
			c.sock.doWork()
		})
		return nil
	}
	return errors.New("create error")
}

func (c *TcpConnectorNew) connected() {
	if c.OnConnected != nil {
		c.OnConnected()
	}
}

func (c *TcpConnectorNew) SendData(data []byte) error {
	if c.sock != nil {
		return c.sock.SendPacket(data)
	}
	return errors.New("no socket")
}

func (c *TcpConnectorNew) Close() {
	if c.sock != nil {
		c.sock.Close()
	}
}
