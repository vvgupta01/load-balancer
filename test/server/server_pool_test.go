package server_test

import (
	"fmt"
	server "loadbalancer/src/server"
	test "loadbalancer/test"
	"testing"
)

func singleCheck(pool server.ServerPool, i int, expected int) error {
	actual, _ := pool.GetNextAvailable(i)
	if actual != expected {
		return fmt.Errorf(test.ErrorIdx(i, actual, expected))
	}
	return nil
}

func checkPool(pool server.ServerPool, expected []int) error {
	for i := 0; i < len(expected); i++ {
		if err := singleCheck(pool, i, expected[i]); err != nil {
			return err
		}
	}
	return nil
}

func TestGetNextAvailable(t *testing.T) {
	test.Setup()

	t.Run("Empty pool", func(t *testing.T) {
		pool := test.CreateDefaultTestPool(0)
		if err := singleCheck(pool, 0, -1); err != nil {
			t.Error(err)
		}
	})

	t.Run("Negative index", func(t *testing.T) {
		pool := test.CreateDefaultTestPool(3)
		if err := singleCheck(pool, -1, -1); err != nil {
			t.Error(err)
		}
	})

	t.Run("Out-of-bounds index (wrap)", func(t *testing.T) {
		pool := test.CreateDefaultTestPool(3)
		if err := singleCheck(pool, 10, 1); err != nil {
			t.Error(err)
		}
	})

	t.Run("Available pool", func(t *testing.T) {
		pool := test.CreateDefaultTestPool(5)
		expected := []int{0, 1, 2, 3, 4}
		if err := checkPool(pool, expected); err != nil {
			t.Error(err)
		}
	})

	t.Run("Unavailable pool", func(t *testing.T) {
		unavailable := []int{0, 1, 2, 3, 4}
		pool := test.CreateTestPool(5, nil, nil, nil, unavailable)
		expected := []int{-1, -1, -1, -1, -1}

		if err := checkPool(pool, expected); err != nil {
			t.Error(err)
		}
	})

	t.Run("Random available pool", func(t *testing.T) {
		unavailable := []int{0, 2, 3, 6, 8, 9}
		pool := test.CreateTestPool(10, nil, nil, nil, unavailable)
		expected := []int{1, 1, 4, 4, 4, 5, 7, 7, 1, 1}

		if err := checkPool(pool, expected); err != nil {
			t.Error(err)
		}
	})
}
