package events

import types "sl.framework.com/game_server/game/service/type"

/**
 * IListenerJoinOrLeave
 * 玩家进入房间或者离开房间事件监听接口
 *
 * 具体游戏如需要该事件 则需要将该接口注册到框架 框架收到该消息会通知到注册的接口
 */

type IListenerJoinOrLeave interface {
	/*
	 * OnJoinEvent
	 * 接收到玩家进入房间的处理函数
	 * @param traceId string - traceId 用于日志跟踪
	 * @param event types.JoinOrLeaveGameRoom - 玩家进入房间事件
	 */
	OnJoinEvent(traceId string, event types.JoinLeaveGameRoom)

	/*
	 * OnLeaveEvent
	 * 接收到玩家离开房间后的处理函数
	 * @param traceId string - traceId 用于日志跟踪
	 * @param event types.JoinLeaveGameRoom - 玩家离开房间事件
	 */
	OnLeaveEvent(traceId string, event types.JoinLeaveGameRoom)
}
