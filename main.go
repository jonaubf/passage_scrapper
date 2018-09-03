package main

import (
	"github.com/spf13/pflag"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func assertError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	pflag.Parse()
	module, err := ReadBQTModule(*modulePath)
	assertError(err)

	tbs, err := parseFile(*srcPath)
	assertError(err)

	generatePDF(module, tbs)
}
