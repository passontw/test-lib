package types

/*Redis相关的Models*/

type RoomCardStatOp string

const (
	RoomCardNumUpdate RoomCardStatOp = "Update"
	RoomCardNumReset  RoomCardStatOp = "Reset"
)

type (
	// LimitInfo 限红信息结构 个人限红和房间限红均可使用该结构
	LimitInfo struct {
		Currency  string
		MinAmount float64
		MaxAmount float64
	}

	// OddInfo 房间赔率信息
	OddInfo struct {
		Odds float32
	}

	//PlayerInfo 玩家信息
	PlayerInfo struct {
		UserId   int64
		Currency string
	}

	// PlayerInRoom 玩家进入房间 离开房间的时候使用
	PlayerInRoom struct {
		PlayerInfoSet []PlayerInfo
	}

	// RoomCardStat 统计当前房间当前靴下牌的数量
	RoomCardStat struct {
		GameRoomId  int64  //房间号
		GameRoundNo string //当前局号 可能存在多节点更新牌的数量所以需要判断局号
		CardNum     int    //当前房间当前靴下牌的数量
	}
)
