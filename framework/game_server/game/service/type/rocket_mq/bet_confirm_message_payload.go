package rocket_mq

import types "sl.framework.com/game_server/game/service/type"

type BetConfirmMessagePayload struct {
	GameId      int64                    `json:"gameId"`      //游戏Id
	GameRoomId  int64                    `json:"gameRoomId"`  //房间Id
	GameRoundId int64                    `json:"gameRoundId"` //局Id
	UserInfo    []types.UserCurrencyInfo `json:"userInfo"`    //用户Id 列表
}
