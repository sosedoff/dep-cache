package main

import (
	"fmt"
	"os"
)

func debug(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

func fatal(message interface{}) {
	fmt.Println("error:", message)
	os.Exit(0)
}

func getEnvVar(name string) string {
	val := os.Getenv(name)
	if val == "" {
		fatal("Please set " + name + " environment variable")
	}
	return val
}
