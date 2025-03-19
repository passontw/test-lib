package worker

import (
	errcode "sl.framework.com/game_server/error_code"
	"sl.framework.com/game_server/game/service/type/VO"
	"sl.framework.com/game_server/game/service/type/const_type"
	"sl.framework.com/game_server/game/service/type/dto"
	"sl.framework.com/game_server/redis/cache"
	rpcreq "sl.framework.com/game_server/rpc_client"
	"sl.framework.com/trace"
	"strconv"
)

/**
 * 游戏主播上播或下播
 * 1.尝试创建现场员工数据
 * 2.房间主播上播
 */

func SignOn(traceId string, workerSignVO *VO.WorkerSignVO) int {
	trace.Info("[主播登录] traceId: %s,workerSignVO:%+v", traceId, workerSignVO)
	ret := errcode.ErrorOk
	trace.Info("[主播登录] 查询游戏房间缓存 gameRoomCache traceId:%v,gameRoomId:%v", traceId, workerSignVO.GameRoomId)
	gameRoomCache := cache.GameRoomCache{RoomId: workerSignVO.GameRoomId}
	gameRoomCache.Get()
	roomDetailedInfo := gameRoomCache.Data
	if roomDetailedInfo == nil {
		trace.Warning("[主播登录] 游戏主播上播或下播,游戏房间数据查询不存在警告,workerLoginVO:%+v", workerSignVO)
		return errcode.RedisErrorGet
	}
	//请求中台创建荷官信息
	workerDTO := new(dto.WorkerDTO)
	ret = rpcreq.WorkerClientBuild(traceId, workerSignVO.UserName, roomDetailedInfo.GameId, string(const_type.WorkerTypeDealer), workerDTO)
	if ret == errcode.ErrorOk {
		trace.Info("[主播登录] 现场员工信息  traceId:%v,gameRoomId:%v workerDTO:%+v", traceId, workerSignVO.GameRoomId, *workerDTO)
	} else {
		trace.Error("[主播登录] 创建员工信息失败 traceId:%v,workerSignVO:%+v", traceId, workerSignVO)
		return errcode.HttpErrorPlatFormBuildWorkerFailed
	}
	//上播
	workId, _ := strconv.ParseInt(workerDTO.Id, 10, 64)
	ret = rpcreq.AnchorSignOn(traceId, workerSignVO.GameRoomId, workId)
	return ret
}
