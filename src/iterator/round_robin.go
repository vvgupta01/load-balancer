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
		curr: 0,
	}
}

func (iter *RoundRobin) Next() *server.ServerInterface {
	iter.mux.Lock()
	defer iter.mux.Unlock()

	iter.curr = (iter.curr + 1) % iter.pool.Len()
	_, srv := iter.pool.GetNextAvailable(iter.pool.DefaultOrder, iter.curr)
	return srv
}
