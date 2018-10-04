package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func fileChecksum(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func archive(dir string, output string) error {
	cmd := exec.Command("tar", "-czf", output, ".")
	cmd.Stderr = os.Stderr
	cmd.Dir = dir
	return cmd.Run()
}

func extract(filename string, path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}
	return exec.Command("tar", "-xzf", filename, "-C", path).Run()
}
