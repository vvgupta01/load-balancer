package iterator

import (
	server "load-balancer/src/server"
	"math"
	"sync/atomic"
)

type RoundRobinIterator struct {
	pool *server.ServerPool
	curr uint64
}

func NewRoundRobinIterator(pool *server.ServerPool) Iterator {
	return &RoundRobinIterator{
		pool: pool,
		curr: math.MaxUint64,
	}
}

func (iter *RoundRobinIterator) Next() *server.ServerInterface {
	next := atomic.AddUint64(&iter.curr, uint64(1)) % uint64(iter.pool.Len())
	return iter.pool.GetNextAvailable(iter.pool.Order, int(next))
}
