package rpcreq

import (
	"fmt"
	"sl.framework.com/game_server/conf"
	"sl.framework.com/game_server/game/service/type/dto"
)

/*
	GetUserClientInfo 根据用户id获取用户信息
	调用接口:/feign/user/getOne/userId/{userId}}
*/

func GetUserClientInfo(traceId string, userId int64) (*dto.UserDto, int) {
	url := fmt.Sprintf("%v/feign/user/getOne/userId/%v",
		conf.GetPlatformInfoUrl(), userId)
	msg := fmt.Sprintf("GetUserClientInfo traceId=%v, userId=%v, url=%v",
		traceId, userId, url)

	userDto := new(dto.UserDto)
	ret := runHttpGet(traceId, msg, url, userDto)
	return userDto, ret
}

/*
	GetUserClientInfoList 根据用户id集合批量查询用户信息，默认 500条
	调用接口:/feign/user/list}
*/

func GetUserClientInfoList(traceId string) ([]*dto.UserDto, int) {
	url := fmt.Sprintf("%v/feign/user/list",
		conf.GetPlatformInfoUrl())
	msg := fmt.Sprintf("GetUserClientInfo traceId=%v, url=%v",
		traceId, url)

	userDto := make([]*dto.UserDto, 0)
	ret := runHttpGet(traceId, msg, url, userDto)
	return userDto, ret
}
