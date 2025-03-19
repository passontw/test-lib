package impl

import (
	errcode "sl.framework.com/game_server/error_code"
	"sl.framework.com/game_server/game/service"
	"sl.framework.com/game_server/game/service/sign/sign_dto"
	types "sl.framework.com/game_server/game/service/type"
	"sl.framework.com/game_server/game/service/type/dto"
	rpcreq "sl.framework.com/game_server/rpc_client"
	"sl.framework.com/trace"
	"strconv"
)

type SignContextImpl struct {
	TraceId string       `json:"trace_id"`
	GameId  types.GameId `json:"game_id"`
}

/**
 * Sign
 * 投注批量签名
 *
 * @param betDTOList  []*types.BetOrderV2- 注单列表
 * @return []*types.BetOrderV2 - 签名后投注信息
 */

func (s *SignContextImpl) Sign(betDTOList []*dto.BetDTO) []*dto.BetDTO {
	trace.Info("Sign traceId:%v betDTOList:%v", s.TraceId, betDTOList)
	var (
		signor            = service.GetSignor(s.TraceId, s.GameId)
		signResultDTOList []*sign_dto.SignResultDTO
		resultList        = make([]*dto.BetDTO, 0)
		resultBetSignMap  = make(map[int64]string)
	)
	defer service.PutSignor(s.GameId, signor)
	if signor == nil {
		trace.Error("Sign traceId:%v betDTOList:%v get signor failed.", s.TraceId, betDTOList)
		return resultList
	}

	signTextDTOList := signor.Sign(betDTOList)
	trace.Info("Sign traceId:%v from game server signTextDTOList:%v", s.TraceId, signTextDTOList)
	signClient := rpcreq.SignClient{}
	if ok := signClient.Sign(s.TraceId, &signTextDTOList, &signResultDTOList); ok != errcode.ErrorOk {
		trace.Error("Sign traceId:%v signTextDTOList:%v request from agent failed.", s.TraceId, signTextDTOList)
	}
	if len(signResultDTOList) == 0 {
		trace.Notice("Sign traceId:%v signTextDTOList:%v request from agent empty.", s.TraceId, signTextDTOList)
		return resultList
	}
	for _, signResultDTO := range signResultDTOList {
		betId, _ := strconv.ParseInt(signResultDTO.Id, 10, 64)
		resultBetSignMap[betId] = signResultDTO.Sign
	}

	for _, betDTO := range betDTOList {
		betDTO.Md5 = resultBetSignMap[betDTO.Id]
	}
	trace.Info("Sign traceId:%v betDTOList:%v success.", s.TraceId, signTextDTOList)
	return betDTOList
}

/**
 * Verify
 * 验证签名,返回延签通过的 游戏记录id集合
 *
 * @param gameRecordsList  []*types.DrawOrder- 游戏记录集合
 * @return []int64 -返回延签通过的 游戏记录id集合
 */

func (s *SignContextImpl) Verify(gameRecordsList []*dto.BetDTO) []int64 {
	trace.Info("Verify traceId:%v gameRecordsList:%v", s.TraceId, gameRecordsList)
	var (
		signor          = service.GetSignor(s.TraceId, s.GameId)
		signTestDTOList []*sign_dto.SignVerifyDTO
		resultList      = make([]int64, 0)
		//resultBetSignMap = make([]*sign_dto.SignVerifyResultDTO, 0)
	)
	defer service.PutSignor(s.GameId, signor)
	if signor == nil {
		trace.Error("Sign traceId:%v gameRecordsList:%v get signor failed.", s.TraceId, gameRecordsList)
		return resultList
	}
	signTestDTOList = signor.PrePareVerifyData(gameRecordsList)
	trace.Info("Verify traceId:%v from game server signTestDTOList:%v", s.TraceId, signTestDTOList)
	//signClient := rpcreq.SignClient{}
	//if ok := signClient.Verify(s.TraceId, &signTestDTOList, &resultBetSignMap); ok != errcode.ErrorOk {
	//	trace.Error("Sign traceId:%v resultBetSignMap:%v request from agent failed.", s.TraceId, resultBetSignMap)
	//}
	//if len(signTestDTOList) == 0 {
	//	trace.Notice("Verify traceId:%v gameRecordsList:%v request from game empty.", s.TraceId, signTestDTOList)
	//	return resultList
	//}
	//for _, signItem := range resultBetSignMap {
	//	if signItem.Ok {
	//		betId, _ := strconv.ParseInt(signItem.Id, 10, 64)
	//		resultList = append(resultList, betId)
	//	}
	//}
	for _, dto := range signTestDTOList {
		orderId, _ := strconv.ParseInt(dto.Id, 10, 64)
		resultList = append(resultList, orderId)
	}
	trace.Info("Verify traceId:%v gameRecordsList:%v success failed list :%v.", s.TraceId, gameRecordsList, resultList)
	return resultList
}
