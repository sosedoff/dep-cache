package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"runtime"
	"strings"
)

var (
	reEnvVar = regexp.MustCompile(`(?i)\$([\w\d\_]+)`)
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

func replaceEnvVars(input string) string {
	matches := reEnvVar.FindAllStringSubmatch(input, -1)
	if len(matches) > 0 {
		for _, m := range matches {
			input = strings.Replace(input, m[0], os.Getenv(m[1]), -1)
		}
	}
	return input
}

func readConfig(path string) (*Config, error) {
	input, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	newInput := replaceEnvVars(string(input))

	cfg := &Config{}
	if err := json.Unmarshal([]byte(newInput), cfg); err != nil {
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
