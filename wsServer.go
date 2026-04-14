package winterSocket

import (
	"bufio"
	"crypto/rand"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/gobwas/ws"
	"io"
	"math"
	"math/big"
	"net"
	"reflect"
	"sync"
	"time"
)

type WsServer struct {
	dispatcher         func(Cmd, *WsConn, []byte) (bool, []byte)
	tracking           WsTrackingInterface
	route              map[string]reflect.Value
	jsonRouteSeparator byte // 路由 key 长度
	maxDisId           *big.Int
	AcceptSem          int64 // 限流
	MaxMessageSize     int64 // 最大消息
	ReadDeadlineSecond int64
}

func NewWsServer() *WsServer {
	wsServerInstance := new(WsServer)
	wsServerInstance.route = make(map[string]reflect.Value)
	wsServerInstance.jsonRouteSeparator = 38
	wsServerInstance.maxDisId = big.NewInt(math.MaxInt32)
	wsServerInstance.MaxMessageSize = 1 << 20 // 1MB（按业务调整）
	wsServerInstance.AcceptSem = 10240
	wsServerInstance.ReadDeadlineSecond = 90
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

func (that *WsServer) SetDispatcher(dispatcher_ func(Cmd, *WsConn, []byte) (bool, []byte)) {
	that.dispatcher = dispatcher_
}

func (that *WsServer) Disconnect(conn *WsConn) {
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
	disId, _ := rand.Int(rand.Reader, big.NewInt(math.MaxInt32))
	cmd.DisId = disId.Int64()
	cmd.DisTime = time.Now().Unix()

	defer func() {
		if err := recover(); nil != err {
			if nil != that.tracking {
				that.tracking.RecoverError(conn, cmd, jsonDataByte, err)
			} else {
				pError("(that *webSocketServer) handle", err)
			}
			return
		}
	}()

	if nil != that.tracking && !that.tracking.DispatcherBefore(conn, cmd, jsonDataByte) {
		return false
	}

	isResultOk := false
	var resultData []byte

	if nil == that.dispatcher {
		isResultOk, resultData = that.wsJSONDispatcher(cmd, conn, jsonDataByte)

	} else {
		isResultOk, resultData = that.dispatcher(*cmd, conn, jsonDataByte)
	}
	if nil != that.tracking {
		that.tracking.DispatcherAfter(conn, cmd, jsonDataByte, resultData)
	}
	return isResultOk
}

// handleConnection 监听新连接数据
func (that *WsServer) handleConnection(conn *WsConn) {
	defer that.Disconnect(conn)

	if nil != that.tracking {
		that.tracking.Connect(conn)
	}

	// 读取 WebSocket 消息
	for {
		_ = conn.Conn.SetReadDeadline(time.Now().Add(time.Duration(that.ReadDeadlineSecond) * time.Second))

		header, err := ws.ReadHeader(&conn.Conn)
		if nil != err {
			pError("", err)
			return
		}

		if header.Length > that.MaxMessageSize {
			pError("payload too large", header.Length)
			return
		}

		payload := make([]byte, header.Length)
		_, err = io.ReadFull(&conn.Conn, payload)
		if nil != err {
			pError("", err)
			return
		}
		if header.Masked {
			ws.Cipher(payload, header.Mask, 0)
		}

		switch header.OpCode {
		case ws.OpText:
			that.handleJSON(conn, payload)
		case ws.OpBinary:
			// ignore or handle
		case ws.OpContinuation:
			// 不支持分片
			return
		case ws.OpClose:
			_ = conn.WriteFrame(ws.NewCloseFrame(payload))
			return
		case ws.OpPing:
			err = conn.WriteFrame(ws.NewPongFrame(payload))
			if nil != err {
				pError("write pong failed", err)
				return
			}
		case ws.OpPong:
			conn.LastPong = time.Now()
			_ = conn.Conn.SetReadDeadline(time.Now().Add(time.Duration(that.ReadDeadlineSecond) * time.Second))
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
	sem := make(chan struct{}, that.AcceptSem)
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
		buffConn := BufferedConn{Conn: conn, r: bufio.NewReader(io.MultiReader(rawData, conn))}

		remoteAddr := headers.Get("RemoteAddr")
		if "" == remoteAddr {
			headers.Set("RemoteAddr", conn.RemoteAddr().String())
		}
		clientIp := GetRealClientIp(headers)

		// WebSocket 握手
		_, err = ws.Upgrade(&buffConn)
		if nil != err {
			pError("ws.Upgrade(conn)", err)
			if e := conn.Close(); nil != e {
				pError("conn.Close()", e)
			}
			continue
		}

		con := WsConn{
			Conn:     buffConn,
			Header:   headers,
			ClientIp: clientIp,
			LastPong: time.Now(),
			mutex:    sync.Mutex{},
			once:     sync.Once{},
		}

		sem <- struct{}{}
		// 处理链接数据
		go func() {
			defer func() { <-sem }()
			defer func() {
				if e := recover(); nil != e {
					pError("", e)
					_ = con.Close()
					if nil != that.tracking {
						that.tracking.RecoverError(&con, nil, nil, e)
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
