package rpcreq

import (
	"fmt"
	"sl.framework.com/game_server/conf"
	"sl.framework.com/game_server/game/service/type/VO"
	types "sl.framework.com/game_server/game/service/type/const_type"
	"sl.framework.com/game_server/game/service/type/dto"
	"sl.framework.com/game_server/rpc_client/config"
	"sl.framework.com/trace"
	"strconv"
)

/**
 * BetRequest
 * 投注
 *
 * @param traceId string - traceId 用于日志跟踪
 * @param currency string - currency 玩家货币类型
 * @param gameRoomId int64 - 游戏房间Id
 * @param gameRoundId int64 - 游戏局Id
 * @param userId int64 - 玩家Id
 * @param betList *[]*Bet -  玩家下注列表
 * @return int - 请求返回码
 */

func BetRequest(traceId, currency, gameRoomId, gameRoundId, userId string,
	betList *[]*dto.BetDTO) int {
	url := fmt.Sprintf(config.BetRequestURL, conf.GetPlatformInfoUrl(), gameRoomId, gameRoundId, userId, currency)
	msg := fmt.Sprintf("BetRequest traceId=%v, gameRoomId=%v, gameRoundId=%v, userId=%v, currency=%+v, url=%v", traceId,
		gameRoomId, gameRoundId, userId, currency, url)
	betDTOList := make([]*VO.BetOrderVO, 0)
	for _, val := range *betList {
		item := new(VO.BetOrderVO)
		item.Id = strconv.FormatInt(val.Id, 10)

		item.UserId = strconv.FormatInt(val.UserId, 10)
		item.Username = val.Username
		item.SiteUsername = val.SiteUsername
		item.Nickname = val.Nickname
		item.GroupId = strconv.FormatInt(val.GroupId, 10)
		item.OrderNo = strconv.FormatInt(val.OrderNo, 10)
		item.GameRoomId = strconv.FormatInt(val.GameRoomId, 10)
		item.GameRoundId = strconv.FormatInt(val.GameRoundId, 10)

		item.GameRoundNo = val.GameRoundNo
		item.GameCategoryId = strconv.FormatInt(val.GameCategoryId, 10)
		item.GameId = strconv.FormatInt(val.GameId, 10)
		item.GameWagerId = strconv.FormatInt(val.GameWagerId, 10)
		item.Currency = val.Currency
		item.Num = val.Num
		item.BetOdds = val.BetOdds
		item.DrawOdds = val.DrawOdds

		item.Type = string(types.Bet)
		item.ClientType = val.ClientType
		item.BetAmount = val.BetAmount
		item.WinAmount = val.WinAmount
		item.AvailableStatus = val.AvailableStatus
		item.ClientStatus = val.ClientStatus
		item.BetStatus = val.BetStatus
		item.WinLostStatus = val.WinLostStatus
		item.PostStatus = val.PostStatus
		item.BetDoneTime = val.BetDoneTime

		item.ManualOn = val.ManualOn
		item.TrialOn = val.TrialOn
		item.Sort = val.Sort
		item.Md5 = val.Md5

		betDTOList = append(betDTOList, item)
		trace.Info("BetRequest traceId=%v url=%v order=%+v", traceId, url, item)
	}
	ret := runHttpPost(traceId, msg, url, &betDTOList, &betDTOList)
	return ret
}
