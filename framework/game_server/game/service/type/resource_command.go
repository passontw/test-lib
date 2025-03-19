package types

// GameEventCommand 游戏事件命令类型
type GameEventCommand string

const (
	GameEventCommandBetStart    GameEventCommand = "Bet_Start"    //游戏开始 对下局个人限红 房间限红等做预缓存
	GameEventCommandBetStop     GameEventCommand = "Bet_Stop"     //游戏结束 如果Bet_Start没有缓存成功则在此消息中对个人限红 房间限红做预缓存
	GameEventCommandGameDraw    GameEventCommand = "Game_Draw"    //游戏开奖
	GameEventCommandGamePause   GameEventCommand = "Game_Pause"   //游戏换靴 游戏服务器接收到该协议后，房间内牌的数量重置为0
	GameEventCommandGameData    GameEventCommand = "Game_Data"    //荷官发牌 游戏服务器接收到该协议暂不做任何处理
	GameEventCommandGameEnd     GameEventCommand = "Game_End"     //游戏局结束 游戏服务器接收到该协议暂不做任何处理
	GameEventCommandChangeDeck  GameEventCommand = "Change_Deck"  //换牌 游戏服务器接收到该协议暂不做任何处理
	GameEventCommandInvalid     GameEventCommand = "Invalid"      //无效的处理命令
	GameEventCommandCancelRound GameEventCommand = "Cancel_Round" //取消局
	GameEventCommandCancelBet   GameEventCommand = "Bet_Cancel"   //取消下注
	GameEventCommandBet         GameEventCommand = "Bet"          //取消下注
	GameEventCommandBetReceipt  GameEventCommand = "Bet_Receipt"  //结算投注小票
)
