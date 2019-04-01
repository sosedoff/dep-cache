package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	s3session *session.Session
	s3service *s3.S3
	s3bucket  string
)

func setupS3(config *Config) {
	sess, err := session.NewSession()
	if err != nil {
		fatal(err)
	}

	// Load region from metadata
	meta := ec2metadata.New(sess)
	if region, err := meta.Region(); err == nil {
		// Only overwrite region if it's not set
		if config.S3.Region == "" {
			config.S3.Region = region
			sess.Config.Region = &region
		}
	}

	if config.S3.Region == "" {
		fatal("s3 region is reqiored")
	}
	if config.S3.Bucket == "" {
		fatal("s3 bucket name is required")
	}

	if config.S3.Key != "" && config.S3.Secret != "" {
		sess.Config.Credentials = credentials.NewStaticCredentials(
			config.S3.Key,
			config.S3.Secret,
			"",
		)
	}

	s3session = sess
	s3service = s3.New(s3session)
	s3bucket = config.S3.Bucket
}

func perform(cache *Cache, command string) {
	if err := cache.prepare(); err != nil {
		log.Println("error:", err)
		return
	}

	switch command {
	case "upload":
		if err := upload(cache); err != nil {
			fmt.Println("error:", err)
		}
	case "download":
		if err := download(cache); err != nil {
			fmt.Println("error:", err)
		}
	}
}

func main() {
	args, opts := initOptions()
	if opts == nil {
		return
	}

	config, err := readConfig(opts.Config)
	if err != nil {
		fatal(err)
	}

	setupS3(config)

	if len(args) < 1 {
		fatal("command required")
	}

	command := args[0]
	if !(command == "upload" || command == "download") {
		fatal("invalid command:" + command)
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(config.Caches))

	for i := range config.Caches {
		go func(c *Cache) {
			defer wg.Done()
			perform(c, command)
		}(&config.Caches[i])
	}

	wg.Wait()
}

func upload(cache *Cache) error {
	debug("checking %s", cache.Key)
	exists, err := s3exists(cache.Key)
	if err != nil {
		return err
	}
	if exists {
		debug("cache %s exists, skipping upload", cache.Key)
		return nil
	}

	debug("preparing %s", cache.Key)
	archivePath := filepath.Join("/tmp", cache.Key)
	defer os.Remove(archivePath)

	if err := archive(cache.Path, archivePath); err != nil {
		return err
	}

	debug("uploading %s", cache.Key)
	return s3upload(cache.Key, archivePath)
}

func download(cache *Cache) error {
	archivePath := filepath.Join("/tmp", cache.Key)
	defer os.Remove(archivePath)

	debug("checking %s", cache.Key)
	exists, err := s3exists(cache.Key)
	if err != nil {
		return err
	}
	if !exists {
		debug("cache %s not found", cache.Key)
		return nil
	}

	debug("downloading %s", cache.Key)
	if err := s3download(cache.Key, archivePath); err != nil {
		return err
	}

	debug("extracting %s", cache.Key)
	return extract(archivePath, cache.Path)
}
