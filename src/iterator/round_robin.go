package iterator

import (
	server "load-balancer/src/server"
	"sync"
)

type RoundRobin struct {
	pool *server.ServerPool
	curr int
	mux  sync.RWMutex
}

func NewRoundRobin(pool *server.ServerPool) Iterator {
	return &RoundRobin{
		pool: pool,
		curr: -1,
	}
}

func (iter *RoundRobin) Next() ([]int, int) {
	iter.mux.Lock()
	defer iter.mux.Unlock()

	if iter.pool.Len() == 0 {
		return nil, -1
	}
	iter.curr = (iter.curr + 1) % iter.pool.Len()
	return iter.pool.DefaultOrder, iter.curr
}

func (iter *RoundRobin) NextAvailable() *server.ServerInterface {
	order, i := iter.Next()
	_, srv := iter.pool.GetNextAvailable(order, i)
	return srv
}
