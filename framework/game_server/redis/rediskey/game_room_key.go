package rediskey

import (
	"fmt"
	"sl.framework.com/game_server/redis/redis_tool"
	"sl.framework.com/game_server/redis/types"
	"strconv"
	"time"
)

//房间相关的redis key
/*
	1.玩家进入房间 离开房间
	2.房间内一个牌靴下牌的数量
*/

const (
	// roomInfoFileKeyPrefix 避免key重复 每个文件中要使用一个key与其他文件区别
	roomInfoFileKeyPrefix = "RoomInfo"
)

// 玩家进入房间 离开房间
const (
	roomActionPrefix     = "JoinLeaveGameRoom"
	roomActionLockPrefix = "RoomActionLock"
)

// GetRoomActionRedisInfo 房间用户key信息
// 如:{serverRedisKeyPrefix}:JoinLeaveGameRoom:JoinLeaveGameRoom
func GetRoomActionRedisInfo(roomId int64) *types.RedisInfo {
	return redistool.BuildRedisInfo(
		time.Duration(15)*time.Minute,
		roomInfoFileKeyPrefix,
		roomActionPrefix,
		strconv.FormatInt(roomId, 10),
	)
}

// GetRoomActionLockRedisInfo 房间用户key信息锁 没有必要加锁
func GetRoomActionLockRedisInfo(roomId int64) *types.RedisLockInfo {
	return redistool.BuildRedisLockInfo(
		types.RedisLockExpireDuration,
		roomInfoFileKeyPrefix,
		roomActionLockPrefix,
		strconv.FormatInt(roomId, 10),
	)
}

const (
	roomCardNumPrefix     = "RoomCardNum"
	roomCardNumLockPrefix = "RoomCardNumLock"
)

// GetGameCardNumRedisInfo 下局局号信息
func GetGameCardNumRedisInfo(gameRoomId int64) *types.RedisInfo {
	return redistool.BuildRedisInfo(
		time.Duration(15)*time.Minute,
		roomInfoFileKeyPrefix,
		roomCardNumPrefix,
		strconv.FormatInt(gameRoomId, 10),
	)
}

// GetGameCardNumLockRedisInfo 下局局号信息key信息锁
func GetGameCardNumLockRedisInfo(gameRoomId int64) *types.RedisLockInfo {
	return redistool.BuildRedisLockInfo(
		types.RedisLockExpireDuration,
		roomInfoFileKeyPrefix,
		roomCardNumLockPrefix,
		strconv.FormatInt(gameRoomId, 10),
	)
}

/*
1.赔率倍率
	表名:	HRoomOddsInfo{GM054247311DA}
	field:	{GameRoomId}{gameId}{gameWagerId}
	value:	{Odds:"23"}
2.房间限红信息
	表名:	HRoomLimitInfo{GM054247311DA}
	field:	{GameRoomId}{gameId}{gameWagerId}
	value:	{currency:"",minAmount:"",maxAmount:""}
*/

/*
房间详情信息 todo:暂时未使用到

	1.每局结束或者开始前从能力中心获取并设置该key
	2.每局开始再从该key中获取数据填充房间赔率信息和房间限红信息
*/
const (
	roomDetailedInfoKeyPrefix = "RoomDetailedInfo"
	roomDetailedLockKeyPrefix = "RoomDetailedInfoLock"
)

// GetRoomDetailedInfoRedisInfo 房间详情信息
// 如: hash表名
func GetRoomDetailedInfoRedisInfo(roomId int64) *types.RedisInfo {
	return redistool.BuildRedisInfo(
		redistool.GetRedisExpireDuration(),
		roomInfoFileKeyPrefix,
		roomDetailedInfoKeyPrefix,
		strconv.FormatInt(roomId, 10),
	)
}

// GetRoomDetailedLockHRedisInfo 房间信息详情表分布式锁key
func GetRoomDetailedLockHRedisInfo(roomId int64) *types.RedisLockInfo {
	return redistool.BuildRedisLockInfo(
		types.RedisLockExpireDuration,
		roomInfoFileKeyPrefix,
		roomDetailedLockKeyPrefix,
		strconv.FormatInt(roomId, 10),
	)
}

/*
赔率倍率
	表名:	HRoomOddsInfo{GM054247311DA}
	field:	{GameRoomId}{gameId}{gameWagerId}
	value:	{Odds:"23"}
*/

/*房间赔率信息*/
const (
	roomOddInfoKeyPrefix     = "HRoomOddInfo"
	roomOddInfoLockKeyPrefix = "HRoomOddInfoLock"
)

// GetRoomOddHRedisInfo 房间信息赔率表 返回表名以及过期时间
// 如: hash表名
func GetRoomOddHRedisInfo(gameRoundId int64) *types.HashRedisInfo {
	return redistool.BuildHashRedisInfo(
		redistool.GetRedisExpireDuration(),
		"",
		roomInfoFileKeyPrefix,
		roomOddInfoKeyPrefix,
		strconv.FormatInt(gameRoundId, 10),
	)
}

// GetRoomOddHRedisInfoEx 房间信息赔率表 返回hash表名 filed 过期时间
// 如: hash表名
func GetRoomOddHRedisInfoEx(gameRoundId, roomId, gameId, gameWagerId int64) *types.HashRedisInfo {
	return redistool.BuildHashRedisInfo(
		redistool.GetRedisExpireDuration(),
		fmt.Sprintf("%v:%v:%v", roomId, gameId, gameWagerId),
		roomInfoFileKeyPrefix,
		roomOddInfoKeyPrefix,
		strconv.FormatInt(gameRoundId, 10),
	)
}

// GetRoomOddLockHRedisInfo 房间信息赔率表分布式锁名
func GetRoomOddLockHRedisInfo(gameRoundId int64) *types.RedisLockInfo {
	return redistool.BuildRedisLockInfo(
		types.RedisLockExpireDuration,
		roomInfoFileKeyPrefix,
		roomOddInfoLockKeyPrefix,
		strconv.FormatInt(gameRoundId, 10),
	)
}

/*
房间限红信息
	表名:	HRoomLimitInfo{GM054247311DA}
	field:	{GameRoomId}{gameId}{gameWagerId}
	value:	{currency:"",minAmount:"",maxAmount:""}
*/

/*房间限红信息*/
const (
	hRoomLimitInfoKeyPrefix     = "HRoomLimitInfo"
	hRoomLimitInfoLockKeyPrefix = "HRoomLimitInfoLock"
)

// GetRoomLimitHRedisInfo 房间限红信息 获取表名和过期时间
// 桌台限红以及桌台赔率信息是一次性设置进去所以HTable只使用gameRoundId
func GetRoomLimitHRedisInfo(gameRoundId int64) *types.HashRedisInfo {
	return redistool.BuildHashRedisInfo(
		redistool.GetRedisExpireDuration(),
		"",
		roomInfoFileKeyPrefix,
		hRoomLimitInfoKeyPrefix,
		strconv.FormatInt(gameRoundId, 10),
	)
}

// GetRoomLimitHRedisInfoEx 房间限红信息 获取表名 Field 过期时间
// 如: hash表名
func GetRoomLimitHRedisInfoEx(gameRoundId, roomId, gameId, gameWagerId int64, currency string) *types.HashRedisInfo {
	return redistool.BuildHashRedisInfo(
		redistool.GetRedisExpireDuration(),
		fmt.Sprintf("%v:%v:%v:%v", roomId, gameId, gameWagerId, currency),
		roomInfoFileKeyPrefix,
		hRoomLimitInfoKeyPrefix,
		strconv.FormatInt(gameRoundId, 10),
	)
}

// GetRoomLimitLockHRedisInfo 房间信息限红信息表分布式锁名
func GetRoomLimitLockHRedisInfo(gameRoundId int64) *types.RedisLockInfo {
	return redistool.BuildRedisLockInfo(
		types.RedisLockExpireDuration,
		roomInfoFileKeyPrefix,
		hRoomLimitInfoLockKeyPrefix,
		strconv.FormatInt(gameRoundId, 10),
	)
}

/*
房间限红信息
	List:	LRoomLimitInfo{gameRoundId}{currency}
*/

const (
	lRoomLimitInfoKeyPrefix = "LRoomLimitInfo"
)

// GetRoomLimitLRedisInfo 房间总限红信息 获取list表名和过期时间
func GetRoomLimitLRedisInfo(gameRoundId int64, currency string) *types.RedisInfo {
	return redistool.BuildRedisInfo(
		redistool.GetRedisExpireDuration(),
		roomInfoFileKeyPrefix,
		lRoomLimitInfoKeyPrefix,
		strconv.FormatInt(gameRoundId, 10),
		currency,
	)
}
