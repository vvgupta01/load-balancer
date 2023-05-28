package iterator

import (
	"loadbalancer/src/server"
	"sync"
)

type WeightedRoundRobin struct {
	pool     *server.ServerPool
	curr     int
	curr_req int32
	mux      sync.RWMutex
}

func NewWeightedRoundRobin(pool *server.ServerPool) Iterator {
	return &WeightedRoundRobin{
		pool:     pool,
		curr:     0,
		curr_req: 0,
	}
}

func (iter *WeightedRoundRobin) Next() ([]int, int) {
	iter.mux.Lock()
	defer iter.mux.Unlock()

	if iter.pool.Len() == 0 {
		return nil, -1
	}

	if iter.curr_req >= iter.pool.Get(iter.curr).Weight {
		iter.curr_req = 0
		iter.curr = (iter.curr + 1) % iter.pool.Len()
	}
	iter.curr_req++
	return iter.pool.DefaultOrder, iter.curr
}

func (iter *WeightedRoundRobin) NextAvailable() *server.ServerInterface {
	order, i := iter.Next()
	_, srv := iter.pool.GetNextAvailable(order, i)
	return srv
}
