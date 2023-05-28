package test

import (
	"fmt"
	"io/ioutil"
	"load-balancer/src/iterator"
	"load-balancer/src/server"
	"log"
	"net/url"
	"os"
)

func Setup() {
	os.Setenv("HEALTH_INTERVAL", "0")
	os.Setenv("HEALTH_TIMEOUT", "0")
	log.SetFlags(0)
	log.SetOutput(ioutil.Discard)
}

func ErrorIdx(i int, actual int, expected int) string {
	return fmt.Sprintf("i=%d: Returned %d; Expected %d", i, actual, expected)
}

func CreateTestPool(n int) *server.ServerPool {
	interfaces := make([]*server.ServerInterface, n)
	for i := 0; i < n; i++ {
		interfaces[i] = server.NewServerInterface(&url.URL{}, 1, 1)
	}
	return server.NewServerPool(interfaces)
}

func CreateTestLoadPool(loads []int32) *server.ServerPool {
	pool := CreateTestPool(len(loads))
	for i := 0; i < len(loads); i++ {
		pool.Get(i).Health.SetLoad(loads[i])
	}
	return pool
}

func CreateTestWeightPool(weights []int32) *server.ServerPool {
	pool := CreateTestPool(len(weights))
	for i := 0; i < len(weights); i++ {
		pool.Get(i).Weight = weights[i]
	}
	return pool
}

func CheckOrder(actual []int, expected []int) error {
	if len(actual) != len(expected) {
		return fmt.Errorf("len: Returned %d; Expected %d", len(actual), len(expected))
	}

	for i := 0; i < len(actual); i++ {
		if actual[i] != expected[i] {
			return fmt.Errorf(ErrorIdx(i, actual[i], expected[i]))
		}
	}
	return nil
}

func CheckIterator(iter iterator.Iterator, expected []int) error {
	for i := 0; i < len(expected); i++ {
		if _, actual := iter.Next(); actual != expected[i] {
			return fmt.Errorf(ErrorIdx(i, actual, expected[i]))
		}
	}
	return nil
}
