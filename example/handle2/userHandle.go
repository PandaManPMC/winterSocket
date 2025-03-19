package handle2

import (
	"fmt"
	"github.com/PandaManPMC/winterSocket/example/cm2"
	"github.com/PandaManPMC/winterSocket/example/proto2"
	"github.com/PandaManPMC/winterSocket/example/util2"
	"golang.org/x/net/websocket"
	"sync"
	"time"
)

type userHandle struct {
	baseHandle
}

var userHandleInstance userHandle

func GetInstanceByUserHandle() *userHandle {
	return &userHandleInstance
}

func init() {
	userHandleInstance.lock = new(sync.Mutex)
}

// Login 登录
func (that *userHandle) Login(conn *websocket.Conn, params proto2.LoginReq) *proto2.Response {
	xIp := util2.GetRequestIp(conn.Request())

	println(xIp)

	// 同步
	that.locked()
	defer that.unLocked()

	now := time.Now().Unix()
	uc := cm2.UserConn{
		ConnBase: cm2.ConnBase{
			Conn:        conn,
			ConnectTime: now,
			LastTime:    now,
		},
		IdMember:     uint64(time.Now().UnixNano()),
		UserToken:    params.UserToken,
		SerialNumber: fmt.Sprintf("%d", time.Now().UnixNano()),
		ConnectState: cm2.ConnectState1,
	}
	uc.TimerCallBack = that.ChangeRecordSession // 隔断时间会回调并且
	if isOk := cm2.GetInstanceByConnManager().Login(&uc); isOk {
		// 更新或记录会话
		println(isOk)
	}
	return proto2.NewResponseByCode(proto2.LoginSucceed)
}

// Ping 客户端 ping
func (that *userHandle) Ping(conn *websocket.Conn) {
	println("Ping LastTime")
	cm2.GetInstanceByConnManager().LastTime(conn)
}

// ChangeRecordSession 更新或记录会话
func (that *userHandle) ChangeRecordSession(uc *cm2.UserConn, onLine bool) {
	conn := uc.Conn
	xIp := util2.GetRequestIp(conn.Request())
	println(fmt.Sprintf("更新或记录会话 %d-%s ip=%s onLine=%v", uc.IdMember, uc.UserToken, xIp, onLine))

	// 判断 token 是否过期，未过期更新会话，过期则发下线通知

	// token 已经过期
	rsp := proto2.Response{
		Code: proto2.OffLine,
		Msg:  "登录已失效，请重新登录。",
		Data: nil,
	}

	if e := uc.Send(rsp); nil != e {
		println(e)
	}

	_ = cm2.GetInstanceByConnManager().OffLine(uc.Conn)
}
