package interfaces

import (
	"sl.framework.com/game_server/game/service/type/dto"
)

type ISignContext interface {

	/**
	* Sign 投注批量签名
	@param betList []*types.BetOrderV2 下注注单列表
	@return  []*types.BetOrderV2 签名数据列表
	*/

	Sign(betList []*dto.BetDTO) []*dto.BetDTO

	/**
	* Sign 验证签名
	@param gameRecordsList []*types.BetOrderV2 游戏记录集合
	@return  []int64 验证结果
	*/

	Verify(gameRecordsList []*dto.BetDTO) []int64
}
