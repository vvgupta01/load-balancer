package utils

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func DisableLogOutput() {
	log.SetFlags(0)
	log.SetOutput(ioutil.Discard)
}

func ToInt(i_str string) int {
	val, err := strconv.Atoi(i_str)
	if err != nil {
		log.Fatal(err)
	}
	return val
}

func GetIntEnv(name string) int {
	return ToInt(os.Getenv(name))
}

func GetTimeEnv(name string) time.Duration {
	str := os.Getenv(name)
	if strings.HasSuffix(str, "ms") {
		return time.Millisecond * time.Duration(ToInt(str[:len(str)-2]))
	} else if strings.HasSuffix(str, "s") {
		return time.Second * time.Duration(ToInt(str[:len(str)-1]))
	}
	log.Fatal("Invalid time suffix: must be 'ms' or 's'")
	return time.Duration(0)
}
