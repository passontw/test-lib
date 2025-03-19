package const_type

// 派奖状态
type PostStatus string

const PostStatusCreate PostStatus = "Create" //创建

const PostStatusReady PostStatus = "Ready" //待派奖

const PostStatusDoing PostStatus = "Doing" //派彩中

const PostStatusInvalid PostStatus = "Invalid" //作废

const PostStatusPaid PostStatus = "Paid" //已派奖

const PostStatusRefund PostStatus = "Refund" //已退款

const PostStatusFailed PostStatus = "Failed" //失败

const PostStatusResettle PostStatus = "Resettle" //重新结算
