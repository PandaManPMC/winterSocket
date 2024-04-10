package handle

import (
	"fmt"
	"golang.org/x/net/websocket"
	"os"
	"strconv"
	"sync/atomic"
	"time"
	"winterSocket/example/cm"
	"winterSocket/example/proto"
	"winterSocket/example/util"
)

type SocketTracking struct {
}

var serialNumber *atomic.Uint64

func init() {
	serialNumber = new(atomic.Uint64)
}

const (
	DispatcherNo        = "DispatcherNo"
	DispatcherMethod    = "DispatcherMethod"
	DispatcherBeginTime = "DispatcherBeginTime"
)

// Connect 新连接
func (*SocketTracking) Connect(conn *websocket.Conn) {
	xIp := util.GetRequestIp(conn.Request())
	println(fmt.Sprintf("webSocketServer ws 新连接%s", xIp))
	cm.GetInstanceByConnManager().RegisterTempConn(conn)
}

// RecoverError 出现 panic 被捕获
func (*SocketTracking) RecoverError(conn *websocket.Conn, err any) {
	no := conn.Request().Header.Get(DispatcherNo)
	method := conn.Request().Header.Get(DispatcherMethod)

	println(err)
	println(fmt.Sprintf("RecoverError %s,%s", method, no))

	// 响应一个系统错误
	Write(conn, proto.NewResponseError())
}

// DispatcherBefore 之前
func (*SocketTracking) DispatcherBefore(conn *websocket.Conn, method, jsonDataStr string) bool {
	xIp := util.GetRequestIp(conn.Request())
	serialNumber.Add(1)
	no := fmt.Sprintf("%d_%d", os.Geteuid(), serialNumber.Load())
	conn.Request().Header.Set(DispatcherNo, no)
	conn.Request().Header.Set(DispatcherMethod, method)
	conn.Request().Header.Set(DispatcherBeginTime, fmt.Sprintf("%d", time.Now().Unix()))
	println(fmt.Sprintf("DispatcherBefore ws %s-%s：%s", xIp, no, jsonDataStr))
	return true
}

// DispatcherAfter 之后
func (*SocketTracking) DispatcherAfter(conn *websocket.Conn) {
	no := conn.Request().Header.Get(DispatcherNo)
	method := conn.Request().Header.Get(DispatcherMethod)
	beginTime, _ := strconv.ParseInt(conn.Request().Header.Get(DispatcherBeginTime), 10, 64)
	println(fmt.Sprintf("DispatcherAfter %s-%s,耗时 %d", method, no, time.Now().Unix()-beginTime))
}

// Disconnect 关闭连接
func (*SocketTracking) Disconnect(conn *websocket.Conn) {
	_ = cm.GetInstanceByConnManager().OffLine(conn)
}

// Dispatcher404 资源未找到
func (*SocketTracking) Dispatcher404(conn *websocket.Conn) {
	// 404
	Write(conn, proto.NewResponseByCode(proto.MethodNotFound))
}

// ParameterError 参数错误
func (*SocketTracking) ParameterError(conn *websocket.Conn, msg string) {
	Write(conn, proto.NewResponseByCodeMsg(proto.ParameterError, msg))
}

// ParameterUnmarshalError 数据解析失败
func (*SocketTracking) ParameterUnmarshalError(conn *websocket.Conn) {
	Write(conn, proto.NewResponseByCode(proto.ParameterUnmarshalError))
}
