package main

import "testing"

func TestLaunch(t *testing.T) {
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
			ParamRunner:    Runner{},
		},
		{
			Description: "run option with no arguments",
			Expected:    1,

			ParamArguments: []string{"-r"},
			ParamRunner:    Runner{},
		},
		{
			Description: "run option with a letter",
			Expected:    2,

			ParamArguments: []string{"-r", "a"},
			ParamRunner:    Runner{},
		},
		{
			Description: "scan option with wrong dir",
			Expected:    3,

			ParamArguments: []string{"-s", "/ii"},
			ParamRunner:    Runner{},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Description, func(t *testing.T) {
			var gotten = Launch(testCase.ParamRunner, testCase.ParamArguments)

			if testCase.Expected != gotten {
				ErrorExpGot(t, testCase.Expected, gotten, false)
			}
		})
	}
}
