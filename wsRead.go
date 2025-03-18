package winterSocket

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"net/http"
	"strings"
)

// readAndCopyHeaders **读取 HTTP 头，同时保留数据**
func readAndCopyHeaders(r io.Reader) (http.Header, *bytes.Buffer, error) {
	headers := http.Header{}
	var rawHeaders bytes.Buffer

	// **创建 TeeReader：读取数据的同时存到 rawHeaders**
	tr := io.TeeReader(r, &rawHeaders)
	br := bufio.NewReader(tr)

	for {
		line, err := br.ReadString('\n')
		if err != nil {
			return nil, nil, err
		}
		if line == "\r\n" { // 头部结束
			break
		}

		// **解析 HTTP 头**
		key, value, found := strings.Cut(line, ":")
		if found {
			headers.Set(strings.TrimSpace(key), strings.TrimSpace(value))
		}
	}

	// **返回读取到的 Headers 和原始数据**
	return headers, &rawHeaders, nil
}

// bufferedConn **包装连接，确保 `ws.Upgrade()` 仍能读取握手数据**
type bufferedConn struct {
	net.Conn
	r *bufio.Reader
}

func (bc *bufferedConn) Read(p []byte) (int, error) {
	return bc.r.Read(p)
}
