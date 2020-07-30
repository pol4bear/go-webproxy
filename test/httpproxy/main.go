package main

import (
	"log"
	"net/http"

	"github.com/pol4bear/cve-2018-11714/pkg/webproxy"
)

var a int

func test(request *http.Request) {
	log.Printf("New Request: %d\n", a)
	a++
}

func main() {
	a = 1
	httpProxy := webproxy.NewHTTPProxy()
	//httpProxy.SetRequestHook(test)
	httpProxy.Start("", 3000)
}
