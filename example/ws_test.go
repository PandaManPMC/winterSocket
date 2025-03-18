package example

import (
	"github.com/PandaManPMC/winterSocket"
	"github.com/PandaManPMC/winterSocket/example/handle"
	"testing"
)

func TestInitSocket(t *testing.T) {
	winterSocket.SetLog(func(s string) {
		println(s)
	}, func(s string) {
		println(s)
	}, func(s string, a any) {
		println(s)
		println(a)
	})
	//	设置回调
	winterSocket.SetTracking(new(handle.SocketTracking))

	winterSocket.PutRoute("login", handle.GetInstanceByUserHandle().Login)
	winterSocket.PutRoute("ping", handle.GetInstanceByUserHandle().Ping)

	if e := winterSocket.GetInstanceByWebSocketServer().Listener(uint16(19999), winterSocket.Dispatcher); nil != e {
		panic(e)
	}
}
