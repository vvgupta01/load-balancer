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
	addr 	  *url.URL
	Interface *ServerInterface

	mux  sync.RWMutex
	stop chan struct{}
}

func NewServer(addr *url.URL, proxy_addr *url.URL, weight int32, capacity int32) *Server {
	return &Server{
		addr: 	   addr,
		Interface: NewServerInterface(addr, weight, capacity),
	}
}

func (server *Server) HTTPHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Server %s: Received request from %s\n", server.Interface.Addr, r.Host)
	fmt.Fprintf(w, "ACK")
}

func (server *Server) ServerRoutine() {
	http_server := http.Server{
		Addr:    fmt.Sprintf(":%s", server.addr.Port()),
		Handler: http.HandlerFunc(server.HTTPHandler),
	}

	go func() {
		if err := http_server.ListenAndServe(); err != http.ErrServerClosed {
			log.Println(err)
		}
	}()
	log.Printf("Started server on %s (weight = %d, capacity = %d)\n", server.addr,
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
	log.Printf("Stopped server on %s...\n", server.addr)
}
