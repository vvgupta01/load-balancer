package iterator_test

import (
	"load-balancer/src/iterator"
	"load-balancer/test"
	"testing"
)

func TestLeastConnectionsNext(t *testing.T) {
	test.Setup()

	t.Run("Index check", func(t *testing.T) {
		pool := test.CreateTestPool(10)
		iter := iterator.NewLeastConnections(pool)
		expected := make([]int, 10)

		if err := test.CheckIterator(iter, expected); err != nil {
			t.Error(err)
		}
	})

	t.Run("Empty pool", func(t *testing.T) {
		pool := test.CreateTestPool(0)
		iter := iterator.NewLeastConnections(pool)
		expected := []int{-1}

		if err := test.CheckIterator(iter, expected); err != nil {
			t.Error(err)
		}
	})

	t.Run("Default order", func(t *testing.T) {
		loads := []int32{1, 2, 3, 4, 5}
		pool := test.CreateTestLoadPool(loads)
		iter := iterator.NewLeastConnections(pool)
		expected := []int{0, 1, 2, 3, 4}

		actual, _ := iter.Next()
		if err := test.CheckOrder(actual, expected); err != nil {
			t.Error(err)
		}
	})

	t.Run("Reverse order", func(t *testing.T) {
		loads := []int32{5, 4, 3, 2, 1}
		pool := test.CreateTestLoadPool(loads)
		iter := iterator.NewLeastConnections(pool)
		expected := []int{4, 3, 2, 1, 0}

		actual, _ := iter.Next()
		if err := test.CheckOrder(actual, expected); err != nil {
			t.Error(err)
		}
	})

	t.Run("Random order", func(t *testing.T) {
		loads := []int32{1, 3, 5, 2, 4}
		pool := test.CreateTestLoadPool(loads)
		iter := iterator.NewLeastConnections(pool)
		expected := []int{0, 3, 1, 4, 2}

		actual, _ := iter.Next()
		if err := test.CheckOrder(actual, expected); err != nil {
			t.Error(err)
		}
	})

	t.Run("Stable order", func(t *testing.T) {
		loads := []int32{2, 2, 2, 1, 1}
		pool := test.CreateTestLoadPool(loads)
		iter := iterator.NewLeastConnections(pool)
		expected := []int{3, 4, 0, 1, 2}

		actual, _ := iter.Next()
		if err := test.CheckOrder(actual, expected); err != nil {
			t.Error(err)
		}
	})
}
