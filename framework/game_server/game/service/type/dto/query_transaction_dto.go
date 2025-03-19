package dto

type QueryTransactionDTO struct {
	UserName     string  `json:"username"`     //站点用户名
	BeginTime    int64   `json:"beginTime"`    //开始时间
	EndTime      int64   `json:"endTime"`      //结束时间
	OrderNoList  []int64 `json:"orderNoList"`  //订单号集合
	SerialNoList []int64 `json:"serialNoList"` //序列号集合
}
