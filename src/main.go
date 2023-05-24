package main

import (
	"fmt"
	balancer "load-balancer/src/balancer"
	iterator "load-balancer/src/iterator"
	server "load-balancer/src/server"
	"net/url"
	"time"
	// "math/rand"
)

func main() {
	port := 3000
	proxy_addr, _ := url.Parse(fmt.Sprintf("http://localhost:%d", port))

	// Initialize multiple servers
	num_servers := 3
	var servers []*server.Server
	var interfaces []*server.ServerInterface
	for i := 1; i <= num_servers; i++ {
		s := server.NewServer(uint16(port+i), proxy_addr, int32(i), 100)

		servers = append(servers, s)
		interfaces = append(interfaces, s.Interface)

		s.Start()
	}
	pool := server.NewServerPool(interfaces)

	// Initialize single server
	// s := server.NewServer(uint16(3001), proxy_addr, 1, 100)
	// s.Start()
	// interfaces := []*server.ServerInterface{s.Interface}
	// pool := server.NewServerPool(interfaces)

	// iter := iterator.NewRandom(rand.Seed(time.Now().UnixNano(), pool))
	// iter := iterator.NewRoundRobin(pool)
	iter := iterator.NewWeightedRoundRobin(pool)
	load_balancer := balancer.NewLoadBalancer(iter, uint16(port))
	go load_balancer.Start()

	// Test client
	time.Sleep(time.Second)
	client := balancer.NewClient(1, proxy_addr)
	client.Start()

	// // Test running server/health
	// time.Sleep(3 * time.Second)

	// // Test stopped server/running health
	// s.Stop()
	// time.Sleep(3 * time.Second)

	// // Test stopped server/stopped health
	// interfaces[0].Health.Stop()
	// time.Sleep(3 * time.Second)

	// // Test running server/stopped health
	// s.Start()
	// time.Sleep(3 * time.Second)

	// interfaces[0].Health.Start()
	// time.Sleep(3 * time.Second)
}
