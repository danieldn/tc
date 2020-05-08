// Package mockpkg is a mock package for testing testcolor. It resembles a generic
// public package you might find in the wild
package mockpkg

func PlusOne(x int) int {
	return x + 1
}

func BuggyPlusOne(x int) int {
	return x + 2
}
