package tools

import (
	"os"
	"strconv"
)

func GetEnvBool(name string) bool {
	value, find := os.LookupEnv(name)
	if !find {
		return false
	}
	result, err := strconv.ParseBool(value)
	if err != nil {
		return false
	}
	return result
}

func GetEnvDefault(name, defaultValue string) string {
	value, find := os.LookupEnv(name)
	if !find {
		return defaultValue
	}
	return value
}
