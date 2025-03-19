package dto

type UserTransactionDTO struct {
	SerialNo     string  `json:"serialNo"`     //流水号
	SiteSerialNo string  `json:"siteSerialNo"` //站点流水号
	OrderNo      string  `json:"orderNo"`      //订单号
	Currency     string  `json:"currency"`     //币种编号
	Before       float64 `json:"before"`       //变更前金额
	Change       float64 `json:"change"`       //变更金额
	After        float64 `json:"after"`        //变更后金额
	Direction    string  `json:"direction"`    //账变方向：出 Out，入 In
	Title        string  `json:"title"`        //会计科目
	Status       string  `json:"status"`       //状态:创建 Create,入账 Recorded,回滚中 Try_Rollback，已回滚 Rollback
	Timestamp    string  `json:"timestamp"`    //交易完成时间戳
	Version      string  `json:"version"`      //版本号
	Summary      string  `json:"summary"`      //说明
}
