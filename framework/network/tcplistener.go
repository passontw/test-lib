package network

import (
	"fmt"
	"net"
	"sl.framework.com/async"
	"sl.framework.com/trace"
)

type FnAcceptNewConn func(uint16, *net.TCPConn)

type TcpListener struct {
	fnAddClient FnAcceptNewConn
}

func NewTcpListener(f FnAcceptNewConn) *TcpListener {
	return &TcpListener{
		fnAddClient: f,
	}
}

func (t *TcpListener) StartListenPorts(ports []uint16) error {
	for _, v := range ports {
		t.StartListen(v)
	}
	return nil
}

func (t *TcpListener) StartListen(port uint16) error {
	trace.Notice("listening on port %d", port)
	str := fmt.Sprintf(":%d", port)
	addr, err := net.ResolveTCPAddr("tcp", str)
	if err != nil {
		trace.Notice("StartListen, ResolveTCPAddr %v failed, err:%v", str, err)
		return err
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		trace.Error("listen on port %v failed, err:%v", port, err)
		return err
	}

	async.AsyncRunCoroutine(func() {
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				trace.Error("accept error %v", err.Error())
				continue
			}
			if t.fnAddClient != nil {
				t.fnAddClient(port, conn)
			}
		}
	})
	return nil
}
