package server

type ServerPool struct {
	servers []*ServerInterface
}

func NewServerPool(servers []*ServerInterface) *ServerPool {
	return &ServerPool{
		servers: servers,
	}
}

func (pool *ServerPool) GetNext(idx int) *ServerInterface {
	for i := 0; i < pool.Len(); i++ {
		try_idx := (idx + i) % pool.Len()
		server := pool.Get(try_idx)

		if server.alive {
			return server
		}
	}
	return nil
}

func (pool *ServerPool) Len() int {
	return len(pool.servers)
}

func (pool *ServerPool) Get(i int) *ServerInterface {
	return pool.servers[i]
}
