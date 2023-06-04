package iterator

import (
	"loadbalancer/src/server"
	"sync"
)

type WeightedRoundRobin struct {
	pool     server.ServerPool
	curr     int
	curr_req int32
	mux      sync.RWMutex
}

func NewWeightedRoundRobin(pool server.ServerPool) Iterator {
	return &WeightedRoundRobin{
		pool:     pool,
		curr:     0,
		curr_req: 0,
	}
}

func (iter *WeightedRoundRobin) Next() int {
	if iter.pool.Len() == 0 {
		return -1
	}

	if iter.curr_req >= iter.pool[iter.curr].Weight {
		iter.curr_req = 0
		iter.curr = (iter.curr + 1) % iter.pool.Len()
	}
	iter.curr_req++
	return iter.curr
}

func (iter *WeightedRoundRobin) NextAvailable() (int, *server.ServerInterface) {
	iter.mux.Lock()
	defer iter.mux.Unlock()

	i := iter.Next()
	avail_i, srv := iter.pool.GetNextAvailable(i)
	if srv == nil {
		return -1, nil
	}

	if avail_i != i {
		iter.curr = (iter.curr + 1) % iter.pool.Len()
		iter.curr_req = 1
	}
	srv.Health.AddLoad(1)
	return i, srv
}

func (iter *WeightedRoundRobin) DoneCallback(i int) {
	iter.pool[i].Health.AddLoad(-1)
}
