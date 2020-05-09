package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
)

var (
	// Target patterns for text searching and editing
	//
	// We define these as byte slices since we use byte slices throughout this
	// program for efficiency (e.g. writing to a buffer instead of concat
	// strings). Alternatively, we could have found ways to use libraries such
	// as src/strings/builder.go but choose to implement similar functions for
	// our specific task
	dashDashDashPass   = []byte("--- PASS")
	dashDashDashFail   = []byte("--- FAIL")
	dashDashDashSkip   = []byte("--- SKIP")
	equalEqualEqualRun = []byte("=== RUN")
	underscoreTest     = []byte("_test.go")
	ok                 = []byte("ok")
	pass               = []byte("PASS")
	fail               = []byte("FAIL")
	colon              = []byte(":")
	question           = []byte("?")
	fourSpaces         = []byte("    ")
	tenSpaces          = []byte("          ")
	// ANSI color escape codes
	boldRed    = []byte("\033[1;31m")
	boldGreen  = []byte("\033[1;32m")
	boldYellow = []byte("\033[1;33m")
	boldCyan   = []byte("\033[1;36m")
	red        = []byte("\033[31m")
	green      = []byte("\033[32m")
	yellow     = []byte("\033[33m")
	cyan       = []byte("\033[36m")
	grey       = []byte("\033[90m")
	reset      = []byte("\033[0m")
)

type options struct {
	nofmt   bool
	nocolor bool
}

type parser interface {
	parseBuffer([]byte) []byte
}

// newParser returns an object that implements parser based on given options
func newParser(opt *options) parser {
	if opt.nofmt || opt.nocolor {
		return &customP{opt}
	}

	return &normalP{}
}

// normalP implements a normal parser
type normalP struct{}

// customP implements a custom parser
type customP struct {
	opt *options
}

func runTestColor(p parser, rd io.Reader, wr io.Writer) {
	// Allocate main buffer
	buffer := bufio.NewReader(rd)

	// We loop to read and parse input from stdin, line by line until EOF
	for {
		// ReadSlice returns of slice of buffer (does not allocate new memory)
		// and WILL include '\n'
		buf, err := buffer.ReadSlice('\n')
		if err != nil {
			if err == io.EOF {
				return
			}

			log.Fatalf("testcolor: Unhandled error reading from stdin! %v", err)
		}

		// we ditch '\n' in order to annotate properly
		editedBuf := p.parseBuffer(buf[:len(buf)-1])

		fmt.Fprintf(wr, "%s\n", editedBuf)
	}
}

func (p normalP) parseBuffer(buf []byte) []byte {
	var annotatedSlice []byte

	trimmedBuf := bytes.TrimLeft(buf, " ")

	switch {
	case bytes.HasPrefix(trimmedBuf, equalEqualEqualRun):
		// === RUN
		annotatedSlice = buf
	case bytes.HasPrefix(trimmedBuf, dashDashDashPass):
		// --- PASS
		annotatedSlice = annotateWithColorAndIntent(buf, boldGreen)
	case bytes.HasPrefix(trimmedBuf, pass):
		// PASS
		annotatedSlice = annotateWithColor(buf, boldGreen)
	case bytes.HasPrefix(trimmedBuf, ok):
		// ok
		annotatedSlice = annotateWithColor(buf, green)
	case bytes.HasPrefix(trimmedBuf, dashDashDashFail):
		// --- FAIL
		annotatedSlice = annotateWithColorAndIntent(buf, boldRed)
	case bytes.HasPrefix(trimmedBuf, fail):
		// FAIL
		annotatedSlice = annotateWithColor(buf, boldRed)
	case bytes.HasPrefix(trimmedBuf, dashDashDashSkip):
		// --- SKIP
		annotatedSlice = annotateWithColorAndIntent(buf, boldYellow)
	case bytes.HasPrefix(trimmedBuf, question):
		// ?
		annotatedSlice = annotateWithColor(buf, cyan)
	case bytes.Contains(trimmedBuf, underscoreTest):
		// For lines containing '_test.go' like
		//
		//		TestMockSucceed: mock_test.go:6: Checking if mocking succeed works
		//
		// we transform to
		//
		//		mock_test.go:6:
		//			Checking if mocking succeed works
		//
		var edited bytes.Buffer
		//    TestMockSucceed: mock_test.go:6: Checking if mocking succeed works
		//                     ^
		//					   cut everything before this
		//
		colons := search(trimmedBuf, colon, 1)
		if len(colons) < 1 {
			panic("testcolor: unexpected nil from search function")
		}

		edited.Write(tenSpaces)
		edited.Write(annotateWithColor(trimmedBuf[colons[0]+2:], grey))

		annotatedSlice = edited.Bytes()
	default:
		annotatedSlice = buf
	}

	return annotatedSlice
}

func (p customP) parseBuffer(buf []byte) []byte {
	var annotatedSlice []byte

	trimmedBuf := bytes.TrimLeft(buf, " ")

	// We have to handle 3 option possibilities
	// 		tc -nofmt
	// 		tc -nocolor
	// 		tc -nofmt -nocolor
	switch {
	case bytes.HasPrefix(trimmedBuf, equalEqualEqualRun):
		// === RUN
		annotatedSlice = buf
	case bytes.HasPrefix(trimmedBuf, dashDashDashPass):
		// --- PASS
		if p.opt.nofmt && !p.opt.nocolor {
			annotatedSlice = annotateWithColor(buf, boldGreen)
			break
		}

		if !p.opt.nofmt && p.opt.nocolor {
			annotatedSlice = annotateWithIndent(buf)
			break
		}

		annotatedSlice = buf
	case bytes.HasPrefix(trimmedBuf, pass):
		// PASS
		if p.opt.nofmt && !p.opt.nocolor {
			annotatedSlice = annotateWithColor(buf, boldGreen)
			break
		}

		// No formatting is done here so if nocolor is given,
		// just return
		annotatedSlice = buf
	case bytes.HasPrefix(trimmedBuf, ok):
		// ok
		if p.opt.nofmt && !p.opt.nocolor {
			annotatedSlice = annotateWithColor(buf, boldGreen)
			break
		}

		// No formatting is done here so if nocolor is given,
		// just return
		annotatedSlice = buf
	case bytes.HasPrefix(trimmedBuf, dashDashDashFail):
		// --- FAIL
		if p.opt.nofmt && !p.opt.nocolor {
			annotatedSlice = annotateWithColor(buf, boldRed)
			break
		}

		if !p.opt.nofmt && p.opt.nocolor {
			annotatedSlice = annotateWithIndent(buf)
			break
		}

		annotatedSlice = buf
	case bytes.HasPrefix(trimmedBuf, fail):
		// FAIL
		if p.opt.nofmt && !p.opt.nocolor {
			annotatedSlice = annotateWithColor(buf, boldRed)
			break
		}

		// No formatting is done here so if nocolor is given,
		// just return
		annotatedSlice = buf
	case bytes.HasPrefix(trimmedBuf, dashDashDashSkip):
		// --- SKIP
		if p.opt.nofmt && !p.opt.nocolor {
			annotatedSlice = annotateWithColor(buf, boldRed)
			break
		}

		if !p.opt.nofmt && p.opt.nocolor {
			annotatedSlice = annotateWithIndent(buf)
			break
		}

		annotatedSlice = buf
	case bytes.HasPrefix(trimmedBuf, question):
		// ?
		if p.opt.nofmt && !p.opt.nocolor {
			annotatedSlice = annotateWithColor(buf, cyan)
			break
		}

		// No formatting is done here so if nocolor is given,
		// just return
		annotatedSlice = buf
	case bytes.Contains(trimmedBuf, underscoreTest):
		var (
			edited bytes.Buffer
			colons []int
		)

		if p.opt.nofmt && !p.opt.nocolor {
			edited.Write(fourSpaces)
			edited.Write(annotateWithColor(trimmedBuf, grey))
			annotatedSlice = edited.Bytes()

			break
		}

		if !p.opt.nofmt && p.opt.nocolor {
			colons = search(trimmedBuf, colon, 1)
			if len(colons) < 1 {
				panic("testcolor: unexpected nil from search function")
			}

			edited.Write(annotateWithIndent(trimmedBuf[colons[0]+2:]))
			annotatedSlice = edited.Bytes()

			break
		}

		annotatedSlice = buf
	default:
		annotatedSlice = buf
	}

	return annotatedSlice
}

// annotate builds a new slice of bytes with the ANSI escape code corresponding
// to color followed by text and ending with the reset escape code
func annotateWithColor(text []byte, color []byte) []byte {
	var edited bytes.Buffer

	// We don't check return val since err returned from Write is always nil
	// will panic if fails
	edited.Write(color)
	edited.Write(text)
	edited.Write(reset)

	return edited.Bytes()
}

func annotateWithIndent(text []byte) []byte {
	var edited bytes.Buffer

	edited.Write(tenSpaces)
	edited.Write(text)

	return edited.Bytes()
}

func annotateWithColorAndIntent(text, color []byte) []byte {
	var edited bytes.Buffer

	edited.Write(tenSpaces)
	edited.Write(color)
	edited.Write(text)
	edited.Write(reset)

	return edited.Bytes()
}

// search performs a naive pattern search on s for pattern p. It returns a slice
// of indices where p begins in s, up to n occurrences (or all if n<0). It
// returns nil if no occurrences are found.
//
// TODO(Daniel): The naive implementation seems good enough for our use case.
// But we could be awesome and shoot for linear time later on with something
// like KMP
func search(s, p []byte, n int) []int {
	var (
		seen   int
		result []int
		found  bool
	)

	// Naive pattern search
	for i := 0; i < len(s); i++ {
		if seen == n {
			break
		}

		// Check that we're still inbounds
		if i+len(p) > len(s) {
			break
		}

		// Lets assume s[i] is start of pattern. We check if s[i] to s[i+len(p)]
		// matches p
		found = true

		for j := 0; j < len(p); j++ {
			if s[i+j] != p[j] {
				found = false
				break
			}
		}

		if found {
			result = append(result, i)
			seen++
		}
	}

	return result
}
