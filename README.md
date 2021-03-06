# testcolor
[![Build Status](https://travis-ci.org/danieldn/tc.svg?branch=master)](https://travis-ci.org/danieldn/tc) [![Go Report Card](https://goreportcard.com/badge/github.com/danieldn/tc)](https://goreportcard.com/report/github.com/danieldn/tc) [![codecov](https://codecov.io/gh/danieldn/tc/branch/master/graph/badge.svg)](https://codecov.io/gh/danieldn/tc)

tc pretty prints your 'go test' output 

Before:

![before](before_image.png)

After:

![after](after_image.png)

Install using


```bash
go get github.com/danieldn/tc
```

Then pipe your `go test` output

```bash
go test -v ./... | tc 
```

Disable features you don't want with flags
```bash
Usage:
	go test -v ./... | tc [flags]

Optional flags:
  -nocolor
    	Disables color (default false)
  -nofmt
    	Disables formatting (default false)
```

Add `tc` to your travis.yml
```
before_install:
- go get github.com/danieldn/tc

script:
- go test -v ./... | tc
- go build
```

## How It Works

This tool was made just using the standard Go library. It parses input to a
buffer and performs a pattern search to insert ANSI color escape codes and text
formatting where we want. 

## Performance

This tool relies on Go's [bytes.Buffer](https://golang.org/src/bytes/buffer.go)
to minimize memory allocation when reading and editing text input. 

## Inspired By

* https://github.com/rakyll/gotest
* https://github.com/fatih/color
* https://github.com/logrusorgru/aurora
* https://github.com/mitchellh/colorstring
