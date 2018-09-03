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

	tbs, err := parseFile(*srcPath)
	assertError(err)

	fmt.Println(tbs)

	for _, tb := range tbs {
		text, err := module.GetScripture(tb)
		if err != nil {
			fmt.Printf("[WARNING] %s\n", err)
		}

		fmt.Println(tb.String())
		fmt.Println(text)
	}
}
