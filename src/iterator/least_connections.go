package iterator

import (
	server "load-balancer/src/server"
	"sort"
)

type LeastConnectionsIterator struct {
	pool *server.ServerPool
}

func NewLeastConnectionsIterator(pool *server.ServerPool) Iterator {
	return &LeastConnectionsIterator{
		pool: pool,
	}
}

func (iter* LeastConnectionsIterator) Next() *server.ServerInterface {
	min_order := make([]int, iter.pool.Len())
	copy(min_order, iter.pool.Order)

	sort.SliceStable(min_order, func (i int, j int) bool {
		return iter.pool.Get(i).Health.GetLoad() < iter.pool.Get(j).Health.GetLoad()
	})
	return iter.pool.GetNextAvailable(min_order, 0)
}