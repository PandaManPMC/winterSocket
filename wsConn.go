package winterSocket

import (
	"github.com/gobwas/ws/wsutil"
	"net"
	"net/http"
)

type WsConn struct {
	Conn     *net.Conn
	Header   http.Header
	ClientIp string
}

func (that *WsConn) WriteServerText(buff []byte) error {
	if e := wsutil.WriteServerText(*that.Conn, buff); nil != e {
		pError("winterSocket WriteBuff", e)
		return e
	}
	return nil
}
