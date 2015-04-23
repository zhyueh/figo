package toolkit

import (
	"net/http"
	"strings"
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
