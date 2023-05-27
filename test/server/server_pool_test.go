package server_test

import (
	"fmt"
	"io/ioutil"
	server "load-balancer/src/server"
	"log"
	"net/url"
	"os"
	"testing"
)

func setup() {
	os.Setenv("HEALTH_INTERVAL", "0")
	os.Setenv("HEALTH_TIMEOUT", "0")
	log.SetFlags(0)
	log.SetOutput(ioutil.Discard)
}

func createTestPool(n int) *server.ServerPool {
	interfaces := make([]*server.ServerInterface, n)
	for i := 0; i < n; i++ {
		interfaces[i] = server.NewServerInterface(&url.URL{}, 1, 1)
	}
	return server.NewServerPool(interfaces)
}

func singleCheck(pool *server.ServerPool, i int, expected int) error {
	actual, _ := pool.GetNextAvailable(pool.DefaultOrder, i)
	if actual != expected {
		return fmt.Errorf("i=%d: Returned %d; Expected %d", i, actual, expected)
	}
	return nil
}

func checkPool(n int, unavailable []int, expected []int) error {
	if len(expected) != n {
		return fmt.Errorf("Invalid test: len(pool)=%d does not equal len(expected)=%d", n, len(expected))
	}

	pool := createTestPool(n)
	for _, i := range unavailable {
		pool.Get(i).Health.SetAlive(false)
	}

	for i := 0; i < n; i++ {
		if err := singleCheck(pool, i, expected[i]); err != nil {
			return err
		}
	}
	return nil
}

func TestGetNextAvailable(t *testing.T) {
	setup()

	t.Run("Empty pool", func(t *testing.T) {
		pool := createTestPool(0)
		if err := singleCheck(pool, 0, -1); err != nil {
			t.Error(err)
		}
	})

	t.Run("Single-server pool", func(t *testing.T) {
		pool := createTestPool(1)
		if err := singleCheck(pool, 0, 0); err != nil {
			t.Error(err)
		}
	})

	t.Run("Out-of-bounds index", func(t *testing.T) {
		pool := createTestPool(3)
		if err := singleCheck(pool, 10, 1); err != nil {
			t.Error(err)
		}
	})

	t.Run("Available small pool", func(t *testing.T) {
		unavailable := []int{}
		expected := []int{0, 1, 2, 3, 4}
		if err := checkPool(5, unavailable, expected); err != nil {
			t.Error(err)
		}
	})

	t.Run("Unavailable small pool", func(t *testing.T) {
		unavailable := []int{0, 1, 2, 3, 4}
		expected := []int{-1, -1, -1, -1, -1}
		if err := checkPool(5, unavailable, expected); err != nil {
			t.Error(err)
		}
	})

	t.Run("Alternating available pool", func(t *testing.T) {
		unavailable := []int{1, 3, 5, 7, 9}
		expected := []int{0, 2, 2, 4, 4, 6, 6, 8, 8, 0}
		if err := checkPool(10, unavailable, expected); err != nil {
			t.Error(err)
		}
	})

	t.Run("Half available pool", func(t *testing.T) {
		unavailable := []int{0, 1, 2, 3, 4}
		expected := []int{5, 5, 5, 5, 5, 5, 6, 7, 8, 9}
		if err := checkPool(10, unavailable, expected); err != nil {
			t.Error(err)
		}
	})

	t.Run("Test large pool", func(t *testing.T) {
		unavailable := []int{0, 1, 4, 9, 16, 25, 36, 49, 64, 81, 99}
		next_avail := []int{2, 2, 5, 10, 17, 26, 37, 50, 65, 82, 2}
		
		expected := make([]int, 100)
		for i := 0; i < 100; i++ {
			expected[i] = i
		}
		for i := 0; i < len(unavailable); i++ {
			expected[unavailable[i]] = next_avail[i]
		}

		if err := checkPool(100, unavailable, expected); err != nil {
			t.Error(err)
		}
	})
}