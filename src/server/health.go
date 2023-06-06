package server

import (
	"log"
	"net"
	"net/url"
	"os"
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

	alive_check func(*HealthService) bool
	stop        chan struct{}
	alive_mux   sync.RWMutex
	stop_mux    sync.RWMutex
}

func NewHealthService(addr *url.URL, interval time.Duration, timeout time.Duration, capacity int32) *HealthService {
	return &HealthService{
		addr:        addr.JoinPath(os.Getenv("HEALTH_ENDPOINT")),
		alive_check: DefaultAliveCheck,
		alive:       true,
		interval:    interval,
		timeout:     timeout,
		load:        0,
		capacity:    capacity,
	}
}

func DefaultAliveCheck(service *HealthService) bool {
	conn, err := net.DialTimeout("tcp", service.addr.Host, service.timeout)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}

func (service *HealthService) HealthCheck() {
	alive := service.alive_check(service)
	if !alive {
		service.SetLoad(0)
	}
	service.SetAlive(alive)

	log.Printf("%s status - alive = %t, load = %d/%d\n", service.addr,
		service.IsAlive(), service.GetLoad(), service.GetCapacity())
}

func (service *HealthService) Run(notify chan int) {
	t := time.NewTicker(service.interval)
	ticks := 0
	for {
		select {
		case <-t.C:
			service.HealthCheck()
			if notify != nil {
				ticks++
				notify <- ticks
			}
		case <-service.stop:
			t.Stop()
			return
		}
	}
}

func (service *HealthService) Start(notify chan int) {
	service.stop_mux.Lock()
	defer service.stop_mux.Unlock()

	if service.stop == nil {
		service.stop = make(chan struct{})
		go service.Run(notify)

		log.Printf("Started health service for %s (interval = %s, timeout = %s)\n",
			service.addr, service.interval, service.timeout)
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

func (service *HealthService) SetAliveCheck(alive_check func(*HealthService) bool) {
	service.alive_check = alive_check
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
