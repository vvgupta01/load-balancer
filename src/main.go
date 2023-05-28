package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"loadbalancer/src/balancer"
	"loadbalancer/src/iterator"
	"loadbalancer/src/server"
	utils "loadbalancer/src/utils"
	"log"
	"net/url"
	"os"
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
	// Parse command-line arguments
	env_file := flag.String("env", "default", "Name of .env file specifying load balancer config")
	log_file := flag.String("log", "", "Name of .log file to log output")
	verbose := flag.Bool("v", false, "Verbose")
	flag.Parse()

	// Load .env file
	if *env_file != "" {
		env_path := fmt.Sprintf("config/%s.env", *env_file)
		if err := godotenv.Load(env_path); err != nil {
			log.Fatal(err)
		}
	}

	// Set log output
	if !*verbose && *log_file == "" {
		utils.DisableLogOutput()
	} else if *log_file != "" {
		if err := os.MkdirAll("logs", os.ModePerm); err != nil {
			log.Fatal(err)
		}

		log_path := fmt.Sprintf("logs/%s.log", *log_file)
		log_file, err := os.OpenFile(log_path, os.O_CREATE | os.O_RDWR, 0755)
		if err != nil {
			log.Fatal(err)
		}
		defer log_file.Close()

		if *verbose {
			mw := io.MultiWriter(os.Stderr, log_file)
			log.SetOutput(mw)
		} else {
			log.SetOutput(log_file)
		}
	}	

	// Load load balancer config
	proxy_port := utils.GetIntEnv("PROXY_PORT")
	proxy_addr, err := url.Parse(fmt.Sprintf("http://localhost:%d", proxy_port))
	if err != nil {
		log.Fatal(err)
	}

	// Load server config
	file, err := ioutil.ReadFile("config/server.json")
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

		// Comment this block if not creating test servers
		servers[i] = server.NewServer(addr, proxy_addr, conf.Weight, conf.Capacity)
		interfaces[i] = servers[i].Interface
		servers[i].Start()
		
		// Uncomment this line if using servers running elsewhere
		// interfaces[i] = server.NewServerInterface(addr, conf.Weight, conf.Capacity)

		// Optional health check
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

	// Initialize test client to ping load balancer
	time.Sleep(time.Second)
	client := balancer.NewClient(1, proxy_addr)
	client.Start()
}
