package winterSocket

type Method struct {
	Method string `json:"method"`
}

type Cmd struct {
	Cmd     string `json:"cmd"`     // 指令
	DisId   int64  `json:"disId"`   // 分发ID
	DisTime int64  `json:"disTime"` // 分发时间
}
