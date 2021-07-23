package main

import (
	"os"
	"testing"
)

func TestRunFromList(t *testing.T) {
	os.MkdirAll(TestDir, 0755)
	defer os.RemoveAll(TestDir)

	var testTable = []struct {
		Description  string
		ExpectedErr  error
		ParamRunner  Runner
		ParamRunProg int
		ParamFork  bool
	}{
		{
			Description: "fork launch first in list with wine and its params",
			ExpectedErr: nil,
			ParamRunner: Runner{
				Program:     "wine",
				ProgramArgs: "",
				List: []exe{
					{1, "rufus", PathJoin(TestDir, "rufus.exe")},
				},
			},
			ParamRunProg: 1,
			ParamFork: true,
		},
		// {
		// 	Description: "fork launch a number out of the range of list",
		// 	ExpectedErr: fmt.Errorf("exe number %d: not in list", 5),
		// 	ParamRunner: Runner{
		// 		Program:     "wine",
		// 		ProgramArgs: "",
		// 		List: []exe{
		// 			{1, "PS", PathJoin(TestDir, "PS.exe")},
		// 		},
		// 	},
		// 	ParamRunProg: 5,
		// 	ParamFork: true,
		// },
	}

	for _, testCase := range testTable {
		t.Run(testCase.Description, func(t *testing.T) {
			// // copy
			// var readData, _ = ioutil.ReadFile("~/Downloads/rufus.exe")
			// // paste
			// ioutil.WriteFile(testCase.ParamRunner.List[0].Path, readData, os.FileMode(0755))
			// // remove later
			// defer os.Remove(testCase.ParamRunner.List[0].Path)

			var gottenErr = testCase.ParamRunner.RunFromList(testCase.ParamRunProg, testCase.ParamFork)

			if EqualErrorList(t, []error{testCase.ExpectedErr}, []error{gottenErr}) == false {
				ErrorExpGot(t, testCase.ExpectedErr, gottenErr, true)
			}
		})
	}
}

func TestDisplayList(t *testing.T) {
	var testTable = []struct {
		Description string
		Expected    string
		ParamRunner Runner
	}{
		{
			Description: "display a list of two",
			Expected:    "1 sr\n2 lon\n",
			ParamRunner: Runner{
				Program:     "wine",
				ProgramArgs: "",
				List: []exe{
					{1, "sr", PathJoin(TestDir, "sr.exe")},
					{2, "lon", PathJoin(TestDir, "lon.exe")},
				},
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Description, func(t *testing.T) {
			var gotten = testCase.ParamRunner.DisplayList()

			if testCase.Expected != gotten {
				ErrorExpGot(t, testCase.Expected, gotten, false)
			}
		})
	}

}
