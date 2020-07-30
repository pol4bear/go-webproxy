// Package webproxy provides http and https proxy with request and response hooks
package webproxy

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
)

// RequestHookFunc typed functions provided to proxy will modify request before sent to the server.
type RequestHookFunc func(request *http.Request)

// ResponseHookFunc typed functions provided to proxy will modify response before sent to the client.
type ResponseHookFunc func(response *http.Response)

// HTTPProxy will store some items that are used in HTTP proxy and state.
type HTTPProxy struct {
	server       *http.Server
	IsStarted    bool
	requestHook  RequestHookFunc
	responseHook ResponseHookFunc
}

// ProxyManager provide methods to manage HTTPProxy and HTTPSProxy.
type ProxyManager interface {
	Start(host string, port uint16)
	Stop()
	GetAddr()
	SetRequestHook(hook RequestHookFunc)
	SetResponseHook(hook ResponseHookFunc)
}

// NewHTTPProxy creates default HTTPProxy struct.
func NewHTTPProxy() *HTTPProxy {
	return &HTTPProxy{
		server:       new(http.Server),
		IsStarted:    false,
		requestHook:  nil,
		responseHook: nil,
	}
}

// Start serves HTTPProxy on host:port.
func (httpProxy *HTTPProxy) Start(host string, port uint16) {
	if httpProxy.IsStarted {
		log.Println("[HTTP] The proxy is already running.")
		return
	} else if len(host) != 0 && net.ParseIP(host) == nil {
		log.Fatal("[HTTP] Inserted host is not valid.")
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	httpProxy.server.Handler = getHandler(httpProxy.requestHook, httpProxy.responseHook)
	httpProxy.server.Addr = addr

	log.Printf("[HTTP] The proxy will be served on %s:%d\n", host, port)
	if err := httpProxy.server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("[HTTP] The proxy ListenAndServer: %v\n", err)
	}
}

// Stop will stop HTTPProxy immediately.
func (httpProxy *HTTPProxy) Stop() {
	if !httpProxy.IsStarted {
		log.Println("[HTTP] The proxy is not running.")
		return
	}
	httpProxy.server.Shutdown(context.Background())
}

// GetAddr will return HTTPProxy's host:port.
func (httpProxy *HTTPProxy) GetAddr() string {
	return httpProxy.server.Addr
}

// SetRequestHook setup request hook into the proxy
func (httpProxy *HTTPProxy) SetRequestHook(hook RequestHookFunc) {
	if httpProxy.IsStarted {
		log.Println("[HTTP] You can't set request hook while the proxy is working.")
		return
	}
	httpProxy.requestHook = hook
}

// SetResponseHook setup response hook into the proxy.
func (httpProxy *HTTPProxy) SetResponseHook(hook ResponseHookFunc) {
	if httpProxy.IsStarted {
		log.Println("[HTTP] You can't set response hook while the proxy is working.")
		return
	}
	httpProxy.responseHook = hook
}
