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
		Description      string
		Expected         Runner
		ParamRunner      Runner
		ParamConfigStart string
		ParamConfigAfter string
	}{
		{
			Description: "read a config with program set as wine staging",
			Expected: Runner{
				Program:     "wine-staging",
				ProgramArgs: "",
				List:        []Exe{},
				ConfigFile:  inTestDir("winelarc"),
			},
			ParamRunner: Runner{
				ConfigFile: inTestDir("winelarc"),
			},
			ParamConfigStart: "Program = wine-staging\n" +
				"ProgramArgs = ",
			ParamConfigAfter: "",
		},
		{
			Description: "read a config after it gets edited",
			Expected: Runner{
				Program:     "wine",
				ProgramArgs: "",
				List:        []Exe{},
				ConfigFile:  inTestDir("winelarc"),
			},
			ParamRunner: Runner{
				ConfigFile: inTestDir("winelarc"),
			},
			ParamConfigStart: "Program = wine-staging\n" + "ProgramArgs = ",
			ParamConfigAfter: "Program = wine\n" + "ProgramArgs = ",
		},
		{
			Description: "read a config with left values",
			Expected: Runner{
				Program:     "",
				ProgramArgs: "",
				List:        []Exe{},
				ConfigFile:  inTestDir("winelarc"),
			},
			ParamRunner: Runner{
				ConfigFile: inTestDir("winelarc"),
			},
			ParamConfigStart: "Prog = wine\n" + "Something = ",
			ParamConfigAfter: "",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Description, func(t *testing.T) {
			// write
			var _ = ioutil.WriteFile(testCase.ParamRunner.ConfigFile, []byte(testCase.ParamConfigStart), os.FileMode(0755))

			// first read
			testCase.ParamRunner.runnerReadConfig()

			// test straight away if there is no config after (to modify to)
			// else do modification and read again then test
			if testCase.ParamConfigAfter == "" {
				// test write read
				if fmt.Sprint(testCase.ParamRunner) != fmt.Sprint(testCase.Expected) {
					errorExpGot(t, testCase.Expected, testCase.ParamRunner, false)
				}
			} else {
				// edit
				var _ = ioutil.WriteFile(testCase.ParamRunner.ConfigFile, []byte(testCase.ParamConfigAfter), os.FileMode(0755))
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

func TestRunnerWriteConfig(t *testing.T) {
	os.MkdirAll(TestDir, 0755)
	defer os.RemoveAll(TestDir)

	var testTable = []struct {
		Description string
		Expected    string
		ExpectedErr error

		ParamRunner Runner
	}{
		{
			Description: "write a full regular config",
			Expected:    "Program = wine\nArgs = \n",

			ParamRunner: Runner{
				Program:     "wine",
				ProgramArgs: "",
				ConfigFile:  inTestDir("winelarc"),
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Description, func(t *testing.T) {
			// write using function
			testCase.ParamRunner.RunnerWriteConfig()

			// read as a single string
			var data, gottenErr = ioutil.ReadFile(testCase.ParamRunner.ConfigFile)
			var gotten = string(data)

			if testCase.Expected != gotten {
				errorExpGot(t, testCase.Expected, gotten, false)
			}

			if equalErrorList(t, []error{testCase.ExpectedErr}, []error{gottenErr}) == false {
				errorExpGot(t, testCase.ExpectedErr, gottenErr, true)
			}
		})
	}
}

func TestRunFromList(t *testing.T) {
	os.MkdirAll(TestDir, 0755)
	defer os.RemoveAll(TestDir)

	var testTable = []struct {
		Description string
		ExpectedErr error

		ParamRunner  Runner
		ParamRunProg int
		ParamFork    bool
	}{
		{
			Description: "fork launch a number out of the range of list",
			ExpectedErr: fmt.Errorf("exe number %d: not in list", 5),
			ParamRunner: Runner{
				Program:     "wine",
				ProgramArgs: "",
				List: []Exe{
					{"PS", inTestDir("PS.exe")},
				},
			},
			ParamRunProg: 5,
			ParamFork:    true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Description, func(t *testing.T) {
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
					{"sr", inTestDir("sr.exe")},
					{"lon", inTestDir("lon.exe")},
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
