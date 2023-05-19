package server

import (
	"fmt"
	"net/http"
	"net/url"
)

type Server struct {
	port             uint16
	s_interface *ServerInterface
}

func NewServer(port uint16, proxy_addr *url.URL) *Server {
	addr, _ := url.Parse(fmt.Sprintf("http://localhost:%d", port))
	return &Server{
		port:             port,
		s_interface: NewServerInterface(addr),
	}
}

func (server *Server) HTTPHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Server %s: Received request from %s\n", server.s_interface.addr, r.Host)
	fmt.Fprintf(w, "ack")
}

func (server *Server) Start() {
	http_server := http.Server{
		Addr:    fmt.Sprintf(":%d", server.port),
		Handler: http.HandlerFunc(server.HTTPHandler),
	}

	fmt.Printf("Server: Running on %s...\n", server.s_interface.addr)
	err := http_server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}

func (server *Server) GetInterface() *ServerInterface {
	return server.s_interface
}
