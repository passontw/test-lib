package types

type (
	/*
		以下机构用于推送能力中心使用 属于统一结构
		所有具体游戏结构最终都要转为GameDrawOutput接头推送出去
	*/
	//Heads 游戏结果公共头
	Heads struct {
		Token       string `json:"-"`           //token
		GameRoomId  string `json:"gameRoomId"`  //房间Id
		GameRoundId string `json:"gameRoundId"` //局号Id
		GameRoundNo string `json:"gameRoundNo"` //局号
		CardNum     int8   `json:"cardNum"`     //牌型列表中牌的数量
	}

	// Payload 推送到能力中心使用
	Payload struct {
		//具体游戏的结果 不同游戏结构体不同所以使用string类型接收不同游戏的数据
		Data any `json:"data"`
		//中奖的游戏玩法 发给平台中心
		Result any `json:"result"`
	}

	/*
		GameDrawOutput
		游戏服务器对游戏结果进行初步解析 将解析结果发送给能力中心
	*/
	GameDrawOutput struct {
		// 游戏结果头
		Headers *Heads `json:"headers"`
		// 游戏结果数据
		Payload *Payload `json:"payload"`
	}
)
