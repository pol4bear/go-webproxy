package webproxy

import (
	"bytes"
	"net/http"
	"sync"
)

type proxyHandler struct {
	requestHook  RequestHookFunc
	responseHook ResponseHookFunc
}

var initialize sync.Once
var hopByHopHeaders []string

func getHandler(requestHookIn RequestHookFunc, responseHookIn ResponseHookFunc) http.Handler {
	initialize.Do(func() {
		hopByHopHeaders = []string{
			"Connection",
			"Keep-Alive",
			"Public",
			"Proxy-Authenticate",
			"Transfer-Encoding",
			"Upgrade",
		}
	})
	return &proxyHandler{requestHook: requestHookIn, responseHook: responseHookIn}
}

// ServeHttp
func (handler *proxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if handler.requestHook != nil {
		handler.requestHook(r)
	}

	for _, header := range hopByHopHeaders {
		r.Header.Del(header)
	}

	res, _ := http.Get("http://test.pol4.dev")
	body := new(bytes.Buffer)
	body.ReadFrom(res.Body)
	w.Write(body.Bytes())
}
