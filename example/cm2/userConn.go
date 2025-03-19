package cm2

import (
	"fmt"
	"github.com/PandaManPMC/winterSocket/example/proto2"
	"time"
)

// UserConn 会员连接
type UserConn struct {
	ConnBase
	IdMember      uint64                // 会员编号
	UserToken     string                // 鉴权 Token
	SerialNumber  string                // 连接序号
	ConnectState  uint8                 // 连接状态:1@正常;2@半连接;3@关闭
	TimerCallBack func(*UserConn, bool) // 定时回调
}

const (
	_             = iota
	ConnectState1 // 1@正常
	ConnectState2 // 2@半连接
	ConnectState3 // 3@关闭
)

// Ping 客户端
func (that *UserConn) Ping() error {
	if e := that.Send(proto2.NewPing()); nil != e {
		println(e)
		println(fmt.Sprintf("userConn ping %d-%s", that.IdMember, that.UserToken))
		return e
	}
	that.LastTime = time.Now().Unix()
	return nil
}

// Listener 监听
func (that *UserConn) Listener() {
	defer func() {
		e := recover()
		if nil != e {
			println(e)
			println(fmt.Sprintf("会员 %d-%s Listener 出现异常", that.IdMember, that.UserToken))
			that.Listener()
		}
	}()

	for {
		select {
		case <-time.After(30 * time.Second):
			if ConnectState1 != that.ConnectState {
				println(fmt.Sprintf("会员 %d-%s 连接关闭 停止 Listener", that.IdMember, that.UserToken))
				that.TimerCallBack(that, false)
				return
			}
			// ping 客户端
			println(fmt.Sprintf("ping  客户端 %s", that.UserToken))
			_ = that.Ping()
		case <-time.After(3 * time.Minute):
			if nil != that.TimerCallBack {
				if ConnectState1 != that.ConnectState {
					that.TimerCallBack(that, false)
					return
				} else {
					that.TimerCallBack(that, true)
				}
			}
		}
	}
}
