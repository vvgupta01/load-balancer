package iterator_test

import (
	"loadbalancer/src/iterator"
	"loadbalancer/test"
	"math/rand"
	"testing"
)

func TestWeightedRoundRobinNext(t *testing.T) {
	test.Setup()

	t.Run("Order check", func(t *testing.T) {
		pool := test.CreateTestPool(10)
		iter := iterator.NewWeightedRoundRobin(pool)

		order, _ := iter.Next()
		if err := test.CheckOrder(order, pool.DefaultOrder); err != nil {
			t.Error(err)
		}
	})

	t.Run("Empty pool", func(t *testing.T) {
		pool := test.CreateTestPool(0)
		iter := iterator.NewWeightedRoundRobin(pool)
		expected := []int{-1}

		if err := test.CheckIterator(iter, expected); err != nil {
			t.Error(err)
		}
	})

	t.Run("Equal weight pool single iteration", func(t *testing.T) {
		weights := []int32{1, 1, 1, 1, 1}
		pool := test.CreateTestWeightPool(weights)
		iter := iterator.NewWeightedRoundRobin(pool)
		expected := []int{0, 1, 2, 3, 4}

		if err := test.CheckIterator(iter, expected); err != nil {
			t.Error(err)
		}
	})

	t.Run("Unequal weight pool single iteration", func(t *testing.T) {
		weights := []int32{1, 2, 3, 2, 1}
		pool := test.CreateTestWeightPool(weights)
		iter := iterator.NewWeightedRoundRobin(pool)
		expected := []int{0, 1, 1, 2, 2, 2, 3, 3, 4}

		if err := test.CheckIterator(iter, expected); err != nil {
			t.Error(err)
		}
	})

	t.Run("Dominant weight pool", func(t *testing.T) {
		weights := []int32{10, 1, 1}
		pool := test.CreateTestWeightPool(weights)
		iter := iterator.NewWeightedRoundRobin(pool)
		expected := []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2}

		if err := test.CheckIterator(iter, expected); err != nil {
			t.Error(err)
		}
	})

	t.Run("Unequal weight multiple iterations", func(t *testing.T) {
		weights := []int32{2, 1, 2}
		pool := test.CreateTestWeightPool(weights)
		iter := iterator.NewWeightedRoundRobin(pool)
		expected := []int{0, 0, 1, 2, 2, 0, 0, 1, 2, 2}

		if err := test.CheckIterator(iter, expected); err != nil {
			t.Error(err)
		}
	})
}

func BenchmarkWeightedRoundRobinNext(b *testing.B) {
	rand.Seed(test.SEED)

	b.Run("100% available pool", func(b *testing.B) {
		pool := test.CreateRandomTestPool(10000, 0)
		iter := iterator.NewWeightedRoundRobin(pool)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			test.IterNext(iter)
		}
	})

	b.Run("90% available pool", func(b *testing.B) {
		pool := test.CreateRandomTestPool(10000, 1000)
		iter := iterator.NewWeightedRoundRobin(pool)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			test.IterNext(iter)
		}
	})

	b.Run("50% available pool", func(b *testing.B) {
		pool := test.CreateRandomTestPool(10000, 5000)
		iter := iterator.NewWeightedRoundRobin(pool)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			test.IterNext(iter)
		}
	})
}
