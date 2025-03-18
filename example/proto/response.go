package proto

import (
	"encoding/json"
	"time"
)

type Response struct {
	Code int    `json:"code"` // 状态码
	Msg  string `json:"msg"`  // 提示
	Data any    `json:"data"` // 数据
}

func NewResponse(code int, msg string, data any) *Response {
	return &Response{Code: code, Msg: msg, Data: data}
}

func NewResponseError() *Response {
	return &Response{Code: SystemError}
}

func NewResponseByCode(code int) *Response {
	return &Response{Code: code}
}

func NewResponseByCodeMsg(code int, msg string) *Response {
	return &Response{Code: code, Msg: msg}
}

func NewPing() Response {
	return Response{Code: Ping, Data: time.Now().Unix()}
}

func (that *Response) Bytes() []byte {
	buf, _ := json.Marshal(that)
	return buf
}
