package proto2

type LoginReq struct {
	UserToken string `json:"userToken" required:"true"` // 鉴权
}
