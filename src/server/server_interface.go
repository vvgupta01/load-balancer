package server

import (
	"loadbalancer/src/utils"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type ServerInterface struct {
	Addr   *url.URL
	proxy  *httputil.ReverseProxy
	Index  int
	Health *HealthService
	Weight int32

	handler func(*ServerInterface, http.ResponseWriter, *http.Request)
}

func NewServerInterface(addr *url.URL, index int, weight int32, capacity int32) *ServerInterface {
	health := NewHealthService(addr, utils.GetTimeEnv("HEALTH_INTERVAL"), utils.GetTimeEnv("HEALTH_TIMEOUT"), capacity)
	return &ServerInterface{
		Addr:    addr,
		proxy:   httputil.NewSingleHostReverseProxy(addr),
		Index:   index,
		Health:  health,
		Weight:  weight,
		handler: ServeHTTP,
	}
}

func (Interface *ServerInterface) HandleRequest(w http.ResponseWriter, r *http.Request, done chan int) {
	Interface.handler(Interface, w, r)
	if done != nil {
		done <- Interface.Index
	}
}

func (Interface *ServerInterface) SetHandler(handler func(*ServerInterface, http.ResponseWriter, *http.Request)) {
	Interface.handler = handler
}

func ServeHTTP(Interface *ServerInterface, w http.ResponseWriter, r *http.Request) {
	Interface.proxy.ServeHTTP(w, r)
}

func (Interface *ServerInterface) StartHealthCheck(notify chan int) {
	Interface.Health.Start(notify)
}

func (Interface *ServerInterface) StopHealthCheck() {
	Interface.Health.Stop()
}
