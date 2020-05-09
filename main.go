// main runs the testcolor cli, tc
//
// You can install tc by running
//
// 		go get github.com/danieldn/testcolor/tc
//
// Alternatively, clone this repo and install with
//
//		cd <path-to-tc>
//		go install
//
package main

import (
	"flag"
	"fmt"
	"os"
)

// version of testcolor to show in the cli usage message. We assign the value at
// build time using a linker flag like go build -ldflags="-X 'main.version=0.1.0'"
var version = "0.1.0"

var usageMessage = `testcolor v%s

tc pretty prints your 'go test' output

Usage:
	go test -v ./... | tc [flags]

Optional flags:
`

// main starts the testcolor cli
func main() {
	nofmt := flag.Bool("nofmt", false, "Disables formating (default false)")
	nocolor := flag.Bool("nocolor", false, "Disables color (default false)")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), usageMessage, version)
		flag.PrintDefaults()
	}

	flag.Parse()

	// Create a parser that knows how to parse text based on user provided
	// options
	p := newParser(&options{
		nofmt:   *nofmt,
		nocolor: *nocolor,
	})

	// Run testcolor with input from os.Stdin and output to os.Stdout
	runTestColor(p, os.Stdin, os.Stdout)
}
