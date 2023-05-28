package server

import (
	"log"
	"net"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

type HealthService struct {
	addr     *url.URL
	alive    bool
	interval time.Duration
	timeout  time.Duration
	load     int32
	capacity int32

	stop      chan struct{}
	alive_mux sync.RWMutex
	stop_mux  sync.RWMutex
}

func NewHealthService(addr *url.URL, interval time.Duration, timeout time.Duration, capacity int32) *HealthService {
	return &HealthService{
		addr:     addr,
		alive:    true,
		interval: interval,
		timeout:  timeout,
		load:     0,
		capacity: capacity,
	}
}

func (service *HealthService) HealthCheck() {
	conn, err := net.DialTimeout("tcp", service.addr.Host, service.timeout)

	if err != nil {
		service.SetAlive(false)
	} else {
		_ = conn.Close()
		service.SetAlive(true)
	}
	log.Printf("%s status - alive = %t, load = %d/%d\n", service.addr,
		service.IsAlive(), service.GetLoad(), service.GetCapacity())
}

func (service *HealthService) HealthRoutine() {
	t := time.NewTicker(service.interval)
	for {
		select {
		case <-t.C:
			service.HealthCheck()
		case <-service.stop:
			t.Stop()
			return
		}
	}
}

func (service *HealthService) Start() {
	service.stop_mux.Lock()
	defer service.stop_mux.Unlock()

	if service.stop == nil {
		service.stop = make(chan struct{})
		log.Printf("Starting health service for %s (interval = %s, timeout = %s)...\n",
			service.addr, service.interval, service.timeout)
		service.HealthCheck()
		go service.HealthRoutine()
	}
}

func (service *HealthService) Stop() {
	service.stop_mux.Lock()
	defer service.stop_mux.Unlock()

	if service.stop != nil {
		service.stop <- struct{}{}
		close(service.stop)
		service.stop = nil
		log.Printf("Stopped health service for %s\n", service.addr)
	}
}

func (service *HealthService) SetAlive(alive bool) {
	service.alive_mux.Lock()
	defer service.alive_mux.Unlock()
	service.alive = alive
}

func (service *HealthService) AddLoad(load int32) {
	atomic.AddInt32(&service.load, load)
}

func (service *HealthService) SetLoad(load int32) {
	atomic.StoreInt32(&service.load, load)
}

func (service *HealthService) IsAlive() bool {
	service.alive_mux.RLock()
	defer service.alive_mux.RUnlock()
	return service.alive
}

func (service *HealthService) IsAvailable() bool {
	return service.IsAlive() && service.GetLoad() < service.GetCapacity()
}

func (service *HealthService) GetLoad() int32 {
	return atomic.LoadInt32(&service.load)
}

func (service *HealthService) GetCapacity() int32 {
	return atomic.LoadInt32(&service.capacity)
}
