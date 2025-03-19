package example

import (
	"github.com/PandaManPMC/winterSocket"
	"github.com/PandaManPMC/winterSocket/example/handle2"
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
	winterSocket.SetTracking(new(handle2.SocketTracking))

	winterSocket.PutRoute("login", handle2.GetInstanceByUserHandle().Login)
	winterSocket.PutRoute("ping", handle2.GetInstanceByUserHandle().Ping)

	if e := winterSocket.GetInstanceByWebSocketServer().Listener(uint16(19999), winterSocket.Dispatcher); nil != e {
		panic(e)
	}
}
