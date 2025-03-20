package proto2

type LoginReq struct {
	UserToken string `json:"userToken" required:"true" max:"5"` // 鉴权
}
