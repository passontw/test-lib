package protocol

/**
 * @Author: M
 * @Date: 2024/7/30 11:23
 * @Desc:协议编码常量
 */

const (
	VLGame       = 14
	VLVID        = 4
	VLDealer     = 8
	VLTable      = 4
	VLShoe       = 16
	VLPackHeader = 12
	VLGMType     = 4
)

const (
	CmdDealerKeepalive    uint32 = 0x000001 // DealerKeepalive
	CmdDealerLogin        uint32 = 0x050001 // DealerLogin
	CmdDealerLoginR       uint32 = 0x060001 // DealerLoginR
	CmdDealerLogout       uint32 = 0x050101 // DealerLogout
	CmdNewCard            uint32 = 0x050002 // NewCard
	CmdNewCardRou         uint32 = 0x050002 // NewCardRou
	CmdNewCardV1          uint32 = 0x05D002 // NewCardV1
	CmdNewCardR           uint32 = 0x060002 // NewCardR
	CmdStartGame          uint32 = 0x050003 // StartGame
	CmdStartGameR         uint32 = 0x060003 // StartGameR
	CmdStopBet            uint32 = 0x085059 // StopBet
	CmdDispatchCard       uint32 = 0x060004 // DispatchCard
	CmdNewShoe            uint32 = 0x050005 // NewShoe
	CmdNewShoeR           uint32 = 0x060005 // NewShoeR
	CmdChangeCard         uint32 = 0x050009 // ChangeCard
	CmdChangeCardR        uint32 = 0x060009 // ChangeCardR
	CmdBacRoundRes        uint32 = 0x06000C // BacRoundRes
	CmdRouRoundRes        uint32 = 0x06000b // BacRoundRes
	CmdDTRoundRes         uint32 = 0x06000d // BacRoundRes
	CmdDealerCloseRound   uint32 = 0x05000E // DealerCloseRound
	CmdDealerCloseRoundR  uint32 = 0x06000E // DealerCloseRoundR
	CmdDealerCancelRound  uint32 = 0x05000F // DealerCancelRound
	CmdDealerCancelRoundR uint32 = 0x06000F // DealerCancelRoundR
	CmdSnapshot           uint32 = 0x010016 // BacSnapshot
	CmdBacSnapshotR       uint32 = 0x020016 // BacSnapshotR
	CmdRouSnapshotR       uint32 = 0x020016 // BacSnapshotR
	CmdDTSnapshotR        uint32 = 0x080016 // BacSnapshotR
	CmdBacResultList      uint32 = 0x020008 // BacResultList
	CmdBacGameRes         uint32 = 0x020011 // BacResultList
	CmdRouGameRes         uint32 = 0x020011 // BacResultList
	CmdDTGameRes          uint32 = 0x020029 // BacResultList
	CmdStartJetton        uint32 = 0x02000b // BacResultList
)

const (
	VIDEO_STATUS_CLOSED = "未激活"
	VIDEO_STATUS_OPENED = "激活"
)

const (
	VISIBLE_STATUS_CLOSED = iota
	VISIBLE_STATUS_OPENED
)

const (
	SUB_STATUS_CLOSED = "禁用"
	SUB_STATUS_OPENED = "启用"
)

const (
	GAME_STATUS_CLOSED    = 0  //游戏关闭
	GAME_STATUS_CAN_BET   = 1  //下注状态
	GAME_STATUS_GAME_DATA = 2  //正在发牌
	GAME_STATUS_NEW_SHOE  = 11 //洗牌
	GAME_STATUS_UNKNOW    = 100
)

// 定义状态映射
var (
	StatusToCode = map[string]int{
		"CLOSED":    GAME_STATUS_CLOSED,
		"CAN_BET":   GAME_STATUS_CAN_BET,
		"GAME_DATA": GAME_STATUS_GAME_DATA,
		"NEW_SHOE":  GAME_STATUS_NEW_SHOE,
		"未知状态":      GAME_STATUS_UNKNOW,
	}

	CodeToStatus = map[int]string{
		GAME_STATUS_CLOSED:    "CLOSED",
		GAME_STATUS_CAN_BET:   "CAN_BET",
		GAME_STATUS_GAME_DATA: "GAME_DATA",
		GAME_STATUS_NEW_SHOE:  "NEW_SHOE",
		GAME_STATUS_UNKNOW:    "未知状态",
	}
)

// GetStatusCode returns the integer code for a given status string
func GetStatusCode(status string) int {
	code, _ := StatusToCode[status]
	return code
}

// GetStatusString returns the status string for a given integer code
func GetStatusString(code int) string {
	status, _ := CodeToStatus[code]
	return status
}
