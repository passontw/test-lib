package interfaces

import (
	sign_dto2 "sl.framework.com/game_server/game/service/sign/sign_dto"
	"sl.framework.com/game_server/game/service/type/dto"
)

type ISignHandler interface {

	/**
	* Sign 投注批量签名
	@param betList []*types.BetOrderV2 下注注单列表
	@return  []*sign_dto.SignTextDTO 签名数据列表
	*/

	Sign(betList []*dto.BetDTO) []*sign_dto2.SignTextDTO

	/**
	*PrePareVerifyData  准备验证签名数据
	@param gameRecordsList []*types.DrawOrder 游戏记录集合
	@return  []*sign_dto.SignVerifyDTO 签名数据列表
	*/

	PrePareVerifyData(gameRecordsList []*dto.BetDTO) []*sign_dto2.SignVerifyDTO
}
