package rpcreq

import (
	"fmt"
	"sl.framework.com/game_server/conf"
	types "sl.framework.com/game_server/game/service/type"
	"sl.framework.com/game_server/game/service/type/dto"
)

/**
 * GetOrBindGameRoundNo
 * 查询局或者绑定局号
 *
 * @param traceId string - traceId 用于日志跟踪
 * @param gameRoundNo string - 局号
 * @param gameRoomId int64 - 房间Id
 * @return *GameRoundInfo - 局信息
 * @return int - 请求返回值
 */

func GetOrBindGameRoundNo(traceId, gameRoundNo string, gameRoomId int64) (*types.GameRoundDTO, int) {
	url := fmt.Sprintf("%v/feign/gameRound/getOne/gameRoomId/%v/gameRoundNo/%v",
		conf.GetPlatformInfoUrl(), gameRoomId, gameRoundNo)
	msg := fmt.Sprintf("GetOrBindGameRoundNo traceId=%v, gameRoomId=%v, gameRoundNo=%v, url=%v", traceId,
		gameRoomId, gameRoundNo, url)

	gameRoundInfo := new(types.GameRoundDTO)
	ret := runHttpGet(traceId, msg, url, gameRoundInfo)
	return gameRoundInfo, ret
}

/**
 * BuildRoundRequest
 * 创建局信息
 *
 * @param traceId string - traceId 用于日志跟踪
 * @param gameRoundInfo *types.GameRoundDTO - 局信息是传入参数也是传出参数
 * @return int - 请求返回码
 */

func BuildRoundRequest(traceId string, gameRoundInfo *types.GameRoundDTO) int {
	url := fmt.Sprintf("%v/feign/gameRound/build", conf.GetPlatformInfoUrl())
	msg := fmt.Sprintf("BuildRoundRequest traceId=%v, gameRoomId=%v, gameRoundNo=%v, url=%v",
		traceId, gameRoundInfo.GameRoomId, gameRoundInfo.RoundNo, url)

	return runHttpPost(traceId, msg, url, gameRoundInfo, gameRoundInfo)
}

/**
 * CloseExceptionRoundRequest
 * 关闭异常局
 *
 * @param traceId string - traceId用于日志跟踪
 * @param gameRoomId int64 - 房间Id
 * @param gameRoundId int64 - 局Id
 * @return int - 请求返回值
 */

func CloseExceptionRoundRequest(traceId string, gameRoomId, gameRoundId int64) int {
	url := fmt.Sprintf("%v/feign/gameRound/closeExceptionRound/%v/%v", conf.GetPlatformInfoUrl(),
		gameRoomId, gameRoundId)
	msg := fmt.Sprintf("CloseExceptionRoundRequest traceId=%v, gameRoomId=%v, gameRoundId=%v, url=%v",
		traceId, gameRoomId, gameRoundId, url)
	return runHttpPut(traceId, msg, url, nil)
}

/**
 * BuildGameEventRequest
 * 创建游戏事件
 *
 * @param traceId string - traceId用于日志跟踪
 * @param gameRoomId int64 - 房间Id
 * @param gameRoundId int64 - 局Id
 * @param gameEvent *cache.GameEvent - 游戏事件指针类型
 * @return int - 请求返回值
 */

func BuildGameEventRequest(traceId string, gameRoomId, gameRoundId int64, gameEvent *dto.GameCommandDTO) int {
	url := fmt.Sprintf("%v/feign/gameCommand/build", conf.GetPlatformInfoUrl())
	msg := fmt.Sprintf("BuildGameEventRequest traceId=%v, gameRoomId=%v, gameRoundId=%v, url=%v",
		traceId, gameRoomId, gameRoundId, url)

	return runHttpPost(traceId, msg, url, gameEvent, nil)
}
