package listenner

import (
	"errors"
	"fmt"
	"sl.framework.com/game_server/error_code"
	"sl.framework.com/game_server/game/dao/redisdb"
	"sl.framework.com/game_server/game/service"
	gameevent "sl.framework.com/game_server/game/service/game_event"
	"sl.framework.com/game_server/game/service/type/VO"
	"sl.framework.com/game_server/game/service/type/dto"
	"sl.framework.com/game_server/redis/cache"

	types "sl.framework.com/game_server/game/service/type"
	rediskey "sl.framework.com/game_server/redis/rediskey"
	rpcreq "sl.framework.com/game_server/rpc_client"
	err "sl.framework.com/resource/error"
	"sl.framework.com/tool"
	"sl.framework.com/trace"
	"strconv"
)

/**
 * DispatchGameEventV2
 * 分发游戏事件，对具体的游戏事件进行处理
 *
 * @param traceId string - 用于日志跟踪
 * @param header types.GameEventMessageHeader - 处理事件的信息头
 * @param result *int - 处理结果
 * @return
 */

func DispatchGameEventV2(parserDto *dto.ControllerParserDTO, event types.GameEventVO, result *int) {

	//参数校验
	if len(parserDto.TraceId) == 0 || event.GameRoomId == 0 || len(event.GameRoundNo) == 0 || len(event.Command) == 0 {
		trace.Error("分派游戏事件 DispatchGameEvent invalid param is invalid traceId=%v,gameRoomId=%v,roundNo=%v,command=%v.",
			parserDto.TraceId, event.GameRoomId, event.GameRoundNo, event.Command)
		return
	}
	if tool.IsEmpty(event.Payload) {
		trace.Info("分派游戏事件 DispatchGameEvent skip traceId=%v, game event=%v payload is empty", parserDto.TraceId, event.Command)
	}
	//分布式锁
	gameEventRedisLockInfo := rediskey.GetGameEventLockRedisInfo(parserDto.RequestId, event.GameRoundNo, string(event.Command))
	if !redisdb.TryLock(gameEventRedisLockInfo) {
		trace.Error("分派游戏事件 DispatchGameEvent traceId=%v, redis lock key=%v, lock failed", parserDto.TraceId, gameEventRedisLockInfo.Key)
		*result = int(errcode.GameErrorGameEventExist)
		return
	}
	//defer redisdb.Unlock(gameEventRedisLockInfo) //同一个局号，同一个消息需要锁住不能释放等过期 数据源是多节点发送同一个消息 这里做幂等
	//查询局信息
	pWatcher := tool.NewWatcher("获取局信息")

	trace.Info("分派游戏事件 调用中台接口查询房间信息 traceId=%v, event:%+v command=%v", parserDto.TraceId, event, event.Command)
	roundDTO := GetRoundInfo(parserDto.TraceId, event.GameRoundNo, event.GameRoomId)
	pWatcher.Stop()
	//尝试关闭异常局

	pWatcher.Start("尝试关闭局")
	gameRoundId, _ := strconv.ParseInt(roundDTO.Id, 10, 64)
	if types.GameEventCommandGameDraw == event.Command {
		trace.Info("分派游戏事件 尝试关闭局 traceId=%v, event:%+v command=%v", parserDto.TraceId, event, event.Command)
		rpcreq.CloseExceptionRoundRequest(parserDto.TraceId, event.GameRoomId, gameRoundId)
	}
	pWatcher.Stop()
	//创建新的局
	pWatcher.Start(fmt.Sprintf("getNextRoundInfo handle traceId=%v", parserDto.TraceId))
	trace.Info("分派游戏事件 创建新的局 traceId=%v, event:%+v command=%v", parserDto.TraceId, event, event.Command)
	nextRoundInfo := getNextRoundInfo(parserDto.TraceId, event.GameRoundNo, parserDto.RequestId, event.GameRoomId)
	if nextRoundInfo == nil {
		trace.Error("分派游戏事件 生成下一局信息失败 traceId=%v, event:%+v command=%v", parserDto.TraceId, event, event.Command)
		*result = int(errcode.GameErrorGameEventExist)
		return
	}
	pWatcher.Stop()

	//获取局信息缓存
	pWatcher.Start("局信息缓存")
	trace.Info("分派游戏事件 获取局信息缓存 traceId=%v, event:%+v command=%v", parserDto.TraceId, event, event.Command)
	gameEventCache := cache.GameEventCache{TraceId: parserDto.TraceId, GameRoomId: strconv.FormatInt(event.GameRoomId, 10), GameRoundId: roundDTO.Id}
	gameEventCache.Data = &event
	//设置缓存并通知客户端
	trace.Info("分派游戏事件 设置缓存并通知客户端 traceId=%v, event:%+v command=%v", parserDto.TraceId, event, event.Command)
	gameEventCache.Notify()
	pWatcher.Stop()

	//游戏事件处理
	pWatcher.Start("游戏事件处理")
	roundId, _ := strconv.ParseInt(roundDTO.Id, 10, 64)
	nextRoundId, _ := strconv.ParseInt(nextRoundInfo.Id, 10, 64)
	gameEventInitVO := &VO.GameEventInitVO{
		TraceId:     parserDto.TraceId,
		RoomId:      event.GameRoomId,
		RoundId:     roundId,
		NextRoundId: nextRoundId,
		RequestId:   parserDto.RequestId,
		Code:        result,
		Time:        event.Time,
		ReceiveTime: event.ReceiveTime,
	}
	GameEvent(event, roundDTO, gameEventInitVO)
	pWatcher.Stop()
}

/**
 * GetRoundInfo
 * 私有函数 查询局信息，如果没有则初始化第一局信息 首字母小写， 内部调用，不再进行参数校验
 *
 * @param traceId string - 用于日志跟踪
 * @param gameRoundNo string -局号
 * @param gameRoomId int64 - 房间号
 * @return *types.GameRoundDTO
 */
func GetRoundInfo(traceId, gameRoundNo string, gameRoomId int64) *types.GameRoundDTO {

	roundInfo, retCode := rpcreq.GetOrBindGameRoundNo(traceId, gameRoundNo, gameRoomId)
	trace.Info("获取局信息 GetRoundInfo traceId=%v, gameRoundNo=%v,gameRoomId=%v roundDTO=%+v", traceId, gameRoundNo, gameRoomId, roundInfo)
	if retCode != err.ERR_OK {
		if initRoundInfo, errMsg := initFirstRound(traceId, gameRoundNo, gameRoomId); errMsg != nil {
			roundInfo = nil
			trace.Error("获取局信息 GetRoundInfo initFirstRound failed, traceId=%v, gameRoundNo=%v,gameRoomId=%v", traceId, gameRoundNo, gameRoomId)
		} else {
			roundInfo = initRoundInfo
			trace.Notice("获取局信息 GetRoundInfo traceId=%v, gameRoundNo=%v,gameRoomId=%v not exist,init.", traceId, gameRoundNo, gameRoomId)
		}

	}
	return roundInfo
}

/**
 * GetNextRoundInfo
 * 生成下一局信息
 *
 * @param traceId string - 用于日志跟踪
 * @param nextRoundNo string -下一局局号
 * @param gameRoomId string - 房间号
 * @return *types.GameRoundDTO
 */
func getNextRoundInfo(traceId, nextRoundNo, requestId string, gameRoomId int64) *types.GameRoundDTO {
	trace.Info("获取下一局信息 getNextRoundInfo  traceId=%v, gameRoundNo=%v,gameRoomId=%v", traceId, nextRoundNo, gameRoomId)

	//参数校验
	if len(traceId) == 0 || len(nextRoundNo) == 0 {
		trace.Notice("获取下一局信息 getNextRoundInfo invalid param")
		return nil
	}
	roundNo, _ := strconv.ParseInt(nextRoundNo, 10, 64)
	roundInfoRedisKey := rediskey.GetNextRoundIdLockRedisInfo(roundNo, requestId)
	if !redisdb.TryLock(roundInfoRedisKey) {
		trace.Error("获取下一局信息 getNextRoundInfo traceid=%v, redis lock key=%v, lock failed", traceId, roundInfoRedisKey.Key)
		return nil
	}
	defer redisdb.Unlock(roundInfoRedisKey)

	roundInfo := &types.GameRoundDTO{}
	roundInfo.GameRoomId = strconv.FormatInt(gameRoomId, 10)
	roundInfo.RoundNo = nextRoundNo
	rpcreq.BuildRoundRequest(traceId, roundInfo)
	trace.Info("获取下一局信息 getNextRoundInfo  traceId=%v, gameRoundNo=%v,gameRoomId=%v,roundInfo=%v success", traceId, nextRoundNo, gameRoomId, roundInfo)

	return roundInfo
}

func initFirstRound(traceId, gameRoundNo string, gameRoomId int64) (*types.GameRoundDTO, error) {
	roundInfo := &types.GameRoundDTO{
		GameRoomId: strconv.FormatInt(gameRoomId, 10),
		RoundNo:    gameRoundNo,
	}
	ret := errors.New("OK")
	if roundInfo == nil {
		trace.Error("DispatchGameEvent invalid traceId=%v, game event=%+v", traceId, gameRoundNo, gameRoomId)
		ret = errors.New(fmt.Sprintf("initFirstRound create failed, traceId=%v, gameRoundNo=%v,gameRoomId=%v", traceId, gameRoundNo, gameRoomId))
	}
	rpcreq.BuildRoundRequest(traceId, roundInfo)
	trace.Info("getNextRoundInfo  traceId=%v, gameRoundNo=%v,gameRoomId=%v,roundInfo=%v success", traceId, gameRoundNo, gameRoomId, roundInfo)

	return roundInfo, ret
}

/**
 * GameEvent
 * 游戏事件处理函数
 *
 * @param traceId string - traceId用于日志跟踪
 * @param roundNo int64 - 局Id
 * @param nextRoundNo int64 - 下一局Id
 * @param event types.GameEventVO - 下一局Id
 * @param result  *int - 结果
 */

func GameEvent(event types.GameEventVO, roundDto *types.GameRoundDTO, gameEventInitVo *VO.GameEventInitVO) {

	instance := gameevent.CreateInstance(event, roundDto, gameEventInitVo)
	instance.HandleRondEvent()
	service.AsyncNotifyGameEventListenerV2(gameEventInitVo.TraceId, event)
}
