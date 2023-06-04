package iterator

import server "loadbalancer/src/server"

type Iterator interface {
	Next() int

	NextAvailable() (int, *server.ServerInterface)

	DoneCallback(int)
}
