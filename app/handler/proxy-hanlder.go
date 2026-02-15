package handler

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/Thashreef45/proxy-server/internal/model"
)

func NewProxyHandler(cfg model.Config) (*httputil.ReverseProxy, error) {

	directorArgs, err := parseUrl(cfg.Backend)
	if err != nil {
		return nil, err
	}

	return &httputil.ReverseProxy{
		Director:      directorHandler(directorArgs),
		ErrorHandler:  errorHandler,
		Transport:     initTransport(cfg),
		FlushInterval: time.Duration(cfg.FlushIntervalMs) * time.Millisecond,
	}, nil
}

type Url struct {
	Schema string
	Host   string
}

func parseUrl(address string) (Url, error) {
	u, err := url.Parse(address)
	if err != nil {
		return Url{}, err
	}

	return Url{
		Schema: u.Scheme,
		Host:   u.Host,
	}, nil
}

func directorHandler(u Url) func(*http.Request) {
	return func(req *http.Request) {
		req.URL.Scheme = u.Schema
		req.URL.Host = u.Host
		req.Host = u.Host
	}
}

func errorHandler(w http.ResponseWriter, req *http.Request, err error) {
	http.Error(w, `Proxy Error : `+err.Error(), http.StatusBadGateway)
}

func initTransport(cfg model.Config) *http.Transport {
	return &http.Transport{
		MaxIdleConnsPerHost: cfg.MaxIdleConnsPerHost,
		MaxIdleConns:        cfg.MaxIdleConns,
		IdleConnTimeout:     time.Duration(cfg.IdleConnTimeoutSec) * time.Second,
	}
}
