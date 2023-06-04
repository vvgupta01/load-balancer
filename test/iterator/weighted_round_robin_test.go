package iterator_test

import (
	"loadbalancer/src/iterator"
	"loadbalancer/test"
	"math/rand"
	"testing"
)

func TestWeightedRoundRobinNext(t *testing.T) {
	test.Setup()

	t.Run("Empty pool", func(t *testing.T) {
		pool := test.CreateDefaultTestPool(0)
		iter := iterator.NewWeightedRoundRobin(pool)
		expected := []int{-1}

		if err := test.CheckIterNextAvailable(iter, expected, expected); err != nil {
			t.Error(err)
		}
	})

	t.Run("Equal weight pool single iteration", func(t *testing.T) {
		weights := []int32{1, 1, 1, 1, 1}
		pool := test.CreateTestPool(5, nil, nil, weights, nil)
		iter := iterator.NewWeightedRoundRobin(pool)
		expected := []int{0, 1, 2, 3, 4}

		if err := test.CheckIterNextAvailable(iter, expected, expected); err != nil {
			t.Error(err)
		}
	})

	t.Run("Dominant weight pool", func(t *testing.T) {
		weights := []int32{10, 1, 1}
		pool := test.CreateTestPool(3, nil, nil, weights, nil)
		iter := iterator.NewWeightedRoundRobin(pool)
		expected := []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2}

		if err := test.CheckIterNextAvailable(iter, expected, expected); err != nil {
			t.Error(err)
		}
	})

	t.Run("Unequal weight multiple iterations", func(t *testing.T) {
		weights := []int32{2, 1, 2}
		pool := test.CreateTestPool(3, nil, nil, weights, nil)
		iter := iterator.NewWeightedRoundRobin(pool)
		expected := []int{0, 0, 1, 2, 2, 0, 0, 1, 2, 2}

		if err := test.CheckIterNextAvailable(iter, expected, expected); err != nil {
			t.Error(err)
		}
	})

	t.Run("Unavailable servers (index update check)", func(t *testing.T) {
		weights := []int32{1, 1, 2, 3}
		unavailable := []int{1, 3}
		pool := test.CreateTestPool(4, nil, nil, weights, unavailable)

		iter := iterator.NewWeightedRoundRobin(pool)
		exp_i := []int{0, 1, 2, 3, 1, 2}
		exp_next := []int{0, 2, 2, 0, 2, 2}

		if err := test.CheckIterNextAvailable(iter, exp_i, exp_next); err != nil {
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
			test.TestNext(iter)
		}
	})

	b.Run("90% available pool", func(b *testing.B) {
		pool := test.CreateRandomTestPool(10000, 1000)
		iter := iterator.NewWeightedRoundRobin(pool)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			test.TestNext(iter)
		}
	})

	b.Run("50% available pool", func(b *testing.B) {
		pool := test.CreateRandomTestPool(10000, 5000)
		iter := iterator.NewWeightedRoundRobin(pool)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			test.TestNext(iter)
		}
	})
}
