package main

import "github.com/spf13/pflag"

var (
	modulePath = pflag.StringP("module", "m", "", "Path to bqt module")
	srcPath    = pflag.StringP("source", "s", "", "Path to source file")
	bgrndPath  = pflag.StringP("background", "b", "static/1.png", "Path to background file. JPG or PNG")
)
