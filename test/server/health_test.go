package server_test

import (
	"fmt"
	"loadbalancer/src/server"
	"loadbalancer/test"
	"net/url"
	"testing"
	"time"
)

func availableCheck(service *server.HealthService, expected bool) error {
	actual := service.IsAvailable()
	if actual != expected {
		return fmt.Errorf("Available: Returned %t; Expected %t", actual, expected)
	}
	return nil
}

func tickCheck(actual int, expected int) error {
	if actual != expected {
		return fmt.Errorf("Tick: Returned %d; Expected %d", actual, expected)
	}
	return nil
}

func createTestServer(addr *url.URL) *server.Server {
	srv := server.NewServer(addr, nil, 0, 1, 1)
	test_check := func(*server.HealthService) bool {
		return srv.IsRunning()
	}
	srv.Interface.Health.SetAliveCheck(test_check)
	return srv
}

func TestHealthService(t *testing.T) {
	test.Setup()

	t.Run("Available server health check", func(t *testing.T) {
		addr, _ := url.Parse("http://localhost:5001")
		srv := createTestServer(addr)
		srv.Start(true)
		defer srv.Stop()

		srv.Interface.Health.HealthCheck()
		if err := availableCheck(srv.Interface.Health, true); err != nil {
			t.Error(err)
		}
	})

	t.Run("Unavailable server (not alive) health check", func(t *testing.T) {
		addr, _ := url.Parse("http://localhost:5002")
		srv := createTestServer(addr)

		srv.Interface.Health.HealthCheck()
		if err := availableCheck(srv.Interface.Health, false); err != nil {
			t.Error(err)
		}
	})

	t.Run("Unavailable server (high load) health check", func(t *testing.T) {
		addr, _ := url.Parse("http://localhost:5003")
		srv := createTestServer(addr)
		srv.Interface.Health.SetLoad(2)

		srv.Start(true)
		defer srv.Stop()

		srv.Interface.Health.HealthCheck()
		if err := availableCheck(srv.Interface.Health, false); err != nil {
			t.Error(err)
		}
	})

	t.Run("Health service check", func(t *testing.T) {
		addr, _ := url.Parse("http://localhost:5004")
		srv := createTestServer(addr)

		notify := make(chan int)
		srv.Interface.Health.Start(notify)
		defer srv.Interface.StopHealthCheck()

		<-notify
		if err := availableCheck(srv.Interface.Health, false); err != nil {
			t.Error(err)
		}

		srv.Start(true)
		defer srv.Stop()

		<-notify
		if err := availableCheck(srv.Interface.Health, true); err != nil {
			t.Error(err)
		}
	})

	t.Run("Health service restart", func(t *testing.T) {
		addr, _ := url.Parse("http://localhost:5005")
		srv := createTestServer(addr)

		notify := make(chan int)
		srv.Interface.StartHealthCheck(notify)

		<-notify
		if err := availableCheck(srv.Interface.Health, false); err != nil {
			t.Error(err)
		}

		srv.Interface.StopHealthCheck()
		srv.Start(true)
		srv.Interface.StartHealthCheck(notify)

		<-notify
		if err := availableCheck(srv.Interface.Health, true); err != nil {
			t.Error(err)
		}
	})

	t.Run("Stale health service", func(t *testing.T) {
		addr, _ := url.Parse("http://localhost:5006")
		srv := createTestServer(addr)
		srv.Start(true)

		notify := make(chan int)
		srv.Interface.StartHealthCheck(notify)

		<-notify
		if err := availableCheck(srv.Interface.Health, true); err != nil {
			t.Error(err)
		}

		srv.Interface.StopHealthCheck()
		srv.Stop()

		if err := availableCheck(srv.Interface.Health, true); err != nil {
			t.Error(err)
		}
	})

	t.Run("Health service multiple check consistency", func(t *testing.T) {
		addr, _ := url.Parse("http://localhost:5007")
		srv := createTestServer(addr)
		srv.Start(true)

		notify := make(chan int)
		srv.Interface.Health.Start(notify)
		defer srv.Interface.StopHealthCheck()

		for i := 1; i <= 10; i++ {
			tick := <-notify
			if err := availableCheck(srv.Interface.Health, true); err != nil {
				t.Error(err)
			} else if err := tickCheck(tick, i); err != nil {
				t.Error(err)
			}
		}

		srv.Stop()

		for i := 1; i <= 10; i++ {
			tick := <-notify
			if err := availableCheck(srv.Interface.Health, false); err != nil {
				t.Error(err)
			} else if err := tickCheck(tick, i+10); err != nil {
				t.Error(err)
			}
		}
	})

	t.Run("Multiple health service check", func(t *testing.T) {
		addr, _ := url.Parse("http://localhost:5008")
		srv := createTestServer(addr)

		service := server.NewHealthService(addr, time.Millisecond, time.Millisecond, 1)
		test_check := func(*server.HealthService) bool {
			return srv.IsRunning()
		}
		service.SetAliveCheck(test_check)

		notify_a, notify_b := make(chan int), make(chan int)
		srv.Interface.Health.Start(notify_a)
		defer srv.Interface.StopHealthCheck()
		service.Start(notify_b)
		defer service.Stop()

		for i := 0; i < 10; i++ {
			select {
			case <-notify_a:
				if err := availableCheck(srv.Interface.Health, false); err != nil {
					t.Error(err)
				}
			case <-notify_b:
				if err := availableCheck(service, false); err != nil {
					t.Error(err)
				}
			}
		}

		srv.Start(true)
		defer srv.Stop()

		for i := 0; i < 10; i++ {
			select {
			case <-notify_a:
				if err := availableCheck(srv.Interface.Health, true); err != nil {
					t.Error(err)
				}
			case <-notify_b:
				if err := availableCheck(service, true); err != nil {
					t.Error(err)
				}
			}
		}
	})
}
