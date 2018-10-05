package main

import (
	"github.com/jessevdk/go-flags"
)

type Options struct {
	Config string `short:"c" long:"config" description:"Config file" required:"true"`
}

func initOptions() ([]string, *Options) {
	opts := &Options{}
	args, err := flags.Parse(opts)
	if err != nil {
		return nil, nil
	}
	return args, opts
}
