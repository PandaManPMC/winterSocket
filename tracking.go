package winterSocket

import "golang.org/x/net/websocket"

var tracking TrackingInterface

func SetTracking(t TrackingInterface) {
	tracking = t
}

type TrackingInterface interface {
	// Connect 新连接
	Connect(*websocket.Conn)
	// RecoverError 出现 panic 被捕获
	RecoverError(*websocket.Conn, any)
	// DispatcherBefore 之前
	DispatcherBefore(*websocket.Conn, string, string) bool
	// DispatcherAfter 之后
	DispatcherAfter(*websocket.Conn)
	// Disconnect 关闭连接
	Disconnect(*websocket.Conn)
	// Dispatcher404 资源未找到
	Dispatcher404(*websocket.Conn)
	// ParameterError 参数错误
	ParameterError(*websocket.Conn, string)
	// ParameterUnmarshalError 数据解析失败
	ParameterUnmarshalError(*websocket.Conn)
}
