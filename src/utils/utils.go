package utils

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
)

func DisableLogOutput() {
	log.SetFlags(0)
	log.SetOutput(ioutil.Discard)
}

func GetIntEnv(name string) int {
	val, err := strconv.Atoi(os.Getenv(name))
	if err != nil {
		log.Fatal(err)
	}
	return val
}

func GetTimeEnv(name string) time.Duration {
	millis := time.Duration(GetIntEnv(name))
	return time.Millisecond * millis
}
