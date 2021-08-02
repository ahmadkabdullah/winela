package main

import (
	"os"
	"testing"
)

func TestLaunch(t *testing.T) {
	// make testing directory
	os.MkdirAll(TestDir, 0755)
	defer os.RemoveAll(TestDir)

	var testTable = []struct {
		Description string
		Expected    int

		ParamArguments []string
		ParamRunner    Runner
	}{
		{
			Description: "list option",
			Expected:    0,

			ParamArguments: []string{"-l"},
			ParamRunner: Runner{
				ListFile: inTestDir("wineladb"),
			},
		},
		{
			Description: "run option with no arguments",
			Expected:    1,

			ParamArguments: []string{"-r"},
			ParamRunner: Runner{
				ListFile: inTestDir("wineladb"),
			},
		},
		{
			Description: "run option with a letter",
			Expected:    2,

			ParamArguments: []string{"-r", "a"},
			ParamRunner: Runner{
				ListFile: inTestDir("wineladb"),
			},
		},
		{
			Description: "scan option with wrong dir",
			Expected:    3,

			ParamArguments: []string{"-s", "/ii"},
			ParamRunner: Runner{
				ListFile: inTestDir("wineladb"),
			},
		},
		{
			Description: "scan option no directory given and no default",
			Expected:    1,

			ParamArguments: []string{"-s"},
			ParamRunner: Runner{
				ListFile: inTestDir("wineladb"),
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Description, func(t *testing.T) {
			var gotten = launch(testCase.ParamRunner, testCase.ParamArguments)

			if testCase.Expected != gotten {
				errorExpGot(t, testCase.Expected, gotten, false)
			}
		})
	}
}
