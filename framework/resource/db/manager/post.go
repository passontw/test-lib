/*
 * @Author: [Malfoy]
 * @Date: 2024-09-19
 * @Last Modified by: [Malfoy]
 * @Last Modified time: 2024-09-19
 * @Description: 该文件负责订阅者状态检查、请求的处理与保存、重试机制及相关的队列操作。
 *               主要功能包括：
 *               1. 通过 `checkSubscriberOnline` 函数定期检查订阅者是否在线，并根据返回的状态更新订阅者状态。
 *               2. 通过 `dispatcherQueueData` 函数处理不同队列中的请求数据，支持将处理失败的请求放入重试队列。
 *               3. 使用 Redis 锁确保并发情况下的请求不会重复处理，保护数据的一致性。
 *               4. 通过 `fasthttp` 库实现高效的 HTTP 请求发送和响应处理。
 *               5. 支持请求失败后的重试逻辑，超过最大重试次数后，记录请求为失败。
 *               6. 所有请求和响应的状态被持久化到 MySQL 数据库，记录详细的请求信息和状态。
 *
 * @Dependencies:
 *               - Beego ORM: 处理数据库的增删改查操作
 *               - fasthttp: 用于高效的 HTTP 请求和响应处理
 *               - Redis: 用于分布式锁操作，防止请求的重复处理
 *               - MySQL: 数据库存储系统，用于存储请求和响应数据
 *               - Go 标准库: 提供如时间、JSON 编解码和并发等功能
 *
 * @Usage:
 *               1. 定期执行订阅者状态检查，确保每个订阅者的在线状态保持最新。
 *               2. 当新的请求到达时，通过队列处理机制进行异步处理，并记录请求的处理结果。
 *               3. 处理请求时使用 Redis 锁机制，确保在分布式环境下的数据一致性。
 *               4. 提供重试机制，在请求失败时能够重试处理，并将失败的请求标记为 "pending" 或 "failed" 状态。
 *
 * @Attention:
 *               - `retryingQueue` 长度有限，如果超出最大长度，重试队列将无法接受更多请求，需要调整系统配置。
 *               - Redis 锁的有效时间设置为 3 秒，可能需要根据实际场景调整。
 *               - 在并发情况下使用锁机制来确保不会有多个节点处理同一请求，避免并发冲突。
 *               - 日志记录使用 `sl.framework.com/trace` 库，记录详细的错误和状态变化。
 */

package mgr

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"sl.framework.com/resource/cache"
	fdb "sl.framework.com/resource/db"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/valyala/fasthttp"
	"sl.framework.com/trace"
)

// CheckRsp 定义从服务器返回的检查响应结构体
type CheckRsp struct {
	Msg  string `json:"msg"`
	Data struct {
		Countdown int64 `json:"countdown"`
	} `json:"data"`
	Code string `json:"code"`
}

// checkSubscriberOnline 检查所有订阅者的在线状态并更新状态
func checkSubscriberOnline() {
	c.RLock() // 读取锁，防止并发读写
	defer c.RUnlock()
	trace.Info("checkMapUrl: %d", len(c.checkMapUrl))

	// 遍历所有订阅者的URL进行状态检查
	for k, v := range c.checkMapUrl {
		check := func() {
			req := fasthttp.AcquireRequest()
			req.SetRequestURI(v)
			req.Header.SetMethod("GET")
			resp := fasthttp.AcquireResponse()

			defer fasthttp.ReleaseResponse(resp)
			defer fasthttp.ReleaseRequest(req)

			// 发送请求并检查状态
			err := fasthttp.Do(req, resp)
			if err != nil || resp.StatusCode() != fasthttp.StatusOK {
				// 错误处理，如果是400错误则解析响应体
				if resp.StatusCode() == fasthttp.StatusBadRequest {
					var x CheckRsp
					if err = json.Unmarshal(resp.Body(), &x); err != nil {
						trace.Error("Request unmarshal failed game_room_id: %d, checkUrl: %s, Method: GET, Error: %v, ResponseCode: 400", k, v, err)
					}
				} else {
					trace.Info("online state check info: game_room_id: %d, checkUrl: %s set to offline", k, v)
					subscriberStatus.Store(k, false) // 设置为离线状态
				}
			} else {
				// 如果请求成功，解析响应体并更新状态
				var x CheckRsp
				if err = json.Unmarshal(resp.Body(), &x); err != nil {
					trace.Error("Request unmarshal failed game_room_id: %d, checkUrl: %s, Method: GET, Error: %v, ResponseCode: 400", k, v, err)
					return
				}
				betTimeMap.Store(k, x.Data.Countdown) // 存储倒计时信息
				trace.Info("online state check info: game_room_id: %d, checkUrl: %s set to online, betime: %d", k, v, x.Data.Countdown)
			}
		}
		check() // 执行检查函数
	}
}

// saveRequest 保存请求信息到数据库
// 参数：Gmcode：游戏代码，traceID：请求唯一标识，requestBody：请求体，endpoint：目标端点，method：请求方法，eid：端点ID
func saveRequest(Gmcode, traceID, requestBody, endpoint, method string, eid int64) error {
	httpRequest := fdb.HttpPostRequests{
		Gmcode:      Gmcode,
		EndpointId:  eid,
		Endpoint:    endpoint,
		Method:      method,
		RequestBody: requestBody,
		RequestTime: time.Now().UTC(),
		TraceId:     traceID,
		Status:      "pending", // 初始状态为 pending
	}
	err := httpRequest.Insert() // 插入数据库
	if err != nil {
		trace.Error("error saving request, TraceID: %s, Error: %v", traceID, err)
		return err
	}
	trace.Info("request saved, TraceID: %s", traceID)
	return nil
}

// dispatcherQueueData 分发队列中的数据到不同地方处理
// 参数：
//   - fromQueue：来源队列
//   - targetQueue：目标队列
func dispatcherQueueData(fromQueue, targetQueue chan *RequestData) {
	for requestData := range fromQueue {
		// 处理请求，如果fromQueue是初次处理(incomingQueue)，保存请求
		if err := processOrSaveRequest(requestData, fromQueue == incomingQueue); err != nil {
			trace.Warning("error processing request, TraceID: %s, Error: %v", requestData.TraceID, err)
			continue
		}
		// 如果fromQueue是初次处理(incomingQueue)，将请求推送到下一个队列，等待请求(请求会在另一个协程中进入上一段代码的请求逻辑)
		if fromQueue == incomingQueue {
			pushToNextQueue(requestData, targetQueue)
		}
	}
}

// processOrSaveRequest 首次保存请求入库 否则发起请求
// 参数：
//   - requestData：请求数据
//   - isFirstReq：是否为初次尝试
func processOrSaveRequest(requestData *RequestData, isFirstReq bool) error {
	if isFirstReq {
		return saveRequest(requestData.GMCode, requestData.TraceID, requestData.Body, requestData.Endpoint, requestData.Method, requestData.EndpointId)
	}
	return processRequest(requestData)
}

// processRequest 处理具体的请求逻辑，包含重试机制
// 参数：
//   - data：请求数据
func processRequest(data *RequestData) error {
	// 检查GMCode
	if data.GMCode == "" {
		return fmt.Errorf("GMCode is null")
	}
	// 检查请求重试次数
	if data.TryTimes > data.TryMaxTimes {
		return fmt.Errorf("tryTimes: %d more than TryMaxTimes: %d", data.TryTimes, data.TryMaxTimes)
	}
	isOnlineAny := isOnline(data.EndpointId) // 检查订阅者是否在线
	if !isOnlineAny {
		trace.Warn("subscriber offline, saving request as pending, TraceID: %s", data.TraceID)
		return recordFailure(data.TraceID, data.GMCode, "pending", "", 0)
	}

	// 抢锁，确保并发处理安全
	key := fmt.Sprintf("%s-traceid-%s", cache.RedisKeyPrefix, data.TraceID)
	lockSuccess, err := cache.Get().SetNX(key, "1", 3*time.Second)
	if err != nil || !lockSuccess {
		return fmt.Errorf("get redis key: [%s] lock failed: %v", key, err)
	}
	defer cache.Get().Delete(key) // 释放锁
	trace.Info("get redis key: [%s] lock success", key)

	// 发送请求
	req, resp := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)
	setupRequest(req, data) // 配置请求
	reqTime := time.Now().UTC()
	if err = fasthttp.Do(req, resp); err != nil || resp.StatusCode() != fasthttp.StatusOK {
		respCode := resp.StatusCode()
		body := resp.Body()
		// 错误处理，如果是400错误则解析响应体
		if respCode == fasthttp.StatusBadRequest {
			trace.Error("body=%s", body)
			isWorking(data.GMCode[1:5])
			var x CheckRsp
			if err = json.Unmarshal(body, &x); err != nil {
				trace.Error("request unmarshal failed, Method: POST, Error: %v, ResponseCode: 400", err)
			}
		} else {
			trace.Info("unexpected exceptions, ResponseCode: %d, err: %v", respCode, err)
		}
		return handleFailure(data, respCode, string(body), err)
	}
	// 成功处理，设置订阅者状态为在线
	subscriberStatus.Store(data.EndpointId, true)
	return saveRequestResult(reqTime, data, resp.StatusCode(), string(resp.Body()))
}

// setupRequest 设置请求参数
// 参数：
//   - req：请求
//   - data：请求数据
func setupRequest(req *fasthttp.Request, data *RequestData) {
	req.SetRequestURI(data.Endpoint)
	req.Header.SetMethod(data.Method)
	req.Header.SetContentType("application/json")
	req.Header.Set("Request-Id", data.TraceID)
	req.SetBody([]byte(data.Body))
}

// handleFailure 处理请求失败情况，记录失败信息并触发重试
// 参数：
//   - data：请求数据
//   - statusCode：响应状态码
//   - err：错误信息
func handleFailure(data *RequestData, statusCode int, body string, err error) error {
	trace.Error("request failed, TraceID: %s, ResponseCode: %d, body: %s, postUrl: %s, error: %v", data.TraceID, statusCode, body, data.Endpoint, err)
	var ok bool
	if err != nil {
		var dnsErr *net.DNSError
		ok = errors.As(err, &dnsErr)
		if ok && dnsErr.IsNotFound {
			statusCode = fasthttp.StatusBadRequest
			body = dnsErr.Error()
		}
	}
	if data.TryTimes == int64(maxTriedTimes) {
		_ = recordFailure(data.TraceID, data.GMCode, "failed", body, statusCode)
	} else {
		data.TryTimes++
		_ = recordFailure(data.TraceID, data.GMCode, "retrying", body, statusCode)
	}
	if ok {
		return nil
	}
	if data.Endpoint != "" {
		return retryRequest(data)
	}
	return nil
}

// retryRequest 将请求推入重试队列
// 参数：
//   - data：请求数据
func retryRequest(data *RequestData) error {
	select {
	case retryingQueue <- data:
		trace.Info("retryingQueue=%d", len(retryingQueue))
		return nil
	default:
		return recordFailure(data.TraceID, data.GMCode, "failed", `{"msg":"重试队列已满","data":null,"code":"9999"}`, 0)
	}
}

// saveRequestResult 保存请求结果到数据库
// 参数：
//   - reqTime：请求时间
//   - data：请求数据
//   - responseCode：响应码
//   - responseBody：响应内容
func saveRequestResult(reqTime time.Time, data *RequestData, responseCode int, responseBody string) error {
	httpRequest := &fdb.HttpPostRequests{
		ResponseCode: responseCode,
		ResponseBody: responseBody,
		RequestTime:  reqTime,
		ResponseTime: time.Now().UTC(),
		TraceId:      data.TraceID,
		Gmcode:       data.GMCode,
		Status:       "success",
	}
	return httpRequest.Update()
}

// recordFailure 记录请求失败信息
// 参数：
//   - traceID：请求标识
//   - status：状态
//   - body：响应体
//   - statusCode：响应码
func recordFailure(traceID, gmcode, status, body string, statusCode int) error {
	httpRequest := &fdb.HttpPostRequests{
		TraceId:      traceID,
		Gmcode:       gmcode,
		Status:       status,
		ResponseCode: statusCode,
		ResponseTime: time.Now().UTC(),
		ResponseBody: body,
	}
	return httpRequest.RecordFailure()
}

// pushToNextQueue 将请求推送到下一个队列，或者保存为pending
// 参数：
//   - data：请求数据
//   - targetQueue：目标队列
func pushToNextQueue(data *RequestData, targetQueue chan *RequestData) {
	select {
	case targetQueue <- data: // 推送成功，继续处理
	default:
		trace.Error("target queue is full, saving to DB. TraceID: %s", data.TraceID)
		_ = recordFailure(data.TraceID, data.GMCode, "pending", "", 0) // 记录失败
	}
}
