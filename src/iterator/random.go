package iterator

import (
	server "load-balancer/src/server"
	"math/rand"
	"time"
)

type Random struct {
	pool *server.ServerPool
}

func NewRandom(seed func(), pool *server.ServerPool) Iterator {
	seed()
	return &Random{
		pool: pool,
	}
}

func DefaultSeed() func() {
	return func() {
		rand.Seed(time.Now().UnixNano())
	}
}

func (iter *Random) Next() ([]int, int) {
	if iter.pool.Len() == 0 {
		return nil, -1
	}

	i := rand.Intn(iter.pool.Len())
	return iter.pool.DefaultOrder, i
}

func (iter *Random) NextAvailable() *server.ServerInterface {
	order, i := iter.Next()
	_, srv := iter.pool.GetNextAvailable(order, i)
	return srv
}

