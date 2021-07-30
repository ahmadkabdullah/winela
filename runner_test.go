package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestRunnerReadConfig(t *testing.T) {
	os.MkdirAll(TestDir, 0755)
	defer os.RemoveAll(TestDir)

	var testTable = []struct {
		Description       string
		Expected          Runner
		ParamRunner       Runner
		ParamConfigString string
		ParamConfigFile   string
		ParamModifyTo     string
	}{
		{
			Description: "read a config with program set us wine staging",
			Expected: Runner{
				Program:     "wine-staging",
				ProgramArgs: "",
				List:        []Exe{},
				ConfigFile:  pathJoin(TestDir, "winelarc"),
			},
			ParamRunner: Runner{
				ConfigFile: pathJoin(TestDir, "winelarc"),
			},
			ParamConfigString: "Program : wine-staging\n" +
				"ProgramArgs : ",
			ParamConfigFile: pathJoin(TestDir, "winelarc"),
			ParamModifyTo:   "",
		},
		{
			Description: "read a config with program set us wine staging",
			Expected: Runner{
				Program:     "wine",
				ProgramArgs: "",
				List:        []Exe{},
				ConfigFile:  pathJoin(TestDir, "winelarc"),
			},
			ParamRunner: Runner{
				ConfigFile: pathJoin(TestDir, "winelarc"),
			},
			ParamConfigString: "Program : wine-staging\n" + "ProgramArgs : ",
			ParamConfigFile:   pathJoin(TestDir, "winelarc"),
			ParamModifyTo:     "Program : wine\n" + "ProgramArgs : ",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Description, func(t *testing.T) {
			// write
			var _ = ioutil.WriteFile(testCase.ParamConfigFile, []byte(testCase.ParamConfigString), os.FileMode(0755))

			// first read
			testCase.ParamRunner.runnerReadConfig()

			// test straight away if there is no value to modify to
			// else do modification and read again then test
			if testCase.ParamModifyTo == "" {
				// test write read
				if fmt.Sprint(testCase.ParamRunner) != fmt.Sprint(testCase.Expected) {
					errorExpGot(t, testCase.Expected, testCase.ParamRunner, false)
				}
			} else {
				// edit
				var _ = ioutil.WriteFile(testCase.ParamConfigFile, []byte(testCase.ParamModifyTo), os.FileMode(0755))
				// second read
				testCase.ParamRunner.runnerReadConfig()
				// test modified read
				if fmt.Sprint(testCase.ParamRunner) != fmt.Sprint(testCase.Expected) {
					errorExpGot(t, testCase.Expected, testCase.ParamRunner, false)
				}
			}
		})
	}
}

func TestRunFromList(t *testing.T) {
	os.MkdirAll(TestDir, 0755)
	defer os.RemoveAll(TestDir)

	var testTable = []struct {
		Description  string
		ExpectedErr  error
		ParamRunner  Runner
		ParamRunProg int
		ParamFork    bool
	}{
		// {
		// 	Description: "fork launch first in list with wine and its params",
		// 	ExpectedErr: nil,
		// 	ParamRunner: Runner{
		// 		Program:     "wine",
		// 		ProgramArgs: "",
		// 		List: []exe{
		// 			{1, "rufus", PathJoin(TestDir, "rufus.exe")},
		// 		},
		// 	},
		// 	ParamRunProg: 1,
		// 	ParamFork: true,
		// },
		{
			Description: "fork launch a number out of the range of list",
			ExpectedErr: fmt.Errorf("exe number %d: not in list", 5),
			ParamRunner: Runner{
				Program:     "wine",
				ProgramArgs: "",
				List: []Exe{
					{"PS", pathJoin(TestDir, "PS.exe")},
				},
			},
			ParamRunProg: 5,
			ParamFork:    true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Description, func(t *testing.T) {
			// // copy
			// var readData, _ = ioutil.ReadFile("~/Downloads/rufus.exe")
			// // paste
			// ioutil.WriteFile(testCase.ParamRunner.List[0].Path, readData, os.FileMode(0755))
			// // remove later
			// defer os.Remove(testCase.ParamRunner.List[0].Path)

			var gottenErr = testCase.ParamRunner.runFromList(testCase.ParamRunProg, testCase.ParamFork)

			if equalErrorList(t, []error{testCase.ExpectedErr}, []error{gottenErr}) == false {
				errorExpGot(t, testCase.ExpectedErr, gottenErr, true)
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
				List: []Exe{
					{"sr", pathJoin(TestDir, "sr.exe")},
					{"lon", pathJoin(TestDir, "lon.exe")},
				},
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Description, func(t *testing.T) {
			var gotten = testCase.ParamRunner.displayList()

			if testCase.Expected != gotten {
				errorExpGot(t, testCase.Expected, gotten, false)
			}
		})
	}

}
