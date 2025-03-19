package types

//这里定义 从gameSvr-> ws->client 的指令

type WSMessageCommand string

const (
	WSMessageCommandDynamicOdds WSMessageCommand = "Dynamic_Odds"
)
