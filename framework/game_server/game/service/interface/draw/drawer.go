package draw

import (
	types "sl.framework.com/game_server/game/service/type"
)

// IGameDrawer 游戏结算接口 不同游戏实现不同逻辑
type IGameDrawer interface {
	/*
		Init初始化函数 用于初始化接口 目前只是用来设置接口的TraceId
	*/
	Init(traceId string)

	/*
		ParseGameResult 设置游戏结果
		gameDrawInput *types.GameDrawInput:游戏结果数据 统一结构
		*types.GameDrawOutput:游戏结果数据 返回给能力中心
		对游戏结果进行解析 将统一的接收结构转化为特定的游戏类型需要的结构
		计算出结算时需要的数据 保存到对象中以待结算使用
		组装能力中心需要的结构作为返回值
	*/
	ParseGameResult(gameDrawInput *types.EventDTO) *types.GameRoundResultDTO

	/*
		SettleAvailableAmount 对注单进行有效投注计算
		drawOrder *types.BetOrderV2:需要结算的注单
		该函数根据对象中保存的中间结果对订单BetOrderV2进行结算
		更新注单中的字段 AvailableBetAmount
	*/
	//SettleAvailableAmount(drawOrder []*dto.BetDTO)

	/*
			SettleOrder 对注单批次列表进行结算
		    traceId string  跟踪id
			result *types.GameRoundResultDTO:局结果
			gameDrawDTO   *types.GameDrawDataDTO 结算对象
			orders     *[]int64 结算注单列表
			wg  *sync.WaitGroup 等待
	*/
	SettleOrder(traceId string, result *types.GameRoundResultDTO, gameDrawDTO *types.GameDrawDataDTO, orders *[]int64) []*types.SettleDTO

	/*
			AfterCompletion 结算完成之后的处理逻辑
		    traceId     跟踪Id
			gameRoomId  int64 房间Id
			gameRoundId int64 局Id
			drawOrder *types.SettleDTO:需要结算的注单
	*/
	AfterCompletion(traceId string, gameRoomId, gameRoundId int64, drawOrder []*types.SettleDTO)
}
