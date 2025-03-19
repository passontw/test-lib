package cache

import (
	"encoding/json"
	"sl.framework.com/game_server/game/dao/redisdb"
	types "sl.framework.com/game_server/game/service/type"
	"sl.framework.com/game_server/redis/rediskey"
	"sl.framework.com/trace"
)

type SettleCache struct {
	*BaseCache //基础缓存结构 封装解析与反解析函数

	Data []*types.GameRoundResultDTO //缓存中实际存储的结构

	TraceId string //用于日志跟踪
	RoomId  int64
}

/**
 * Set
 * 保存游戏局结果
 * @param sign_dto *types.GameRoundResultDTO   - 游戏结果结构体
 * @return
 */

func (c *SettleCache) Set(dto *types.GameRoundResultDTO) {
	trace.Info("settle_cache put gameRoomId:%v,sign_dto:%v", c.RoomId, dto)
	cacheKey := rediskey.GetGameRoundResultCacheKey(c.RoomId)
	if !redisdb.SetList[types.GameRoundResultDTO](cacheKey.Key, cacheKey.Expire, *dto) {
		trace.Error("cache put key:%v,value:%v", cacheKey.Key, dto)
	}
	if len(c.Data) == 0 {
		c.Data = make([]*types.GameRoundResultDTO, 0)
	}
	c.Data = append(c.Data, dto)
	//检查存储长度
	listLen := redisdb.GetListSize(cacheKey.Key)
	if listLen >= types.GameDrawResultCacheSize {
		trace.Notice("cache put key:%v,value:%v over size:%v", cacheKey.Key, dto, listLen)
		if err := redisdb.RemoveTopFromList(cacheKey.Key); err != nil {
			trace.Error("cache remove key:%v failed,error:%v", cacheKey.Key, err.Error())
		}
	}
}

/**
 * Get
 * 获取游戏局结果
 * @param
 * @return bool
 */

func (c *SettleCache) Get() (success bool) {
	trace.Info("settle_cache 获取结算结果 gameRoomId:%v", c.RoomId)
	var (
		results []string
		ok      bool
		err     error
	)
	if len(c.Data) == 0 {
		c.Data = make([]*types.GameRoundResultDTO, 0)
	}
	cacheKey := rediskey.GetGameRoundResultCacheKey(c.RoomId)
	if results, ok = redisdb.LAllMember(cacheKey.Key); !ok {
		trace.Notice("结算缓存获取失败 key:%v", cacheKey.Key)
		return
	}

	for _, val := range results {
		item := new(types.GameRoundResultDTO)
		if err = json.Unmarshal([]byte(val), item); nil != err {
			trace.Error("结算缓存 json unmarshal failed, key=%v, val=%v, err=%v", cacheKey.Key,
				val, err.Error())
			return
		}
		c.Data = append(c.Data, item)
	}
	return true
}
