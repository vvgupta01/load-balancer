package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"load-balancer/src/balancer"
	"load-balancer/src/iterator"
	"load-balancer/src/server"
	utils "load-balancer/src/utils"
	"log"
	"net/url"
	"time"

	"github.com/joho/godotenv"
	// "math/rand"
)

type ServerConfig struct {
	Addr     string
	Weight   int32
	Capacity int32
}

func main() {
	godotenv.Load("config.env")
	proxy_port := utils.GetIntEnv("PROXY_PORT")
	proxy_addr, err := url.Parse(fmt.Sprintf("http://localhost:%d", proxy_port))
	if err != nil {
		log.Fatal(err)
	}

	file, err := ioutil.ReadFile("config_server.json")
	if err != nil {
		log.Fatal(err)
	}

	var configs []ServerConfig
	if err := json.Unmarshal(file, &configs); err != nil {
		log.Fatal(err)
	}

	// Initialize servers/interfaces
	servers := make([]*server.Server, len(configs))
	interfaces := make([]*server.ServerInterface, len(configs))

	for i, conf := range configs {
		addr, err := url.Parse(conf.Addr)
		if err != nil {
			log.Fatal(err)
		}

		servers[i] = server.NewServer(addr, proxy_addr, conf.Weight, conf.Capacity)
		interfaces[i] = servers[i].Interface
		servers[i].Start()

		// interfaces[i] = server.NewServerInterface(addr, conf.Weight, conf.Capacity)

		interfaces[i].StartHealthCheck()
	}
	pool := server.NewServerPool(interfaces)

	// Initialize iterator/load-balancer
	// iter := iterator.NewRandom(iterator.DefaultSeed, pool)
	// iter := iterator.NewRoundRobin(pool)
	// iter := iterator.NewWeightedRoundRobin(pool)
	iter := iterator.NewLeastConnections(pool)
	load_balancer := balancer.NewLoadBalancer(iter, uint16(proxy_port))
	go load_balancer.Start()

	// Initialize client
	time.Sleep(time.Second)
	client := balancer.NewClient(1, proxy_addr)
	client.Start()
}
