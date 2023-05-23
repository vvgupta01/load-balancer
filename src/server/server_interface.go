package server

import (
	"math"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
	"time"
)

type ServerInterface struct {
	Addr   *url.URL
	proxy  *httputil.ReverseProxy
	Health *HealthService
}

func NewServerInterface(addr *url.URL) *ServerInterface {
	health := NewHealthService(addr, time.Second*2, time.Second*2, math.MaxInt32)
	health.Start()

	return &ServerInterface{
		Addr:   addr,
		proxy:  httputil.NewSingleHostReverseProxy(addr),
		Health: health,
	}
}

func (Interface *ServerInterface) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	atomic.AddInt32(&Interface.Health.load, 1)
	defer atomic.AddInt32(&Interface.Health.load, -1)

	Interface.proxy.ServeHTTP(w, r)
}
