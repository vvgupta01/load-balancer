package test

import (
	"fmt"
	"loadbalancer/src/iterator"
	"loadbalancer/src/server"
	"loadbalancer/src/utils"
	"math/rand"
	"net/url"
	"os"
)

const SEED = 0

func Setup() {
	os.Setenv("HEALTH_INTERVAL", "1")
	os.Setenv("HEALTH_TIMEOUT", "1")
	utils.DisableLogOutput()
}

func ErrorIdx(i int, actual int, expected int) string {
	return fmt.Sprintf("i=%d: Returned %d; Expected %d", i, actual, expected)
}

func IterNext(iter iterator.Iterator) {
	srv := iter.NextAvailable()
	srv.Health.AddLoad(1)
}

func CreateTestPool(n int) *server.ServerPool {
	interfaces := make([]*server.ServerInterface, n)
	for i := range interfaces {
		interfaces[i] = server.NewServerInterface(&url.URL{}, 1, 1)
	}
	return server.NewServerPool(interfaces)
}

func CreateTestAlivePool(n int, unavailable []int) *server.ServerPool {
	pool := CreateTestPool(n)
	for _, i := range unavailable {
		pool.Get(i).Health.SetAlive(false)
	}
	return pool
}

func CreateTestLoadPool(loads []int32) *server.ServerPool {
	pool := CreateTestPool(len(loads))
	for i := range loads {
		pool.Get(i).Health.SetLoad(loads[i])
	}
	return pool
}

func CreateTestWeightPool(weights []int32) *server.ServerPool {
	pool := CreateTestPool(len(weights))
	for i := range weights {
		pool.Get(i).Weight = weights[i]
	}
	return pool
}

func CreateRandomTestPool(n int, n_unavail int) *server.ServerPool {
	interfaces := make([]*server.ServerInterface, n)
	for i := range interfaces {
		weight := int32(rand.Intn(10) + 1)
		capacity := int32(rand.Intn(9000) + 1000)
		interfaces[i] = server.NewServerInterface(&url.URL{}, weight, capacity)
	}
	pool := server.NewServerPool(interfaces)

	perm := rand.Perm(n)
	for i := 0; i < n_unavail; i++ {
		pool.Get(perm[i]).Health.SetAlive(false)
	}
	return pool
}

func CheckOrder(actual []int, expected []int) error {
	if len(actual) != len(expected) {
		return fmt.Errorf("len: Returned %d; Expected %d", len(actual), len(expected))
	}

	for i := range actual {
		if actual[i] != expected[i] {
			return fmt.Errorf(ErrorIdx(i, actual[i], expected[i]))
		}
	}
	return nil
}

func CheckIterator(iter iterator.Iterator, expected []int) error {
	for i := range expected {
		if _, actual := iter.Next(); actual != expected[i] {
			return fmt.Errorf(ErrorIdx(i, actual, expected[i]))
		}
	}
	return nil
}
