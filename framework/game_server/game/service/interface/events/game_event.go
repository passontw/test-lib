package events

import types "sl.framework.com/game_server/game/service/type"

type IEventHandler interface {
	HandleRondEvent()
}

/**
 * IListenerGameEvent
 * 游戏事件监听接口
 * 游戏相关事件 包括 Bet_Start Bet_Stop Game_Data Game_End Change_Deck Game_Pause
 * Bet_Start: 游戏开始 对下局个人限红 房间限红等做预缓存
 * Bet_Stop: 游戏结束 如果Bet_Start没有缓存成功则在此消息中对个人限红 房间限红做预缓存
 * Game_Data: 荷官发牌 游戏服务器框架接收到该协议暂不做任何处理
 * Game_End: 游戏局结束 游戏服务器框架接收到该协议暂不做任何处理
 * Change_Deck: 换牌 游戏服务器接框架收到该协议暂不做任何处理
 * Game_Pause: 游戏换靴 游戏服务器框架接收到该协议后，房间内牌的数量重置为0
 *
 * 具体游戏如需要该事件 则需要将该接口注册到框架 框架收到该消息会通知到注册的接口
 */

type IListenerGameEvent interface {
	/*
	 * OnGameEventV2
	 * 接收到游戏事件后的处理函数
	 * @param traceId string - traceId 用于日志跟踪
	 * @param event types.GameEventVO - 游戏事件
	 */
	OnGameEvent(traceId string, event types.GameEventVO)

	/*
	 * OnPreEvent
	 * OnGameEvent执行之前执行
	 * @param traceId string - traceId 用于日志跟踪
	 * @param gameRoomId int64 -房间Id
	 * @param gameRoundNo string -局Id
	 */
	OnPreEvent(traceId string, gameRoomId int64, gameRoundNo string)
}
