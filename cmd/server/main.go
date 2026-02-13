package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

func main() {
	proxyHandler, err := newProxy()
	if err != nil {
		fmt.Println("error creating proxy :", err)
		return
	}
	fmt.Println("Starting proxy on :8080")
	http.ListenAndServe(":8080", proxyHandler)
}

func newProxy() (*httputil.ReverseProxy, error) {

	address := "http://localhost:8000"

	parsedUrl, urlParseError := url.Parse(address)
	if urlParseError != nil {
		return nil, urlParseError
	}

	return &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = parsedUrl.Scheme
			req.URL.Host = parsedUrl.Host
			req.Host = parsedUrl.Host
		},
		Transport: &http.Transport{
			MaxIdleConns:        1000,
			MaxIdleConnsPerHost: 1000,
			IdleConnTimeout:     time.Second * 60,
		},
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, `Proxy Error : `+err.Error(), http.StatusBadGateway)
		},
		FlushInterval: time.Millisecond * 50,
	}, nil

}
