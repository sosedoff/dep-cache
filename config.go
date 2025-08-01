package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strings"
)

var (
	reEnvVar = regexp.MustCompile(`(?i)\$([\w\d\_]+)`)
)

const (
	downloadPolicyDefault      = "default"
	downloadPolicySkipNotEmpty = "skip_not_empty"
)

type Config struct {
	S3     S3Config `json:"s3"`
	Caches []Cache  `json:"cache"`
}

type S3Config struct {
	Key            string `json:"key"`
	Secret         string `json:"secret"`
	Region         string `json:"region"`
	Bucket         string `json:"bucket"`
	Endpoint       string `json:"endpoint"`
	ForcePathStyle bool   `json:"force_path_style"`
}

type Cache struct {
	Manifest       string `json:"manifest"`
	Path           string `json:"path"`
	Prefix         string `json:"prefix"`
	DownloadPolicy string `json:"download_policy"`
	Key            string `json:"-"`
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
	input, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	newInput := replaceEnvVars(string(input))

	cfg := &Config{}
	if err := json.Unmarshal([]byte(newInput), cfg); err != nil {
		return nil, err
	}

	if cfg.S3.Key != "" {
		if err := loadFromFile(cfg.S3.Key, &cfg.S3.Key); err != nil {
			return nil, err
		}
	}
	if cfg.S3.Secret != "" {
		if err := loadFromFile(cfg.S3.Secret, &cfg.S3.Secret); err != nil {
			return nil, err
		}
	}

	return cfg, nil
}

func (c *Cache) validate() error {
	if c.DownloadPolicy == "" {
		c.DownloadPolicy = downloadPolicyDefault
	}

	switch c.DownloadPolicy {
	case downloadPolicyDefault, downloadPolicySkipNotEmpty:
		return nil
	default:
		return fmt.Errorf("invalid download policy: %v", c.DownloadPolicy)
	}
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
