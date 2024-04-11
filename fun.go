package winterSocket

import (
	"fmt"
	"golang.org/x/net/websocket"
)

func WriteBuff(buff []byte, conn *websocket.Conn) {
	if c, e := conn.Write(buff); nil != e {
		pError("winterSocket WriteBuff", e)
	} else {
		pDebug(fmt.Sprintf("winterSocket WriteBuff %d", c))
	}
}
