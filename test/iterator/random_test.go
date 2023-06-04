package iterator_test

import (
	"loadbalancer/src/iterator"
	"loadbalancer/test"
	"math/rand"
	"testing"
)

func TestRandomNext(t *testing.T) {
	test.Setup()

	seed := func() {
		rand.Seed(test.SEED)
	}

	t.Run("Empty pool", func(t *testing.T) {
		pool := test.CreateDefaultTestPool(0)
		iter := iterator.NewRandom(seed, pool)
		expected := []int{-1}

		if err := test.CheckIterNextAvailable(iter, expected, expected); err != nil {
			t.Error(err)
		}
	})

	t.Run("Seed check", func(t *testing.T) {
		pool := test.CreateDefaultTestPool(10)
		iter := iterator.NewRandom(seed, pool)

		r := rand.New(rand.NewSource(test.SEED))
		expected := make([]int, 100)
		for i := 0; i < len(expected); i++ {
			expected[i] = r.Intn(pool.Len())
		}

		if err := test.CheckIterNextAvailable(iter, expected, expected); err != nil {
			t.Error(err)
		}
	})
}

func BenchmarkRandomNext(b *testing.B) {
	seed := func() {
		rand.Seed(test.SEED)
	}

	b.Run("100% available pool", func(b *testing.B) {
		pool := test.CreateRandomTestPool(10000, 0)
		iter := iterator.NewRandom(seed, pool)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			test.TestNext(iter)
		}
	})

	b.Run("90% available pool", func(b *testing.B) {
		pool := test.CreateRandomTestPool(10000, 1000)
		iter := iterator.NewRandom(seed, pool)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			test.TestNext(iter)
		}
	})

	b.Run("50% available pool", func(b *testing.B) {
		pool := test.CreateRandomTestPool(10000, 5000)
		iter := iterator.NewRandom(seed, pool)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			test.TestNext(iter)
		}
	})
}
