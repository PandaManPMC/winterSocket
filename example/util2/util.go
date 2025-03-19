package util2

import (
	"net/http"
	"strings"
)

// GetRequestIp 获取客户端 ip，绕过代理，会去除端口号
func GetRequestIp(req *http.Request) string {
	ip := func(req *http.Request) string {
		// 优先使用 X-Forwarded-For
		fIp := req.Header.Get("X-Forwarded-For")
		if "" != fIp && !strings.Contains(fIp, "[") {
			// x.x.x.x,xx.xx.x.x,x.x.x.xx ...
			if strings.Contains(fIp, ",") {
				ips := strings.Split(fIp, ",")
				return ips[0]
			}
			return fIp
		}

		rIp := req.RemoteAddr
		// RemoteAddr=[::1] or 127.0.0.1
		if "" != rIp && !strings.Contains(rIp, "[") && !strings.HasPrefix(rIp, "127.") {
			return rIp
		}

		xIp := req.Header.Get("X-Real-IP")
		if "" != xIp && !strings.Contains(xIp, "[") {
			return xIp
		}

		remoteAddr := req.Header.Get("Remote_addr")
		if "" != remoteAddr && !strings.Contains(remoteAddr, "[") {
			return remoteAddr
		}

		return req.RemoteAddr
	}(req)
	if strings.Contains(ip, ":") {
		return strings.Split(ip, ":")[0]
	}
	return strings.Trim(ip, " ")
}
