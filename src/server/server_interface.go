package server

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type ServerInterface struct {
	Addr   *url.URL
	proxy  *httputil.ReverseProxy
	Health *HealthService
}

func NewServerInterface(addr *url.URL) *ServerInterface {
	health := NewHealthService(addr, time.Second*2, time.Second*2, -1)
	health.Start()

	return &ServerInterface{
		Addr:   addr,
		proxy:  httputil.NewSingleHostReverseProxy(addr),
		Health: health,
	}
}

func (si *ServerInterface) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	si.proxy.ServeHTTP(w, r)
}
