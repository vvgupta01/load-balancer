package iterator

import (
	server "load-balancer/src/server"
	"math/rand"
)

type Random struct {
	pool *server.ServerPool
}

func NewRandom(seed func(), pool *server.ServerPool) Iterator {
	return &Random{
		pool: pool,
	}
}

func (iter *Random) Next() *server.ServerInterface {
	i := rand.Intn(iter.pool.Len())
	return iter.pool.GetNextAvailable(iter.pool.DefaultOrder, i)
}
