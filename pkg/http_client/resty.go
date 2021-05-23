
package http_client

import (
	"net"
	"net/http"
	"net/http/cookiejar"

	"golang.org/x/net/publicsuffix"
)

func New() *HttpClient {
	cookieJar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	return createClient(&http.Client{
		Jar: cookieJar,
	})
}

func NewWithClient(hc *http.Client) *HttpClient {
	return createClient(hc)
}

func NewWithLocalAddr(localAddr net.Addr) *HttpClient {
	cookieJar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	return createClient(&http.Client{
		Jar:       cookieJar,
		Transport: createTransport(localAddr),
	})
}
