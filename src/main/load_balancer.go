package main

import (
	"fmt"
	iterator "loadbalancer/src/iterator"
	"log"
	"net/http"
)

type LoadBalancer struct {
	iter         iterator.Iterator
	port         uint16
	done         chan int
	transactions int64
}

func NewLoadBalancer(iter iterator.Iterator, port uint16) *LoadBalancer {
	return &LoadBalancer{
		iter: iter,
		port: port,
		done: make(chan int),
	}
}

func (balancer *LoadBalancer) HTTPHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Load balancer: Received request from %s\n", r.RemoteAddr)

	_, srv := balancer.iter.NextAvailable()
	if srv != nil {
		log.Printf("Load balancer: Forwarding request to %s...\n", srv.Addr)
		srv.HandleRequest(w, r, balancer.done)
	} else {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
	}
}

func (balancer *LoadBalancer) requestCallback() {
	for {
		i := <-balancer.done
		balancer.iter.DoneCallback(i)
		balancer.transactions++
	}
}

func (balancer *LoadBalancer) Start() {
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", balancer.port),
		Handler: http.HandlerFunc(balancer.HTTPHandler),
	}

	go balancer.requestCallback()

	log.Printf("Load balancer: Forwarding on http://localhost:%d...\n", balancer.port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
