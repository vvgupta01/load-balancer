package utils

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func GetIntEnv(name string) int {
	val, err := strconv.Atoi(os.Getenv(name))
	if err != nil {
		fmt.Println(err)
	}
	return val
}

func GetTimeEnv(name string) time.Duration {
	millis := time.Duration(GetIntEnv(name))
	return time.Millisecond * millis
}
