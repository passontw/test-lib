package rpcreq

import (
	"fmt"
	"sl.framework.com/game_server/conf"
	"sl.framework.com/game_server/game/service/type/dto"
	"sl.framework.com/trace"
)

/*
	BalanceResponse
	请求玩家余额回包
*/

type BalanceResponse struct {
	UserId           string  `json:"userId"`           //数据库字段:user_id 用户id
	FinancialAccount string  `json:"financialAccount"` //数据库字段:financial_account 财务账号,全局唯一
	Currency         string  `json:"currency"`         //数据库字段:currency 币种编码
	Balance          float64 `json:"balance"`          //数据库字段:balance 余额
	Status           string  `json:"status"`           //数据库字段:status 状态:启用 Enable，停用 Disable
}

/**
 * BalanceRequest
 * 获取玩家余额信息
 *
 * @param traceId string - traceId 用于日志跟踪
 * @param currency string - 玩家货币类型
 * @param userId string - 玩家Id
 * @return *BalanceResponse - 返回的玩家余额信息
 * @return int - 请求返回码
 */

func BalanceRequest(traceId, currency, userId string) (*BalanceResponse, int) {
	url := fmt.Sprintf("%v/feign/wallet/get/balance/%v/%v",
		conf.GetPlatformInfoUrl(), userId, currency)
	msg := fmt.Sprintf("BalanceRequest traceId=%v, userId=%v, currency=%v, url=%v", traceId, userId, currency, url)

	balance := new(BalanceResponse)
	ret := runHttpGet(traceId, msg, url, balance)
	return balance, ret
}

/**
 * GetTransactionList
 * 查询交易记录操作 才结算是校验是否已经扣钱成功，只有扣钱成功的注单才进行派彩
 *
 * @param traceId string - traceId 用于日志跟踪
 * @param currency string - 玩家货币类型
 * @param userId string - 玩家Id
 * @param queryDto *dto.QueryTransactionDTO -查询请求体
 * @return []*dto.UserTransactionDTO - 返回的玩家交易记录
 * @return int - 请求返回码
 */

func GetTransactionList(traceId string, queryDto *dto.QueryTransactionDTO) ([]*dto.UserTransactionDTO, int) {
	url := fmt.Sprintf("%v/feign/wallet/list/transaction",
		conf.GetPlatformInfoUrl())
	msg := fmt.Sprintf("GetTransactionList traceId=%v, url=%v", traceId, url)

	transactionDTO := make([]*dto.UserTransactionDTO, 0)
	ret := runHttpPost(traceId, msg, url, queryDto, &transactionDTO)
	for _, item := range transactionDTO {
		trace.Debug("开奖 校验注单 返回值结构体 item=%+v", item)
	}
	return transactionDTO, ret
}
