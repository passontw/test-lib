package rpcreq

import (
	"fmt"
	"sl.framework.com/game_server/conf"
	types "sl.framework.com/game_server/game/service/type"
)

/*
	GetRoomInfoRequest 查询房间详情,包括房间限红规则和游戏玩法，以及赔率
	调用接口:/v1/gamecenter/feign/gameRoom/getDetail/gameRoomId/{gameRoomId}
*/

func GetRoomInfoRequest(traceId string, gameRoomId, gameRoundId int64) (*types.RoomDetailedInfo, int) {
	url := fmt.Sprintf("%v/feign/gameRoom/getDetail/gameRoomId/%v", conf.GetPlatformInfoUrl(), gameRoomId)
	msg := fmt.Sprintf("GetRoomInfoRequest traceId=%v, gameRoundId=%v, gameRoomId=%v, url=%v",
		traceId, gameRoundId, gameRoomId, url)

	roomDetailedInfo := new(types.RoomDetailedInfo)
	ret := runHttpGet(traceId, msg, url, roomDetailedInfo)
	return roomDetailedInfo, ret
}
