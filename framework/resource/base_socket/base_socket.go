package basesk

import (
	"net"
	"sl.framework.com/base"
	"sl.framework.com/network"
	"sl.framework.com/trace"
	"sync"
	"time"
)

type ProtocolHandler func(ds any, packet []byte)

var protocolHandlers = make(map[uint32]ProtocolHandler)

// RegisterProtocolHandler 注册协议处理函数
func RegisterProtocolHandler(protocolID uint32, handler ProtocolHandler) {
	protocolHandlers[protocolID] = handler
	trace.Notice("[√] register protocol: [0x%06x] size: %d", protocolID, len(protocolHandlers))
}

type Socket struct {
	conn            *network.InboundConnection
	bConn           base.AtomBool
	headSize        uint32
	lastLoginTime   int64
	receiveDateTime int64
	createTime      int64
	onceFirstPacket sync.Once
	hasValidPacket  base.AtomBool
	stableMode      int8

	businessSocket any // 应该是业务socket
	loginCommand   uint32
}

// New 创建新的 Socket 并初始化 packetHandlers
func New(conn *net.TCPConn, stableMode int8, headSize uint32) *Socket {
	p := &Socket{
		conn:       network.NewInboundConnection(conn, headSize),
		stableMode: stableMode,
		createTime: time.Now().Unix(),
		headSize:   headSize,
	}
	p.hasValidPacket.Set(false)
	p.conn.OnCheckPacket = p.CheckPacket
	p.conn.OnRecvPacket = p.ReceivePacket
	return p
}

func (r *Socket) SetBusinessSocket(bsk any) {
	r.businessSocket = bsk
}

func (r *Socket) SetLoginCommand(cmd uint32) {
	r.loginCommand = cmd
}

func (r *Socket) SetClose(fn func()) {
	r.conn.OnClosed = fn
}

func (r *Socket) Endpoint() string {
	return r.conn.Endpoint()
}

func (r *Socket) Run() {
	r.conn.Run()
}

func (r *Socket) IsValidPacket() bool {
	return r.hasValidPacket.Get()
}

func (r *Socket) CreateTime() int64 {
	return r.createTime
}

func (r *Socket) ReceiveDateTime() int64 {
	return r.receiveDateTime
}

func (r *Socket) SetBConn(f bool) {
	r.bConn.Set(f)
}

func (r *Socket) Close() {
	r.conn.Close()
}

func (r *Socket) SendPacket(data []byte) error {
	return r.conn.SendPacket(data)
}

func (r *Socket) BConnGet() bool {
	return r.bConn.Get()
}

func (r *Socket) LastLoginTimeGet() int64 {
	return r.lastLoginTime
}

func (r *Socket) SetLastLoginTime(t int64) {
	r.lastLoginTime = t
}

func (r *Socket) CheckPacket(buf []byte) (uint32, error) {
	header := new(PacketHeader)
	err := base.UnSerializeFromBytes(buf, header)
	if err != nil {
		trace.Error("[1] addr %v, failed to serialize packet header: %v", r.conn.Endpoint(), err)
		return 0, err
	}
	// 使用 map 中注册的处理函数来处理数据
	if _, exists := protocolHandlers[header.Cmd]; exists {
		trace.Notice("[1] addr %s,     check packet, size=[%d], cmd=[0x%06x]", r.conn.Endpoint(), len(buf), header.Cmd)
		return header.Size, nil
	} else {
		return handleInvalidPacket(r, header.Cmd)
	}

}

func (r *Socket) ReceivePacket(buf []byte) {
	trace.Notice("[2] addr %v,   receive packet, size=[%d]", r.conn.Endpoint(), len(buf))

	header := new(PacketHeader)
	_ = base.UnSerializeFromBytes(buf, header)

	// 检查连接状态
	if !r.bConn.Get() && header.Cmd != r.loginCommand {
		trace.Error("[3] addr %s, received cmd [0x%06x] before login", r.conn.Endpoint(), int(header.Cmd))
		return
	}

	// 标记首次接收有效数据
	r.onceFirstPacket.Do(func() { r.hasValidPacket.Set(true) })
	r.receiveDateTime = time.Now().Unix()

	// 提取数据部分
	data := buf[r.headSize:]

	// 使用 map 中注册的处理函数来处理数据
	if handler, exists := protocolHandlers[header.Cmd]; exists {
		trace.Notice("[3] addr %s, deal with packet, size=[%d], cmd=[0x%06x]", r.conn.Endpoint(), len(buf), header.Cmd)
		handler(r.businessSocket, data)
	} else {
		trace.Warn("[3] addr %s, unsupported packet, size=[%d], cmd=[0x%06x]", r.conn.Endpoint(), len(buf), header.Cmd)
	}
}

func handleInvalidPacket(r *Socket, cmd uint32) (uint32, error) {
	msg := "keep socket"
	if r.stableMode <= 0 {
		msg = "close socket"
		r.conn.Close()
	}
	trace.Error("[1] addr %v, invalid packet, cmd=[0x%06x], %s", r.conn.Endpoint(), int(cmd), msg)
	return 0, nil
}
