package const_type

type ClientStatus string

// 显示状态：已支付 Paid，已结算 Settled，取消 Cancel，结算失败 Settled_Failed

const ClientStatusPaid ClientStatus = "Paid"

const ClientStatusSettled ClientStatus = "Settled"

const ClientStatusCancel ClientStatus = "Cancel"

const ClientStatusSettledFailed ClientStatus = "Settled_Failed"
