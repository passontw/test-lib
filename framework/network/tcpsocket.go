package network

import (
	"errors"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"sl.framework.com/async"
	"sl.framework.com/base"
	"sl.framework.com/trace"
)

const (
	SocketType_normal        int32 = 0 //普通模式
	SocketType_noreadtimeout int32 = 1 //read没有超时,client的时候,可能长时间不会收到任何包(服务器不发keeplive包)
)

type packetRawData []byte

type fnRecvCallback func([]byte)
type fnCheckPacketCallback func([]byte) (uint32, error)
type fnClosedCallback func()
type tcpsocket struct {
	conn *net.TCPConn
	//buff       *streamBuffer
	headerSize    uint32
	inPacket      chan packetRawData
	outPacket     chan packetRawData
	alive         base.AtomBool
	OnRecvPacket  fnRecvCallback
	OnClosed      fnClosedCallback
	OnCheckPacket fnCheckPacketCallback
	remoteAddr    string
	bStart        base.AtomBool
	syncWait      *sync.WaitGroup
	quitRecvChan  chan bool
	quitSendChan  chan bool
	socketType    int32 //0-normal 1-read时没有超时
}

func newTcpSocket(conn *net.TCPConn, headerSize uint32) *tcpsocket {
	initConn(conn)
	p := &tcpsocket{
		conn: conn,
		//buff:       newStreamBuffer(),
		headerSize:   headerSize,
		remoteAddr:   conn.RemoteAddr().String(),
		inPacket:     make(chan packetRawData, 10240),
		outPacket:    make(chan packetRawData, 10240),
		syncWait:     new(sync.WaitGroup),
		quitRecvChan: make(chan bool, 1),
		quitSendChan: make(chan bool, 1),
		socketType:   SocketType_noreadtimeout,
	}
	p.bStart.Set(false)
	return p
}

var (
	Conf_AliveTime     = time.Second * 600 /* keepalive 时间 */
	Conf_ReadBlockTime = time.Second * 900 /* tcp 读数据超时时间 */
	Conf_SendBlockTime = time.Second * 900 /* tcp 写数据超时时间 */
	Conf_ReadBufLen    = 1024 * 8          /* tcp 接收缓冲区 */
	Conf_SendBufLen    = 1024 * 8          /* tcp 发送缓冲区 */
	Conf_Nodelay       = true              /* tcp nodelay */
	Conf_Interval      = time.Second * 10  /* 线路检测周期 */
	Conf_Timeout       = time.Second * 10  /* 线路检测超时时间 */
)

func initConn(conn *net.TCPConn) bool {
	if err := conn.SetKeepAlivePeriod(Conf_AliveTime); err != nil {
		trace.Error("initConn, SetKeepAlivePeriod failed, err=%v, addr=%v", err, conn.RemoteAddr().String())
		return false
	}

	if err := conn.SetReadBuffer(Conf_ReadBufLen); err != nil {
		trace.Error("initConn, SetReadBuffer failed, err=%v, addr=%v", err, conn.RemoteAddr().String())
		return false
	}
	if err := conn.SetWriteBuffer(Conf_SendBufLen); err != nil {
		trace.Error("initConn, SetWriteBuffer failed, err=%v, addr=%v", err, conn.RemoteAddr().String())
		return false
	}
	if err := conn.SetNoDelay(Conf_Nodelay); err != nil {
		trace.Error("initConn, SetNoDelay failed, err=%v, addr=%v", err, conn.RemoteAddr().String())
		return false
	}
	trace.Info("initConn %v success", conn.RemoteAddr().String())
	return true
}

func (t *tcpsocket) SendPacket(buf []byte) error {
	if t.isAlive() {
		if t.outPacket != nil {
			t.outPacket <- buf
			return nil
		}
	}
	return errors.New("socket not alive")
}

func (t *tcpsocket) SetAlive(flag bool) {
	t.alive.Set(true)
}

func (t *tcpsocket) readPacket() (packetRawData, error) {
	t.conn.SetReadDeadline(time.Now().Add(Conf_ReadBlockTime))
	msgHead := make([]byte, t.headerSize)
	if _, err := io.ReadFull(t.conn, msgHead); err != nil {
		return nil, err
	}
	pktLen, err := t.OnCheckPacket(msgHead)
	if err != nil {
		return nil, err
	}
	msgBuf := msgHead

	//attention: pktLen is uint32, when minus will be large number cause out of memory. exsample 0-16=4294967200
	dataSize := int32(pktLen - t.headerSize)
	if dataSize > 0 {
		msgBody := make([]byte, dataSize)
		if _, err := io.ReadFull(t.conn, msgBody); err != nil {
			return nil, err
		}
		msgBuf = append(msgBuf, msgBody...)
	}
	return msgBuf, nil
}

func (t *tcpsocket) exitThread() {
	t.syncWait.Wait()
	t.alive.Set(false)
	close(t.inPacket)
	close(t.outPacket)
	t.conn.Close()
	trace.Notice("close socket %v", t.conn.RemoteAddr().String())
	if t.OnClosed != nil {
		t.OnClosed()
	}
}

func (t *tcpsocket) doSendPacket() {
	for {
		select {
		case ch, _ := <-t.outPacket:
			{
				if ch != nil {
					totalLen := len(ch)
					for totalLen > 0 {
						t.conn.SetWriteDeadline(time.Now().Add(Conf_SendBlockTime))
						n, err := t.conn.Write(ch)
						if err != nil {
							trace.Error("write to %v error, err=%v", t.conn.RemoteAddr(), err)
							t.Close()
							return
						}
						totalLen = totalLen - n
						ch = ch[n:]
					}
				}
			}
		case <-t.quitSendChan:
			{
				t.syncWait.Done()
				trace.Info("doSendPacket, recv quit, addr=%v", t.conn.RemoteAddr())
				return
			}
		}
	}
}

func (t *tcpsocket) doRecvPacket() {
	for {
		select {
		case ch, _ := <-t.inPacket:
			{
				if t.OnRecvPacket != nil {
					t.OnRecvPacket(ch)
				}
			}
		case <-t.quitRecvChan:
			{
				trace.Info("doRecvPacket, recv quit, addr=%v", t.conn.RemoteAddr())
				t.syncWait.Done()
				return
			}
		}
	}
}

func (t *tcpsocket) doReadPacket() {
	for {
		if !t.isAlive() {
			break
		}
		pkt, err := t.readPacket()
		if err != nil {
			trace.Info("doReadData failed, err=%v, addr=%v", err, t.conn.RemoteAddr())
			t.Close()
			break
		}
		if t.inPacket != nil {
			t.inPacket <- pkt
		} else {
			break
		}
	}
	t.syncWait.Done()
	trace.Info("doReadPacket quit, addr=%v", t.conn.RemoteAddr())
}

func (t *tcpsocket) doWork() {
	if t.bStart.Get() {
		trace.Error("doWork failed, alread start")
		return
	}
	t.bStart.Set(true)
	t.SetAlive(true)

	t.syncWait.Add(3)
	async.AsyncRunCoroutine(func() {
		t.exitThread()
	})

	// send data
	async.AsyncRunCoroutine(func() {
		t.doSendPacket()
	})

	async.AsyncRunCoroutine(func() {
		t.doRecvPacket()
	})

	async.AsyncRunCoroutine(func() {
		t.doReadPacket()
	})
}

func (t *tcpsocket) Close() {
	trace.Info("Close addr=%v", t.conn.RemoteAddr())
	t.alive.Set(false)
	t.quitRecvChan <- true
	t.quitSendChan <- true
	t.conn.Close()
}

func (t *tcpsocket) isAlive() bool {
	return t.alive.Get()
}

func (t *tcpsocket) SetSocketType(sockettype int32) {
	atomic.StoreInt32(&t.socketType, sockettype)
}

func (t *tcpsocket) IsCloseWhenReadTimeout() bool {
	if SocketType_noreadtimeout == atomic.LoadInt32(&t.socketType) {
		return false
	} else {
		return true
	}
}
