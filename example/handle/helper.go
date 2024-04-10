package handle

import (
	"encoding/json"
	"golang.org/x/net/websocket"
	"winterSocket/example/proto"
)

func Write(conn *websocket.Conn, msg *proto.Response) {
	buf, _ := json.Marshal(msg)
	if _, e := conn.Write(buf); nil != e {
		println(e)
	}
}
