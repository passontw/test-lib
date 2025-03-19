package gameevent

import (
	"encoding/json"
	"fmt"
	"sl.framework.com/async"
	"sl.framework.com/game_server/conf"
	errcode "sl.framework.com/game_server/error_code"
	"sl.framework.com/game_server/game/dao/redisdb"
	types "sl.framework.com/game_server/game/service/type"
	"sl.framework.com/game_server/game/service/type/VO"
	"sl.framework.com/game_server/game/service/type/const_type"
	"sl.framework.com/game_server/game/service/type/rocket_mq"
	"sl.framework.com/game_server/mq"
	"sl.framework.com/game_server/redis/cache"
	"sl.framework.com/game_server/redis/rediskey"
	"sl.framework.com/tool"
	"sl.framework.com/trace"
	"strconv"
	"time"
)

type GameDataEvent struct {
	types.EventBase
}

/**
 * NewGameData
 * 创建发牌实例
 *
 * @param roomId int64 - 房间id
 * @param roundId int64 - 局id
 * @param nextRoundId int64 - 下一局id
 * @param event types.GameEventVO - 来自数据源的事件信息
 * @param roundDto *types.GameRoundDTO - 的游戏局信息
 * @return RETURN - 返回游戏事件实例
 */

func NewGameData(event types.GameEventVO, roundDto *types.GameRoundDTO, gameEventInitVo *VO.GameEventInitVO) *GameDataEvent {
	return &GameDataEvent{
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

func (e *GameDataEvent) HandleRondEvent() {
	//直接设置为成功
	trace.Info("[游戏发牌] GameData %v 局信息%+v", e.MsgHeader, e.Dto.Payload)
	*e.RetHandleEvent = errcode.ErrorOk
	EventCommonSet(&e.EventBase, string(types.GameEventCommandGameData), string(types.GameEventCommandGameData))
	//发第一张牌的时候，提交为提交的注单
	redisLockInfo := rediskey.GetBetConfirmedGameDataLockRedisInfo(strconv.FormatInt(e.Dto.GameRoomId, 10),
		strconv.FormatInt(e.Dto.GameRoundId, 10))
	if !redisdb.TryLock(redisLockInfo) {
		trace.Notice("[游戏发牌] traceId=%v, 首次发牌已经提交过一次注单 redis lock failed, lock info=%+v", e.TraceId, redisLockInfo)
		return
	}
	//从redis中把当前局的注单全部拿出来遍历然后提交注单
	betDtoList := cache.GetOrders(e.TraceId, strconv.FormatInt(e.Dto.GameRoomId, 10),
		strconv.FormatInt(e.Dto.GameRoundId, 10))
	if len(betDtoList) == 0 {
		trace.Notice("[游戏发牌] 当前局注单数量为0，不需要confirm TraceId=%v,GameRoomId=%v,GameRoundId=%v", e.TraceId, e.Dto.GameRoomId, e.Dto.GameRoundId)
	} else {
		trace.Debug("[游戏发牌] 当前局注单数量为%v，提交确认中 TraceId=%v,GameRoomId=%v,GameRoundId=%v", len(betDtoList), e.TraceId, e.Dto.GameRoomId, e.Dto.GameRoundId)

		userInfoMap := make(map[string]types.UserCurrencyInfo)
		userInfoList := make([]types.UserCurrencyInfo, 0)
		//把list转为map
		for _, item := range betDtoList {
			//只把没有提交的注单选出来
			if item.PostStatus == string(const_type.PostStatusCreate) {
				strUserId := strconv.FormatInt(item.UserId, 10)
				userInfo := types.UserCurrencyInfo{UserId: strUserId, Currency: item.Currency}
				if _, ok := userInfoMap[strUserId]; !ok {
					userInfoMap[strUserId] = userInfo
					userInfoList = append(userInfoList, userInfo)
				}
			}
		}
		if len(userInfoList) == 0 {
			trace.Info("[游戏发牌] [提交注单分片] traceId:%v 没有需要提交的注单", e.TraceId)
			return
		}
		pWatcher := tool.NewWatcher("提交注单分片")
		//3.获取分片大小
		patchSize := conf.ServerConf.Common.BetConfirmSIze
		trace.Info("[游戏发牌] [提交注单分片] traceId:%v patchSize:%v settleOrderList:%+v", e.TraceId, patchSize, userInfoList)
		patches := tool.SplitList[types.UserCurrencyInfo](userInfoList, patchSize)
		//遍历
		for _, row := range patches {
			if len(row) == 0 {
				continue
			}
			betConfirmPayload := rocket_mq.BetConfirmMessagePayload{
				GameId:      e.Dto.GameId,
				GameRoomId:  e.Dto.GameRoomId,
				GameRoundId: e.Dto.GameRoundId,
				UserInfo:    row,
			}

			messageStr, err := json.Marshal(betConfirmPayload)
			//id := tool.GenerateRandomString(32)
			trace.Info("[游戏发牌] [提交注单分片] traceId:%v 分片数组大小:%v patchSize:%v messageStr:%v", e.TraceId, len(patches), patchSize, messageStr)
			if err != nil {
				trace.Error("[游戏发牌] [提交注单分片] traceId:%v  序列化messageDto=%+v 失败.", e.TraceId, betConfirmPayload)
			} else {
				topic := generateBetConfirmTopic()
				createTime := strconv.FormatInt(time.Now().Unix(), 10)
				fn := func() {
					trace.Info("[游戏发牌] [提交注单分片] traceId:%v 异步发送到Mq topic:%v messageStr:%v", e.TraceId, topic, messageStr)
					mq.SendMessage(topic, strconv.FormatInt(e.Dto.GameId, 10), e.TraceId, createTime, string(messageStr))
				}
				async.AsyncRunCoroutine(fn)
			}

		}
		pWatcher.Stop()

	}

	return
}

func generateBetConfirmTopic() string {
	var result string
	rndInt := tool.GenerateRandomRange(1, 9999) % 2
	switch rndInt {
	case 0:
		result = string(mq.TopicBetConfirmOut0)
	case 1:
		result = string(mq.TopicBetConfirmOut1)
	default:
		result = string(mq.TopicBetConfirmOut0)
	}
	return result
}
