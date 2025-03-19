package const_type

type AvailableStatus string

const AvailableStatusAvailable = AvailableStatus("Available") //有效

const AvailableStatusCancel = AvailableStatus("Cancel") //取消

const AvailableStatusResettle = AvailableStatus("Resettle") //重新结算
