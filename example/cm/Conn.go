package cm

import (
	"encoding/json"
	"golang.org/x/net/websocket"
	"winterSocket/example/proto"
)

type ConnBase struct {
	Conn        *websocket.Conn // 连接
	ConnectTime int64           // 连接时间
	LastTime    int64           // 最后通讯时间
}

// Send 发送消息
func (that *ConnBase) Send(response proto.Response) error {
	b, err := json.Marshal(response)
	if nil != err {
		return err
	}
	if _, e := that.Conn.Write(b); nil != e {
		return e
	}
	return nil
}
