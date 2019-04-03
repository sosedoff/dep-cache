package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
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

	awsmeta := false

	// When AWS keys are not set, load session from environment
	if config.S3.Key == "" && config.S3.Secret == "" {
		debug("api keys not provider, fetching aws metadata...")
		region, err := ec2metadata.New(sess).Region()
		if err == nil {
			config.S3.Region = region
			awsmeta = true
		} else {
			debug("aws metadata fetch error: " + err.Error())
		}
	}

	if config.S3.Region == "" {
		fatal("region is not set")
	} else {
		sess.Config.Region = aws.String(config.S3.Region)
	}

	if config.S3.Bucket == "" {
		fatal("bucket name is not set")
	}

	if !awsmeta {
		if config.S3.Key == "" {
			fatal("access key is not set")
		}
		if config.S3.Secret == "" {
			fatal("secret key is not set")
		}

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
		fmt.Println("error:", err)
		return
	}

	switch command {
	case "status":
		if err := status(cache); err != nil {
			fmt.Println("error:", err)
		}
	case "upload":
		if err := upload(cache); err != nil {
			fmt.Println("error:", err)
		}
	case "download":
		if err := download(cache); err != nil {
			fmt.Println("error:", err)
		}
	case "reset":
		if err := reset(cache); err != nil {
			fmt.Println("error:", err)
		}
	}
}

func main() {
	args, opts, err := initOptions()
	if err != nil {
		fatal(err)
	}
	if opts == nil {
		return
	}

	if len(args) < 1 {
		fatal("command required")
	}

	command := args[0]
	commands := map[string]bool{"upload": true, "download": true, "reset": true, "status": true}

	if _, ok := commands[command]; !ok {
		fatal("invalid command:" + command)
	}

	config, err := readConfig(opts.Config)
	if err != nil {
		fatal(err)
	}

	if len(config.Caches) == 0 {
		fatal("no cache manifests found")
	}

	setupS3(config)

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

func status(cache *Cache) error {
	debug("checking %s", cache.Key)
	exists, err := s3exists(cache.Key)
	if err != nil {
		return err
	}
	if exists {
		debug("cache %s exists", cache.Key)
	} else {
		debug("cache %s does not exist", cache.Key)
	}
	return nil
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

func reset(cache *Cache) error {
	debug("deleteing %s", cache.Key)
	return s3delete(cache.Key)
}
