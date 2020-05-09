package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// Usages Tests
//
// Basic idea: We want to test running tc with its different options like
//  tc
//  tc -nocolor
//  tc -nofmt
//  tc -nocolor -nofmt
//
// To do this, we mock os.Stdin and os.Stdout by creating an io.Reader from
// fake_test_output.txt and io.Writer from actual_output.txt. Run testcolor()
// with this reader and capture the output to writer. Do a file diff between
// actual_output.txt and expected_output.txt and assert no differences
func TestUsage(t *testing.T) {
	var (
		fakeOutputFile     = "fake.txt"
		actualOutputFile   = "actual.txt"
		expectedOutputFile = "expected.txt"
	)

	cases := []struct {
		name string
		dir  string
		opt  *options
	}{
		{
			name: "default",
			dir:  "testdata/default",
			opt:  &options{nofmt: false, nocolor: false},
		},
		{
			name: "nocolor",
			dir:  "testdata/nocolor",
			opt:  &options{nofmt: false, nocolor: true},
		},
		{
			name: "nofmt",
			dir:  "testdata/nofmt",
			opt:  &options{nofmt: true, nocolor: false},
		},
		{
			name: "nofmtnocolor",
			dir:  "testdata/nofmtnocolor",
			opt:  &options{nofmt: true, nocolor: true},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// create our parser
			p := newParser(tc.opt)

			// mock os.Stdin
			mockStdin, err := os.Open(filepath.Join(tc.dir, fakeOutputFile))
			if err != nil {
				t.Error(err)
			}
			defer close(t, mockStdin)

			// mock os.Stdout
			mockStdout, err := os.Create(filepath.Join(tc.dir, actualOutputFile))
			if err != nil {
				t.Error(err)
			}

			// Run output will be written to actual_output.txt and immediately
			// close file once done
			runTestColor(p, mockStdin, mockStdout)
			close(t, mockStdout)

			// Assert no difference
			if filesEqual(t,
				filepath.Join(tc.dir, actualOutputFile),
				filepath.Join(tc.dir, expectedOutputFile),
			) != 0 {
				t.Error("Actual output file does not match expected")
			}
		})
	}

	// Cleanup function to delete files that get created in subtests, i.e. /testdata/*/actual.txt
	t.Cleanup(func() {
		for _, tc := range cases {
			file := filepath.Join(tc.dir, actualOutputFile)
			t.Log("Removing created file", file)
			os.Remove(file)
		}
	})
}

func filesEqual(t *testing.T, pathA, pathB string) int {
	a, err := os.Open(pathA)
	if err != nil {
		t.Error(err)
	}

	b, err := os.Open(pathB)
	if err != nil {
		t.Error(err)
	}

	bytesA, _ := ioutil.ReadAll(a)
	bytesB, _ := ioutil.ReadAll(b)

	return bytes.Compare(bytesA, bytesB)
}

func close(t *testing.T, c io.Closer) {
	err := c.Close()
	if err != nil {
		t.Error(err)
	}
}

func TestSearch(t *testing.T) {
	cases := []struct {
		data     []byte
		pattern  []byte
		n        int
		expected []int
	}{
		{
			data:     []byte(""),
			pattern:  []byte("bb"),
			n:        1,
			expected: nil,
		},
		{
			data:     []byte("aaaaa"),
			pattern:  []byte("c"),
			n:        1,
			expected: nil,
		},
		{
			data:     []byte("aaabbcc"),
			pattern:  []byte("bb"),
			n:        1,
			expected: []int{3},
		},
		{
			data:     []byte("aaabbb"),
			pattern:  []byte("bb"),
			n:        2,
			expected: []int{3, 4},
		},
		{
			data:     []byte("aaabbb"),
			pattern:  []byte("bb"),
			n:        1,
			expected: []int{3},
		},
		{
			data:     []byte("AABAACAADAABAABA"),
			pattern:  []byte("AABA"),
			n:        -1,
			expected: []int{0, 9, 12},
		},
		{
			data:     []byte("    TestMockSubTest/two: mock_test.go:47: Here's another occurrence of _test.go"),
			pattern:  []byte("_test.go"),
			n:        1,
			expected: []int{29},
		},
		{
			data:     []byte("_test.go:47: bunch : of : colons:::"),
			pattern:  []byte(":"),
			n:        2,
			expected: []int{8, 11},
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("Test %v", i), func(t *testing.T) {
			t.Logf("Searching text \"%s\" for pattern \"%s\"\n", tc.data, tc.pattern)
			res := search(tc.data, tc.pattern, tc.n)

			// We must deep equal to compare slices
			if !reflect.DeepEqual(res, tc.expected) {
				t.Errorf("Got %v Expected %v\n", res, tc.expected)
			}
		})
	}
}

// Test search function with larger input data. Read in testdata/moby.txt and
// check it is the number of bytes we expect
//
//  $ stat moby.txt
//  size    1276201
//  ...
//
// Then run our pattern search function and check it is the number we expect
//
//  $ cat moby.txt | grep  -o "good" | wc -l
//  201
func TestSearchWithLargerInput(t *testing.T) {
	testFilePath := "testdata/moby.txt"
	expectedSize := 1276201
	expectedMatches := 201

	data, err := ioutil.ReadFile(testFilePath)
	if err != nil {
		t.Fatal(err)
	}

	if len(data) != expectedSize {
		t.Errorf("Got %v Expected %v\n", len(data), expectedSize)
	}

	res := search(data, []byte("good"), -1)
	if len(res) != expectedMatches {
		t.Errorf("Got %v Expected %v\n", len(res), expectedMatches)
	}
}

func BenchmarkSearch(b *testing.B) {
	testFilePath := "testdata/moby.txt"
	pattern := []byte("good")

	data, err := ioutil.ReadFile(testFilePath)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		search(data, pattern, -1)
	}
}
