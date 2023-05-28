package server_test

import (
	"fmt"
	server "load-balancer/src/server"
	test "load-balancer/test"
	"testing"
)

func singleCheck(pool *server.ServerPool, order []int, i int, expected int) error {
	actual, _ := pool.GetNextAvailable(order, i)
	if actual != expected {
		return fmt.Errorf(test.ErrorIdx(i, actual, expected))
	}
	return nil
}

func checkPool(pool *server.ServerPool , order []int, expected []int) error {
	for i := 0; i < len(expected); i++ {
		if err := singleCheck(pool, order, i, expected[i]); err != nil {
			return err
		}
	}
	return nil
}

func TestGetNextAvailable(t *testing.T) {
	test.Setup()

	t.Run("Empty pool", func(t *testing.T) {
		pool := test.CreateTestPool(0)
		if err := singleCheck(pool, pool.DefaultOrder, 0, -1); err != nil {
			t.Error(err)
		}
	})

	t.Run("Out-of-bounds index", func(t *testing.T) {
		pool := test.CreateTestPool(3)
		if err := singleCheck(pool, pool.DefaultOrder, 10, 1); err != nil {
			t.Error(err)
		}
	})

	t.Run("Available pool", func(t *testing.T) {
		expected := []int{0, 1, 2, 3, 4}
		pool := test.CreateTestAlivePool(len(expected), []int{})
		if err := checkPool(pool, pool.DefaultOrder, expected); err != nil {
			t.Error(err)
		}
	})

	t.Run("Unavailable pool", func(t *testing.T) {
		unavailable := []int{0, 1, 2, 3, 4}
		expected := []int{-1, -1, -1, -1, -1}
		pool := test.CreateTestAlivePool(len(expected), unavailable)

		if err := checkPool(pool, pool.DefaultOrder, expected); err != nil {
			t.Error(err)
		}
	})

	t.Run("Alternating available pool", func(t *testing.T) {
		unavailable := []int{1, 3, 5, 7, 9}
		expected := []int{0, 2, 2, 4, 4, 6, 6, 8, 8, 0}
		pool := test.CreateTestAlivePool(len(expected), unavailable)

		if err := checkPool(pool, pool.DefaultOrder, expected); err != nil {
			t.Error(err)
		}
	})

	t.Run("Available pool, reverse order", func(t *testing.T) {
		order := []int{9, 8, 7, 6, 5, 4, 3, 2, 1, 0}
		expected := []int{9, 8, 7, 6, 5, 4, 3, 2, 1, 0}
		pool := test.CreateTestAlivePool(len(expected), []int{})

		if err := checkPool(pool, order, expected); err != nil {
			t.Error(err)
		}
	})

	t.Run("Semi-available pool, random order", func(t *testing.T) {
		unavailable := []int{7, 8, 9, 3, 4}
		order := []int{7, 8, 9, 3, 4, 5, 6, 0, 1, 2}
		expected := []int{5, 5, 5, 5, 5, 5, 6, 0, 1, 2}
		pool := test.CreateTestAlivePool(len(expected), unavailable)

		if err := checkPool(pool, order, expected); err != nil {
			t.Error(err)
		}
	})

	t.Run("Random available pool, random order", func(t *testing.T) {
		unavailable := []int{2, 8, 4, 9, 6, 1}
		order := []int{2, 0, 8, 4, 7, 3, 9, 5, 6, 1}
		expected := []int{0, 0, 7, 7, 7, 3, 5, 5, 0, 0}
		pool := test.CreateTestAlivePool(len(expected), unavailable)

		if err := checkPool(pool, order, expected); err != nil {
			t.Error(err)
		}
	})
}