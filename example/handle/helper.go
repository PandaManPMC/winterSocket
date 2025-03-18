package handle

import (
	"encoding/json"
	"github.com/PandaManPMC/winterSocket/example/proto"
	"golang.org/x/net/websocket"
)

func Write(conn *websocket.Conn, msg *proto.Response) {
	buf, _ := json.Marshal(msg)
	if _, e := conn.Write(buf); nil != e {
		println(e)
	}
}
