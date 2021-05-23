
package http_client

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
)

type (

	RedirectPolicy interface {
		Apply(req *http.Request, via []*http.Request) error
	}

	RedirectPolicyFunc func(*http.Request, []*http.Request) error
)

func (f RedirectPolicyFunc) Apply(req *http.Request, via []*http.Request) error {
	return f(req, via)
}

func NoRedirectPolicy() RedirectPolicy {
	return RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error {
		return errors.New("auto redirect is disabled")
	})
}

func FlexibleRedirectPolicy(noOfRedirect int) RedirectPolicy {
	return RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error {
		if len(via) >= noOfRedirect {
			return fmt.Errorf("stopped after %d redirects", noOfRedirect)
		}
		checkHostAndAddHeaders(req, via[0])
		return nil
	})
}

func DomainCheckRedirectPolicy(hostnames ...string) RedirectPolicy {
	hosts := make(map[string]bool)
	for _, h := range hostnames {
		hosts[strings.ToLower(h)] = true
	}

	fn := RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error {
		if ok := hosts[getHostname(req.URL.Host)]; !ok {
			return errors.New("redirect is not allowed as per DomainCheckRedirectPolicy")
		}

		return nil
	})

	return fn
}


func getHostname(host string) (hostname string) {
	if strings.Index(host, ":") > 0 {
		host, _, _ = net.SplitHostPort(host)
	}
	hostname = strings.ToLower(host)
	return
}

func checkHostAndAddHeaders(cur *http.Request, pre *http.Request) {
	curHostname := getHostname(cur.URL.Host)
	preHostname := getHostname(pre.URL.Host)
	if strings.EqualFold(curHostname, preHostname) {
		for key, val := range pre.Header {
			cur.Header[key] = val
		}
	}
}
