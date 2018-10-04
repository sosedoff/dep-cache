package main

import (
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	s3session *session.Session
	s3service *s3.S3
	s3bucket  string
)

func setupS3(opts *Options) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(opts.S3Region),
		Credentials: credentials.NewStaticCredentials(opts.S3Key, opts.S3Secret, ""),
	})
	if err != nil {
		fatal(err)
	}
	s3session = sess
	s3service = s3.New(s3session)
	s3bucket = opts.S3Bucket
}

func main() {
	opts := initOptions()
	if opts == nil {
		return
	}

	if len(os.Args) < 2 {
		fatal("command required")
	}

	command := os.Args[1]

	if err := prepare(opts); err != nil {
		fatal(err)
	}

	setupS3(opts)

	switch command {
	case "upload":
		if err := upload(opts); err != nil {
			fatal(err)
		}
	case "download":
		if err := download(opts); err != nil {
			fatal(err)
		}
	default:
		fatal("invalid command")
	}
}

func upload(opts *Options) error {
	debug("checking cache %s", opts.Key)
	exists, err := s3exists(opts.Key)
	if err != nil {
		return err
	}
	if exists {
		debug("cache %s already exists, skipping upload", opts.Key)
		return nil
	}

	debug("preparing cache")
	archivePath := filepath.Join("/tmp", opts.Key)
	defer os.Remove(archivePath)

	if err := archive(opts.Path, archivePath); err != nil {
		return err
	}

	debug("uploading cache")
	return s3upload(opts.Key, archivePath)
}

func download(opts *Options) error {
	archivePath := filepath.Join("/tmp", opts.Key)
	defer os.Remove(archivePath)

	debug("checking for cache %s", opts.Key)
	exists, err := s3exists(opts.Key)
	if err != nil {
		return err
	}
	if !exists {
		debug("cache does not exist")
		return nil
	}

	debug("downloading cache")
	if err := s3download(opts.Key, archivePath); err != nil {
		return err
	}

	debug("extracting cache")
	return extract(archivePath, opts.Path)
}
