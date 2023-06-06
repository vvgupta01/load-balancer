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
	os.Setenv("HEALTH_INTERVAL", "1ms")
	os.Setenv("HEALTH_TIMEOUT", "1ms")
	utils.DisableLogOutput()
}

func ErrorIdx(i int, actual int, expected int) string {
	return fmt.Sprintf("i=%d: Returned %d; Expected %d", i, actual, expected)
}

func TestNext(iter iterator.Iterator) {
	_, srv := iter.NextAvailable()
	iter.DoneCallback(srv.Index)
}

func CreateTestPool(n int, loads []int32, capacities []int32, weights []int32, unavailable []int) server.ServerPool {
	pool := make(server.ServerPool, n)
	for i := range pool {
		capacity, weight := int32(100), int32(1)
		if capacities != nil {
			capacity = capacities[i]
		}

		if weights != nil {
			weight = weights[i]
		}

		addr, _ := url.Parse("")
		pool[i] = server.NewServerInterface(addr, i, weight, capacity)

		if loads != nil {
			pool[i].Health.SetLoad(loads[i])
		}
	}

	if unavailable != nil {
		for _, i := range unavailable {
			pool[i].Health.SetAlive(false)
		}
	}
	return pool
}

func CreateDefaultTestPool(n int) server.ServerPool {
	return CreateTestPool(n, nil, nil, nil, nil)
}

func CreateRandomTestPool(n int, n_unavail int) server.ServerPool {
	pool := make(server.ServerPool, n)
	for i := range pool {
		weight := int32(rand.Intn(10) + 1)
		capacity := int32(rand.Intn(9000) + 1000)
		pool[i] = server.NewServerInterface(&url.URL{}, i, weight, capacity)
	}

	perm := rand.Perm(n)
	for i := 0; i < n_unavail; i++ {
		pool[perm[i]].Health.SetAlive(false)
	}
	return pool
}

func CheckIterNextAvailable(iter iterator.Iterator, exp_i []int, exp_next []int) error {
	for i := range exp_i {
		act_i, act_next := iter.NextAvailable()

		if act_i != exp_i[i] {
			return fmt.Errorf("Index: %s", ErrorIdx(i, act_i, exp_i[i]))
		} else if act_next == nil && exp_next[i] != -1 {
			return fmt.Errorf(ErrorIdx(i, -1, exp_next[i]))
		} else if act_next != nil && act_next.Index != exp_next[i] {
			return fmt.Errorf("Next: %s", ErrorIdx(i, act_next.Index, exp_next[i]))
		}
	}
	return nil
}
