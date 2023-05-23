package server

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Server struct {
	port        uint16
	Interface *ServerInterface

	mux  sync.RWMutex
	stop chan struct{}
}

func NewServer(port uint16, proxy_addr *url.URL) *Server {
	addr, _ := url.Parse(fmt.Sprintf("http://localhost:%d", port))
	return &Server{
		port:        port,
		Interface: NewServerInterface(addr),
	}
}

func (server *Server) HTTPHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Server %s: Received request from %s\n", server.Interface.Addr, r.Host)
	fmt.Fprintf(w, "ACK")
}

func (server *Server) ServerRoutine() {
	http_server := http.Server{
		Addr:    fmt.Sprintf(":%d", server.port),
		Handler: http.HandlerFunc(server.HTTPHandler),
	}

	go func() {
		if err := http_server.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Println(err)
		}
	}()
	fmt.Printf("Started server on %s...\n", server.Interface.Addr)
	
	<-server.stop
	ctx, cancel := context.WithTimeout(context.Background(), 2 * time.Second)
	defer cancel()

	if err := http_server.Shutdown(ctx); err != nil {
		fmt.Println(err)
	}
}

func (server *Server) Start() {
	server.mux.Lock()
	defer server.mux.Unlock()

	if server.stop == nil {
		server.stop = make(chan struct{})
		go server.ServerRoutine()
	}
}

func (server *Server) Stop() {
	server.mux.Lock()
	defer server.mux.Unlock()
	
	if server.stop != nil {
		server.stop <- struct{}{}
		close(server.stop)
		server.stop = nil
	}
	fmt.Printf("Stopped server on %s...\n", server.Interface.Addr)
}