package rpcreq

import (
	"fmt"
	"sl.framework.com/game_server/conf"
	types "sl.framework.com/game_server/game/service/type"
)

/*
	GetUserLimitRequest 查询用户个人限红
	调用接口:/v1/user/feign/userBetLimitRule/getDetail/userId/{userId}/{currency}
*/

func GetUserLimitRequest(traceId, currency string, userId, gameRoundId int64) (*types.UserBetLimitInfo, int) {
	url := fmt.Sprintf("%v/feign/userBetLimitRule/getDetail/userId/%v/%v",
		conf.GetPlatformInfoUrl(), userId, currency)
	msg := fmt.Sprintf("GetUserLimitRequest traceId=%v, gameRoundId=%v, userId=%v, currency=%v, url=%v",
		traceId, gameRoundId, userId, currency, url)

	userBetLimitInfo := new(types.UserBetLimitInfo)
	ret := runHttpGet(traceId, msg, url, userBetLimitInfo)
	return userBetLimitInfo, ret
}

/**
 * GetUserLimitBatchRequest
 * 批量获取玩家限红信息
 *
 * @param traceId string - traceId 用于日志跟踪
 * @param currency string - currency 玩家货币类型
 * @param userId int64 - userId 玩家Id
 * @param userInfoList []types.UserBetLimitBatchRequest - userInfoList 玩家信息列表
 * @return []*types.UserBetLimitInfo - 玩家限红信息集合
 * @return int - 请求返回码
 */

func GetUserLimitBatchRequest(traceId, gameRoundNo string, gameRoundId int64,
	userInfoList types.UserBetLimitBatchRequest) ([]*types.UserBetLimitInfo, int) {
	url := fmt.Sprintf("%v/feign/userBetLimitRule/list/userBetLimitRule", conf.GetPlatformInfoUrl())
	msg := fmt.Sprintf("GetUserLimitBatchRequest traceId=%v, gameRoundId=%v, gameRoundNo=%v, userInfoList=%+v, "+
		"url=%v", traceId, gameRoundId, gameRoundNo, userInfoList, url)

	userBetLimitList := make([]*types.UserBetLimitInfo, 0)
	ret := runHttpPost(traceId, msg, url, userInfoList, &userBetLimitList)
	return userBetLimitList, ret
}
