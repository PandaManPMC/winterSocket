package handle

import (
	"fmt"
	"github.com/PandaManPMC/winterSocket/example/cm"
	"github.com/PandaManPMC/winterSocket/example/proto"
	"github.com/PandaManPMC/winterSocket/example/util"
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
func (that *userHandle) Login(conn *websocket.Conn, params proto.LoginReq) *proto.Response {
	xIp := util.GetRequestIp(conn.Request())

	println(xIp)

	// 同步
	that.locked()
	defer that.unLocked()

	now := time.Now().Unix()
	uc := cm.UserConn{
		ConnBase: cm.ConnBase{
			Conn:        conn,
			ConnectTime: now,
			LastTime:    now,
		},
		IdMember:     uint64(time.Now().UnixNano()),
		UserToken:    params.UserToken,
		SerialNumber: fmt.Sprintf("%d", time.Now().UnixNano()),
		ConnectState: cm.ConnectState1,
	}
	uc.TimerCallBack = that.ChangeRecordSession // 隔断时间会回调并且
	if isOk := cm.GetInstanceByConnManager().Login(&uc); isOk {
		// 更新或记录会话
		println(isOk)
	}
	return proto.NewResponseByCode(proto.LoginSucceed)
}

// Ping 客户端 ping
func (that *userHandle) Ping(conn *websocket.Conn) {
	println("Ping LastTime")
	cm.GetInstanceByConnManager().LastTime(conn)
}

// ChangeRecordSession 更新或记录会话
func (that *userHandle) ChangeRecordSession(uc *cm.UserConn, onLine bool) {
	conn := uc.Conn
	xIp := util.GetRequestIp(conn.Request())
	println(fmt.Sprintf("更新或记录会话 %d-%s ip=%s onLine=%v", uc.IdMember, uc.UserToken, xIp, onLine))

	// 判断 token 是否过期，未过期更新会话，过期则发下线通知

	// token 已经过期
	rsp := proto.Response{
		Code: proto.OffLine,
		Msg:  "登录已失效，请重新登录。",
		Data: nil,
	}

	if e := uc.Send(rsp); nil != e {
		println(e)
	}

	_ = cm.GetInstanceByConnManager().OffLine(uc.Conn)
}
