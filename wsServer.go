package winterSocket

import (
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/gobwas/ws"
	"io"
	"net"
	"reflect"
)

type WsServer struct {
	dispatcher         func(Cmd, *WsConn, []byte) bool
	tracking           WsTrackingInterface
	route              map[string]reflect.Value
	jsonRouteSeparator byte // 路由 key 长度
}

func NewWsServer() *WsServer {
	wsServerInstance := new(WsServer)
	wsServerInstance.route = make(map[string]reflect.Value)
	wsServerInstance.jsonRouteSeparator = 38
	return wsServerInstance
}

// SetJsonRouteSeparator 设置 路由 key 长度
func (that *WsServer) SetJsonRouteSeparator(jsonRouteSeparator_ byte) {
	that.jsonRouteSeparator = jsonRouteSeparator_
}

func (that *WsServer) PutRoute(cmd string, fun any) {
	if _, isOk := route[cmd]; isOk {
		panic(errors.New(fmt.Sprintf("route %s repetition", cmd)))
	}
	route[cmd] = reflect.ValueOf(fun)
	return
}

func (that *WsServer) GetRoute(cmd string) (reflect.Value, bool) {
	v, is := route[cmd]
	return v, is
}

func (that *WsServer) SetTracking(tracking_ WsTrackingInterface) {
	that.tracking = tracking_
}

func (that *WsServer) SetDispatcher(dispatcher_ func(Cmd, *WsConn, []byte) bool) {
	that.dispatcher = dispatcher_
}

func (that *WsServer) Disconnect(conn *net.Conn) {
	if nil != that.tracking {
		that.tracking.Disconnect(conn, nil)
	} else {
		if e := (*conn).Close(); nil != e {
			pError("", e)
		}
	}
}

// handleJSON 处理
func (that *WsServer) handleJSON(conn *WsConn, jsonDataByte_ []byte) bool {

	// 分隔符
	separator := -1
	for i, v := range jsonDataByte_ {
		if that.jsonRouteSeparator == v {
			separator = i
			break
		}
	}

	// 指令
	var command []byte
	var jsonDataByte []byte
	if -1 == separator {
		command = jsonDataByte_
	} else {
		command = jsonDataByte_[:separator]
		jsonDataByte = jsonDataByte_[separator+1:]
	}

	cmd := new(Cmd)
	cmd.Cmd = string(command)
	//if e := json.Unmarshal(jsonDataByte, cmd); nil != e {
	//	pError("", e)
	//	if nil != that.tracking {
	//		that.tracking.ParameterUnmarshalError(conn, jsonDataByte)
	//	}
	//	return false
	//}

	if nil != that.tracking {
		if !that.tracking.DispatcherBefore(conn, cmd, jsonDataByte) {
			return false
		}
	}

	defer func() {
		if err := recover(); nil != err {
			if nil != that.tracking {
				that.tracking.RecoverError(conn, jsonDataByte, err)
			} else {
				pError("(that *webSocketServer) handle", err)
			}
			return
		}
		if nil != that.tracking {
			that.tracking.DispatcherAfter(conn)
		}
	}()

	if nil == that.dispatcher {
		return that.wsJSONDispatcher(cmd, conn, jsonDataByte)
	}
	return that.dispatcher(*cmd, conn, jsonDataByte)
}

// handleConnection 监听新连接数据
func (that *WsServer) handleConnection(conn *WsConn) {
	defer that.Disconnect(conn.Conn)

	if nil == that.tracking {
		that.tracking.Connect(conn)
	}

	// 读取 WebSocket 消息
	for {
		header, err := ws.ReadHeader(*conn.Conn)
		if nil != err {
			pError("", err)
			return
		}

		switch header.OpCode {
		case ws.OpContinuation:
		case ws.OpText:
			payload := make([]byte, header.Length)
			_, err = io.ReadFull(*conn.Conn, payload)
			if nil != err {
				pError("", err)
				return
			}
			if header.Masked {
				ws.Cipher(payload, header.Mask, 0)
			}
			that.handleJSON(conn, payload)
		case ws.OpBinary:
		case ws.OpClose:
			return
		case ws.OpPing:
		case ws.OpPong:
		}
	}
}

// Listener 开始监听
func (that *WsServer) Listener(port uint16) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if nil != err {
		pError("", err)
		return err
	}
	defer func() {
		if e := listener.Close(); nil != e {
			pError("listener.Close()", e)
		}
	}()

	pInfo(fmt.Sprintf("🚀 WebSocket 服务器运行 ws://_:%d", port))

	that.listenerConnect(listener)
	return nil
}

// listenerConnect 处理连接
func (that *WsServer) listenerConnect(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if nil != err {
			pError("listener.Accept()", err)
			continue
		}

		headers, rawData, err := readAndCopyHeaders(conn)
		if nil != err {
			if e := conn.Close(); nil != e {
				pError("", e)
			}
			continue
		}

		// **确保 `ws.Upgrade()` 可以继续读取数据**
		buffConn := &bufferedConn{Conn: conn, r: bufio.NewReader(io.MultiReader(rawData, conn))}

		remoteAddr := headers.Get("RemoteAddr")
		if "" == remoteAddr {
			headers.Set("RemoteAddr", conn.RemoteAddr().String())
		}
		clientIp := GetRealClientIp(headers)

		// WebSocket 握手
		_, err = ws.Upgrade(buffConn)
		if nil != err {
			pError("ws.Upgrade(conn)", err)
			if e := conn.Close(); nil != e {
				pError("conn.Close()", e)
			}
			continue
		}

		con := WsConn{
			Conn:     &buffConn.Conn,
			Header:   headers,
			ClientIp: clientIp,
		}

		// 处理链接数据
		go func() {
			defer func() {
				if e := recover(); nil != e {
					pError("", e)
					if nil != that.tracking {
						that.tracking.RecoverError(&con, nil, e)
					}
				}
			}()
			that.handleConnection(&con)
		}()
	}
}

// ListenerWSS 开始监听 WSS
func (that *WsServer) ListenerWSS(port uint16, certPath, keyPath string) error {
	pInfo(fmt.Sprintf("WSS ListenerWSS %d", port))

	// ** TLS 证书**
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if nil != err {
		pError("tls.LoadX509KeyPair", err)
		return err
	}

	// **TLS 配置**
	tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}

	// **监听 TLS 端口**
	listener, err := tls.Listen("tcp", fmt.Sprintf(":%d", port), tlsConfig)
	if nil != err {
		pError("tls.Listen", err)
		return err
	}

	defer func() {
		if e := listener.Close(); nil != e {
			pError("tls.Listen", e)
		}
	}()

	pInfo(fmt.Sprintf("🚀 WebSocket 服务器运行 wss://_:%d", port))

	that.listenerConnect(listener)
	return nil
}
