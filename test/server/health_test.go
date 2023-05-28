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

func TestHealthService(t *testing.T) {
	test.Setup()
	addr, _ := url.Parse("http://localhost:3001")

	t.Run("Available server health check", func(t *testing.T) {
		srv := server.NewServer(addr, nil, 1, 1)
		srv.Start()
		defer srv.Stop()

		time.Sleep(10 * time.Millisecond)
		
		srv.Interface.Health.HealthCheck()
		if err := availableCheck(srv.Interface.Health, true); err != nil {
			t.Error(err)
		}
	})

	t.Run("Unavailable server (not alive) health check", func(t *testing.T) {
		srv := server.NewServer(addr, nil, 1, 1)
		time.Sleep(10 * time.Millisecond)

		srv.Interface.Health.HealthCheck()
		if err := availableCheck(srv.Interface.Health, false); err != nil {
			t.Error(err)
		}
	})

	t.Run("Unavailable server (high load) health check", func(t *testing.T) {
		srv := server.NewServer(addr, nil, 1, 0)
		srv.Start()
		defer srv.Stop()

		time.Sleep(10 * time.Millisecond)

		srv.Interface.Health.HealthCheck()
		if err := availableCheck(srv.Interface.Health, false); err != nil {
			t.Error(err)
		}
	})

	t.Run("Health service check", func(t *testing.T) {
		srv := server.NewServer(addr, nil, 1, 1)
		srv.Start()
		srv.Interface.StartHealthCheck()
		defer srv.Stop()
		defer srv.Interface.StopHealthCheck()

		time.Sleep(20 * time.Millisecond)
		
		if err := availableCheck(srv.Interface.Health, true); err != nil {
			t.Error(err)
		}
	})

	t.Run("Server restart, health service check", func(t *testing.T) {
		srv := server.NewServer(addr, nil, 1, 1)
		srv.Start()
		srv.Interface.StartHealthCheck()
		defer srv.Interface.StopHealthCheck()
		defer srv.Stop()

		time.Sleep(20 * time.Millisecond)
		if err := availableCheck(srv.Interface.Health, true); err != nil {
			t.Error(err)
		}

		srv.Stop()

		time.Sleep(20 * time.Millisecond)
		if err := availableCheck(srv.Interface.Health, false); err != nil {
			t.Error(err)
		}

		srv.Start()

		time.Sleep(20 * time.Millisecond)
		if err := availableCheck(srv.Interface.Health, true); err != nil {
			t.Error(err)
		}
	})

	t.Run("Health service restart", func(t *testing.T) {
		srv := server.NewServer(addr, nil, 1, 1)
		srv.Start()
		srv.Interface.StartHealthCheck()
		defer srv.Interface.StopHealthCheck()

		time.Sleep(20 * time.Millisecond)
		if err := availableCheck(srv.Interface.Health, true); err != nil {
			t.Error(err)
		}

		srv.Interface.StopHealthCheck()
		srv.Stop()
		srv.Interface.StartHealthCheck()

		time.Sleep(20 * time.Millisecond)
		if err := availableCheck(srv.Interface.Health, false); err != nil {
			t.Error(err)
		}
	})
}