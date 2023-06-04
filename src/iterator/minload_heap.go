package iterator

import "loadbalancer/src/server"

type MinLoadHeap []*server.ServerInterface

// Based on https://pkg.go.dev/container/heap
func (h MinLoadHeap) Len() int {
	return len(h)
}

func (h MinLoadHeap) Less(i, j int) bool {
	i_avail, j_avail := h[i].Health.IsAvailable(), h[j].Health.IsAvailable()
	if i_avail != j_avail {
		return i_avail
	}

	i_load, j_load := h[i].Health.GetLoad(), h[j].Health.GetLoad()
	if i_load != j_load {
		return i_load < j_load
	}
	return h[i].Index < h[j].Index
}

func (h MinLoadHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h MinLoadHeap) Push(x any) {
	h = append(h, x.(*server.ServerInterface))
}

func (h MinLoadHeap) Pop() any {
	n := h.Len()
	x := h[n-1]
	h = h[0 : n-1]
	return x
}
