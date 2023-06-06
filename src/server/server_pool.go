package server

type ServerPool []*ServerInterface

type PoolStatus struct {
	Total_load         int32
	Percent_load       float32
	Avg_load           float32
	Available_capacity int32
	Available          int
	Percent_available  float32
	Transactions       int64
}

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

func (pool ServerPool) GetStatus() PoolStatus {
	total_load, total_capacity := int32(0), int32(0)
	available := 0
	for _, srv := range pool {
		total_load += srv.Health.GetLoad()
		if srv.Health.IsAvailable() {
			available++
			total_capacity += srv.Health.GetCapacity()
		}
	}

	return PoolStatus{
		Total_load:         total_load,
		Percent_load:       float32(total_load) / float32(total_capacity),
		Avg_load:           float32(total_load) / float32(available),
		Available_capacity: total_capacity,
		Available:          available,
		Percent_available:  float32(available) / float32(pool.Len()),
	}
}
