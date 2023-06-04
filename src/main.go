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
	Host     string
	Port     uint16
	Weight   int32
	Capacity int32
}

type ClientConfig struct {
	Rate float32
}

func main() {
	// Parse command-line arguments
	env_file := flag.String("env", "default", "Name of env file in config/load_balancer/ specifying load balancer config")
	server_file := flag.String("server", "default", "Name of JSON file in config/server/ specifying server config")
	client_file := flag.String("client", "", "Name of JSON file in config/client/ specifying test client config, empty for none")

	test_server := flag.Bool("t", false, "Indicates if local test servers should be instantiated or if preexisting servers should be used")
	log_file := flag.String("log", "", "Name of log file to log output, empty to disable logging")
	verbose := flag.Bool("v", false, "Verbose")
	flag.Parse()

	// Load env file
	if *env_file != "" {
		env_path := fmt.Sprintf("config/load_balancer/%s.env", *env_file)
		if err := godotenv.Load(env_path); err != nil {
			log.Fatal(err)
		}
	}

	// Load load balancer config
	proxy_port := utils.GetIntEnv("PROXY_PORT")
	proxy_addr, err := url.Parse(fmt.Sprintf("http://localhost:%d", proxy_port))
	if err != nil {
		log.Fatal(err)
	}

	// Load server config
	server_path := fmt.Sprintf("config/server/%s.json", *server_file)
	server_json, err := ioutil.ReadFile(server_path)
	if err != nil {
		log.Fatal(err)
	}

	var server_configs []ServerConfig
	if err := json.Unmarshal(server_json, &server_configs); err != nil {
		log.Fatal(err)
	}

	// Load test client config
	var client_configs []ClientConfig
	if *client_file != "" {
		client_path := fmt.Sprintf("config/client/%s.json", *client_file)
		client_json, err := ioutil.ReadFile(client_path)
		if err != nil {
			log.Fatal(err)
		}
		if err := json.Unmarshal(client_json, &client_configs); err != nil {
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
		log_file, err := os.OpenFile(log_path, os.O_CREATE|os.O_RDWR, 0755)
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

	// Initialize servers/interfaces
	servers := make([]*server.Server, len(server_configs))
	pool := make(server.ServerPool, len(server_configs))

	for i, conf := range server_configs {
		addr, err := url.Parse(fmt.Sprintf("http://%s:%d", conf.Host, conf.Port))
		if err != nil {
			log.Fatal(err)
		}

		if *test_server {
			servers[i] = server.NewServer(addr, proxy_addr, i, conf.Weight, conf.Capacity)
			pool[i] = servers[i].Interface
			servers[i].Start(false)
		} else {
			pool[i] = server.NewServerInterface(addr, i, conf.Weight, conf.Capacity)
		}

		// Optional: Start concurrent health check
		pool[i].StartHealthCheck(nil)
	}

	// Initialize test clients to periodically ping load balancer
	time.Sleep(time.Second)
	for _, conf := range client_configs {
		client := balancer.NewClient(conf.Rate, proxy_addr)
		go client.Start()
	}

	// Initialize iterator/load-balancer
	// iter := iterator.NewRandom(iterator.DefaultSeed, pool)
	// iter := iterator.NewRoundRobin(pool)
	// iter := iterator.NewWeightedRoundRobin(pool)
	iter := iterator.NewLeastConnections(pool)
	load_balancer := balancer.NewLoadBalancer(iter, uint16(proxy_port))
	load_balancer.Start()
}
