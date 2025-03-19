package VO

import types "sl.framework.com/game_server/game/service/type"

type GameResultVO struct {
	GameRoundId string        `json:"gameRoundId"`
	Timestamp   int64         `json:"timestamp"`
	Payload     types.Payload `json:"payload"`
}
