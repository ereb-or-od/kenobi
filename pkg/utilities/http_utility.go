package utilities

import (
	"fmt"
	"github.com/xiam/to"
	"net"
	"net/http"
)

func GetForwardedIP(r *http.Request) string {
	return r.Header.Get("X-Forwarded-For")
}
func GetIP(r *http.Request) (string, error) {
	ip := to.String(r.Context().Value("ip"))

	ip = r.Header.Get("X-Real-IP")
	if len(ip) > 0 {
		return ip, nil
	}

	// no nginx reverse proxy?
	// get IP old fashioned way
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", fmt.Errorf("%q is not IP:port", r.RemoteAddr)
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		return "", fmt.Errorf("%q is not IP:port", r.RemoteAddr)
	}
	return userIP.String(), nil
}



