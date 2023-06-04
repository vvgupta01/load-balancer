package iterator_test

import (
	"loadbalancer/src/iterator"
	"loadbalancer/test"
	"math/rand"
	"testing"
)

func TestRoundRobinNext(t *testing.T) {
	test.Setup()

	t.Run("Empty pool", func(t *testing.T) {
		pool := test.CreateDefaultTestPool(0)
		iter := iterator.NewRoundRobin(pool)
		expected := []int{-1}

		if err := test.CheckIterNextAvailable(iter, expected, expected); err != nil {
			t.Error(err)
		}
	})

	t.Run("Single iteration", func(t *testing.T) {
		pool := test.CreateDefaultTestPool(10)
		iter := iterator.NewRoundRobin(pool)
		expected := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

		if err := test.CheckIterNextAvailable(iter, expected, expected); err != nil {
			t.Error(err)
		}
	})

	t.Run("Multiple iterations", func(t *testing.T) {
		pool := test.CreateDefaultTestPool(5)
		iter := iterator.NewRoundRobin(pool)
		expected := []int{0, 1, 2, 3, 4, 0, 1, 2, 3, 4}

		if err := test.CheckIterNextAvailable(iter, expected, expected); err != nil {
			t.Error(err)
		}
	})

	t.Run("Unavailable servers (index update check)", func(t *testing.T) {
		unavailable := []int{0, 2, 4}
		pool := test.CreateTestPool(5, nil, nil, nil, unavailable)
		iter := iterator.NewRoundRobin(pool)
		exp_i := []int{0, 2, 4, 2}
		exp_next := []int{1, 3, 1, 3}

		if err := test.CheckIterNextAvailable(iter, exp_i, exp_next); err != nil {
			t.Error(err)
		}
	})
}

func BenchmarkRoundRobinNext(b *testing.B) {
	rand.Seed(test.SEED)

	b.Run("100% available pool", func(b *testing.B) {
		pool := test.CreateRandomTestPool(10000, 0)
		iter := iterator.NewRoundRobin(pool)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			test.TestNext(iter)
		}
	})

	b.Run("90% available pool", func(b *testing.B) {
		pool := test.CreateRandomTestPool(10000, 1000)
		iter := iterator.NewRoundRobin(pool)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			test.TestNext(iter)
		}
	})

	b.Run("50% available pool", func(b *testing.B) {
		pool := test.CreateRandomTestPool(10000, 5000)
		iter := iterator.NewRoundRobin(pool)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			test.TestNext(iter)
		}
	})
}
