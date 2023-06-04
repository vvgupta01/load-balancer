package iterator

import (
	server "loadbalancer/src/server"
	"sync"
)

type RoundRobin struct {
	pool server.ServerPool
	curr int
	mux  sync.RWMutex
}

func NewRoundRobin(pool server.ServerPool) Iterator {
	return &RoundRobin{
		pool: pool,
		curr: -1,
	}
}

func (iter *RoundRobin) Next() int {
	if iter.pool.Len() == 0 {
		return -1
	}
	iter.curr = (iter.curr + 1) % iter.pool.Len()
	return iter.curr
}

func (iter *RoundRobin) NextAvailable() (int, *server.ServerInterface) {
	iter.mux.Lock()
	defer iter.mux.Unlock()

	i := iter.Next()
	avail_i, srv := iter.pool.GetNextAvailable(i)
	if srv == nil {
		return -1, nil
	}

	iter.curr = avail_i
	srv.Health.AddLoad(1)
	return i, srv
}

func (iter *RoundRobin) DoneCallback(i int) {
	iter.pool[i].Health.AddLoad(-1)
}
