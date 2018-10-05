package main

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
)

type Cache struct {
	Manifest string `json:"manifest"`
	Path     string `json:"path"`
	Prefix   string `json:"prefix"`
	Key      string `json:"-"`
}

type Config struct {
	S3 struct {
		Key    string `json:"key"`
		Secret string `json:"secret"`
		Region string `json:"region"`
		Bucket string `json:"bucket"`
	} `json:"s3"`
	Caches []Cache `json:"cache"`
}

func readConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	cfg := &Config{}
	if err := json.NewDecoder(file).Decode(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Cache) prepare() error {
	checksum, err := fileChecksum(c.Manifest)
	if err != nil {
		return err
	}

	c.Key = fmt.Sprintf(
		"%s_%s_%s_%s.tar.gz",
		c.Prefix,
		checksum,
		runtime.GOOS,
		runtime.GOARCH,
	)

	return nil
}
