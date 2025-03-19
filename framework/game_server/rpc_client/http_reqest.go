package rpcreq

import (
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/client/httplib"
	"sl.framework.com/game_server/conf"
	snowflaker "sl.framework.com/game_server/conf/snow_flake_id"
	"sl.framework.com/game_server/error_code"
	"sl.framework.com/tool"
	"sl.framework.com/trace"
	"strconv"
	"strings"
	"time"
)

type MessageCmd int

const (
	/*
		聊天消息指令
	*/

	MessageCmdChat MessageCmd = 1 //聊天消息 指令

	/*
		游戏流程相关指令
	*/

	MessageCmdBetStart    MessageCmd = 2  //投注开始 指令 (开局)
	MessageCmdBetStop     MessageCmd = 3  //投注结束 指令
	MessageCmdGameDraw    MessageCmd = 4  //开奖 指令 (局结束)
	MessageCmdGameData    MessageCmd = 8  //游戏数据，如：发牌，骰宝开牌，轮盘结果等等
	MessageCmdGamePause   MessageCmd = 9  //游戏暂停,如换靴，维护等
	MessageCmdCancelRound MessageCmd = 10 //取消局
	MessageCmdChangeCard  MessageCmd = 11 //改牌
	MessageCmdGameEnd     MessageCmd = 25 //游戏结束，本局结束

	/*
		用户游戏行为相关指令
	*/

	MessageCmdBet        MessageCmd = 12 //投注 指令
	MessageCmdBetCancel  MessageCmd = 14 //取消投注 指令
	MessageCmdBetReceipt MessageCmd = 26 //投注小票 指令
)

/**
 * ToString
 * 消息命令转化为字符串
 *
 * @param
 * @return string - 消息命令字符串
 */

func (c MessageCmd) ToString() string {
	cmd := "unknown message command"
	switch c {
	case MessageCmdChat:
		cmd = "MessageCmdChat"
	case MessageCmdBetStart:
		cmd = "MessageCmdBetStart"
	case MessageCmdBetStop:
		cmd = "MessageCmdBetStop"
	case MessageCmdGameDraw:
		cmd = "MessageCmdGameDraw"
	case MessageCmdGameData:
		cmd = "MessageCmdGameData"
	case MessageCmdGamePause:
		cmd = "MessageCmdGamePause"
	case MessageCmdCancelRound:
		cmd = "MessageCmdCancelRound"
	case MessageCmdChangeCard:
		cmd = "MessageCmdChangeCard"
	case MessageCmdGameEnd:
		cmd = "MessageCmdGameEnd"
	case MessageCmdBet:
		cmd = "MessageCmdBet"
	case MessageCmdBetCancel:
		cmd = "MessageCmdBetCancel"
	case MessageCmdBetReceipt:
		cmd = "MessageCmdBetReceipt"
	}
	return cmd
}

// HttpHeaderTag Http头类型
type HttpHeaderTag string

const (
	httpHeaderTagContentType HttpHeaderTag = "content-type"
	httpHeaderTagTraceId     HttpHeaderTag = "traceId"
	httpHeaderTagGameId      HttpHeaderTag = "Game-Id"
	httpHeaderToken          HttpHeaderTag = "token"
	httpHeaderRequestId      HttpHeaderTag = "Request-Id"
)

// beegoServerVersion beego服务版本
const (
	beegoServerVersion = "beegoServer-v2.3.1"
	contentTypeJson    = "application/json"
)

func httpSet(traceId, requestId string, request *httplib.BeegoHTTPRequest) *httplib.BeegoHTTPRequest {
	request.Header(string(httpHeaderTagContentType), contentTypeJson)
	request.Header(string(httpHeaderTagTraceId), traceId)
	request.Header(string(httpHeaderTagGameId), strconv.FormatInt(conf.GetGameId(), 10))
	request.Header(string(httpHeaderToken), conf.ServerConf.AgentToken)
	if len(requestId) > 0 {
		request.Header(string(httpHeaderRequestId), requestId)
	}
	request.SetUserAgent(beegoServerVersion)
	return request.SetTimeout(conf.GetHttpConnectTimeout(), conf.GetHttpReadWriteTimeout())
}

// 封装http get方法,统一设置http头 超时时间
func httpGet(traceId, requestId, url string) *httplib.BeegoHTTPRequest {
	return httpSet(traceId, requestId, httplib.Get(url))
}

// 封装http post方法,统一设置http头 超时时间
func httpPost(traceId, requestId string, url string) *httplib.BeegoHTTPRequest {
	return httpSet(traceId, requestId, httplib.Post(url))
}

// 封装http put方法,统一设置http头 超时时间
func httpPut(traceId, requestId string, url string) *httplib.BeegoHTTPRequest {
	return httpSet(traceId, requestId, httplib.Put(url))
}

// exceptionResponse 能力中心异常结果结构
type exceptionResponse struct {
	code string
	data string
	msg  string
}

// 检测能力中心返回包是否有异常
func checkExceptionResponse(traceId string, resp []byte) int {
	var (
		err           error
		strResp       = string(resp)
		exceptionResp exceptionResponse
	)
	msgHeader := fmt.Sprintf("checkExceptionResponse traceId=%v", traceId)

	if strings.Contains(strResp, "msg") && strings.Contains(strResp, "data") &&
		strings.Contains(strResp, "code") {
		if err = json.Unmarshal(resp, &exceptionResp); nil != err {
			trace.Error("%v, unmarshal failed, error=%v", msgHeader, err.Error())
			return errcode.JsonErrorUnMarshal
		} else {
			retCode := errcode.HttpErrorDataFailed
			switch exceptionResp.code {
			case "1406": //错误码1406:The draw result repeat 平台错误码是1406时不重试
				retCode = errcode.ErrorOk
			}
			trace.Error("%v, code=%v, msg=%v strResp=%v", msgHeader, exceptionResp.code, exceptionResp.msg, strResp)
			return retCode
		}
	}
	return errcode.ErrorOk
}

type FuncOnTry func() int

// http重试部分的封装
func httpRunOnRetry(fn FuncOnTry) (ret int) {
	retryTime, retryInterval := conf.GetHttpRetryInfo()
	for i := 0; i < retryTime; i++ {
		if ret = fn(); errcode.ErrorOk != ret {
			time.Sleep(time.Duration(retryInterval) * time.Millisecond)
			continue
		}
		break
	}
	return
}

/**
 * runHttpGet
 * 发起http get请求，其中对http request进行封装，错误重试，数据解析等
 *
 * @param traceId string - traceId 用于日志跟踪
 * @param msg string - msg 打印消息
 * @param url string - http请求的url
 * @param receiver interfaces{} - 数据反序列化接受者 必须为指针类型
 * @return int - 请求结果 成功则返回errcode.ErrorOk
 */

func runHttpGet(traceId, msg, url string, receiver interface{}) int {
	var err error
	pDog := tool.NewWatcher(msg)

	fn := func() int {
		var data []byte
		requestId := strconv.FormatInt(snowflaker.GetSnowFlakeInstance().GetUniqueId(), 10)
		req := httpGet(traceId, requestId, url)
		if data, err = req.Bytes(); nil != err {
			trace.Error("runHttpGet %v, http get error=%v", msg, err.Error())
			return errcode.HttpErrorDataFailed
		}
		if len(data) <= 0 {
			trace.Notice("runHttpGet %v no data", msg)
			return errcode.ErrorOk
		}
		if ret := checkExceptionResponse(traceId, data); errcode.ErrorOk != ret {
			trace.Error("runHttpGet %v, platform internal error, ret code=%v", msg, ret)
			return ret
		}

		if err = json.Unmarshal(data, receiver); nil != err {
			trace.Error("runHttpGet %v json unmarshal failed, data=%v, error=%v", msg, string(data), err.Error())
			return errcode.JsonErrorUnMarshal
		}
		trace.Info("runHttpGet %v get data success", msg)

		return errcode.ErrorOk
	}
	code := httpRunOnRetry(fn)

	pDog.Stop()
	return code
}

/**
 * runHttpPut
 * 发起http put请求，其中对http request进行封装，错误重试，数据解析等
 *
 * @param traceId string - traceId 用于日志跟踪
 * @param msg string - msg 打印消息
 * @param url string - http请求的url
 * @param sender interfaces{} - http发送的信息 放在http body中, 如果为nil则不设置http body
 * @return int - 请求结果 成功则返回errcode.ErrorOk
 */

func runHttpPut(traceId, msg, url string, sender interface{}) int {
	var (
		data []byte
		err  error
	)
	pDog := tool.NewWatcher(msg)

	if sender != nil {
		//序列化发送的数据
		switch t := sender.(type) {
		case string: //字符串原样发送 否则使用marshal函数序列化会有转义字符
			data = []byte(t)
		default:
			if data, err = json.Marshal(sender); err != nil {
				trace.Error("runHttpPut %v, json marshal failed, error=%v", msg, err.Error())
				return errcode.JsonErrorMarshal
			}
		}
	}

	//发送数据
	fn := func() int {
		var respData []byte
		requestId := strconv.FormatInt(snowflaker.GetSnowFlakeInstance().GetUniqueId(), 10)
		req := httpPut(traceId, requestId, url)
		if len(data) > 0 {
			req.Body(data) //有数据则放入body中
		}

		//发送请求并得到回包
		if respData, err = req.Bytes(); nil != err {
			trace.Error("runHttpPut %v, http post failed, error=%v, response data=%v",
				msg, err.Error(), string(respData))
			return errcode.HttpErrorDataFailed
		}
		if ret := checkExceptionResponse(traceId, respData); errcode.ErrorOk != ret {
			trace.Notice("runHttpPut %v, platform internal error, response data=%v, code=%v",
				msg, ret, string(respData))
			return ret
		}

		trace.Info("runHttpPut %v put data success, len(data)=%v, data=%+v", msg, len(data), string(data))
		return errcode.ErrorOk
	}
	code := httpRunOnRetry(fn)

	pDog.Stop()
	return code
}

/**
 * runHttpPost
 * 发起http post请求，其中对http request进行封装，错误重试，数据解析等
 *
 * @param traceId string - traceId 用于日志跟踪
 * @param msg string - msg 打印消息
 * @param url string - http请求的url
 * @param sender interfaces{} - http发送的信息 放在http body中, 如果为nil则不设置http body
 * @param receive interfaces{} - http post返回消息的接收接口 必须为指针类型 如果没有返回值或不解析返回值则传nil
 * @return int - 请求结果 成功则返回errcode.ErrorOk
 */

func runHttpPost(traceId, msg, url string, sender interface{}, receiver interface{}) int {
	var (
		data []byte
		err  error
	)
	pDog := tool.NewWatcher(msg)

	if sender != nil {
		//序列化发送的数据
		switch t := sender.(type) {
		case string: //字符串原样发送 否则使用marshal函数序列化会有转义字符
			data = []byte(t)
		default:
			if data, err = json.Marshal(sender); err != nil {
				trace.Error("runHttpPost %v, json marshal failed, error=%v", msg, err.Error())
				return errcode.JsonErrorMarshal
			}
		}
	}

	//发送数据
	fn := func() int {
		var respData []byte
		requestId := strconv.FormatInt(snowflaker.GetSnowFlakeInstance().GetUniqueId(), 10)
		req := httpPost(traceId, requestId, url)
		if len(data) > 0 {
			req.Body(data) //有数据则放入body中
		}

		//发送请求并得到回包
		if respData, err = req.Bytes(); nil != err {
			trace.Error("runHttpPost %v, http post failed, error=%v, response data=%v",
				msg, err.Error(), string(respData))
			return errcode.HttpErrorDataFailed
		}
		if ret := checkExceptionResponse(traceId, respData); errcode.ErrorOk != ret {
			trace.Notice("runHttpPost %v, platform internal error, response data=%v, code=%v send=%+v",
				msg, ret, string(respData), sender)
			return ret
		}

		//解析回包中的数据
		if receiver != nil {
			if err = json.Unmarshal(respData, receiver); nil != err {
				trace.Error("runHttpPost %v json unmarshal failed, error=%v", msg, err.Error())
				return errcode.JsonErrorUnMarshal
			}
		}

		trace.Info("runHttpPost %v post data response success, len(respData)=%v",
			msg, len(respData))
		trace.Debug("runHttpPost %v post data response success, len(respData)=%v,respData=%v",
			msg, len(respData), string(respData))
		return errcode.ErrorOk
	}
	code := httpRunOnRetry(fn)

	pDog.Stop()
	return code
}
