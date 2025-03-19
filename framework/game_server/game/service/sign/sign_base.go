package sign

import (
	sign_dto2 "sl.framework.com/game_server/game/service/sign/sign_dto"
	types "sl.framework.com/game_server/game/service/type"
)

// 此类适用于进行下注校验的父类，游戏需要根据自己的订单数据继承此类，并重写对应的函数
type SignBase struct {
}

/**

* Sign 投注批量签名
@param betList []*types.BetOrderV2 下注注单列表
@return  []*sign_dto.SignTextDTO 签名数据列表
*/

func (s *SignBase) Sign(betList []*types.BetOrderV2) []*sign_dto2.SignTextDTO {
	//TODO
	return make([]*sign_dto2.SignTextDTO, 0)
}

/**
*PrePareVerifyData  准备验证签名数据
@param gameRecordsList []*types.BetOrderV2 游戏记录集合
@return  []*sign_dto.SignTextDTO 签名数据列表
*/

func (s *SignBase) PrePareVerifyData(gameRecordsList []*any) []*sign_dto2.SignVerifyDTO {
	//TODO
	return make([]*sign_dto2.SignVerifyDTO, 0)
}
