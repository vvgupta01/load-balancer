package server

type ServerPool struct {
	servers      []*ServerInterface
	DefaultOrder []int
}

func NewServerPool(servers []*ServerInterface) *ServerPool {
	order := make([]int, len(servers))
	for i := range servers {
		order[i] = i
	}
	return &ServerPool{
		servers:      servers,
		DefaultOrder: order,
	}
}

func (pool *ServerPool) GetNextAvailable(order []int, idx int) (int, *ServerInterface) {
	if order == nil || idx < 0 {
		return -1, nil
	}

	for i := range order {
		try_idx := (idx + i) % pool.Len()
		server := pool.Get(order[try_idx])

		if server.Health.IsAvailable() {
			return order[try_idx], server
		}
	}
	return -1, nil
}

func (pool *ServerPool) Len() int {
	return len(pool.servers)
}

func (pool *ServerPool) Get(i int) *ServerInterface {
	return pool.servers[i]
}
