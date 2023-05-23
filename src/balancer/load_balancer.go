package balancer

import (
	"fmt"
	iterator "load-balancer/src/iterator"
	"net/http"
)

type LoadBalancer struct {
	iter iterator.Iterator
	port uint16
}

func NewLoadBalancer(iter iterator.Iterator, port uint16) *LoadBalancer {
	return &LoadBalancer{
		iter: iter,
		port: port,
	}
}

func (balancer *LoadBalancer) HTTPHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Load balancer: Received request from %s\n", r.RemoteAddr)

	server := balancer.iter.Next()
	if server != nil {
		fmt.Printf("Load balancer: Forwarding request to %s...\n", server.Addr)
		server.ServeHTTP(w, r)
	} else {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
	}
}

func (balancer *LoadBalancer) Start() {
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", balancer.port),
		Handler: http.HandlerFunc(balancer.HTTPHandler),
	}

	fmt.Printf("Load balancer: Running on http://localhost:%d...\n", balancer.port)
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
