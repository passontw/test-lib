package network

import (
	"errors"
	"net"
	"reflect"
	"strings"
	"sync"
	"time"

	pb "github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/runtime/protoiface"
	"sl.framework.com/async"
	"sl.framework.com/base"
	"sl.framework.com/trace"
)

const (
	DEFAULT_PACKET_HEAD_LEN   = 16
	DEFAULT_CONNECTING_PERIOD = 5 * time.Second
	DEFAULT_LOGINING_PERIOD   = 5 * time.Second
	DEFAULT_KEEPALIVE_PERIOD  = 5 * time.Second
	default_connect_period    = 20 * time.Second
)

const (
	ConnectState_IDLE = iota
	ConnectState_CONNECTING
	ConnectState_CONNECTED
	ConnectState_LOGINING
	ConnectState_LOGIN
)

type Default_Packet_Header struct {
	Cmd     uint32 // 协议号
	Size    uint32 // 协议大小
	Seq     uint32 // 保留
	Session uint16 // 保留
	Version uint16 // 保留
}

type fnCallback func()
type msgCallback func(interface{})

type PbMsgInfo struct {
	msgType    reflect.Type
	msgHandler interface{}
	msgID      string
}

type TcpConnector struct {
	sock          *tcpsocket
	OnSendLogin   fnCallback
	OnKeepAlive   fnCallback
	OnConnected   fnCallback
	OnClosed      fnClosedCallback
	OnCheckPacket fnCheckPacketCallback
	OnRecvPacket  fnRecvCallback
	IsLogin       bool
	HeaderSize    uint32
	Addr          string

	// private variable
	msgInfoes       map[uint32]*PbMsgInfo
	onceFirstPacket sync.Once
	tmRecv          int64
	tmCreate        int64
	hasValidPacket  base.AtomBool
	connState       uint32
	refreshTime     time.Time
	exit            bool
	exitWait        sync.WaitGroup
	socketType      int32
}

func NewTcpConnector(dest string) *TcpConnector {
	connector := &TcpConnector{
		Addr:        dest,
		HeaderSize:  base.StreamSizeof(Default_Packet_Header{}),
		msgInfoes:   make(map[uint32]*PbMsgInfo),
		connState:   ConnectState_IDLE,
		refreshTime: time.Now(),
		exit:        false,
		socketType:  SocketType_noreadtimeout,
	}
	return connector
}

func (c *TcpConnector) RegisterPbMsgHandler(cmd uint32, msg interface{}, f interface{}) error {
	msgType := reflect.TypeOf(msg)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		return errors.New("header pointer required")
	}
	msgID := msgType.Elem().Name()
	if msgID == "" {
		return errors.New("unamed pb msg")
	}
	if _, ok := c.msgInfoes[cmd]; ok {
		return errors.New("pb message " + msgID + "is already register")
	}

	i := new(PbMsgInfo)
	i.msgType = msgType
	i.msgHandler = f
	i.msgID = msgID
	c.msgInfoes[cmd] = i
	return nil
}

func (c *TcpConnector) Connect() error {
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
		c.sock.SetSocketType(c.socketType)
		c.sock.OnCheckPacket = c.cbCheckPacket
		c.sock.OnRecvPacket = c.cbRecvPacket
		c.sock.OnClosed = c.cbClosed
		c.connected()
		c.connState = ConnectState_CONNECTED
		async.AsyncRunCoroutine(func() {
			c.sock.doWork()
		})
		return nil
	}
	return errors.New("create error")
}

func (c *TcpConnector) connected() {
	if c.OnConnected != nil {
		c.OnConnected()
	}
}

func (c *TcpConnector) SendData(data []byte) error {
	if c.sock != nil {
		return c.sock.SendPacket(data)
	}
	return errors.New("no socket")
}

func (c *TcpConnector) GetIp() uint32 {
	arr := strings.Split(c.sock.remoteAddr, ":")
	if len(arr) > 0 {
		return uint32(base.IpToInt(arr[0]))
	}
	return 0
}

func (c *TcpConnector) Endpoint() string {
	return c.sock.remoteAddr
}

// 主动close socket
func (c *TcpConnector) Close() {
	c.sock.Close() //触发cbClose,set to ConnectState_IDLE
	trace.Notice("TcpConnector initiative to close, addr=%s", c.Addr)
}

func (c *TcpConnector) SendPacket(buf []byte) error {
	return c.sock.SendPacket(buf)
}

func (c *TcpConnector) LoginRet(stat bool) {
	if stat {
		c.connState = ConnectState_LOGIN
		c.refreshTime = time.Now()
		if c.OnKeepAlive != nil {
			c.OnKeepAlive()
		}
	} else {
		c.Close() //login failed
	}
}

// connect->connected->sendlogin->login->onkeeplive
// socket close后会重连
func (c *TcpConnector) Run() {

	async.AsyncRunCoroutine(func() {
		isFirst := true
		for {
			if c.exit {
				c.exitWait.Done()
				return
			}
			if c.connState == ConnectState_IDLE {
				time_now := time.Now()
				diff_time := time_now.Sub(c.refreshTime)
				if diff_time >= default_connect_period || isFirst {
					c.connState = ConnectState_CONNECTING
					c.Connect()
					c.refreshTime = time.Now()
					isFirst = false
				}
			} else if c.connState == ConnectState_CONNECTING {
				time_now := time.Now()
				diff_time := time_now.Sub(c.refreshTime)
				if diff_time >= DEFAULT_CONNECTING_PERIOD {
					c.connState = ConnectState_IDLE
				}
			} else if c.connState == ConnectState_CONNECTED {
				if c.OnSendLogin != nil {
					c.OnSendLogin()
					c.connState = ConnectState_LOGINING
					c.refreshTime = time.Now()
					trace.Notice("send login to connect %s ...", c.Addr)
				}
			} else if c.connState == ConnectState_LOGINING {
				time_now := time.Now()
				diff_time := time_now.Sub(c.refreshTime)
				if diff_time >= DEFAULT_LOGINING_PERIOD {
					trace.Notice("long time no response login")
					c.Close()
					trace.Notice("long time no response login, close socket")
				}
			} else if c.connState == ConnectState_LOGIN {
				time_now := time.Now()
				diff_time := time_now.Sub(c.refreshTime)
				if diff_time >= DEFAULT_KEEPALIVE_PERIOD {
					c.refreshTime = time.Now()
					if c.OnKeepAlive != nil {
						c.OnKeepAlive()
					}
				}
			}
			time.Sleep(time.Duration(50) * time.Millisecond)
		}

		// c.sock.doWork()
	})
}

func (c *TcpConnector) Stop() {
	c.exitWait.Add(1)
	c.exit = true
	c.exitWait.Wait()
}

func (c *TcpConnector) cbRecvPacket(buf []byte) {
	if c.OnRecvPacket != nil {
		c.OnRecvPacket(buf)
	} else {
		header := &Default_Packet_Header{}
		base.UnSerializeFromBytes(buf, header)

		//if p.gs != nil{
		//	trace.Info("cbRecvPacket, cmd=%0x, size=%v, seq=%v, sid=%v, addr=%v",
		//		header.Cmd, header.Size, header.Seq,  p.gs.sid, p.connection.Endpoint())
		//}else {
		//	trace.Info("cbRecvPacket, cmd=%0x, size=%v, seq=%v, addr=%v",
		///		header.Cmd, header.Size, header.Seq, p.connection.Endpoint())
		//}

		c.onceFirstPacket.Do(func() {
			c.hasValidPacket.Set(true)
		})
		c.tmRecv = time.Now().Unix()

		if i, ok := c.msgInfoes[header.Cmd]; ok {
			pkMsg := reflect.New(i.msgType.Elem()).Interface().(protoiface.MessageV1)
			protoMsgBuf := buf[DEFAULT_PACKET_HEAD_LEN:]
			err := pb.Unmarshal(protoMsgBuf, pkMsg)
			if err != nil {
				trace.Error("Unmarshal failed, err=%v", err)
			} else {
				async.AsyncRunCoroutine(func() {
					i.msgHandler.(func(interface{}))(pkMsg)
				})
			}
		} else {
			trace.Notice("packet %0x not support", header.Cmd)
		}
	}
}

func (c *TcpConnector) cbCheckPacket(buf []byte) (uint32, error) {
	if c.OnCheckPacket != nil {
		return c.OnCheckPacket(buf)
	} else {
		header := &Default_Packet_Header{}
		err := base.UnSerializeFromBytes(buf, header)
		/*
			if p.gs != nil{
				trace.Info("cbCheckPacket, cmd=%0x, size=%v, seq=%v, sid=%v, addr=%v",
					header.Cmd, header.Size, header.Seq, p.gs.sid, p.connection.Endpoint())
			}else {
				trace.Info("cbCheckPacket, cmd=%0x, size=%v, seq=%v, addr=%v",
					header.Cmd, header.Size, header.Seq, p.connection.Endpoint())
			}*/
		valid := false
		nLen := header.Size
		if nil == err {
			if _, ok := c.msgInfoes[header.Cmd]; ok {
				valid = nLen >= base.StreamSizeof(header)
			}
			// trace.Info("check packet, cmd=%0x, size=%v", header.Cmd, header.Size)
		}
		if valid {
			return header.Size, nil
		} else {
			trace.Error("invalid packet, cmd=%08x, addr=%v, close socket",
				int(header.Cmd), c.Endpoint())
			c.Close()
		}
		return 0, err
	}
}

// socket close的時候会被调用
func (c *TcpConnector) cbClosed() {
	c.connState = ConnectState_IDLE
	if c.OnClosed != nil {
		c.OnClosed()
	}
	trace.Notice("TcpConnector cbClosed, addr=%s", c.Addr)
}

func (c *TcpConnector) HasValidPacket() bool {
	return c.hasValidPacket.Get()
}
