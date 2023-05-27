package iterator

import (
	server "load-balancer/src/server"

	"github.com/mkmik/argsort"
)

type LeastConnections struct {
	pool *server.ServerPool
}

func NewLeastConnections(pool *server.ServerPool) Iterator {
	return &LeastConnections{
		pool: pool,
	}
}

func (iter *LeastConnections) Next() ([]int, int) {
	if iter.pool.Len() == 0 {
		return nil, -1
	}

	loads := make([]int32, iter.pool.Len())
	for i := 0; i < iter.pool.Len(); i++ {
		loads[i] = iter.pool.Get(i).Health.GetLoad()
	}

	order := argsort.SortSlice(loads, func (i int, j int) bool {
		return loads[i] < loads[j]
	})
	return order, 0
}

func (iter *LeastConnections) NextAvailable() *server.ServerInterface {
	order, i := iter.Next()
	_, srv := iter.pool.GetNextAvailable(order, i)
	return srv
}
