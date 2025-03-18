package wshandle

import (
	"fmt"
	"github.com/PandaManPMC/winterSocket"
	"github.com/PandaManPMC/winterSocket/example/proto"
	"net"
	"os"
	"sync/atomic"
	"time"
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
func (*SocketTracking) Connect(conn *winterSocket.WsConn) {
	println(fmt.Sprintf("webSocketServer ws 新连接%s", conn.Header))
	//cm.GetInstanceByConnManager().RegisterTempConn(conn)
}

// RecoverError 出现 panic 被捕获
func (*SocketTracking) RecoverError(conn *winterSocket.WsConn, jsonDataByte []byte, err any) {
	no := conn.Header[DispatcherNo]
	method := conn.Header[DispatcherMethod]

	fmt.Println(err)
	println(fmt.Sprintf("RecoverError %s,%s,%s", method, no, string(jsonDataByte)))

	// 响应一个系统错误
	conn.WriteServerText(proto.NewResponseError().Bytes())
}

// DispatcherBefore 之前
func (*SocketTracking) DispatcherBefore(conn *winterSocket.WsConn, cmd *winterSocket.Cmd, jsonDataByte []byte) bool {
	//xIp := util.GetRequestIp(conn.Request())
	xIp := ""
	serialNumber.Add(1)
	no := fmt.Sprintf("%d_%d", os.Geteuid(), serialNumber.Load())
	//conn.Header[DispatcherNo] = no
	//conn.Header[DispatcherMethod] = cmd
	//conn.Header[DispatcherBeginTime] = fmt.Sprintf("%d", time.Now().Unix())
	println(fmt.Sprintf("DispatcherBefore ws %s-%s：%s", xIp, no, jsonDataByte))
	return true
}

// DispatcherAfter 之后
func (*SocketTracking) DispatcherAfter(conn *winterSocket.WsConn) {
	no := conn.Header[DispatcherNo]
	method := conn.Header[DispatcherMethod]
	//beginTime, _ := strconv.ParseInt(conn.Header[DispatcherBeginTime], 10, 64)
	println(fmt.Sprintf("DispatcherAfter %s-%s,耗时 %d", method, no, time.Now().Unix()))
}

// Disconnect 关闭连接
func (*SocketTracking) Disconnect(conn *net.Conn, err any) {
	fmt.Println(fmt.Sprintf("Disconnect 关闭连接 close %v - err=%s", conn, err))
	//_ = cm.GetInstanceByConnManager().OffLine(conn)
}

// Dispatcher404 资源未找到
func (*SocketTracking) Dispatcher404(conn *winterSocket.WsConn, route *winterSocket.Cmd, jsonDataByte []byte) {
	// 404
	fmt.Println("Dispatcher404：" + route.Cmd + " ::: " + string(jsonDataByte))
	conn.WriteServerText(proto.NewResponseByCode(proto.MethodNotFound).Bytes())
}

// ParameterError 参数错误
func (*SocketTracking) ParameterError(conn *winterSocket.WsConn, msg string) {
	conn.WriteServerText(proto.NewResponseByCodeMsg(proto.ParameterError, msg).Bytes())
}

// ParameterUnmarshalError 数据解析失败
func (*SocketTracking) ParameterUnmarshalError(conn *winterSocket.WsConn, cmd *winterSocket.Cmd, jsonDataByte []byte) {
	fmt.Println(cmd.Cmd)
	fmt.Println(string(jsonDataByte))
	conn.WriteServerText(proto.NewResponseByCode(proto.ParameterUnmarshalError).Bytes())
}
