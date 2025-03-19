package bet

import (
	"sl.framework.com/game_server/game/service/type/dto"
)

/*
	IGameBettor 游戏下注校验接口 不同游戏实现不同逻辑
	接口函数中每个功能都是相互独立里 游戏框架会使用任务编排原语WaitGroup进行并发调用
	接口函数中的参数使用指针类型传入 每个接口函数只使用指针下数据不可修改指针下的数据
*/

type IGameBettor interface {
	/*
		Init初始化函数 用于初始化接口 目前只是用来设置接口的TraceId
	*/
	Init(traceId string)

	/*
		ValidateUserLimit 个人限红校验
		order *types.BetOrder 玩家下注订单
		返回值:校验结果 校验通过返回errcode.ErrorOk 校验失败返回相应错误码
		个人限红开关关闭则直接返回验证通过 个人限红开启则从redis缓存中获取个人限红信息并对个人单次下注进行校验
		并校验玩家在当前币种所有该玩法下注的限红
		如果没有个人限红校验则直接返回errcode.ErrorOk

		个人限红规则:
			1.单次下注小于等于个人限红最大值
			2.单次下注大于个人限红最小值
			3.个人下注总额小于等于个人限红最大值
	*/
	ValidateUserLimit(order *dto.BetDTO) int

	/*
		ValidateRoomLimit 房间限红校验
		order *types.BetOrder 玩家下注订单
		返回值:校验结果 校验通过返回errcode.ErrorOk 校验失败返回相应错误码
		房间限红校验 从redis缓存中获取房间限红信息并对个人单次下注进行房间限红校验
		并校验房间内玩家在当前币种所有该玩法下注的房间限红
		如果没有房间限红校验则直接返回errcode.ErrorOk

		房间限红规则:
			1.单次限红不超过房间限红最大值
			2.单次限红小于房间限红最小值
			3.庄闲差值不超过房间限红最大值 包含本次下注
	*/
	ValidateRoomLimit(order *dto.BetDTO) int

	/*
		ValidatePlayType 玩法规则校验
		order *types.BetOrder 玩家下注订单
		返回值:校验结果 校验通过返回errcode.ErrorOk 校验失败返回相应错误码
		如果没有玩法校验则直接返回 errcode.ErrorOk
	*/
	ValidatePlayType(order *dto.BetDTO) int

	/*
		ValidateExtraRule 其他规则校验
		order *types.BetOrder 玩家下注订单
		返回值:校验结果 校验通过返回errcode.ErrorOk 校验失败返回相应错误码
		其他规则校验用于以上几个校验未包含的规则
		如果没有则直接返回 errcode.ErrorOk
	*/
	ValidateExtraRule(order *dto.BetDTO) int

	/*
			AfterBetComplete 下注完成回调函数
			gameRoomId int64 房间ID
			gameRoundId int64 局id
		    userId	int64 用户id
			order *types.BetOrder 玩家下注订单
			返回值:调用结果 成功返回errcode.ErrorOk 失败返回相应错误码
			默认返回 errcode.ErrorOk
	*/
	AfterBetComplete(gameRoomId, gameRoundId, userId int64, drawOrder []*dto.BetDTO) int

	/*
			AfterCancelComplete 下注取消完成回调函数
			gameRoomId int64 房间ID
			gameRoundId int64 局id
		    userId	int64 用户id
			order *types.BetOrder 玩家下注订单
			返回值:调用结果 成功返回errcode.ErrorOk 失败返回相应错误码
			默认返回 errcode.ErrorOk
	*/
	AfterCancelComplete(gameRoomId, gameRoundId, userId int64, drawOrder []*dto.BetDTO) int

	/*
			AfterConfirmedComplete 下注确认完成回调函数
			gameRoomId int64 房间ID
			gameRoundId int64 局id
		    userId	int64 用户id
			order *types.BetOrder 玩家下注订单
			返回值:调用结果 成功返回errcode.ErrorOk 失败返回相应错误码
			默认返回 errcode.ErrorOk
	*/
	AfterConfirmedComplete(gameRoomId, gameRoundId, userId int64, drawOrder []*dto.BetDTO) int
}
