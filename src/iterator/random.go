package iterator

import (
	server "load-balancer/src/server"
	"math/rand"
)

type RandomIterator struct {
	pool *server.ServerPool
}

func NewRandomIterator(seed func(), pool *server.ServerPool) Iterator {
	return &RandomIterator{
		pool: pool,
	}
}

func (iter *RandomIterator) Next() *server.ServerInterface {
	i := rand.Intn(iter.pool.Len())
	return iter.pool.GetNext(i)
}
