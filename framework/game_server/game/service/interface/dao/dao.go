package dao

import (
	"sl.framework.com/game_server/game/service/type/dto"
)

/**
 * IGameDB
 * 保存订单到数据库
 * 不同游戏需要实现保存数据库的接口
 */

type IGameDB interface {
	/*
	 * SaveDBBatch
	 * 批量保存数据到数据库
	 * @param traceId string - traceId 用于日志跟踪
	 * @param gameRoomId int64 - gameRoomId 房间Id
	 * @param gameRoundId int64 - gameRoundId 局Id
	 * @param betList []dto.BetDTO - 需要保存的订单
	 *
	 * 该函数使用beego自带orm对数据库进行操作 需要从框架内获取orm对象然后再对数据库操作
	 */
	SaveDBBatch(traceId string, gameRoomId, gameRoundId int64, betList *[]dto.BetDTO)

	/**
	 * GetOrderNoList
	 * 从数据库获取未结算订单列表
	 *
	 * @param traceId string -  traceId 用于日志跟踪
	 * @param gameRoomId int64 - gameRoomId 房间Id
	 * @param gameRoundId int64 - gameRoundId 局Id
	 * @param gameRoundNo string - 局号
	 * @return []int64 - 未结算订单号列表
	 * 该函数使用beego自带orm对数据库进行操作 需要从框架内获取orm对象然后再对数据库操作
	 */

	GetOrderNoList(traceId string, gameRoomId, gameRoundId int64, gameRoundNo string) []int64

	/**
	 * UpdateOrders
	 * 更新注单信息
	 *
	 * @param traceId - 跟踪id
	 * @param gameRoomId - 房间id
	 * @param gameRoundId - 局id
	 * @param betList *[]*dto.BetDTO - 注单信息
	 * @return RETURN - 返回值说明
	 * 该函数使用beego自带orm对数据库进行操作 需要从框架内获取orm对象然后再对数据库操作
	 */

	UpdateOrders(traceId string, gameRoomId, gameRoundId int64, betList *[]*dto.BetDTO)
}
