package rpcreq

import (
	"fmt"
	"sl.framework.com/game_server/conf"
	"sl.framework.com/game_server/rpc_client/config"
)

/*
*
  - 游戏房间主播 上播

调用接口:调用接口:/feign/gameRoomAnchor/signOn/build/{gameRoomId}/{workerId}
*/
func AnchorSignOn(traceId string, gameRoomId, workerId int64) int {
	url := fmt.Sprintf(config.AnchorSignOn, conf.GetPlatformInfoUrl(), gameRoomId, workerId)
	msg := fmt.Sprintf("[游戏房间主播上播] traceId=%v, gameRoomId=%v, workerId=%v, url=%v", traceId,
		gameRoomId, workerId, url)

	ret := runHttpGet(traceId, msg, url, nil)
	return ret
}
