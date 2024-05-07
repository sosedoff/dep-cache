package main

import (
	"fmt"
	"os"
	"strings"
)

func debug(format string, args ...interface{}) {
	fmt.Printf("[dep-cache] "+format+"\n", args...)
}

func fatal(message interface{}) {
	fmt.Println("error:", message)
	os.Exit(0)
}

func loadFromFile(path string, dst *string) error {
	if !strings.HasPrefix(path, "file://") {
		return nil
	}

	str, err := os.ReadFile(strings.TrimPrefix(path, "file://"))
	if err != nil {
		return err
	}

	*dst = strings.TrimSpace(string(str))
	return nil
}
