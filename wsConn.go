package winterSocket

import (
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"net/http"
	"sync"
	"time"
)

type WsConn struct {
	Conn     BufferedConn
	Header   http.Header
	ClientIp string
	LastPong time.Time
	mutex    sync.Mutex
	once     sync.Once
}

func (that *WsConn) WriteFrame(frame ws.Frame) error {
	that.mutex.Lock()
	defer that.mutex.Unlock()

	_ = that.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	return ws.WriteFrame(that.Conn, frame)
}

func (that *WsConn) WriteServerText(buff []byte) error {
	that.mutex.Lock()
	defer that.mutex.Unlock()

	_ = that.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	if e := wsutil.WriteServerText(that.Conn, buff); e != nil {
		pError("winterSocket WriteBuff", e)
		return e
	}
	return nil
}

func (that *WsConn) WriteServerTextMust(buff []byte) {
	_ = that.WriteServerText(buff)
}

func (that *WsConn) CloseMust() {
	_ = that.Close()
}

func (that *WsConn) Close() error {
	var err error
	that.once.Do(func() {
		that.mutex.Lock()
		defer that.mutex.Unlock()

		_ = that.Conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
		_ = ws.WriteFrame(that.Conn, ws.NewCloseFrame(nil))

		err = that.Conn.Close()
	})
	return err
}

func (that *WsConn) Ping() error {
	err := that.WriteFrame(ws.NewPingFrame(nil))
	if nil != err {
		pError("write ping failed", err)
		return err
	}
	return nil
}
