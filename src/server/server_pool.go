package server

type ServerPool struct {
	servers    []*ServerInterface
	Order 	   []int
}

func NewServerPool(servers []*ServerInterface) *ServerPool {
	order := make([]int, len(servers))
	for i := 0; i < len(servers); i++ {
		order[i] = i
	}

	return &ServerPool{
		servers: servers,
		Order: order,
	}
}

func (pool *ServerPool) GetNextAvailable(order []int, idx int) *ServerInterface {
	for i := 0; i < pool.Len(); i++ {
		try_idx := (idx + i) % pool.Len()
		server := pool.Get(order[try_idx])

		if server.Health.IsAvailable() {
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
