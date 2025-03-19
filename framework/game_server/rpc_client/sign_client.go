package rpcreq

import (
	"fmt"
	"sl.framework.com/game_server/conf"
	sign_dto2 "sl.framework.com/game_server/game/service/sign/sign_dto"
)

type SignClient struct {
}

/**
 * Sign
 * 投注签名
 *
 * @param traceId string - traceId 用于日志跟踪
 * @param signTestDTOList *[]*sign_dto.SignTextDTO -  玩家下注列表
 * @return *[]*sign_dto.SignResultDTO - 请求返回码
 */

func (c *SignClient) Sign(traceId string, signTestDTOList *[]*sign_dto2.SignTextDTO, signResultDTOList *[]*sign_dto2.SignResultDTO) int {
	url := fmt.Sprintf("%v/feign/sign/batch/sign", conf.GetPlatformInfoUrl())
	msg := fmt.Sprintf("Sign Request traceId=%v, url=%v", traceId, url)

	ret := runHttpPost(traceId, msg, url, signTestDTOList, signResultDTOList)
	return ret
}

/**
 * Verify
 * 验证签名数据 获取验证结果列表
 *
 * @param traceId string - traceId 用于日志跟踪
 * @param signVerifyDTOList *[]*sign_dto.SignVerifyDTO -  玩家下注列表
 * @return *[]*sign_dto.SignVerifyResultDTO - 请求返回码
 */

func (c *SignClient) Verify(traceId string, signVerifyDTOList *[]*sign_dto2.SignVerifyDTO, signVerifyResultDTOList *[]*sign_dto2.SignVerifyResultDTO) int {
	url := fmt.Sprintf("%v/feign/sign/verify", conf.GetPlatformInfoUrl())
	msg := fmt.Sprintf("Sign Request traceId=%v, url=%v", traceId, url)

	ret := runHttpPost(traceId, msg, url, signVerifyDTOList, signVerifyResultDTOList)
	return ret
}
