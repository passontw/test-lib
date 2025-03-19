package rpcreq

import (
	"fmt"
	"sl.framework.com/game_server/conf"
	"sl.framework.com/game_server/game/service/type/dto"
	"sl.framework.com/game_server/rpc_client/config"
)

/*
*
  - 创建现场员工信息

调用接口:/feign/worker/build/{username}/{gameId}/{workerType}
*/
func WorkerClientBuild(traceId, userName, gameId, workerType string, workerDTO *dto.WorkerDTO) int {
	url := fmt.Sprintf(config.WorkerClientBuild, conf.GetPlatformInfoUrl(), userName, gameId, workerType)
	msg := fmt.Sprintf("[创建现场员工信息] traceId=%v, username=%v, gameId=%v, workerType=%v, url=%v", traceId,
		userName, gameId, workerType, url)

	ret := runHttpPost(traceId, msg, url, nil, workerDTO)
	return ret
}
