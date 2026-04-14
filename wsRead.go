package winterSocket

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"net/http"
	"strings"
)

func readAndCopyHeaders(r io.Reader) (http.Header, *bytes.Buffer, error) {
	headers := http.Header{}
	var rawHeaders bytes.Buffer

	tr := io.TeeReader(r, &rawHeaders)
	br := bufio.NewReader(io.LimitReader(tr, 8<<10)) // 限制 8KB

	var requestLine string

	for {
		line, err := br.ReadString('\n')
		if err != nil {
			return nil, nil, err
		}

		if line == "\r\n" {
			break
		}

		line = strings.TrimRight(line, "\r\n")

		// 处理请求行
		if !strings.Contains(line, ":") && requestLine == "" {
			requestLine = line
			continue
		}

		key, value, found := strings.Cut(line, ":")
		if found {
			key = http.CanonicalHeaderKey(strings.TrimSpace(key))
			headers.Add(key, strings.TrimSpace(value))
		}
	}

	return headers, &rawHeaders, nil
}

// BufferedConn **包装连接，确保 `ws.Upgrade()` 仍能读取握手数据**
type BufferedConn struct {
	net.Conn
	r *bufio.Reader
}

func (bc *BufferedConn) Read(p []byte) (int, error) {
	return bc.r.Read(p)
}
