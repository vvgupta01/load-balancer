package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Server struct {
	port      uint16
	Interface *ServerInterface

	mux  sync.RWMutex
	stop chan struct{}
}

func NewServer(port uint16, proxy_addr *url.URL, weight int32, capacity int32) *Server {
	addr, _ := url.Parse(fmt.Sprintf("http://localhost:%d", port))
	return &Server{
		port:      port,
		Interface: NewServerInterface(addr, weight, capacity),
	}
}

func (server *Server) HTTPHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Server %s: Received request from %s\n", server.Interface.Addr, r.Host)
	fmt.Fprintf(w, "ACK")
}

func (server *Server) ServerRoutine() {
	http_server := http.Server{
		Addr:    fmt.Sprintf(":%d", server.port),
		Handler: http.HandlerFunc(server.HTTPHandler),
	}

	go func() {
		if err := http_server.ListenAndServe(); err != http.ErrServerClosed {
			log.Println(err)
		}
	}()
	log.Printf("Started server on %s (weight = %d, capacity = %d)\n", server.Interface.Addr,
		server.Interface.Weight, server.Interface.Health.capacity)

	<-server.stop
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := http_server.Shutdown(ctx); err != nil {
		log.Println(err)
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
	log.Printf("Stopped server on %s...\n", server.Interface.Addr)
}
