package main

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"
)

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

func BenchmarkSearch(b *testing.B) {
	data := []byte("aaabbcc")
	pattern := []byte("bb")

	for n := 0; n < b.N; n++ {
		search(data, pattern, 2)
	}
}

func TestSearchWithLargerInput(t *testing.T) {
	// Read in testdata/moby.txt and check it is the number of bytes we expect
	//
	//  $ stat moby.txt
	//  size    1276201
	//  ...
	//
	// Then run our pattern search function and check it is the number we expect
	//
	//  $ cat moby.txt | grep  -o "good" | wc -l
	//  201
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
