package main

import "github.com/spf13/pflag"

var (
	modulePath = pflag.StringP("module", "m", "", "Path to bqt module")
)
