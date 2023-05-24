package iterator

import (
	"load-balancer/src/server"
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

func (iter *WeightedRoundRobin) Next() *server.ServerInterface {
	iter.mux.Lock()
	defer iter.mux.Unlock()

	if iter.curr_req >= iter.pool.Get(iter.curr).Weight {
		iter.curr_req = 0
		iter.curr = (iter.curr + 1) % iter.pool.Len()
	}
	iter.curr_req++
	return iter.pool.GetNextAvailable(iter.pool.DefaultOrder, iter.curr)
}
