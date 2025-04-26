package utils

import (
	"fmt"
	"net/http"
	"os"
)

// MustError is a helper function that panics if an error occurs
func MustError[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}
	return val
}

// MustEnv reads an enviroment variable and panics if it's missing
func MustEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Sprint("Enviroment variable %s is required", key))
	}
	return val
}

// MustQueryParam reads a query parameter from the URL and panics if it's missing
func MustQueryParam(r *http.Request, key string) string {
	val := r.URL.Query().Get(key)
	if val == "" {
		panic(fmt.Sprint("Query parameter %s is required", key))
	}
	return val
}
