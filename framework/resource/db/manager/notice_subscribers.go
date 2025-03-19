/**
 * @Author: M
 * @Date: 2024/07/30
 * @Description: 该文件实现了向订阅者发送游戏状态更新通知的功能，涵盖了不同游戏状态的推送，确保所有订阅者能够及时接收相关通知。
 *               主要功能：
 *               1. 向订阅者推送不同游戏状态（如开始、发牌、结束、换鞋等）的通知。
 *               2. 构建相应的 JSON 数据结构，并通过 HTTP 请求将通知发送至订阅者的终端。
 *               3. 使用队列（incomingQueue）确保消息顺序发送，并支持重试机制，确保通知最终成功到达。
 *
 * @Dependencies:
 *               - `github.com/google/uuid`: 用于生成唯一标识符（traceID），确保每次请求都有唯一追踪ID。
 *               - `sl.framework.com/trace`: 用于记录日志，跟踪消息推送的过程和错误。
 *               - `solidleisure.com/dice-prize-resource/gamelogic`: 用于处理游戏逻辑，例如获取下一轮游戏代码。
 *
 * @Usage:
 *               1. `noticeSubStartRound(vid, currentGMCode)`: 通知订阅者游戏开始。
 *               2. `noticeSubNewCard(vid, currentGMCode, dicept)`: 通知订阅者发新牌，指定牌值。
 *               3. `noticeSubCloseRound(vid, currentGMCode, dicept)`: 通知订阅者当前回合结束。
 *               4. `noticeSubCancelRound(vid, currentGMCode)`: 通知订阅者取消当前回合。
 *               5. `noticeSubNewShoe(vid, currentGMCode)`: 通知订阅者更换新的一副牌。
 *               6. `noticeSubStopRound(vid, currentGMCode)`: 通知订阅者停止当前回合。
 *               7. `noticeSubChangeCard(vid, dealer, currentGMCode, shoe)`: 通知订阅者换牌。
 */

package mgr

import (
	"fmt"
	"github.com/google/uuid"
	"sl.framework.com/trace"
	"strings"
)

const (
	baseId        = 1
	maxRetryTimes = 5 // 使用常量代替硬编码值

	// 命令模板
	baseTemplate = `{"command":"%s","gameRoomId":%d,"roundNo":"%s","nextRoundNo":"%s"%s}`
	// 推送通知类型
	startRoundTemplate = `Bet_Start`    // 开局
	stopRoundTemplate  = `Bet_Stop`     // 停局
	newCardCommand     = `Game_Data`    // 发牌
	closeRoundCommand  = `Game_Draw`    // 结算
	endRoundCommand    = `Game_End`     // 关局
	changeCardCommand  = `Change_Deck`  // 改牌
	cancelRoundCommand = `Cancel_Round` // 取消局
	newShoeCommand     = `Game_Pause`   // 换靴

)

// RequestData 包含了向订阅者发送请求的相关数据
// 字段：
//   - EndpointId: 订阅者的唯一标识符，用于区分不同的订阅者。
//   - Endpoint: 请求的目标 URL，表示通知将发送到的地址。
//   - Method: HTTP 请求方法，例如 POST 或 GET。
//   - Body: 请求的 JSON 数据内容，表示发送给订阅者的消息体。
//   - TraceID: 用于追踪请求的唯一标识符，用于在日志或重试中跟踪请求。
//   - GMCode: 当前的游戏代码，表示游戏的当前状态。
//   - TryMaxTimes: 最大重试次数，表示通知失败时最大允许的重试次数。
//   - TryTimes: 当前的重试次数，用于记录已经尝试了多少次通知。
type RequestData struct {
	EndpointId    int64
	Endpoint      string
	Method        string
	Body          string
	TraceID       string
	GMCode        string
	TryMaxTimes   int64
	TryTimes      int64
	isDNSNotFound bool
}

type Notice struct {
	Vid           string
	CurrentGMCode string
	NextGMCode    string
	CardValue     string
	WhosCard      int
	StartTime     int64   // 开始时间戳
	BetSpanTime   int     // 下注剩余秒数
	Number        int8    // 虎
	Dict          [3]byte // 骰宝 骰子点数
}

func BaseNotice(vid, currentGMCode, nextGMCode string) *Notice {
	return &Notice{
		Vid:           vid,
		CurrentGMCode: currentGMCode,
		NextGMCode:    nextGMCode,
	}
}

// sendNotification 向订阅者发送通知的通用方法
// 参数：
//   - vid: 视频 ID，用于标识具体的游戏视频。
//   - currentGMCode: 当前的游戏代码。
//   - nextGMCode: 下局的游戏代码。
//   - command: 通知类型，例如 startRoundTemplate stopRoundTemplate newCardCommand closeRoundCommand changeCardCommand cancelRoundCommand newShoeCommand 等。
//   - payload: 附加的 JSON 数据，例如牌面数据或轮次信息。
//
// 功能：
//   - 根据订阅者信息构建请求数据，并通过队列发送通知。
func sendNotification(vid, currentGMCode, nextGMCode, command, payload string) {
	subscribers := GetSubscribersByVid(vid)
	if len(subscribers) <= 0 {
		trace.Warning("Notifying round for vid: %s gmcode: %s failed, the video has no subscribers", vid, currentGMCode)
		return
	}
	for _, sub := range subscribers {
		traceID := strings.ReplaceAll(uuid.New().String(), "-", "")
		gameRoomId := sub.GameRoomId

		endp, method := EndpointAndMethod(baseId, sub.Endpoint)

		bodyJSON := fmt.Sprintf(baseTemplate, command, gameRoomId, currentGMCode, nextGMCode, payload)
		data := &RequestData{
			EndpointId:  gameRoomId,
			Endpoint:    endp,
			Method:      method,
			Body:        bodyJSON,
			TraceID:     traceID,
			GMCode:      currentGMCode,
			TryMaxTimes: maxRetryTimes,
		}

		trace.Info("Pushing notification: %s(%s) for vid: %s into incomingQueue(%d)", command, bodyJSON, vid, len(incomingQueue))
		incomingQueue <- data
	}
	trace.Info("Notifying round for vid: %s gmcode: %s success, the video has [%d] subscribers", vid, currentGMCode, len(subscribers))
}

// noticeSubStartRound 通知订阅者游戏开始
// 参数：
//   - vid: 视频 ID。
//   - currentGMCode: 当前的游戏代码。
//
// 功能：
//   - 构建游戏开始的 JSON 数据结构，并将其推送给订阅者。
func noticeSubStartRound(vid, currentGMCode, nextGMCode string) {
	trace.Info("Notifying start round for vid: %s gmcode: %s", vid, currentGMCode)
	sendNotification(vid, currentGMCode, nextGMCode, startRoundTemplate, "")
}

// noticeSubNewCard 通知订阅者发新牌
// 参数：
//   - vid: 视频 ID。
//   - currentGMCode: 当前的游戏代码。
//   - dicept: 当前回合结束时的牌值（以 byte 数组表示）。
//
// 功能：
//   - 根据牌的归属构建 JSON 数据结构，并通知订阅者发新牌的结果。
func noticeSubNewCard(vid, currentGMCode, nextGMCode string, dicept [3]byte) {
	trace.Info("Notifying new card for vid: %s gmcode: %s dice: %v", vid, currentGMCode, dicept)

	var diceList = make([]string, 0, 3)
	for _, dice := range dicept {
		diceList = append(diceList, fmt.Sprintf(`"%d"`, dice))
	}

	payload := fmt.Sprintf(`,"payload":{"diceList":[%s]}`, strings.Join(diceList, ","))
	sendNotification(vid, currentGMCode, nextGMCode, newCardCommand, payload)
}

// noticeSubCloseRound 通知订阅者当前回合结束
// 参数：
//   - vid: 视频 ID。
//   - currentGMCode: 当前的游戏代码。
//   - dicept: 当前回合结束时的牌值（以 byte 数组表示）。
//
// 功能：
//   - 构建结束回合的 JSON 数据结构，并将结果推送给订阅者。
func noticeSubCloseRound(vid, currentGMCode, nextGMCode string, dicept [3]byte) {
	trace.Info("Notifying close round for vid: %s gmcode: %s dice: %v", vid, currentGMCode, dicept)

	var diceList = make([]string, 0, 3)
	for _, dice := range dicept {
		diceList = append(diceList, fmt.Sprintf(`"%d"`, dice))
	}

	payload := fmt.Sprintf(`,"payload":{"diceList":[%s]}`, strings.Join(diceList, ","))
	sendNotification(vid, currentGMCode, nextGMCode, closeRoundCommand, payload)
}

// noticeSubEndRound 通知订阅者当前回合结束
// 参数：
//   - vid: 视频 ID。
//   - currentGMCode: 当前的游戏代码。
//   - dragon: 龙牌值列表。
//   - tiger: 虎牌值列表。
//
// 功能：
//   - 构建结束回合的 JSON 数据结构，并将结果推送给订阅者。
func noticeSubEndRound(vid, currentGMCode, nextGMCode string, dicept [3]byte) {
	trace.Info("Notifying end round for vid: %s gmcode: %s dice: %v", vid, currentGMCode, dicept)
	var diceList = make([]string, 0, 3)
	for _, dice := range dicept {
		diceList = append(diceList, fmt.Sprintf(`"%d"`, dice))
	}

	payload := fmt.Sprintf(`,"payload":{"diceList":[%s]}`, strings.Join(diceList, ","))
	sendNotification(vid, currentGMCode, nextGMCode, endRoundCommand, payload)
}

// noticeSubCancelRound 通知订阅者取消当前回合
// 参数：
//   - vid: 视频 ID。
//   - currentGMCode: 当前的游戏代码。
//
// 功能：
//   - 构建取消回合的 JSON 数据结构，并将其推送给订阅者。
func noticeSubCancelRound(vid, currentGMCode, nextGMCode string) {
	trace.Info("Notifying cancel round for vid: %s gmcode: %s", vid, currentGMCode)
	sendNotification(vid, currentGMCode, nextGMCode, cancelRoundCommand, "")
}

// noticeSubNewShoe 通知订阅者换鞋（新的一副牌）
// 参数：
//   - vid: 视频 ID。
//   - currentGMCode: 当前的游戏代码。
//
// 功能：
//   - 构建换鞋 newShoeCommand 的 JSON 数据结构，并推送给订阅者。
func noticeSubNewShoe(vid, currentGMCode, nextGMCode string) {
	trace.Info("Notifying new shoe for vid: %s gmcode: %s", vid, currentGMCode)
	sendNotification(vid, currentGMCode, nextGMCode, newShoeCommand, "")
}

// NoticeSubStopRound 通知订阅者停止当前回合
// 参数：
//   - vid: 视频 ID。
//   - currentGMCode: 当前的游戏代码。
//
// 功能：
//   - 构建停止当前回合的 JSON 数据结构，并推送给订阅者。
func NoticeSubStopRound(vid, currentGMCode, nextGMCode string) {
	trace.Info("Notifying stop round for vid: %s gmcode: %s", vid, currentGMCode)
	sendNotification(vid, currentGMCode, nextGMCode, stopRoundTemplate, "")
}

// noticeSubChangeCard 通知订阅者换牌
// 参数：
//   - vid: 视频 ID。
//   - dealer: 荷官名称。
//   - currentGMCode: 当前的游戏代码。
//   - shoe: 新牌的数量。
//
// 功能：
//   - 构建换牌的 JSON 数据结构，并推送给订阅者。
func noticeSubChangeCard(vid, dealer, currentGMCode, nextGMCode string, shoe int) {
	trace.Info("Notifying change card for vid: %s gmcode: %s shoe: %d", vid, currentGMCode, shoe)
	sendNotification(vid, currentGMCode, nextGMCode, changeCardCommand, "")
}

// Dispatcher 函数，根据传入的类型调用不同的处理函数
func Dispatcher(typ string, param *Notice) {
	vid := param.Vid
	gmCode := param.CurrentGMCode
	if gmCode == "" {
		return
	}
	switch typ {
	case "start":
		go noticeSubStartRound(vid, gmCode, param.NextGMCode)
	case "stop":
		go NoticeSubStopRound(vid, gmCode, param.NextGMCode)
	case "newCard":
		go noticeSubNewCard(vid, gmCode, param.NextGMCode, param.Dict)
	case "newShoe":
		go noticeSubNewShoe(vid, gmCode, param.NextGMCode)
	case "cancel":
		go noticeSubCancelRound(vid, gmCode, param.NextGMCode)
	case "close":
		go noticeSubCloseRound(vid, gmCode, param.NextGMCode, param.Dict)
	case "end":
		go noticeSubEndRound(vid, gmCode, param.NextGMCode, param.Dict)
	default:
		trace.Warning("Unknown type: %s", typ)
	}
}
