package main

import (
	"fmt"
	"loadbalancer/src/server"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ServerConfig struct {
	Host      string
	Port      uint16
	Weight    int32
	Load      int32
	Capacity  int32
	Available bool
}

type ServerManager struct {
	port     uint16
	balancer *LoadBalancer
	servers  []ServerConfig
	pool     server.ServerPool
}

func NewServerManager(balancer *LoadBalancer, port uint16) *ServerManager {
	return &ServerManager{
		port:     port,
		balancer: balancer,
	}
}

func (manager *ServerManager) SetServers(servers []ServerConfig) {
	manager.servers = servers
}

func (manager *ServerManager) SetPool(pool server.ServerPool) {
	manager.pool = pool
}

func (manager *ServerManager) Start() {
	router := gin.Default()
	router.GET("/servers", manager.getServers)
	router.GET("/status", manager.getStatus)

	log.Printf("Manager: Running on http://localhost:%d...\n", manager.port)
	router.Run(fmt.Sprintf("localhost:%d", manager.port))
}

func (manager *ServerManager) getServers(c *gin.Context) {
	for i, srv := range manager.pool {
		manager.servers[i].Load = srv.Health.GetLoad()
		manager.servers[i].Available = srv.Health.IsAvailable()
	}
	c.IndentedJSON(http.StatusOK, manager.servers)
}

func (manager *ServerManager) getStatus(c *gin.Context) {
	status := manager.pool.GetStatus()
	status.Transactions = manager.balancer.transactions
	c.IndentedJSON(http.StatusOK, status)
}
