package config

// 下注请求路径
// 调用接口:/feign/order/bet/{gameRoomId}/{gameRoundId}/{userId}/{currency}
const BetRequestURL = "%v/feign/order/bet/%v/%v/%v/%v"

// 结算
// 调用接口:/feign/orderDraw/result/{siteId}/{gameRoomId}/{gameRoundId}
const SettleURL = "%v/feign/settle/settle/%v/%v"

/*
DrawResultPost 推送游戏结果
调用接口:/feign/game/drawResult/{gameRoomId}/{gameRoundId}
*/
const DrawResultPostURL = "%v/feign/game/drawResult/%v/%v"

/*
*
*创建现场员工信息
调用接口:/feign/worker/build/{username}/{gameId}/{workerType}
*/
const WorkerClientBuild = "%v/feign/worker/build/%v/%v/%v"

/*
*
*游戏房间主播 上播接口
调用接口:/feign/gameRoomAnchor/signOn/build/{gameRoomId}/{workerId}
*/
const AnchorSignOn = "%v/feign/worker/build/%v/%v/%v"
