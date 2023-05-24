package iterator

import (
	server "load-balancer/src/server"
	"sort"
)

type LeastConnections struct {
	pool *server.ServerPool
}

func NewLeastConnections(pool *server.ServerPool) Iterator {
	return &LeastConnections{
		pool: pool,
	}
}

func (iter *LeastConnections) Next() *server.ServerInterface {
	min_order := make([]int, iter.pool.Len())
	copy(min_order, iter.pool.DefaultOrder)

	sort.SliceStable(min_order, func(i int, j int) bool {
		return iter.pool.Get(i).Health.GetLoad() < iter.pool.Get(j).Health.GetLoad()
	})
	return iter.pool.GetNextAvailable(min_order, 0)
}
