package server

type ServerPool []*ServerInterface

func (pool ServerPool) GetNextAvailable(idx int) (int, *ServerInterface) {
	if idx < 0 {
		return -1, nil
	}

	for i := range pool {
		try_idx := (idx + i) % pool.Len()
		srv := pool[try_idx]

		if srv.Health.IsAvailable() {
			return try_idx, srv
		}
	}
	return -1, nil
}

func (pool ServerPool) Len() int {
	return len(pool)
}
