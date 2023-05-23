package server

import (
	"fmt"
	"net"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

type HealthService struct {
	addr     *url.URL
	alive    bool
	period   time.Duration
	timeout  time.Duration
	load     int32
	capacity int32

	stop 	   chan struct{}
	alive_mux  sync.RWMutex
	stop_mux   sync.RWMutex
}

func NewHealthService(addr *url.URL, period time.Duration, timeout time.Duration, capacity int32) *HealthService {
	return &HealthService{
		addr:     addr,
		alive:    true,
		period:   period,
		timeout:  timeout,
		load:     0,
		capacity: capacity,
	}
}

func (service *HealthService) HealthCheck() {
	conn, err := net.DialTimeout("tcp", service.addr.Host, service.timeout)

	service.alive_mux.Lock()
	defer service.alive_mux.Unlock()

	if err != nil {
		service.alive = false
	} else {
		_ = conn.Close()
		service.alive = true
	}
	fmt.Printf("%s status - alive = %t, load = %d\n", service.addr, service.alive, service.load)
}

func (service *HealthService) HealthRoutine() {
	t := time.NewTicker(service.period)
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
		go service.HealthRoutine()
		fmt.Printf("Started health service for %s\n", service.addr)
	}
}

func (service *HealthService) Stop() {
	service.stop_mux.Lock()
	defer service.stop_mux.Unlock()

	if service.stop != nil {
		service.stop <- struct{}{}
		close(service.stop)
		service.stop = nil
		fmt.Printf("Stopped health service for %s\n", service.addr)
	}
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
