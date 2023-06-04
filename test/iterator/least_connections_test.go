package iterator_test

import (
	"loadbalancer/src/iterator"
	"loadbalancer/test"
	"math/rand"
	"testing"
)

func TestLeastConnectionsNext(t *testing.T) {
	test.Setup()

	t.Run("Empty pool", func(t *testing.T) {
		pool := test.CreateDefaultTestPool(0)
		iter := iterator.NewLeastConnections(pool)
		expected := []int{-1}

		if err := test.CheckIterNextAvailable(iter, expected, expected); err != nil {
			t.Error(err)
		}
	})

	t.Run("Unavailable pool", func(t *testing.T) {
		unavailable := []int{0, 1, 2, 3, 4}
		pool := test.CreateTestPool(5, nil, nil, nil, unavailable)
		iter := iterator.NewLeastConnections(pool)
		expected := []int{-1, -1, -1, -1, -1}

		if err := test.CheckIterNextAvailable(iter, expected, expected); err != nil {
			t.Error(err)
		}
	})

	t.Run("All available, equal load", func(t *testing.T) {
		pool := test.CreateDefaultTestPool(5)
		iter := iterator.NewLeastConnections(pool)
		exp_i := make([]int, 5)
		exp_next := []int{0, 1, 2, 3, 4}

		if err := test.CheckIterNextAvailable(iter, exp_i, exp_next); err != nil {
			t.Error(err)
		}
	})

	t.Run("Random available, equal load", func(t *testing.T) {
		unavailable := []int{0, 1, 4}
		pool := test.CreateTestPool(5, nil, nil, nil, unavailable)
		iter := iterator.NewLeastConnections(pool)

		exp_i := make([]int, 4)
		exp_next := []int{2, 3, 2, 3}

		if err := test.CheckIterNextAvailable(iter, exp_i, exp_next); err != nil {
			t.Error(err)
		}
	})

	t.Run("All available, random load", func(t *testing.T) {
		loads := []int32{0, 1, 2, 1, 0}
		pool := test.CreateTestPool(5, loads, nil, nil, nil)
		iter := iterator.NewLeastConnections(pool)

		exp_i := make([]int, 5)
		exp_next := []int{0, 4, 0, 1, 3, 4, 0, 1, 2, 3, 4, 0}

		if err := test.CheckIterNextAvailable(iter, exp_i, exp_next); err != nil {
			t.Error(err)
		}
	})

	t.Run("Random available, random load", func(t *testing.T) {
		loads := []int32{0, 1, 1, 2, 3}
		capacities := []int32{1, 1, 2, 5, 5}
		pool := test.CreateTestPool(5, loads, capacities, nil, nil)
		iter := iterator.NewLeastConnections(pool)

		exp_i := make([]int, 5)
		exp_next := []int{0, 2, 3, 3, 4, 3, 4, -1}

		if err := test.CheckIterNextAvailable(iter, exp_i, exp_next); err != nil {
			t.Error(err)
		}
	})
}

func BenchmarkLeastConnectionsNext(b *testing.B) {
	rand.Seed(test.SEED)

	b.Run("100% available pool", func(b *testing.B) {
		pool := test.CreateRandomTestPool(10000, 0)
		iter := iterator.NewLeastConnections(pool)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			test.TestNext(iter)
		}
	})

	b.Run("90% available pool", func(b *testing.B) {
		pool := test.CreateRandomTestPool(10000, 1000)
		iter := iterator.NewLeastConnections(pool)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			test.TestNext(iter)
		}
	})

	b.Run("50% available pool", func(b *testing.B) {
		pool := test.CreateRandomTestPool(10000, 5000)
		iter := iterator.NewLeastConnections(pool)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			test.TestNext(iter)
		}
	})
}
