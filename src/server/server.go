package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"sync"
)

type Server struct {
	addr      *url.URL
	Interface *ServerInterface

	mux  sync.RWMutex
	stop chan struct{}
}

func NewServer(addr *url.URL, proxy_addr *url.URL, weight int32, capacity int32) *Server {
	return &Server{
		addr:      addr,
		Interface: NewServerInterface(addr, weight, capacity),
	}
}

func (server *Server) HTTPHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Server %s: Received request from %s\n", server.Interface.Addr, r.Host)
	fmt.Fprintf(w, "ACK")
}

func (server *Server) ServerRoutine(ack chan struct{}) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", server.addr.Port()))
	if err != nil {
		log.Println(err)
	}

	if ack != nil {
		ack <- struct{}{}
	}

	go func() {
		if err := http.Serve(l, http.HandlerFunc(server.HTTPHandler)); err != http.ErrServerClosed {
			log.Println(err)
		}
	}()
	log.Printf("Started server on %s (weight = %d, capacity = %d)\n", server.addr,
		server.Interface.Weight, server.Interface.Health.capacity)

	<-server.stop
	l.Close()
}

func (server *Server) Start(block bool) {
	server.mux.Lock()
	defer server.mux.Unlock()

	if server.stop == nil {
		server.stop = make(chan struct{})

		if block {
			ack := make(chan struct{})
			go server.ServerRoutine(ack)
			<-ack
		} else {
			go server.ServerRoutine(nil)
		}
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

func (server *Server) IsRunning() bool {
	server.mux.Lock()
	defer server.mux.Unlock()

	return server.stop != nil
}
