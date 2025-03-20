package winterSocket

import "net"

var wsTracking WsTrackingInterface

func SetWsTracking(t WsTrackingInterface) {
	wsTracking = t
}

type WsTrackingInterface interface {
	// Connect 新连接
	Connect(*WsConn)
	// RecoverError 出现 panic 被捕获 conn *winterSocket.WsConn, cmd *winterSocket.Cmd, jsonDataByte []byte, err any
	RecoverError(*WsConn, *Cmd, []byte, any)
	// DispatcherBefore 之前
	DispatcherBefore(*WsConn, *Cmd, []byte) bool
	// DispatcherAfter 之后 conn *winterSocket.WsConn, cmd *winterSocket.Cmd, jsonDataByte []byte, resultData []byte
	DispatcherAfter(*WsConn, *Cmd, []byte, []byte)
	// Disconnect 关闭连接
	Disconnect(*net.Conn, any)
	// Dispatcher404 资源未找到 conn *winterSocket.WsConn, route *winterSocket.Cmd, jsonDataByte []byte
	Dispatcher404(*WsConn, *Cmd, []byte)
	// ParameterError 参数错误 conn *winterSocket.WsConn, cmd *winterSocket.Cmd, jsonDataByte []byte, msg string
	ParameterError(*WsConn, *Cmd, []byte, string)
	// ParameterUnmarshalError 数据解析失败 conn *winterSocket.WsConn, cmd *winterSocket.Cmd, jsonDataByte []byte
	ParameterUnmarshalError(*WsConn, *Cmd, []byte)
}
