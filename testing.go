package main

import (
	"path/filepath"
	"testing"
)

const TestDir = "testground"

// path of file to create and its permission
type PairPathPerm struct {
	Path string
	Perm int
}

// shortening of standard func
func PathJoin(strA string, strB string) string {
	return filepath.Join(strA, strB)
}

// compare two exe slices and return true if they are the same
func EqualExeList(t *testing.T, listA []exe, listB []exe) (equal bool) {
	t.Helper()

	if len(listA) != len(listB) {
		return false
	}

	for i := range listA {
		if listA[i].Name != listB[i].Name {
			return false
		} else if listA[i].Path != listB[i].Path {
			return false
		}
	}
	return true
}

// compare two lists of errors
func EqualErrorList(t *testing.T, listA []error, listB []error) (equal bool) {
	t.Helper()

	if len(listA) != len(listB) {
		return false
	}

	for i := range listA {
		if listA[i] != nil && listB[i] != nil {
			// as none of them is nil compare error string
			if listA[i].Error() != listB[i].Error() {
				return false
			}
		} else {
			// otherwise just compare type since one of them
			if listA[i] != listB[i] {
				return false
			}
		}
	}

	return true
}

// print out an error with expected and gotten values
func ErrorExpGot(t *testing.T, expected, gotten interface{}, err bool) {
	t.Helper()

	if err {
		t.Error("\nExpectedErr:", expected, "\nGottenErr:", gotten)
	} else {
		t.Error("\nExpected:", expected, "\nGotten:", gotten)
	}

}
