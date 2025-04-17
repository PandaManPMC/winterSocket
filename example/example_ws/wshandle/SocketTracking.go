package wshandle

import (
	"fmt"
	"github.com/PandaManPMC/winterSocket"
	"github.com/PandaManPMC/winterSocket/example/proto2"
	"sync/atomic"
	"time"
)

type SocketTracking struct {
}

var serialNumber *atomic.Uint64

func init() {
	serialNumber = new(atomic.Uint64)
}

// Connect 新连接
func (*SocketTracking) Connect(conn *winterSocket.WsConn) {
	println(fmt.Sprintf("webSocketServer ws 新连接%s", conn.Header))
	//cm2.GetInstanceByConnManager().RegisterTempConn(conn)
}

// RecoverError 出现 panic 被捕获
func (*SocketTracking) RecoverError(conn *winterSocket.WsConn, cmd *winterSocket.Cmd, jsonDataByte []byte, err any) {
	println(fmt.Sprintf("RecoverError %s,%d,%s", cmd.Cmd, cmd.DisId, string(jsonDataByte)))

	// 响应一个系统错误
	conn.WriteServerText(proto2.NewResponseError().Bytes())
}

// DispatcherBefore 之前
func (*SocketTracking) DispatcherBefore(conn *winterSocket.WsConn, cmd *winterSocket.Cmd, jsonDataByte []byte) bool {
	xIp := ""
	serialNumber.Add(1)
	println(fmt.Sprintf("DispatcherBefore ws %s-%d：%s", xIp, cmd.DisId, jsonDataByte))
	return true
}

// DispatcherAfter 之后
func (*SocketTracking) DispatcherAfter(conn *winterSocket.WsConn, cmd *winterSocket.Cmd, jsonDataByte []byte, resultData []byte) {
	fmt.Println(string(resultData))
	fmt.Println(string(jsonDataByte))
	println(fmt.Sprintf("DispatcherAfter %s-%d,耗时 %d", cmd.Cmd, cmd.DisId, time.Now().Unix()-cmd.DisTime))
}

// Disconnect 关闭连接
func (*SocketTracking) Disconnect(conn *winterSocket.WsConn, err any) {
	fmt.Println(fmt.Sprintf("Disconnect 关闭连接 close %v - err=%s", conn, err))
	//_ = cm2.GetInstanceByConnManager().OffLine(conn)
}

// Dispatcher404 资源未找到s
func (*SocketTracking) Dispatcher404(conn *winterSocket.WsConn, route *winterSocket.Cmd, jsonDataByte []byte) {
	// 404
	fmt.Println("Dispatcher404：" + route.Cmd + " ::: " + string(jsonDataByte))
	conn.WriteServerText(proto2.NewResponseByCode(proto2.MethodNotFound).Bytes())
}

// ParameterError 参数错误
func (*SocketTracking) ParameterError(conn *winterSocket.WsConn, cmd *winterSocket.Cmd, jsonDataByte []byte, msg string) {
	conn.WriteServerText(proto2.NewResponseByCodeMsg(proto2.ParameterError, msg).Bytes())
}

// ParameterUnmarshalError 数据解析失败
func (*SocketTracking) ParameterUnmarshalError(conn *winterSocket.WsConn, cmd *winterSocket.Cmd, jsonDataByte []byte) {
	fmt.Println(cmd.Cmd)
	fmt.Println(string(jsonDataByte))
	conn.WriteServerText(proto2.NewResponseByCode(proto2.ParameterUnmarshalError).Bytes())
}
