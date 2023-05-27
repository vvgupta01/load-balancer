package server

import (
	"load-balancer/src/utils"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type ServerInterface struct {
	Addr   *url.URL
	proxy  *httputil.ReverseProxy
	Health *HealthService
	Weight int32
}

func NewServerInterface(addr *url.URL, weight int32, capacity int32) *ServerInterface {
	health := NewHealthService(addr, utils.GetTimeEnv("HEALTH_INTERVAL"),
		utils.GetTimeEnv("HEALTH_TIMEOUT"), capacity)
	return &ServerInterface{
		Addr:   addr,
		proxy:  httputil.NewSingleHostReverseProxy(addr),
		Health: health,
		Weight: weight,
	}
}

func (Interface *ServerInterface) StartHealthCheck() {
	Interface.Health.Start()
}

func (Interface *ServerInterface) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	Interface.Health.AddLoad(1)
	defer Interface.Health.AddLoad(-1)

	Interface.proxy.ServeHTTP(w, r)
}
