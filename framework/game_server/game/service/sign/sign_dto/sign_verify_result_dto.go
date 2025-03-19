package sign_dto

type SignVerifyResultDTO struct {
	Id string `json:"id"` //数据id
	Ok bool   `json:"ok"` //签名结果    通过：true,不通过：false
}
