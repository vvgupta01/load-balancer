package iterator

import (
	server "loadbalancer/src/server"
	"math/rand"
	"time"
)

type Random struct {
	pool server.ServerPool
}

func NewRandom(seed func(), pool server.ServerPool) Iterator {
	seed()
	return &Random{
		pool: pool,
	}
}

func DefaultSeed() {
	rand.Seed(time.Now().UnixNano())
}

func (iter *Random) Next() int {
	if iter.pool.Len() == 0 {
		return -1
	}
	return rand.Intn(iter.pool.Len())
}

func (iter *Random) NextAvailable() (int, *server.ServerInterface) {
	i := iter.Next()
	_, srv := iter.pool.GetNextAvailable(i)
	if srv == nil {
		return -1, nil
	}

	srv.Health.AddLoad(1)
	return i, srv
}

func (iter *Random) DoneCallback(i int) {
	iter.pool[i].Health.AddLoad(-1)
}
