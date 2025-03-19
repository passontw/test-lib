package const_type

type TransactionStatus string

// 状态:交易中
const TransactionStatusDoing TransactionStatus = "Doing"

// 状态:成功
const TransactionStatusSuccess TransactionStatus = "Success"

// 状态:失败
const TransactionStatusFailed TransactionStatus = "Failed"

// 状态:重试
const TransactionStatusRetry TransactionStatus = "Retry"

// 状态:回滚
const TransactionStatusRollback TransactionStatus = "Rollback"
