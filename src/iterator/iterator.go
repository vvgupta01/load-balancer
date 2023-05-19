package iterator

import server "load-balancer/src/server"

type Iterator interface {
	Next() *server.ServerInterface
}
