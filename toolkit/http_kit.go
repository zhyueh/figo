package toolkit

import (
	"errors"
	"net/http"
	"strings"
	"time"
)

const (
	XForwardedFor = "X-Forwarded-For"
	XRealIp       = "X-Real-Ip"
)

func GetRemoteIP(r *http.Request) (remoteIP string) {
	remoteIP = GetHeaderIP(r, XRealIp)
	if remoteIP == "" {
		remoteIP = GetRemoteAddr(r)
	}

	return
}

func GetRemoteAddr(r *http.Request) (remoteIP string) {
	remoteAddr := strings.Split(r.RemoteAddr, ":")
	return remoteAddr[0]
}

func GetHeaderIP(r *http.Request, header string) (remoteIP string) {
	remoteIP = r.Header.Get(header)
	if remoteIP != "" { //可能出现这种 117.169.143.20, 120.203.215.3
		remoteIPS := strings.Split(remoteIP, ",")
		remoteIP = strings.TrimSpace(remoteIPS[len(remoteIPS)-1])
	}
	return
}

type httpresp struct {
	Resp *http.Response
	Err  error
}

func GetQuery(url string, timeout int) (*http.Response, error) {
	c1 := make(chan httpresp, 1)
	go func() {
		resp, err := http.Get(url)
		c1 <- httpresp{Resp: resp, Err: err}
	}()
	select {
	case res := <-c1:
		return res.Resp, res.Err
	case <-time.After(time.Second * 3):
		return nil, errors.New("timeout")
	}
}
