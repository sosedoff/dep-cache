package main

import (
	"fmt"
	"runtime"

	"github.com/jessevdk/go-flags"
)

type Options struct {
	S3Key        string `long:"s3-key" description:"AWS access key" required:"true"`
	S3Secret     string `long:"s3-secret" description:"AWS secret key" required:"true"`
	S3Region     string `long:"s3-region" description:"AWS region" default:"us-east-1"`
	S3Bucket     string `long:"s3-bucket" description:"S3 bucket name" required:"true"`
	Prefix       string `long:"prefix" description:"S3 files prefix" required:"true"`
	ManifestPath string `long:"manifest" description:"Path to manifest file (Gemfile.lock, etc)" required:"true"`
	Path         string `long:"path" description:"Path to dependencies directory" required:"true"`
	Key          string
}

func prepare(opts *Options) error {
	checksum, err := fileChecksum(opts.ManifestPath)
	if err != nil {
		return err
	}
	opts.Key = fmt.Sprintf("%s_%s_%s.tar.gz", opts.Prefix, checksum, runtime.GOARCH)
	return nil
}

func initOptions() *Options {
	opts := &Options{}

	if _, err := flags.Parse(opts); err != nil {
		return nil
	}
	return opts
}
