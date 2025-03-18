package cm

import (
	"encoding/json"
	"fmt"
	"github.com/PandaManPMC/winterSocket/example/proto"
	"github.com/PandaManPMC/winterSocket/example/util"
	"golang.org/x/net/websocket"
	"sync"
	"sync/atomic"
	"time"
)

type connManager struct {
	lock *sync.Mutex

	tempConnPool  *sync.Map // 临时连接池 *websocket.Conn -> *TempConn
	userConnPool  *sync.Map // 登录会员连接池 *websocket.Conn -> *UserConn
	tokenConnPool *sync.Map // 会员连接池 UserToken -> *UserConn

	allCountConnection  *atomic.Int64 // 总连接数量
	tempCountConnection *atomic.Int64 // 临时总连接数量

	ipOnLineConnCount *sync.Map // 每个IP连接数量 ip - > count
}

var connManagerInstance connManager

func GetInstanceByConnManager() *connManager {
	return &connManagerInstance
}

func init() {
	connManagerInstance.lock = new(sync.Mutex)

	connManagerInstance.allCountConnection = new(atomic.Int64)
	connManagerInstance.tempCountConnection = new(atomic.Int64)
	connManagerInstance.tempConnPool = new(sync.Map)
	connManagerInstance.userConnPool = new(sync.Map)
	connManagerInstance.tokenConnPool = new(sync.Map)
	connManagerInstance.ipOnLineConnCount = new(sync.Map)
}

func (that *connManager) locked() {
	that.lock.Lock()
}

func (that *connManager) unLocked() {
	that.lock.Unlock()
}

// RegisterTempConn 注册临时连接
func (that *connManager) RegisterTempConn(conn *websocket.Conn) {
	if _, isOk := that.tempConnPool.Load(conn); isOk {
		return
	}
	tc := new(TempConn)
	tc.Conn = conn
	now := time.Now().Unix()
	tc.LastTime = now
	tc.ConnectTime = now
	that.tempConnPool.Store(conn, tc)
	that.tempCountConnection.Add(1)
}

// OffLine 踢下线
func (that *connManager) OffLine(conn *websocket.Conn) error {
	_, isOk := that.tempConnPool.Load(conn)
	if isOk {
		// 移除出临时连接池
		that.tempCountConnection.Add(-1)
		that.tempConnPool.Delete(conn)
		return nil
	}

	// 移除出会员连接池
	uc, isOk := that.userConnPool.LoadAndDelete(conn)
	if !isOk {
		return nil
	}
	u := uc.(*UserConn)
	println(fmt.Sprintf("OffLine 用户 %d 下线", u.IdMember))

	u.ConnectState = ConnectState2
	that.tokenConnPool.Delete(u.UserToken)
	that.allCountConnection.Add(-1)

	// 减少 ip 在线用户数量
	that.ReduceIpOnLineCount(u)

	return nil
}

// Login 登录，登记到连接池
func (that *connManager) Login(uc *UserConn) bool {

	that.locked()
	defer that.unLocked()

	println(fmt.Sprintf("connManager Login=%s", uc.UserToken))

	_, isOk := that.tempConnPool.Load(uc.Conn)
	if isOk {
		// 移除出临时连接池
		that.tempCountConnection.Add(-1)
		that.tempConnPool.Delete(uc.Conn)
	}

	_, isOk = that.userConnPool.Load(uc.Conn)
	if isOk {
		println(fmt.Sprintf("Login %d 重连 当前在线数 %d", uc.IdMember, that.allCountConnection))
		// 已经登录，属于重连
		that.userConnPool.Delete(uc.Conn)
		that.userConnPool.Store(uc.Conn, uc)
		that.tokenConnPool.Store(uc.UserToken, uc) // 由于没处理重连的 token 的清理，这个会无限增长，不过内存占用在可控范围，暂不处理。
		return true
	}

	// 记录到连接池
	that.userConnPool.Store(uc.Conn, uc)
	that.tokenConnPool.Store(uc.UserToken, uc)
	that.allCountConnection.Add(1)

	// ip 在线用户数量增加
	that.addIpOnLineCount(uc)

	go uc.Listener()
	return true
}

// addIpOnLineCount 增加 ip 在线用户连接数量
func (that *connManager) addIpOnLineCount(uc *UserConn) {
	xIp := util.GetRequestIp(uc.Conn.Request())
	count, isOk := that.ipOnLineConnCount.Load(xIp)
	if isOk {
		c := count.(int)
		c++
		that.ipOnLineConnCount.Store(xIp, c)
	} else {
		that.ipOnLineConnCount.Store(xIp, 1)
	}
}

// ReduceIpOnLineCount 减少 IP 在线用户连接数量
func (that *connManager) ReduceIpOnLineCount(uc *UserConn) {
	xIp := util.GetRequestIp(uc.Conn.Request())
	count, isOk := that.ipOnLineConnCount.Load(xIp)
	if isOk {
		c := count.(int)
		c--
		that.ipOnLineConnCount.Store(xIp, c)
	}
}

// GetIpOnLineCount 获取 IP 在线用户数量
func (that *connManager) GetIpOnLineCount(xIp string) int {
	c, isOk := that.ipOnLineConnCount.Load(xIp)
	if !isOk {
		return 0
	}
	count := c.(int)
	return count
}

// LastTime 更新连接最后通信时间
func (that *connManager) LastTime(conn *websocket.Conn) {
	u, isOk := that.userConnPool.Load(conn)
	if !isOk {
		return
	}
	uc := u.(*UserConn)
	uc.LastTime = time.Now().Unix()
}

// WriteAll 向所有连接者写消息，内部是异步的
func (that *connManager) WriteAll(response *proto.Response) {
	go func() {
		defer func() {
			err := recover()
			if nil != err {
				println(err)
			}
		}()

		bytes, err := json.Marshal(response)
		if nil != err {
			println(err)
			return
		}

		that.tempConnPool.Range(func(key, value any) bool {
			conn := key.(*websocket.Conn)
			write(conn, bytes)
			return true
		})

		that.userConnPool.Range(func(key, value any) bool {
			conn := key.(*websocket.Conn)
			write(conn, bytes)
			return true
		})

	}()
}

func write(conn *websocket.Conn, bytes []byte) {
	if _, e := conn.Write(bytes); nil != e {
		println(e)
	}
}

// OffLineByUserToken 踢下线
func (that *connManager) OffLineByUserToken(userToken, msg string) {
	println(fmt.Sprintf("OffLineByUserToken 踢下线=%s", userToken))
	ucp, isOk := that.tokenConnPool.Load(userToken)
	if !isOk {
		return
	}

	uc := ucp.(*UserConn)
	rsp := proto.Response{
		Code: proto.OffLine,
		Msg:  msg,
		Data: nil,
	}

	if e := uc.Send(rsp); nil != e {
		println(e)
	}

	if e := that.OffLine(uc.Conn); nil != e {
		println(e)
	}
}
