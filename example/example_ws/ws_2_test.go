package example

import (
	"fmt"
	"github.com/PandaManPMC/winterSocket"
	"github.com/PandaManPMC/winterSocket/example/example_ws/wshandle"
	"testing"
	"time"
)

func TestInitWs(t *testing.T) {
	winterSocket.SetLog(func(s string) {
		fmt.Println(s)
	}, func(s string) {
		fmt.Println(s)
	}, func(s string, a any) {
		fmt.Println(s)
	})

	//	设置回调
	wsSer := winterSocket.NewWsServer()

	wsSer.SetTracking(new(wshandle.SocketTracking))

	wsSer.PutRoute("login", wshandle.GetInstanceByUserHandle().Login)
	wsSer.PutRoute("ping", wshandle.GetInstanceByUserHandle().Ping)

	if e := wsSer.Listener(uint16(19999)); nil != e {
		panic(e)
	}

}

func TestBytes(t *testing.T) {
	buf := []byte("abc123456&")
	t.Log(buf)
	t.Log(string(buf)) // abc123456
	a := buf[:3]
	t.Log(a)
	a[0] = 'm'
	a[2] = 'p'
	t.Log(string(a))
	t.Log(string(buf)) // mbp123456
}

func panicRecover() {
	defer func() {
		e := recover()
		fmt.Println(e)
	}()

	fmt.Println(time.Now().Unix())
	time.Sleep(time.Second * 3)
	panic(1)
}

func TestGo(t *testing.T) {
	go func() {
		for {
			panicRecover()
		}
	}()
	time.Sleep(100 * time.Minute)
}
