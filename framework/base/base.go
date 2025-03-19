package base

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"net"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"runtime/pprof"
	"sl.framework.com/async"
	"sl.framework.com/trace"
	"syscall"
	"time"
	"unsafe"
)

// 打印系统相关的数据
func PrintSysInfo() {
	fmt.Printf("=======================> cpu num:%d\n", runtime.NumCPU())
	fmt.Printf("=======================> go version:%v\n", runtime.Version())
}

func Wait() {
	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT)
	for {
		select {
		case s := <-chSignal:
			trace.Notice("receive os.signal:%v, exited", s)
			time.Sleep(time.Millisecond * 500) // wait for writing trace
			os.Exit(0)
		}
	}
}

// 获取数据的字符流长度
func StreamSizeof(data interface{}) uint32 {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, data)
	if err != nil {
		trace.Notice("StreamSizeof failed err = %v", err)
	}
	return uint32(buf.Len())
}

// 延迟执行
func RunAfter(duration time.Duration, handler interface{}) {
	async.AsyncRunCoroutine(func() {
		timer := time.NewTimer(duration)
		defer timer.Stop() // Ensure timer is stopped
		select {
		case <-timer.C:
			{
				if cb, ok := handler.(func()); ok {
					cb()
				} else {
					trace.Error("interfaces(%v) is not of type func()", handler)
				}
				return
			}
		}
	})
}

// 序列化
func SerializeToBytes(data interface{}) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, data)
	if err != nil {
		panic(err.Error())
	}
	return buf.Bytes()
}

// 反序列化
func UnSerializeFromBytes(buf []byte, data interface{}) error {
	buffer := new(bytes.Buffer)
	buffer.Write(buf)
	err := binary.Read(buffer, binary.BigEndian, data)
	if err != nil {
		return err
	}
	return nil
}

func PProf(file string) {
	f, err := os.Create(file)
	if err != nil {
		panic(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
}

func ArrayToString(data interface{}) string {
	// 使用反射检查类型是否为数组
	value := reflect.ValueOf(data)
	if value.Kind() != reflect.Array {
		panic("data is not array")
	}

	// 创建一个缓冲区，将数组写入其中
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, data)
	if err != nil {
		panic(fmt.Sprintf("binary write failed: %v", err))
	}

	// 获取写入的字符串并去除末尾的所有 \x00 字节
	str := buf.String()
	return string(bytes.TrimRight([]byte(str), "\x00"))
}

func Htonl(num uint32) uint32 {
	buffer := SerializeToBytes(&num)
	buf := bytes.NewReader(buffer)
	binary.Read(buf, binary.LittleEndian, &num)
	return num
}

func Ltonh(num uint32) uint32 {
	buffer := SerializeToBytes(&num)
	buf := bytes.NewReader(buffer)
	binary.Read(buf, binary.BigEndian, &num)
	return num
}

func String(data []byte) string {
	str := *(*string)(unsafe.Pointer(&data))
	return str
}

func ConvertSecsToStr(sec int64) string {
	tm := time.Unix(sec, 0)
	return tm.Format(time.DateTime)
}

func IpToStr(ip int64) string {
	return fmt.Sprintf("%d.%d.%d.%d",
		byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
}

func IpToInt(ip string) int64 {
	if r := net.ParseIP(ip); nil == r {
		return 0
	}
	ret := big.NewInt(0)
	ret.SetBytes(net.ParseIP(ip).To4())
	return ret.Int64()
}

func ToHex(data []byte) string {
	dst := make([]byte, hex.EncodedLen(len(data)))
	hex.Encode(dst, data)
	return fmt.Sprintf("%s", dst)
}

/**
 * 校验入参是否是空
 * CheckEmpty
 *
 * @param data interface{} - 任意类型
 * @return bool - true 参数为空，false 参数非空
 */

func CheckEmpty(data interface{}) bool {
	if nil == data {
		return true
	}
	val := reflect.ValueOf(data)
	switch val.Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
		return val.Len() == 0
	case reflect.Pointer, reflect.Interface:
		return val.IsNil()
	default:
		trace.Warning("[校验是否为空] unhandled default case 类型：%v", val.Kind())
	}
	return false
}
