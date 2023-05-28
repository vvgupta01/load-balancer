package iterator_test

import (
	"loadbalancer/src/iterator"
	"loadbalancer/test"
	"testing"
)

func TestRoundRobinNext(t *testing.T) {
	test.Setup()

	t.Run("Order check", func(t *testing.T) {
		pool := test.CreateTestPool(10)
		iter := iterator.NewRoundRobin(pool)

		order, _ := iter.Next()
		if err := test.CheckOrder(order, pool.DefaultOrder); err != nil {
			t.Error(err)
		}
	})

	t.Run("Empty pool", func(t *testing.T) {
		pool := test.CreateTestPool(0)
		iter := iterator.NewRoundRobin(pool)
		expected := []int{-1}

		if err := test.CheckIterator(iter, expected); err != nil {
			t.Error(err)
		}
	})

	t.Run("Single iteration", func(t *testing.T) {
		pool := test.CreateTestPool(10)
		iter := iterator.NewRoundRobin(pool)
		expected := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

		if err := test.CheckIterator(iter, expected); err != nil {
			t.Error(err)
		}
	})

	t.Run("Multiple iterations", func(t *testing.T) {
		pool := test.CreateTestPool(5)
		iter := iterator.NewRoundRobin(pool)
		expected := []int{0, 1, 2, 3, 4, 0, 1, 2, 3, 4}

		if err := test.CheckIterator(iter, expected); err != nil {
			t.Error(err)
		}
	})
}
