package iterator

import (
	"container/heap"
	server "loadbalancer/src/server"
	"sync"
)

type LeastConnections struct {
	pool      server.ServerPool
	pool_heap MinLoadHeap
	mux       sync.RWMutex
}

func NewLeastConnections(pool server.ServerPool) Iterator {
	h := make(MinLoadHeap, pool.Len())
	for i := range pool {
		h[i] = pool[i]
	}
	heap.Init(h)

	return &LeastConnections{
		pool:      pool,
		pool_heap: h,
	}
}

func (iter *LeastConnections) Next() int {
	if iter.pool.Len() == 0 {
		return -1
	}
	return 0
}

func (iter *LeastConnections) NextAvailable() (int, *server.ServerInterface) {
	iter.mux.Lock()
	defer iter.mux.Unlock()

	i := iter.Next()
	if i == -1 {
		return i, nil
	}

	srv := heap.Pop(iter.pool_heap).(*server.ServerInterface)
	if !srv.Health.IsAvailable() {
		return -1, nil
	}

	srv.Health.AddLoad(1)
	heap.Push(iter.pool_heap, srv)
	return i, srv
}

func (iter *LeastConnections) DoneCallback(i int) {
	iter.mux.Lock()
	defer iter.mux.Unlock()

	srv := heap.Remove(iter.pool_heap, i).(*server.ServerInterface)
	srv.Health.AddLoad(-1)
	heap.Push(iter.pool_heap, srv)
}
