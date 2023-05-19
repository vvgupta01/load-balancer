package server

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type ServerInterface struct {
	addr  *url.URL
	proxy *httputil.ReverseProxy
	alive bool
	load  uint32
}

func NewServerInterface(addr *url.URL) *ServerInterface {
	return &ServerInterface{
		addr:  addr,
		proxy: httputil.NewSingleHostReverseProxy(addr),
		alive: true,
	}
}

func (si *ServerInterface) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	si.proxy.ServeHTTP(w, r)
}

func (si *ServerInterface) GetAddr() *url.URL {
	return si.addr
}
