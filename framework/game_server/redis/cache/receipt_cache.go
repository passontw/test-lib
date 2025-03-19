package cache

import (
	"encoding/json"
	"fmt"
	"sl.framework.com/game_server/game/dao/redisdb"
	types "sl.framework.com/game_server/game/service/type"
	"sl.framework.com/game_server/redis/rediskey"
	"sl.framework.com/trace"
)

/**
 * GetReceipts
 * 获取小票信息
 *
 * @param traceId string - traceId 用于日志跟踪
 * @param gameRoomId int64 - 房间Id
 * @param gameRoundId int64 - 局Id
 * @param usserId int64 - 用户id
 * @return *types.BetReceipt - 小票详情
 */

func GetReceipts(traceId string, gameRoomId, gameRoundId, userId int64) []*types.BetReceipt {

	msgHeader := fmt.Sprintf("GetReceipts traceId=%v, gameRoomId=%v, gameRoundId=%v",
		traceId, gameRoomId, gameRoundId)

	var (
		val string
		err error
	)
	receiptKey := rediskey.GetBetReceiptKey(gameRoomId, gameRoundId)
	if val, err = redisdb.HGet(receiptKey.HTable, receiptKey.Filed); nil != err {
		trace.Error("%v, redis HGet failed, table=%v, field=%v, error=%v", msgHeader, receiptKey.HTable,
			receiptKey.Filed, err.Error())
		return nil
	}
	if len(val) <= 0 {
		trace.Notice("%v, redis HGet no data, table=%v, field=%v", msgHeader, receiptKey.HTable, receiptKey.Filed)
		return nil
	}

	betReceipts := make(map[int64][]*types.BetReceipt, 0)
	if err = json.Unmarshal([]byte(val), &betReceipts); nil != err {
		trace.Error("%v, json unmarshal failed, table=%v, field=%v, val=%v, err=%v", msgHeader, receiptKey.HTable,
			receiptKey.Filed, val, err.Error())
		return nil
	}

	if len(betReceipts) <= 0 {
		trace.Error("%v, GetOrders empty, table=%v, field=%v, val=%v, err=%v", msgHeader, receiptKey.HTable,
			receiptKey.Filed, val, err.Error())
		return nil
	}
	trace.Info("%v, order list=%v", msgHeader, betReceipts)

	return betReceipts[userId]
}

/**
 * PutReceipts
 * 存储小票信息
 *
 * @param traceId string - traceId 用于日志跟踪
 * @param gameRoomId int64 - 房间Id
 * @param gameRoundId int64 - 局Id
 * @param receipts types.BetReceipt -小票信息
 * @return
 */

func PutReceipts(traceId string, gameRoomId, gameRoundId int64, receipts map[int64][]*types.BetReceipt) {

	msgHeader := fmt.Sprintf("PutReceipts traceId=%v, gameRoomId=%v, gameRoundId=%v",
		traceId, gameRoomId, gameRoundId)

	var (
		val  int64
		err  error
		data []byte
	)
	if data, err = json.Marshal(receipts); err != nil {
		trace.Error("%v, json marshal failed, error=%v", msgHeader, err.Error())
		return
	}
	receiptKey := rediskey.GetBetReceiptKey(gameRoomId, gameRoundId)
	if val, err = redisdb.HSet(receiptKey.HTable, receiptKey.Filed, string(data), receiptKey.Expire); nil != err {
		trace.Error("%v, redis HGet failed, table=%v, field=%v, error=%v", msgHeader, receiptKey.HTable,
			receiptKey.Filed, err.Error())
		return
	}

	trace.Info("%v, val=%v, order list=%v", msgHeader, val, receipts)
}
