package main

import (
	"fmt"
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
	fmt.Printf("%+v\n", module)
}
