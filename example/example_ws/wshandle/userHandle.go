package wshandle

import (
	"fmt"
	"github.com/PandaManPMC/winterSocket"
	"github.com/PandaManPMC/winterSocket/example/proto2"
	"sync"
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
func (that *userHandle) Login(conn *winterSocket.WsConn, params proto2.LoginReq) *proto2.Response {
	fmt.Println("Login conn.Header:", conn.Header)
	fmt.Println(conn.ClientIp)

	// 同步
	that.locked()
	defer that.unLocked()

	fmt.Println(params.UserToken)
	//panic(1)

	res := make(map[string]any)
	res["age"] = 100
	res["name"] = "小黑"
	return proto2.NewResponse(proto2.LoginSucceed, "Login 成功", res)
}

// Ping 客户端 ping
func (that *userHandle) Ping(conn *winterSocket.WsConn) {
	fmt.Println("Ping")
}
