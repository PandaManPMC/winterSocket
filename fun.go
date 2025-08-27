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
		if "" != fIp {
			if strings.Contains(fIp, "[") {
				fIp = strings.ReplaceAll(fIp, "[", "")
				fIp = strings.ReplaceAll(fIp, "]", "")
			}

			// x.x.x.x,xx.xx.x.x,x.x.x.xx ...
			if strings.Contains(fIp, ",") {
				ips := strings.Split(fIp, ",")
				return strings.TrimSpace(ips[0])
			}
			return strings.TrimSpace(fIp)
		}

		rIp := req.Get("RemoteAddr")
		if strings.Contains(rIp, "[") {
			rIp = strings.ReplaceAll(rIp, "[", "")
			rIp = strings.ReplaceAll(rIp, "]", "")
		}

		// RemoteAddr=[::1] or 127.0.0.1
		if "" != rIp && !strings.HasPrefix(rIp, "127.") {
			return rIp
		}

		xIp := req.Get("X-Real-IP")
		if strings.Contains(xIp, "[") {
			xIp = strings.ReplaceAll(xIp, "[", "")
			xIp = strings.ReplaceAll(xIp, "]", "")
		}
		if "" != xIp {
			return xIp
		}

		remoteAddr := req.Get("Remote_addr")
		if strings.Contains(remoteAddr, "[") {
			remoteAddr = strings.ReplaceAll(remoteAddr, "[", "")
			remoteAddr = strings.ReplaceAll(remoteAddr, "]", "")
		}
		if "" != remoteAddr {
			return remoteAddr
		}

		return req.Get("RemoteAddr")
	}(req)

	return strings.ReplaceAll(ip, " ", "")
}
