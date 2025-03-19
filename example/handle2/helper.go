package handle2

import (
	"encoding/json"
	"github.com/PandaManPMC/winterSocket/example/proto2"
	"golang.org/x/net/websocket"
)

func Write(conn *websocket.Conn, msg *proto2.Response) {
	buf, _ := json.Marshal(msg)
	if _, e := conn.Write(buf); nil != e {
		println(e)
	}
}
