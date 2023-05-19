package main

import (
	"fmt"
	balancer "load-balancer/src/balancer"
	iterator "load-balancer/src/iterator"
	server "load-balancer/src/server"
	"net/url"
	"time"
)

func main() {
	port, num_clients := 3000, 3
	proxy_addr, _ := url.Parse(fmt.Sprintf("http://localhost:%d", port))

	var servers []*server.Server
	var interfaces []*server.ServerInterface
	for i := 1; i <= num_clients; i++ {
		s := server.NewServer(uint16(port+i), proxy_addr)

		servers = append(servers, s)
		interfaces = append(interfaces, s.GetInterface())

		go s.Start()
	}
	pool := server.NewServerPool(interfaces)

	load_balancer := balancer.NewLoadBalancer(iterator.NewRoundRobinIterator(pool), uint16(port))
	go load_balancer.Start()

	time.Sleep(time.Second)
	client := balancer.NewClient(1, proxy_addr)
	client.Start()
}
