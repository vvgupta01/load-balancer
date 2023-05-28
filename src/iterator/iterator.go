package iterator

import server "loadbalancer/src/server"

type Iterator interface {
	Next() ([]int, int)

	NextAvailable() *server.ServerInterface
}
