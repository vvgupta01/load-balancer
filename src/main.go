package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"load-balancer/src/balancer"
	"load-balancer/src/iterator"
	"load-balancer/src/server"
	utils "load-balancer/src/utils"
	"net/url"
	"time"

	"github.com/joho/godotenv"
	// "math/rand"
)

type ServerConfig struct {
	Port     uint16
	Weight   int32
	Capacity int32
}

func main() {
	godotenv.Load("config.env")
	proxy_port := utils.GetIntEnv("PROXY_PORT")

	proxy_addr, _ := url.Parse(fmt.Sprintf("http://localhost:%d", proxy_port))

	file, err := ioutil.ReadFile("config_server.json")
	if err != nil {
		fmt.Println(err)
	}

	var configs []ServerConfig
	if err := json.Unmarshal(file, &configs); err != nil {
		fmt.Println(err)
	}

	// Initialize servers
	servers := make([]*server.Server, len(configs))
	interfaces := make([]*server.ServerInterface, len(configs))

	for i, conf := range configs {
		srv := server.NewServer(conf.Port, proxy_addr, conf.Weight, conf.Capacity)

		servers[i] = srv
		interfaces[i] = srv.Interface

		srv.Start()
		srv.Interface.StartHealthCheck()
	}
	pool := server.NewServerPool(interfaces)

	// Initialize iterator/load-balancer
	// iter := iterator.NewRandom(rand.Seed(time.Now().UnixNano(), pool))
	// iter := iterator.NewRoundRobin(pool)
	iter := iterator.NewWeightedRoundRobin(pool)
	load_balancer := balancer.NewLoadBalancer(iter, uint16(proxy_port))
	go load_balancer.Start()

	// Initialize client
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
