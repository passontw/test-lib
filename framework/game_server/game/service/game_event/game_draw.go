package gameevent

import (
	"encoding/json"
	"fmt"
	"sl.framework.com/async"
	"sl.framework.com/game_server/conf"
	errcode "sl.framework.com/game_server/error_code"
	"sl.framework.com/game_server/game/service"
	types "sl.framework.com/game_server/game/service/type"
	"sl.framework.com/game_server/game/service/type/VO"
	"sl.framework.com/game_server/game/service/type/const_type"
	"sl.framework.com/game_server/game/service/type/dto"
	"sl.framework.com/game_server/mq"
	"sl.framework.com/game_server/redis/cache"
	rpcreq "sl.framework.com/game_server/rpc_client"
	"sl.framework.com/tool"
	"sl.framework.com/trace"
	"strconv"
	"time"
)

type GameDrawEvent struct {
	types.EventBase
}

/**
 * NewGameDraw
 * 创建开奖实例
 *
 * @param roomId int64 - 房间id
 * @param roundId int64 - 局id
 * @param nextRoundId int64 - 下一局id
 * @param event types.GameEventVO - 来自数据源的事件信息
 * @param roundDto *types.GameRoundDTO - 的游戏局信息
 * @return RETURN - 返回游戏事件实例
 */

func NewGameDraw(event types.GameEventVO, roundDto *types.GameRoundDTO, gameEventInitVo *VO.GameEventInitVO) *GameDrawEvent {
	return &GameDrawEvent{
		EventBase: types.EventBase{
			Dto: &types.EventDTO{
				GameRoomId:      gameEventInitVo.RoomId,
				GameRoundId:     gameEventInitVo.RoundId,
				GameId:          conf.GetGameId(),
				NextGameRoundId: gameEventInitVo.NextRoundId,
				GameRoundNo:     event.GameRoundNo,
				Command:         string(event.Command),
				Time:            event.Time,
				ReceiveTime:     event.ReceiveTime,
				Payload:         event.Payload,
			},
			RoundDTO:       roundDto,
			TraceId:        gameEventInitVo.TraceId,
			RequestId:      gameEventInitVo.RequestId,
			RetHandleEvent: gameEventInitVo.Code,
			MsgHeader: fmt.Sprintf("command=%s  traceId=%v,requestId=%v, roomId=%v, gameRoundId=%v, "+
				"nextGameRoundId=%v", event.Command, gameEventInitVo.TraceId, gameEventInitVo.RequestId, gameEventInitVo.RoomId, gameEventInitVo.RoundId, gameEventInitVo.NextRoundId),
		},
	}
}

/**
 * HandleEvent
 * 处理游戏事件函数
 *
 * @param traceId string - 跟踪id
 * @return RETURN
 */

func (e *GameDrawEvent) HandleRondEvent() {
	//直接设置为成功
	var (
		transactionList []*dto.UserTransactionDTO
		retCode         int

		settleOrderList = make([]int64, 0)
	)
	trace.Info("[游戏开奖] GameDraw %v 局信息%+v", e.MsgHeader, e.Dto.Payload)
	*e.RetHandleEvent = errcode.ErrorOk
	//
	//参数校验
	trace.Info("[游戏开奖] HandleRondEvent traceid=%v gameEvent=%v", e.TraceId, e)
	if e.Dto == nil {
		*e.RetHandleEvent = errcode.HttpErrorInvalidParam
		trace.Error("[游戏开奖] HandleRondEvent traceid=%v invalid param.", e.TraceId)
		return
	}
	pWatcher := tool.NewWatcher("ParseGameResultV2")

	//1计算游戏结果
	drawer := service.GetDrawer(e.TraceId, types.GameId(conf.GetGameId()))
	defer service.PutDrawer(types.GameId(conf.GetGameId()), drawer)
	gameResult := drawer.ParseGameResult(e.Dto)
	if gameResult == nil {
		//*e.RetHandleEvent = errcode.HttpErrorInvalidParam
		trace.Error("[游戏开奖] 解析牌局结果错误 traceid=%v.", e.TraceId)
		return
	}
	gameResult.GameRoundId = strconv.FormatInt(e.Dto.GameRoundId, 10)
	gameResult.Timestamp = tool.Current()
	pWatcher.Stop()
	trace.Info("[游戏开奖] 计算游戏结果 ParseGameResultV2 traceid=%v gameResult=%+v", e.TraceId, gameResult)

	//2.推送游戏开奖结果 token
	pWatcher.Start("推送开奖结果")
	trace.Info("[游戏开奖] 推送游戏开奖结果到中台 traceid=%v gameResult=%+v", e.TraceId, gameResult)
	if errcode.ErrorOk != rpcreq.DrawResultPost(e.TraceId, gameResult) {
		trace.Error("[游戏开奖] DrawResultPost traceid=%v gameResult=%v failed.", e.TraceId, gameResult)
		return
	}
	//答应时间差
	e.PrintTimeOffset()
	//推送游戏开奖结果
	rpcreq.AsyncSendRoundMessage(e.TraceId, strconv.FormatInt(e.Dto.GameRoomId, 10), e.RoundDTO.Id, string(types.GameEventCommandGameDraw), gameResult)
	pWatcher.Stop()

	//写入缓存
	pWatcher.Start("游戏结果缓存")
	gameResultCache := &cache.SettleCache{TraceId: e.TraceId, RoomId: e.Dto.GameRoomId}
	gameResultCache.Set(gameResult)
	pWatcher.Stop()

	//查询待开奖集合，并分组发送mq
	pWatcher.Start("查询待开奖集合")

	//异步调用具体游戏服接口批量入库 避免具体游戏服数据库写入操作耗时太久而阻塞游戏框架流程
	dbGet := service.NewGameDBSaver(e.TraceId, types.GameId(e.Dto.GameId))
	if dbGet == nil {
		//*e.RetHandleEvent = errcode.HttpErrorInvalidParam
		trace.Error("[游戏开奖] 获取游戏数据库失败 NewGameDBSaver %v, new game order saver interfaces failed", e.TraceId)
		return
	}
	trace.Info("[游戏开奖] 查询待开奖集合，并分组发送mq traceId=%v gameRoomId=%v,gameRoundId=%v", e.TraceId, e.Dto.GameRoomId, e.Dto.GameRoundId)
	orderList := dbGet.GetOrderNoList(e.TraceId, e.Dto.GameRoomId, e.Dto.GameRoundId, e.Dto.GameRoundNo)
	if len(orderList) == 0 {
		//*e.RetHandleEvent = errcode.HttpErrorInvalidParam
		trace.Notice("[游戏开奖] 从数据库查询全量注单为空  traceId:%v gameRoomId=%v,gameRoundId=%v", e.TraceId, e.Dto.GameRoomId, e.Dto.GameRoundId)
		return
	}
	//校验注单扣款
	queryTransactionList := &dto.QueryTransactionDTO{BeginTime: time.Now().UnixMilli() - 3600*1000, EndTime: time.Now().UnixMilli(), OrderNoList: orderList}
	trace.Debug("[游戏开奖] 校验注单扣款 入参校验 queryTransactionList OrderNoList=%+v", queryTransactionList)
	if transactionList, retCode = rpcreq.GetTransactionList(e.TraceId, queryTransactionList); retCode != errcode.ErrorOk {
		trace.Error("[游戏开奖] 校验注单是否已经成功扣款失败 traceId:%v 错误码:%v", e.TraceId, retCode)
		//*e.RetHandleEvent = errcode.HttpErrorInvalidParam
		return
	}
	if len(transactionList) == 0 {
		//*e.RetHandleEvent = errcode.HttpErrorInvalidParam
		trace.Error("[游戏开奖] 校验注单是否已经成功扣款成功注单数量为空 traceId:%v", e.TraceId)
		return
	}
	//获取注单缓存
	orderAllList := cache.GetOrders(e.TraceId, strconv.FormatInt(e.Dto.GameRoomId, 10), strconv.FormatInt(e.Dto.GameRoundId, 10))
	if orderAllList == nil || len(orderAllList) == 0 {
		trace.Notice("[游戏开奖] 获取当局全量注单失败，traceId=%v,gameRoomId=%v,gameRoundId=%v", e.TraceId, e.Dto.GameRoomId, e.Dto.GameRoundId)
		*e.RetHandleEvent = errcode.RedisErrorGet
		return
	}
	//把注单放到map中
	betDtoList := make(map[int64]*dto.BetDTO)
	for _, betDTO := range orderAllList {
		betDtoList[betDTO.OrderNo] = betDTO
	}
	//从orderList获取真正扣款的注单用于分票
	for _, transaction := range transactionList {
		nOrderNo, _ := strconv.ParseInt(transaction.OrderNo, 10, 64)
		if transaction.Status == string(const_type.TransactionStatusSuccess) {
			settleOrderList = append(settleOrderList, nOrderNo)
		} else {
			trace.Notice("[游戏开奖] 校验注单扣款 状态异常注单 traceId=%v, gameRoomId=%v, gameRoundId=%v, transaction=%+v",
				e.TraceId, e.Dto.GameRoomId, e.Dto.GameRoundId, transaction)
		}
		betDtoList[nOrderNo].BetStatus = transaction.Status
		trace.Debug("[游戏开奖] 校验注单扣款 traceId=%v gameRoomId=%v,gameRoundId=%v totalBetNum=%v,settleBetNum=%v",
			e.TraceId, e.Dto.GameRoomId, e.Dto.GameRoundId, len(betDtoList), len(settleOrderList))
	}
	pWatcher.Stop()
	//更新完注单状态后重新更新缓存
	pWatcher.Start("开奖更新注单缓存")
	cache.SetOrders(e.TraceId, strconv.FormatInt(e.Dto.GameRoomId, 10), strconv.FormatInt(e.Dto.GameRoundId, 10), orderAllList)
	pWatcher.Stop()

	trace.Debug("[游戏开奖] 查询待开奖集合，并分组发送mq traceId=%v gameRoomId=%v,gameRoundId=%v transactionList=%v settleOrderList=%v",
		e.TraceId, e.Dto.GameRoomId, e.Dto.GameRoundId, transactionList, settleOrderList)

	pWatcher.Start("分片发送MQ")
	//3.获取分片大小
	patchSize := conf.ServerConf.Common.DrawSize
	trace.Info("[游戏开奖] 分片并发送MQ traceId:%v patchSize:%v settleOrderList:%+v", e.TraceId, patchSize, settleOrderList)
	patches := tool.SplitList[int64](settleOrderList, patchSize)
	//遍历
	for _, row := range patches {
		if len(row) == 0 {
			continue
		}
		gameDrawDataDTOItem := types.GameDrawDataDTO{
			GameRoomId:         e.Dto.GameRoomId,
			GameRoundId:        e.Dto.GameRoundId,
			GameId:             e.Dto.GameId,
			GameRoundNo:        e.Dto.GameRoundNo,
			GameRoundResultDTO: *gameResult,
			OrderList:          row,
		}
		messageStr, err := generateGameDrawMessage(e.TraceId, gameDrawDataDTOItem)
		//id := tool.GenerateRandomString(32)
		trace.Info("[游戏开奖] 分片并发送MQ traceId:%v 分片数组大小 :%v patchSize:%v messageStr:%v", e.TraceId, len(patches), patchSize, messageStr)
		if err != nil {
			trace.Error("[游戏开奖]  生成开奖MQ消息失败 generateGameDrawMessage traceId:%v failed.", e.TraceId)
		} else {
			topic := generateTopic()
			createTime := strconv.FormatInt(time.Now().Unix(), 10)
			fn := func() {
				trace.Info("[游戏开奖] 分片并发送MQ  异步发送HandleRondEvent traceId:%v dispatch orders slice to gameSvr topic:%v messageStr:%v", e.TraceId, topic, messageStr)
				mq.SendMessage(topic, strconv.FormatInt(e.Dto.GameId, 10), e.TraceId, createTime, messageStr)
			}
			async.AsyncRunCoroutine(fn)
		}

	}
	pWatcher.Stop()
	return
}

/**
 * generateGameDrawMessage
 * 生成开奖MQ消息
 *
 * @param PARAM - 参数说明
 * @return RETURN - 返回值说明
 */

func generateGameDrawMessage(traceId string, dto types.GameDrawDataDTO) (string, error) {
	var (
		err        error
		messageBuf []byte
		messageStr string
	)
	gameDrawMessage := types.GameDrawMessage{
		GameRoomId:         dto.GameRoomId,
		GameRoundId:        dto.GameRoundId,
		GameRoundNo:        dto.GameRoundNo,
		GameId:             dto.GameId,
		GameRoundResultDTO: dto.GameRoundResultDTO,
		OrderList:          dto.OrderList,
	}

	trace.Info("[生成结算消息] generateGameDrawMessage traceId=%v gameDrawMessage=%+v.", traceId, gameDrawMessage)
	if messageBuf, err = json.Marshal(gameDrawMessage); err != nil {
		trace.Error("[生成结算消息] generateGameDrawMessage traceId=%v 序列化 messageDto=%+v 失败.", traceId, gameDrawMessage)
		return messageStr, err
	}

	return string(messageBuf), err
}

/**
 * FunctionName
 * 函数功能描述
 *
 * @param PARAM - 参数说明
 * @return RETURN - 返回值说明
 */

func generateTopic() string {
	var result string
	rndInt := tool.GenerateRandomRange(1, 9999) % 10
	switch rndInt {
	case 0:
		result = string(mq.TopicGameDrawOut0)
	case 1:
		result = string(mq.TopicGameDrawOut1)
	case 2:
		result = string(mq.TopicGameDrawOut2)
	case 3:
		result = string(mq.TopicGameDrawOut3)
	case 4:
		result = string(mq.TopicGameDrawOut4)
	case 5:
		result = string(mq.TopicGameDrawOut5)
	case 6:
		result = string(mq.TopicGameDrawOut6)
	case 7:
		result = string(mq.TopicGameDrawOut7)
	case 8:
		result = string(mq.TopicGameDrawOut8)
	case 9:
		result = string(mq.TopicGameDrawOut9)
	default:
		result = string(mq.TopicGameDrawOut0)
	}
	return result
}

/**
 * PrintTimeOffset
 * 打印事件时间差
 *
 * @param PARAM - 参数说明
 * @return RETURN - 返回值说明
 */

func (e *GameDrawEvent) PrintTimeOffset() {
	//发送到ws的时间
	sendTimestamp := time.Now().UnixMilli()
	sendTime := tool.FormatTime(sendTimestamp)
	nEventReceiveOffset := e.Dto.ReceiveTime - e.Dto.Time
	eventTime := tool.FormatTime(e.Dto.Time)
	receiveTime := tool.FormatTime(e.Dto.ReceiveTime)
	nSendReceiveOffset := sendTimestamp - e.Dto.ReceiveTime

	trace.Notice("[数据源时间转发ws] traceId:%v,事件类型：%v\r\n,事件接收时间：%v,时间发生时间：%v\r\n,事件发生到接收时间差：%v毫秒\r\n,发送时间：%v,接收和发送时差：%v毫秒\r\n",
		e.TraceId, e.Dto.Command, receiveTime, eventTime, nEventReceiveOffset, sendTime, nSendReceiveOffset)
}
