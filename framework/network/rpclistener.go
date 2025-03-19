package network

import (
	"net"
	"net/rpc"
)

// TODO
type RpcListener struct {
	addr string
}

func NewRpcListener(addr string) *RpcListener {
	return &RpcListener{addr: addr}
}

func (r *RpcListener) Start() error {
	rpc.HandleHTTP()

	_, err := net.Listen("tcp", r.addr)
	if err != nil {
		return err
	}
	return nil
}
