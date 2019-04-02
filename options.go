package main

import (
	"github.com/jessevdk/go-flags"
)

type Options struct {
	Config string `short:"c" long:"config" description:"Config file" default:".dep-cache.json"`
}

func initOptions() ([]string, *Options, error) {
	opts := &Options{}
	args, err := flags.Parse(opts)
	if err != nil {
		return nil, nil, err
	}
	return args, opts, nil
}
