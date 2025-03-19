package sign_dto

type SignVerifyDTO struct {
	Id   string `json:"id"`   //数据id
	Text string `json:"text"` //待签名文本
	Sign string `json:"sign"` //签名字符串
}
