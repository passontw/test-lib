package draw

import (
	"encoding/json"
	"sl.framework.com/game_server/error_code"
	"sl.framework.com/game_server/game/service/type/rocket_mq"
	"sl.framework.com/game_server/mq/handler"
	"sl.framework.com/trace"
)

/**
 * ParseMessage
 * 解析RocketMq 的消息
 *
 * @param traceId string - 跟踪Id
 * @param msgBody []byte 消息体
 * @return int - 返回码
 */

func ParseMessage(traceId string, msgBody []byte) int {
	trace.Info("[解析RocketMq消息] traceId=%v 包体长度=%v", traceId, string(msgBody))

	//解析消息类型
	var rocketMqMessage = new(rocket_mq.RocketMqMessageDTO)
	if err := json.Unmarshal(msgBody, &rocketMqMessage); err != nil {
		trace.Error("[解析RocketMq消息] traceId=%v json unmarshal failed, error=%v, data=%v", traceId, err.Error(), string(msgBody))
		return errcode.HttpErrorOrderParse
	}
	switch rocketMqMessage.MessageType {
	case string(rocket_mq.RocketMessageTypeGameDraw):
		//结算消息
		return handler.OnGameDrawHandler(traceId, rocketMqMessage.MessagePayload)
	case string(rocket_mq.RocketMessageTypeBetConfirm):
		//提交注单消息
		return handler.OnBetConfirmHandler(traceId, rocketMqMessage.MessagePayload)
	default:
		trace.Error("[解析RocketMq消息] traceId=%v 未知的消息类型：%v 包体内容：%v！", traceId, rocketMqMessage.MessageType, string(msgBody))
		return errcode.HttpErrorDataFailed
	}
}
