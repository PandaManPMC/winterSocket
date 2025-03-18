package wshandle

import (
	"fmt"
	"sync"
	"winterSocket"
	"winterSocket/example/proto"
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
func (that *userHandle) Login(conn *winterSocket.WsConn, params proto.LoginReq) *proto.Response {
	fmt.Println("Login conn.Header:", conn.Header)
	fmt.Println(conn.ClientIp)

	// 同步
	that.locked()
	defer that.unLocked()

	fmt.Println(params.UserToken)
	//panic(1)

	return proto.NewResponseByCodeMsg(proto.LoginSucceed, "Login 成功")
}

// Ping 客户端 ping
func (that *userHandle) Ping(conn *winterSocket.WsConn) {
	fmt.Println("Ping")
}
