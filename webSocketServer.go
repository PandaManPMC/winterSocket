package winterSocket

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"net/http"
)

type webSocketServer struct {
	dispatcher func(Method, *websocket.Conn, string) bool
}

var webSocketServerInstance webSocketServer

func GetInstanceByWebSocketServer() *webSocketServer {
	return &webSocketServerInstance
}

func (that *webSocketServer) Disconnect(conn *websocket.Conn) {
	if nil != tracking {
		tracking.Disconnect(conn)
	} else {
		conn.Close()
	}
}

// ws websocket
func (that *webSocketServer) ws() {
	http.Handle("/ws", websocket.Handler(func(conn *websocket.Conn) {
		if nil != tracking {
			tracking.Connect(conn)
		}
		for {
			jsonDataStr := ""
			if e := websocket.Message.Receive(conn, &jsonDataStr); nil != e {
				pError("", e)
				break
			}
			if that.handle(conn, jsonDataStr) {
				continue
			}
		}
		that.Disconnect(conn)
	}))
}

// handle 处理
func (that *webSocketServer) handle(conn *websocket.Conn, jsonDataStr string) bool {
	method := new(Method)
	e := json.Unmarshal([]byte(jsonDataStr), method)
	if nil != e {
		pError("", e)
		return false
	}

	if nil != tracking {
		if !tracking.DispatcherBefore(conn, method.Method, jsonDataStr) {
			return false
		}
	}
	defer func() {
		err := recover()
		if nil != err {
			if nil != tracking {
				tracking.RecoverError(conn, err)
			} else {
				pError("(that *webSocketServer) handle", err)
			}
			return
		}
		if nil != tracking {
			tracking.DispatcherAfter(conn)
		}
	}()

	return that.dispatcher(*method, conn, jsonDataStr)
}

// Listener 开始监听
func (that *webSocketServer) Listener(port uint16, dispatcher func(method Method, conn *websocket.Conn, jsonDataStr string) bool) error {
	that.dispatcher = dispatcher
	that.ws()
	pInfo(fmt.Sprintf("WS Listener %d", port))
	if err := http.ListenAndServe(fmt.Sprintf(
		":%d", port), nil); nil != err {
		pError("ListenAndServe", err)
		return err
	}
	return nil
}

// ListenerWSS 开始监听 WSS
func (that *webSocketServer) ListenerWSS(port uint16, certPath, keyPath string, dispatcher func(method Method, conn *websocket.Conn, jsonDataStr string) bool) error {
	that.dispatcher = dispatcher
	that.ws()
	pInfo(fmt.Sprintf("WSS ListenerWSS %d", port))
	if err := http.ListenAndServeTLS(fmt.Sprintf(
		":%d", port), certPath, keyPath, nil); nil != err {
		pError("WSS ListenAndServe", err)
		return err
	}
	return nil
}

// Handler websocket
// path 路由，如 ws
func (that *webSocketServer) Handler(path string, dispatcher func(method Method, conn *websocket.Conn, jsonDataStr string) bool) {
	that.dispatcher = dispatcher
	http.Handle(path, websocket.Handler(func(conn *websocket.Conn) {
		if nil != tracking {
			tracking.Connect(conn)
		}
		for {
			jsonDataStr := ""
			if e := websocket.Message.Receive(conn, &jsonDataStr); nil != e {
				pError("", e)
				break
			}
			if that.handle(conn, jsonDataStr) {
				continue
			}
		}
		that.Disconnect(conn)
	}))
}
