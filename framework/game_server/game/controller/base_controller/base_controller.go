package base_controller

import (
	"encoding/json"
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"sl.framework.com/base"
	"sl.framework.com/game_server/error_code"
	"sl.framework.com/game_server/game/service/type/dto"
	"sl.framework.com/trace"
)

// Tag 变量
type Tag string

const (
	TagTraceId    Tag = "traceId"     //日志追踪Id
	TagTimestamp  Tag = "timestamp"   // 时间戳
	TagUserId     Tag = "User-Id"     //用户Id
	TagUserName   Tag = "User-Name"   //用户名
	TagUserType   Tag = "User-Type"   //用户类型
	TagLanguage   Tag = "Language"    //用户类型
	TagCurrency   Tag = "currency"    //货币类型
	TagClientType Tag = "Client-Type" //客户端类型
	TagRequestId  Tag = "Request-Id"  //请求Id
)

/*
	HttpResponse http回包统一结构
	所有与平台中心交互的接口都使用此基础类中的Response函数进行回包
*/

type HttpResponse struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// HttpStatusError 如果回包中有错误则将http status设置为400 此值已经与平台中心一起确认
const HttpStatusError = 400

/**
 * ClientResponse
 * 控制器回包
 *
 * @param code - 返回码
 * @param data - 回包数据
 * @return
 */

func (c *BaseController) DataSourceResponse(code int, traceId string, data interface{}) {
	res := HttpResponse{
		Code: fmt.Sprintf("%04d", code), Msg: errcode.GetErrMsg(code), Data: data,
	}

	// 设置 TraceID
	c.Ctx.Output.Header(string(TagTraceId), traceId)
	c.Ctx.Output.Header("Content-Type", "application/json;charset=utf-8")

	if code != errcode.ErrorOk && code != errcode.HttpStatusOK {
		c.Ctx.Output.SetStatus(HttpStatusError) // 400 错误
	}
	dataString, _ := json.Marshal(res) //null
	if err := c.Ctx.Output.Body(dataString); err != nil {
		trace.Error("请求回包报错：%v", err.Error())
	}
	trace.Debug("请求回包：traceId=%v response=%s", traceId, dataString)
}

/**
 * ClientResponse
 * 控制器回包
 *
 * @param code - 返回码
 * @param data - 回包数据
 * @return
 */

func (c *BaseController) ClientResponse(code int, traceId string, data interface{}) {
	//如果 data == nil json在序列化的时候为转为null，java对null的解析不兼容
	if base.CheckEmpty(data) {
		data = ""
	}
	res := HttpResponse{
		Code: fmt.Sprintf("%04d", code), Msg: errcode.GetErrMsg(code), Data: data,
	}

	// 设置 TraceID
	c.Ctx.Output.Header(string(TagTraceId), traceId)
	c.Ctx.Output.Header("Content-Type", "application/json;charset=utf-8")

	if code != errcode.ErrorOk && code != errcode.HttpStatusOK {
		c.Ctx.Output.SetStatus(HttpStatusError) // 400 错误
	}
	dataString, _ := json.Marshal(res) //null
	if err := c.Ctx.Output.Body(dataString); err != nil {
		trace.Error("请求回包报错：%v", err.Error())
	}
	trace.Debug("请求回包：traceId=%v response=%s", traceId, dataString)
}

/**
 * ParserFromDataSource
 * 解析来自数据源的请求header头以及参数
 *
 * @param receiver interfaces{} - http body中data的解析对象 必须为指针类型 如果为nil则不解析数据
 * @return controllerParserDTO dto.ControllerParserDTO - 解析得到的数据
 */

func (c *BaseController) ParserFromDataSource(receiver interface{}) (controllerParserDTO *dto.ControllerParserDTO) {
	controllerParserDTO = new(dto.ControllerParserDTO)
	controllerParserDTO.Code = errcode.ErrorOk

	data := c.Ctx.Input.RequestBody
	controllerParserDTO.TraceId = c.Ctx.Input.Header(string(TagTraceId))
	controllerParserDTO.Timestamp = c.Ctx.Input.Header(string(TagTimestamp))

	//用于debug不返回给函数调用者
	controllerParserDTO.RequestId = c.Ctx.Input.Header(string(TagRequestId))

	if controllerParserDTO.TraceId == "" || len(controllerParserDTO.TraceId) == 0 {
		controllerParserDTO.TraceId = controllerParserDTO.RequestId
	}

	trace.Debug("BaseController parser controllerParserDTO=%+v,"+
		"body=%v,head=%v", controllerParserDTO, string(c.Ctx.Input.RequestBody), c.Ctx.Request.Header)

	//没有接受者则直接返回不解析
	if receiver == nil {
		return
	}

	//解析body中的数据到receiver中
	if err := json.Unmarshal(data, receiver); nil != err {
		controllerParserDTO.Code = errcode.JsonErrorUnMarshal
		trace.Error("BaseController json unmarshal failed,errMsg=%+v controllerParserDTO=%+v", err.Error(),
			controllerParserDTO)
		c.DataSourceResponse(errcode.JsonErrorUnMarshal, controllerParserDTO.TraceId, nil)
		return
	}

	return
}

/**
 * ParserFromClient
 * 解析来自客户端的请求header头以及参数
 *
 * @param receiver interfaces{} - http body中data的解析对象 必须为指针类型 如果为nil则不解析数据
 * @return controllerParserDTO dto.ControllerParserDTO - 解析得到的数据
 */

func (c *BaseController) ParserFromClient(receiver interface{}) (controllerParserDTO *dto.ControllerParserDTO) {
	controllerParserDTO = new(dto.ControllerParserDTO)
	controllerParserDTO.Code = errcode.ErrorOk

	data := c.Ctx.Input.RequestBody
	controllerParserDTO.TraceId = c.Ctx.Input.Header(string(TagTraceId))
	controllerParserDTO.Timestamp = c.Ctx.Input.Header(string(TagTimestamp))

	//用于debug不返回给函数调用者
	userId := c.Ctx.Input.Header(string(TagUserId))
	userName := c.Ctx.Input.Header(string(TagUserName))
	userType := c.Ctx.Input.Header(string(TagUserType))
	language := c.Ctx.Input.Header(string(TagLanguage))
	currency := c.Ctx.Input.Header(string(TagCurrency))
	clientType := c.Ctx.Input.Header(string(TagClientType))
	controllerParserDTO.RequestId = c.Ctx.Input.Header(string(TagRequestId))

	trace.Debug("BaseController parser controllerParserDTO=%+v, userId=%v, userName=%v, userType=%v, language=%v, "+
		"currency=%v, clientType=%v,data=%v,head=%v", controllerParserDTO, userId, userName, userType, language,
		currency, clientType, string(c.Ctx.Input.RequestBody), c.Ctx.Request.Header)

	//没有接受者则直接返回不解析
	if receiver == nil {
		return
	}

	//解析body中的数据到receiver中
	if err := json.Unmarshal(data, receiver); nil != err {
		controllerParserDTO.Code = errcode.JsonErrorUnMarshal
		trace.Error("BaseController json unmarshal failed,errMsg=%+v controllerParserDTO=%+v", err.Error(),
			controllerParserDTO)
		c.ClientResponse(errcode.JsonErrorUnMarshal, controllerParserDTO.TraceId, nil)
		return
	}

	return
}

// BaseController 控制器基础类 所有与平台中心交互的接口都使用此基础类
type BaseController struct {
	beego.Controller
}
