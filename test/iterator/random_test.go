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
		rand.Seed(0)
	}

	t.Run("Order check", func(t *testing.T) {
		pool := test.CreateTestPool(10)
		iter := iterator.NewRandom(seed, pool)

		order, _ := iter.Next()
		if err := test.CheckOrder(order, pool.DefaultOrder); err != nil {
			t.Error(err)
		}
	})

	t.Run("Empty pool", func(t *testing.T) {
		pool := test.CreateTestPool(0)
		iter := iterator.NewRandom(seed, pool)
		expected := []int{-1}

		if err := test.CheckIterator(iter, expected); err != nil {
			t.Error(err)
		}
	})

	t.Run("Seed check", func(t *testing.T) {
		pool := test.CreateTestPool(10)
		iter := iterator.NewRandom(seed, pool)

		r := rand.New(rand.NewSource(0))
		expected := make([]int, 100)
		for i := 0; i < len(expected); i++ {
			expected[i] = r.Intn(pool.Len())
		}

		if err := test.CheckIterator(iter, expected); err != nil {
			t.Error(err)
		}
	})
}
