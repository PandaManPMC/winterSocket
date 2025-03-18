package winterSocket

import "net"

var wsTracking WsTrackingInterface

func SetWsTracking(t WsTrackingInterface) {
	wsTracking = t
}

type WsTrackingInterface interface {
	// Connect 新连接
	Connect(*WsConn)
	// RecoverError 出现 panic 被捕获
	RecoverError(*WsConn, []byte, any)
	// DispatcherBefore 之前
	DispatcherBefore(*WsConn, *Cmd, []byte) bool
	// DispatcherAfter 之后
	DispatcherAfter(*WsConn)
	// Disconnect 关闭连接
	Disconnect(*net.Conn, any)
	// Dispatcher404 资源未找到
	Dispatcher404(*WsConn, *Cmd, []byte)
	// ParameterError 参数错误
	ParameterError(*WsConn, string)
	// ParameterUnmarshalError 数据解析失败
	ParameterUnmarshalError(*WsConn, *Cmd, []byte)
}
