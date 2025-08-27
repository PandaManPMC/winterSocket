package winterSocket

import (
	"net/http"
	"testing"
)

func TestGetRealClientIp(t *testing.T) {
	req := make(http.Header)
	req.Set("Accept-Language", "[zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7]")
	req.Set("Remote_addr", "[172.71.151.215]")
	req.Set("X-Forwarded-For", "[2408:8240:43f:cec4:b1c2:5932:4480:4f92, 172.71.151.215]")
	req.Set("Cf-Connecting-Ip", "[2408:8240:43f:cec4:b1c2:5932:4480:4f92]")
	req.Set("Remoteaddr", "[127.0.0.1:58144]")
	req.Set("X-Real-Ip", "[172.71.151.215]")
	t.Log(GetRealClientIp(req))
}
