package proto

type LoginReq struct {
	UserToken string `json:"userToken" required:"true"` // 鉴权
}
