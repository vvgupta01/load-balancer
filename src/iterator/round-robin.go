package iterator

import (
	server "load-balancer/src/server"
	"sync/atomic"
)

type RoundRobinIterator struct {
	pool *server.ServerPool
	curr int64
}

func NewRoundRobinIterator(pool *server.ServerPool) Iterator {
	return &RoundRobinIterator{
		pool: pool,
		curr: -1,
	}
}

func (iter *RoundRobinIterator) Next() *server.ServerInterface {
	next := atomic.AddInt64(&iter.curr, int64(1)) % int64(iter.pool.Len())
	return iter.pool.GetNext(int(next))
}
