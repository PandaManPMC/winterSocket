package winterSocket

import (
	"fmt"
	"golang.org/x/net/websocket"
	"net/http"
	"strings"
)

func WriteBuff(buff []byte, conn *websocket.Conn) {
	if c, e := conn.Write(buff); nil != e {
		pError("winterSocket WriteBuff", e)
	} else {
		pDebug(fmt.Sprintf("winterSocket WriteBuff %d", c))
	}
}

func GetRealClientIp(req http.Header) string {
	ip := func(req http.Header) string {
		// 优先使用 X-Forwarded-For
		fIp := req.Get("X-Forwarded-For")
		if "" != fIp && !strings.Contains(fIp, "[") {
			// x.x.x.x,xx.xx.x.x,x.x.x.xx ...
			if strings.Contains(fIp, ",") {
				ips := strings.Split(fIp, ",")
				return ips[0]
			}
			return fIp
		}

		rIp := req.Get("RemoteAddr")
		// RemoteAddr=[::1] or 127.0.0.1
		if "" != rIp && !strings.Contains(rIp, "[") && !strings.HasPrefix(rIp, "127.") {
			return rIp
		}

		xIp := req.Get("X-Real-IP")
		if "" != xIp && !strings.Contains(xIp, "[") {
			return xIp
		}

		remoteAddr := req.Get("Remote_addr")
		if "" != remoteAddr && !strings.Contains(remoteAddr, "[") {
			return remoteAddr
		}

		return req.Get("RemoteAddr")
	}(req)
	if strings.Contains(ip, ":") {
		return strings.Split(ip, ":")[0]
	}
	return strings.Trim(ip, " ")
}
