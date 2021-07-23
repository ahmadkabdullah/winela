package main

import "testing"

func TestLaunch(t *testing.T) {
	var testTable = []struct {
		Description  string
		ExpectedErrs []error

		ParamSearchDir string
		ParamWriteFile string
		ParamDirs      []PairPathPerm
		ParamFiles     []PairPathPerm
	}{
		{
			Description: "scan a dir and export result to a file then import back from exported file",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Description, func(t *testing.T) {

		})
	}
}
