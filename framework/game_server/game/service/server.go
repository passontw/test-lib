package service

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
	"sl.framework.com/async"
	"sl.framework.com/game_server/conf"
	"sl.framework.com/game_server/game/service/interface/bet"
	"sl.framework.com/game_server/game/service/interface/dao"
	"sl.framework.com/game_server/game/service/interface/draw"
	"sl.framework.com/game_server/game/service/interface/events"
	"sl.framework.com/game_server/game/service/interface/ws_message"
	"sl.framework.com/game_server/game/service/sign/interfaces"
	types "sl.framework.com/game_server/game/service/type"
	"sl.framework.com/tool"
	"sl.framework.com/trace"
	"sort"
	"strings"
	"sync"
)

/*
	BaseService 下注和结算服务基础类
	其中有TraceId成员和Init函数
	TraceId用于日志跟踪 在该基础类中统一初始化
*/

type BaseService struct {
	/*
		TraceId用于日志跟踪 在该基础类中统一初始化
	*/
	TraceId string
}

/*
	Init 初始化TraceId 用于日志跟踪
*/

func (b *BaseService) Init(traceId string) {
	b.TraceId = traceId
}

/*
	InterfaceType接口类型
*/

type InterfaceType int

/*
	接口类型常量 非常量范围内的不可注册接口
*/

const (
	InterfaceTypeBettor              InterfaceType = iota + 1 //下注接口类型
	InterfaceTypeDrawer                                       //结算接口类型
	InterfaceTypeGameEventListener                            //游戏事件监听者接口类型
	InterfaceTypeJoinOrLeaveListener                          //玩家离开房间或者加入房间接口类型
	InterfaceTypeDBSaver                                      //保存订单到数据库
	InterfaceTypeSign                                         //校验接口服务
	InterfaceTypeWsMessage                                    //来自ws的消息推送
	InterfaceTypeInvalid                                      //无效的接口类型
)

var (
	/*
		变量无需加锁 只是启动的时候插入数据 程序运行时候从变量中读取数据
	*/
	gameBettorMap     = make(map[types.GameId]bet.IGameBettor) // 下注接口map
	gameBettorPoolMap = make(map[types.GameId]sync.Pool)       // 下注接口池map

	gameDrawerMap     = make(map[types.GameId]draw.IGameDrawer) // 结算接口map
	gameDrawerPoolMap = make(map[types.GameId]sync.Pool)        // 结算接口池map

	gameEventListenerMap   = make(map[types.GameId]events.IListenerGameEvent)   // 游戏事件监听者map
	joinOrLeaveListenerMap = make(map[types.GameId]events.IListenerJoinOrLeave) // 玩家进入游戏或者离开游戏监听者map

	gameDBSaveBatchMap = make(map[types.GameId]dao.IGameDB) // 数据库保存接口map

	gameSignorMap     = make(map[types.GameId]interfaces.ISignHandler)
	gameSignerPoolMap = make(map[types.GameId]sync.Pool)

	gameWsMessageMap     = make(map[types.GameId]ws_message.IWsMessageHandler) //处理来自ws消息的handler
	gameWsMessagePoolMap = make(map[types.GameId]sync.Pool)                    //ws消息handler池

	gameIdMap = make(map[types.GameId]struct{}) // GameId接口map
)

/**
 * RegisterService
 * 服务注册函数
 *
 * @param gameId types.GameId - 游戏Id
 * @param serviceType InterfaceType - 注册的服务类型
 * @param service interfaces{} - 服务接口
 * @return
 */

func RegisterService(gameId types.GameId, serviceType InterfaceType, service interface{}) {
	if serviceType >= InterfaceTypeInvalid {
		panic("RegisterService invalid interfaces type")
	}

	switch serviceType {
	case InterfaceTypeBettor:
		bettorRegister(gameId, service)
	case InterfaceTypeDrawer:
		drawerRegister(gameId, service)
	case InterfaceTypeGameEventListener:
		gameEventListenerRegister(gameId, service)
	case InterfaceTypeJoinOrLeaveListener:
		joinOrLeaveListenerRegister(gameId, service)
	case InterfaceTypeSign:
		signorRegister(gameId, service)
	case InterfaceTypeDBSaver:
		gameDBSaverRegister(gameId, service)
	case InterfaceTypeWsMessage:
		websocketMessageResister(gameId, service)
	default:
		panic("RegisterService invalid interfaces type")
	}

	//将游戏Id放入map中供校验使用
	insertGameId(gameId)
}

/**
 * bettorRegister
 * 下注接口注册以及初始化pool
 *
 * @param gameId string - 游戏Id
 * @param service interfaces{} - 下注接口
 * @return
 */

func bettorRegister(gameId types.GameId, service interface{}) {
	//对服务参数进行类型断言 不是下注接口则报错
	bettor, ok := service.(bet.IGameBettor)
	if !ok {
		panic("service not IGameBettor")
	}

	//将下注接口放入内存中
	if _, ok = gameBettorMap[gameId]; !ok {
		gameBettorMap[gameId] = bettor
		trace.Info("bettorRegister bet.handler gameId=%v success", gameId)
	} else {
		trace.Notice("bettorRegister bet.handler gameId=%v already registered", gameId)
	}

	//初始化下注接口对象池
	if _, ok = gameBettorPoolMap[gameId]; !ok {
		gameBettorPoolMap[gameId] = sync.Pool{
			New: func() interface{} { return NewBettor(gameId) },
		}
		trace.Info("bettorRegister bet.handler pool gameId=%v success", gameId)
	} else {
		trace.Notice("bettorRegister bet.handler pool gameId=%v already registered", gameId)
	}
}

/**
 * gameEventListenerRegister
 * 游戏事件监听接口注册
 *
 * @param gameId string - 游戏Id
 * @param service interfaces{} - 游戏事件监听接口
 * @return
 */

func gameEventListenerRegister(gameId types.GameId, service interface{}) {
	//对服务参数进行类型断言 不是游戏事件监听接口则报错
	listener, ok := service.(events.IListenerGameEvent)
	if !ok {
		panic("service not IListenerGameEvent")
	}

	//将游戏事件监听接口放入内存中
	if _, ok = gameEventListenerMap[gameId]; !ok {
		gameEventListenerMap[gameId] = listener
		trace.Info("gameEventListenerRegister gameId=%v success", gameId)
	} else {
		trace.Notice("gameEventListenerRegister gameId=%v already registered", gameId)
	}
}

/**
 * joinOrLeaveListenerRegister
 * 玩家进入房间或者离开房间事件监听接口注册
 *
 * @param gameId string - 游戏Id
 * @param service interfaces{} - 游戏事件监听接口
 * @return
 */

func joinOrLeaveListenerRegister(gameId types.GameId, service interface{}) {
	//对服务参数进行类型断言 不是玩家进入房间或者离开房间事件监听接口则报错
	listener, ok := service.(events.IListenerJoinOrLeave)
	if !ok {
		panic("service not IListenerJoinOrLeave")
	}

	//将玩家进入房间或者离开房间事件监听接口放入内存中
	if _, ok = joinOrLeaveListenerMap[gameId]; !ok {
		joinOrLeaveListenerMap[gameId] = listener
		trace.Info("joinOrLeaveListenerRegister gameId=%v success", gameId)
	} else {
		trace.Notice("joinOrLeaveListenerRegister gameId=%v already registered", gameId)
	}
}

/**
 * gameDBSaverRegister
 * 数据库保存接口注册
 *
 * @param gameId string - 游戏Id
 * @param service interfaces{} - 数据库保存接口
 * @return
 */

func gameDBSaverRegister(gameId types.GameId, service interface{}) {
	//对服务参数进行类型断言 不是数据库保存接口则报错
	saver, ok := service.(dao.IGameDB)
	if !ok {
		panic("service not IGameDB")
	}

	//将数据库保存接口放入内存中
	if _, ok = gameDBSaveBatchMap[gameId]; !ok {
		gameDBSaveBatchMap[gameId] = saver
		trace.Info("gameDBSaverRegister gameId=%v success", gameId)
	} else {
		trace.Notice("gameDBSaverRegister gameId=%v already registered", gameId)
	}
}

/**
 * websocketMessageResister
 * 数据库保存接口注册
 *
 * @param gameId string - 游戏Id
 * @param service interfaces{} - 数据库保存接口
 * @return
 */

func websocketMessageResister(gameId types.GameId, service interface{}) {
	//对服务参数进行类型断言 不是数据库保存接口则报错
	saver, ok := service.(ws_message.IWsMessageHandler)
	if !ok {
		panic("service not IGameDB")
	}

	//将数据库保存接口放入内存中
	if _, ok = gameWsMessageMap[gameId]; !ok {
		gameWsMessageMap[gameId] = saver
		trace.Info("websocketMessageResister gameId=%v success", gameId)
	} else {
		trace.Notice("websocketMessageResister gameId=%v already registered", gameId)
	}
	//初始化下注接口对象池
	if _, ok := gameWsMessagePoolMap[gameId]; !ok {
		gameWsMessagePoolMap[gameId] = sync.Pool{
			New: func() interface{} { return NewWebsocketMessage(gameId) },
		}
		trace.Info("websocketMessageResister websocket message pool gameId=%v success", gameId)
	} else {
		trace.Notice("websocketMessageResister websocket message gameId=%v already registered", gameId)
	}
}

/**
 * drawerRegister
 * 结算接口注册以及初始化pool
 *
 * @param gameId string 游戏Id
 * @param service IGameDrawer 结算接口
 * @return
 */

func drawerRegister(gameId types.GameId, service interface{}) {
	//对结算接口进行类型校验
	drawer, ok := service.(draw.IGameDrawer)
	if !ok {
		panic("service not IGameDrawer")
	}

	//将结算接口放入到内存中
	if _, ok = gameDrawerMap[gameId]; !ok {
		gameDrawerMap[gameId] = drawer
		trace.Info("drawerRegister draw.handler gameId=%v success", gameId)
	} else {
		trace.Notice("drawerRegister draw.handler gameId=%v already registered", gameId)
	}

	//初始化结算对象池
	if _, ok = gameDrawerPoolMap[gameId]; !ok {
		gameDrawerPoolMap[gameId] = sync.Pool{
			New: func() interface{} { return NewDrawer(gameId) },
		}
		trace.Info("drawerRegister draw.handler pool gameId=%v success", gameId)
	} else {
		trace.Notice("drawerRegister draw.handler pool gameId=%v already registered", gameId)
	}
}

/**
 * RegisterConfiguration
 * 游戏服注册配置文件名字包括后缀名
 *
 * @param fileName string - 配置文件名字
 * @return
 */

func RegisterConfiguration(fileName string) {
	conf.SetConfigurationFileName(fileName)
	trace.Info("RegisterConfiguration fileName=%v", fileName)
}

/**
 * validateService
 * 对注册的服务进行校验 每个gameId下必须有下注接口和结算接口
 */

func ValidateService() {
	for gameId, _ := range gameIdMap {
		//校验下注接口是否存在
		if _, ok := gameBettorMap[gameId]; !ok {
			panic(errors.Errorf("gameId=%v no bet.handler registered", gameId))
		}
		trace.Info("ValidateService gameId=%v have bet.handler service", gameId)

		//校验结算接口是否存在
		if _, ok := gameDrawerMap[gameId]; !ok {
			panic(errors.Errorf("gameId=%v no draw.handler registered", gameId))
		}
		trace.Info("ValidateService gameId=%v have draw.handler service", gameId)

		//检查保存数据接口是否注册
		if _, ok := gameDBSaveBatchMap[gameId]; !ok {
			panic(errors.Errorf("gameId=%v no game order saver registered", gameId))
		}
		trace.Info("ValidateService gameId=%v have game order save service", gameId)
	}
}

/**
 * ValidateGameId
 * 校验GameId是否合法
 *
 * @param gameId - 游戏唯一Id
 * @return bool - true:GameId合法 false:GameId不合法
 */

func ValidateGameId(gameId types.GameId) bool {
	isValid := false
	if _, ok := gameIdMap[gameId]; ok {
		isValid = true
	}

	return isValid
}

/**
 * insertGameId
 * 将游戏Id放入map中供校验使用
 *
 * @param gameId - 游戏唯一Id
 * @return 无返回值
 */

func insertGameId(gameId types.GameId) {
	gameIdMap[gameId] = struct{}{}
	trace.Info("insertGameId gameId=%v", gameId)
}

/**
 * BuildMQTags
 * 根据GameIdMap组装mq tag
 *
 * @param
 * @return string - 返回mq tag
 */

func BuildMQTags() string {
	var keys []string
	for key := range gameIdMap {
		keys = append(keys, key.ToString())
	}
	sort.Strings(keys)

	strKeys := strings.Join(keys, "||")
	trace.Info("BuildMQTags keys=%v", strKeys)

	return strKeys
}

/**
 * AsyncNotifyGameEventListener
 * 异步将游戏事件通知到对应接口 避免具体游戏接口耗时过长影响到框架
 *
 * @param traceId string - traceId用于日志跟踪
 * @param event types.GameEventMessageHeader - 游戏事件消息
 * @return
 */

func AsyncNotifyGameEventListener(traceId string, event types.GameEventVO) {
	async.AsyncRunCoroutine(func() { notifyGameEventListener(traceId, event) })
}

/**
 * AsyncNotifyGameEventListener
 * 异步将游戏事件通知到对应接口 避免具体游戏接口耗时过长影响到框架
 *
 * @param traceId string - traceId用于日志跟踪
 * @param event types.GameEventMessageHeader - 游戏事件消息
 * @return
 */

func AsyncNotifyGameEventListenerV2(traceId string, event types.GameEventVO) {
	async.AsyncRunCoroutine(func() { notifyGameEventToListener(traceId, event) })
}

/**
 * notifyGameEventListener
 * 将游戏事件通知到对应接口
 *
 * @param traceId string - traceId用于日志跟踪
 * @param event types.GameEventMessageHeader - 游戏事件消息
 * @return
 */

func notifyGameEventListener(traceId string, event types.GameEventVO) {
	msgHeader := fmt.Sprintf("notifyGameEventListener traceId=%v, event=%+v", traceId, event)
	listener, ok := gameEventListenerMap[types.GameId(conf.GetGameId())]
	if !ok {
		trace.Error("%v, no game event listener registered with game id=%v", msgHeader, conf.GetGameId())
		return
	}
	trace.Info("%v", msgHeader)

	pWatcher := tool.NewWatcher(msgHeader)
	reflectVal := reflect.ValueOf(listener)
	reflectTyp := reflect.Indirect(reflectVal).Type()
	vc := reflect.New(reflectTyp)
	exeListener, okk := vc.Interface().(events.IListenerGameEvent)
	if !okk {
		trace.Error("%v, wrong listener type registered with game id=%v", msgHeader, conf.GetGameId())
		return
	}
	//执行游戏事件前的预处理逻辑
	exeListener.OnPreEvent(traceId, event.GameRoomId, event.GameRoundNo)
	//执行游戏事件相关逻辑
	exeListener.OnGameEvent(traceId, event)
	pWatcher.Stop()
}

/**
 * notifyGameEventToListener
 * 将游戏事件通知到对应接口
 *
 * @param traceId string - traceId用于日志跟踪
 * @param event types.GameEventVO - 游戏事件消息
 * @return
 */

func notifyGameEventToListener(traceId string, event types.GameEventVO) {
	msgHeader := fmt.Sprintf("notifyGameEventListener traceId=%v, event=%+v", traceId, event)
	listener, ok := gameEventListenerMap[types.GameId(conf.GetGameId())]
	if !ok {
		trace.Error("%v, no game event listener registered with game id=%v", msgHeader, conf.GetGameId())
		return
	}
	trace.Info("%v", msgHeader)

	pWatcher := tool.NewWatcher(msgHeader)
	reflectVal := reflect.ValueOf(listener)
	reflectTyp := reflect.Indirect(reflectVal).Type()
	vc := reflect.New(reflectTyp)
	exeListener, okk := vc.Interface().(events.IListenerGameEvent)
	if !okk {
		trace.Error("%v, wrong listener type registered with game id=%v", msgHeader, conf.GetGameId())
		return
	}
	exeListener.OnGameEvent(traceId, event)
	pWatcher.Stop()
}

/**
 * AsyncNotifyJoinLeaveGameRoomListener
 * 将游戏事件通知到对应接口 避免具体游戏接口耗时过长影响到框架
 *
 * @param traceId string - traceId用于日志跟踪
 * @param gameId types.GameId - 游戏Id
 * @param event types.GameEventMessageHeader - 游戏事件消息
 * @return
 */

func AsyncNotifyJoinLeaveGameRoomListener(traceId string, action types.JoinLeaveGameRoom) {
	async.AsyncRunCoroutine(func() { notifyJoinLeaveGameRoomListener(traceId, action) })
}

/**
 * notifyJoinLeaveGameRoomListener
 * 将游戏事件通知到对应接口
 *
 * @param traceId string - traceId用于日志跟踪
 * @param event types.JoinOrLeaveGameRoom - 进入房间和离开房间事件消息
 * @return
 */

func notifyJoinLeaveGameRoomListener(traceId string, event types.JoinLeaveGameRoom) {
	msgHeader := fmt.Sprintf("notifyJoinLeaveGameRoomListener traceId=%v, , event=%+v", traceId, event)
	listener, ok := joinOrLeaveListenerMap[types.GameId(conf.GetGameId())]
	if !ok {
		trace.Error("%v, no game event listener registered with game id=%v", msgHeader, event.GameId)
		return
	}
	trace.Info("%v", msgHeader)

	pWatcher := tool.NewWatcher(msgHeader)
	reflectVal := reflect.ValueOf(listener)
	reflectTyp := reflect.Indirect(reflectVal).Type()
	vc := reflect.New(reflectTyp)
	exeListener, okk := vc.Interface().(events.IListenerJoinOrLeave)
	if !okk {
		trace.Error("%v, wrong listener type registered with game id=%v", msgHeader)
		return
	}

	switch event.Type {
	case types.RoomActionJoin:
		exeListener.OnJoinEvent(traceId, event)
	case types.RoomActionLeave:
		exeListener.OnLeaveEvent(traceId, event)
	default:
		trace.Error("%v, no room action exist.", msgHeader)
	}

	pWatcher.Stop()
}

/**
 * GetBettor
 * 从对象池获取一个下注接口
 *
 * @param - gameId types.GameId 游戏Id
 * @param - traceId string 日志跟踪Id
 * @return - IGameBettor接口
 */

func GetBettor(traceId string, gameId types.GameId) bet.IGameBettor {
	pool, ok := gameBettorPoolMap[gameId]
	if !ok {
		trace.Error("getBettor no bet.handler pool, traceId=%v, gameId=%v", traceId, gameId)
		return nil
	}
	bettor := pool.Get().(bet.IGameBettor)
	bettor.Init(traceId)

	return bettor
}

/**
 * PutBettor
 * 将一个下注接口归还给资源池
 *
 * @param - gameId types.GameId 游戏Id
 * @param - bet.handler IGameBettor接口
 * @return
 */

func PutBettor(gameId types.GameId, bettor bet.IGameBettor) {
	pool, ok := gameBettorPoolMap[gameId]
	if !ok {
		trace.Error("putBettor no bet.handler pool, gameId=%v", gameId)
		return
	}

	pool.Put(bettor)
}

// NewBettor 根据GameId创建新的下注对象
func NewBettor(gameId types.GameId) bet.IGameBettor {
	if bettor, ok := gameBettorMap[gameId]; ok {
		reflectVal := reflect.ValueOf(bettor)
		reflectTyp := reflect.Indirect(reflectVal).Type()
		vc := reflect.New(reflectTyp)
		execBettor, okk := vc.Interface().(bet.IGameBettor)
		if !okk {
			trace.Error("newBettor wrong bet.handler type registered with game id=%v", gameId)
			return nil
		}
		trace.Notice("newBettor new bet.handler for game id=%v", gameId)

		return execBettor
	}

	trace.Error("newBettor no game bet.handler registered with game id=%v", gameId)
	return nil
}

/**
 * GetDrawer
 * 从对象池获取一个下注接口
 *
 * @param - traceId string 日志跟踪Id
 * @param - gameId types.GameId 游戏Id
 * @return - IGameDrawer接口
 */

func GetDrawer(traceId string, gameId types.GameId) draw.IGameDrawer {
	pool, ok := gameDrawerPoolMap[gameId]
	if !ok {
		trace.Error("getDrawer no bet.handler pool, traceId=%v, gameId=%v", traceId, gameId)
		return nil
	}
	drawer := pool.Get().(draw.IGameDrawer)
	drawer.Init(traceId)

	return drawer
}

/**
 * PutDrawer
 * 将一个下注接口归还给资源池
 *
 * @param - gameId types.GameId 游戏Id
 * @param - draw.handler IGameDrawer 结算接口
 * @return
 */

func PutDrawer(gameId types.GameId, drawer draw.IGameDrawer) {
	pool, ok := gameDrawerPoolMap[gameId]
	if !ok {
		trace.Error("putDrawer no draw.handler pool, gameId=%v", gameId)
		return
	}

	pool.Put(drawer)
}

// NewDrawer 根据GameId创建新的结算对象
func NewDrawer(gameId types.GameId) draw.IGameDrawer {
	if drawer, ok := gameDrawerMap[gameId]; ok {
		reflectVal := reflect.ValueOf(drawer)
		reflectTyp := reflect.Indirect(reflectVal).Type()
		vc := reflect.New(reflectTyp)
		execDrawer, okk := vc.Interface().(draw.IGameDrawer)
		if !okk {
			trace.Error("newDrawer wrong draw.handler type registered with game id=%v", gameId)
			return nil
		}
		trace.Notice("newDrawer new draw.handler for game id=%v", gameId)
		return execDrawer
	}

	trace.Error("newDrawer no game draw.handler registered with game id=%v", gameId)
	return nil
}

/**
 * NewGameDBSaver
 * 创建一个数据库保存对象
 *
 * @param - gameId types.GameId 游戏Id
 * @param - traceId string 日志跟踪Id
 * @return - IGameBettor接口
 */

func NewGameDBSaver(traceId string, gameId types.GameId) dao.IGameDB {
	if saver, ok := gameDBSaveBatchMap[gameId]; ok {
		reflectVal := reflect.ValueOf(saver)
		reflectTyp := reflect.Indirect(reflectVal).Type()
		vc := reflect.New(reflectTyp)
		execSaver, okk := vc.Interface().(dao.IGameDB)
		if !okk {
			trace.Error("newGameDBSaver wrong order saver type, registered with game id=%v, traceId=%v", gameId, traceId)
			return nil
		}
		return execSaver
	}

	trace.Error("newGameDBSaver no game order saver registered with game id=%v, traceId=%v", gameId)
	return nil
}

/**
 * signorRegistor
 * 注册校验类
 *
 * @param gameId string - 游戏id
 * @param constructor any - sign的实现类
 * @return  -
 */

func signorRegister(gameId types.GameId, constructor any) {
	conType := reflect.TypeOf(constructor)
	signor, ok := constructor.(interfaces.ISignHandler)
	if !ok {
		panic(fmt.Sprintf("signorRegistor constructor %s does not implement ISignHandler", conType.Name()))
	}

	//将下注接口放入内存中
	if _, ok := gameSignorMap[gameId]; !ok {
		gameSignorMap[gameId] = signor
		trace.Info("signorRegistor gameId=%v success", gameId)
	} else {
		trace.Notice("signorRegistor gameId=%v already registered", gameId)
	}

	//初始化下注接口对象池
	if _, ok := gameSignerPoolMap[gameId]; !ok {
		gameSignerPoolMap[gameId] = sync.Pool{
			New: func() interface{} { return NewSignor(gameId) },
		}
		trace.Info("bettorRegister bet.handler pool gameId=%v success", gameId)
	} else {
		trace.Notice("bettorRegister bet.handler pool gameId=%v already registered", gameId)
	}
}

// NewSignor 根据GameId创建新的校验对象
func NewSignor(gameId types.GameId) interfaces.ISignHandler {
	if signor, ok := gameSignorMap[gameId]; ok {
		reflectVal := reflect.ValueOf(signor)
		reflectTyp := reflect.Indirect(reflectVal).Type()
		vc := reflect.New(reflectTyp)
		execSignor, okk := vc.Interface().(interfaces.ISignHandler)
		if !okk {
			trace.Error("newSignor wrong signor type registered with game id=%v", gameId)
			return nil
		}
		trace.Notice("newSignor new signor for game id=%v", gameId)

		return execSignor
	}

	trace.Error("newSignor no game signor registered with game id=%v", gameId)
	return nil
}

/**
 * GetSignor
 * 从对象池获取一个校验接口
 *
 * @param - traceId string 日志跟踪Id
 * @param - gameId types.GameId 游戏Id
 * @return - ISignHandler接口
 */

func GetSignor(traceId string, gameId types.GameId) interfaces.ISignHandler {
	pool, ok := gameSignerPoolMap[gameId]
	if !ok {
		trace.Error("getSignor no signor pool, traceId=%v, gameId=%v", traceId, gameId)
		return nil
	}
	signor := pool.Get().(interfaces.ISignHandler)
	return signor
}

/**
 * PutSignor
 * 将一个校验接口归还给资源池
 *
 * @param - gameId types.GameId 游戏Id
 * @param - draw.handler IGameDrawer 结算接口
 * @return
 */

func PutSignor(gameId types.GameId, signor interfaces.ISignHandler) {
	pool, ok := gameSignerPoolMap[gameId]
	if !ok {
		trace.Error("putSignor no signor pool, gameId=%v", gameId)
		return
	}

	pool.Put(signor)
}

/**
 * GetWebsocketMessage
 * 从对象池获取一个消息处理接口
 *
 * @param - traceId string 日志跟踪Id
 * @param - gameId types.GameId 游戏Id
 * @return - ISignHandler接口
 */

func GetWebsocketMessage(traceId string, gameId types.GameId) ws_message.IWsMessageHandler {
	pool, ok := gameWsMessagePoolMap[gameId]
	if !ok {
		trace.Error("GetWebsocketMessage no websocket message pool, traceId=%v, gameId=%v", traceId, gameId)
		return nil
	}
	message := pool.Get().(ws_message.IWsMessageHandler)
	return message
}

/**
 * NewWebsocketMessage
 * 创建一个websocket message保存对象
 *
 * @param - gameId types.GameId 游戏Id
 * @return - IGameBettor接口
 */

func NewWebsocketMessage(gameId types.GameId) ws_message.IWsMessageHandler {
	if message, ok := gameWsMessageMap[gameId]; ok {
		reflectVal := reflect.ValueOf(message)
		reflectTyp := reflect.Indirect(reflectVal).Type()
		vc := reflect.New(reflectTyp)
		execMessage, okk := vc.Interface().(ws_message.IWsMessageHandler)
		if !okk {
			trace.Error("newMessage wrong message type registered with game id=%v", gameId)
			return nil
		}
		trace.Notice("newMessage new message for game id=%v", gameId)

		return execMessage
	}

	trace.Error("newMessage no message registered with game id=%v", gameId)
	return nil
}

/**
 * PutMessage
 * 将一个websocket message接口归还给资源池
 *
 * @param - gameId types.GameId 游戏Id
 * @param - draw.handler IGameDrawer 结算接口
 * @return
 */

func PutMessage(gameId types.GameId, message ws_message.IWsMessageHandler) {
	pool, ok := gameWsMessagePoolMap[gameId]
	if !ok {
		trace.Error("PutMessage no message pool, gameId=%v", gameId)
		return
	}

	pool.Put(message)
}
