package proto2

const (
	Ping = 1301 // ping
)

const (
	SystemError             = 1444 // 系统错误
	FaultyOperation         = 1445 // 错误的操作
	ParameterError          = 1454 // 参数错误
	ParameterUnmarshalError = 1455 // 参数解析错误
	MethodNotFound          = 1464 // 未找到 404
)

const (
	LoginSucceed = 1500 // 登录成功
	LoginRep     = 1501 // 重复登录
	OffLine      = 1502 // 踢下线
	TokenExpires = 1503 // 登录信息已过期
)
