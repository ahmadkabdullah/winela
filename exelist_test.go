package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"testing"
)

// integrating tests

func TestExportImport(t *testing.T) {
	// make testing directory
	os.MkdirAll(TestDir, 0755)
	defer os.RemoveAll(TestDir)

	var testTable = []struct {
		Description  string
		Expected     []Exe
		ExpectedErrs []error

		ParamSearchDir string
		ParamWriteFile string
		ParamDirs      []PairPathPerm
		ParamFiles     []PairPathPerm
	}{
		{
			Description: "scan a dir and export result to a file then import back from exported file",
			Expected: []Exe{
				{"flap", pathJoin(TestDir, "games/flap.exe")},
				{"paint", pathJoin(TestDir, "ms/paint.exe")},
				{"pt", pathJoin(TestDir, "pt.exe")},
			},
			ExpectedErrs: []error{},

			ParamSearchDir: TestDir,
			ParamWriteFile: pathJoin(TestDir, "expoFile"),
			ParamDirs: []PairPathPerm{
				{pathJoin(TestDir, "ms"), 0755},
				{pathJoin(TestDir, "games"), 0755},
			},
			ParamFiles: []PairPathPerm{
				{pathJoin(TestDir, "ms/paint.exe"), 0755},
				{pathJoin(TestDir, "games/flap.exe"), 0755},
				{pathJoin(TestDir, "pt.exe"), 0755},
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Description, func(t *testing.T) {
			// make dirs and defer deletion
			for _, dirToMake := range testCase.ParamDirs {
				os.Mkdir(dirToMake.Path, fs.FileMode(dirToMake.Perm))
				defer os.RemoveAll(dirToMake.Path)
			}

			// create files and defer deletion
			for _, fileToMake := range testCase.ParamFiles {
				os.WriteFile(fileToMake.Path, []byte{}, fs.FileMode(fileToMake.Perm))
				defer os.RemoveAll(fileToMake.Path)
			}

			// make var to store errors of all steps
			var gottenErrs []error

			// step 1 run a scan
			var listFromScan, scanErrs = importFromScan(testCase.ParamSearchDir)
			if scanErrs != nil {
				gottenErrs = append(gottenErrs, scanErrs...)
			}

			// step 2 export into a file
			var writeErr = exportToFile(testCase.ParamWriteFile, listFromScan)
			defer os.RemoveAll(testCase.ParamWriteFile)
			if writeErr != nil {
				gottenErrs = append(gottenErrs, writeErr)
			}

			// step 3 read from file
			var listFromFile, readErr = importFromFile(testCase.ParamWriteFile)
			if readErr != nil {
				gottenErrs = append(gottenErrs, readErr)
			}

			// assign step 3 results to gotten results
			var gotten = listFromFile

			// test

			if equalExeList(t, testCase.Expected, gotten) == false {
				errorExpGot(t, testCase.Expected, gotten, false)
			}

			if equalErrorList(t, testCase.ExpectedErrs, gottenErrs) == false {
				errorExpGot(t, testCase.ExpectedErrs, gottenErrs, true)
			}
		})
	}
}

// function tests

func TestImportFromFile(t *testing.T) {
	// make testing directory
	os.MkdirAll(TestDir, 0755)
	defer os.RemoveAll(TestDir)

	// file name for all cases
	var testFileName = pathJoin(TestDir, "testFile")

	var testTable = []struct {
		Description string
		Expected    []Exe
		ExpectedErr error

		ParamContent string
		ParamFile    PairPathPerm
	}{
		{
			Description: "import regular file",
			Expected:    []Exe{{"okay", "~/Downloads/okay.exe"}},
			ExpectedErr: nil,

			ParamContent: "okay => ~/Downloads/okay.exe\n",
			ParamFile:    PairPathPerm{Path: testFileName, Perm: 0755},
		},
		{
			Description: "file can't be read",
			Expected:    []Exe{},
			ExpectedErr: fmt.Errorf("open %s: permission denied", testFileName),

			ParamContent: "okay => ~/Downloads/okay.exe\n" + "yes=> ~/go/bin/yes.exe\n",
			ParamFile:    PairPathPerm{Path: testFileName, Perm: 0111},
		},
		{
			Description: "multiple separators in one line in file",
			Expected: []Exe{
				{"okay", "~/Downloads/okay.exe"},
				{"yes", "~/go/bin/yes.exe"},
			},
			ExpectedErr: nil,

			ParamContent: "okay => ~/Downloads/okay.exe\n" + "hey=>~/hey.exe=>exe\n" + "yes=> ~/go/bin/yes.exe\n",
			ParamFile:    PairPathPerm{Path: testFileName, Perm: 0755},
		},
	}

	// case cycling

	for _, testCase := range testTable {
		t.Run(testCase.Description, func(t *testing.T) {

			// write a file with some data
			os.WriteFile(
				testFileName,
				[]byte(testCase.ParamContent),
				fs.FileMode(testCase.ParamFile.Perm),
			)
			// remove file created by each test to prevent fileMode mixing
			defer os.Remove(testFileName)

			// run function
			var gotten, gottenErr = importFromFile(testFileName)

			// test data
			if equalExeList(t, testCase.Expected, gotten) == false {
				errorExpGot(t, testCase.Expected, gotten, false)
			}

			// test error
			if equalErrorList(t, []error{testCase.ExpectedErr}, []error{gottenErr}) == false {
				errorExpGot(t, testCase.ExpectedErr, gottenErr, true)
			}
		})
	}
}

func TestImportFromScan(t *testing.T) {
	// make testing directory
	os.MkdirAll(TestDir, 0755)
	defer os.RemoveAll(TestDir)

	var testTable = []struct {
		Description  string
		Expected     []Exe
		ExpectedErrs []error

		ParamSearch string
		ParamDirs   []PairPathPerm
		ParamFiles  []PairPathPerm
	}{
		{
			Description: "scanning a regular dir",
			Expected: []Exe{
				{Name: "second", Path: pathJoin(TestDir, "extra/second.exe")},
				{Name: "first", Path: pathJoin(TestDir, "first.exe")},
				{Name: "third", Path: pathJoin(TestDir, "third.exe")},
			},
			ExpectedErrs: []error{},

			ParamSearch: TestDir,
			ParamDirs: []PairPathPerm{
				{pathJoin(TestDir, "extra"), 0755},
			},
			ParamFiles: []PairPathPerm{
				{pathJoin(TestDir, "first.exe"), 0755},
				{pathJoin(TestDir, "extra/second.exe"), 0755},
				{pathJoin(TestDir, "third.exe"), 0755},
			},
		},
		{
			Description: "dir not accessible",
			Expected: []Exe{
				{Name: "first", Path: pathJoin(TestDir, "first.exe")},
				{Name: "third", Path: pathJoin(TestDir, "third.exe")},
			},
			ExpectedErrs: []error{
				// clunky
				fmt.Errorf("open %s: permission denied", pathJoin(TestDir, "extra")),
			},

			ParamSearch: TestDir,
			ParamDirs: []PairPathPerm{
				{pathJoin(TestDir, "extra"), 0111},
			},
			ParamFiles: []PairPathPerm{
				{pathJoin(TestDir, "first.exe"), 0755},
				{pathJoin(TestDir, "extra/second.exe"), 0755},
				{pathJoin(TestDir, "third.exe"), 0755},
			},
		},
		{
			Description: "file not accessible",
			Expected: []Exe{
				{Name: "third", Path: pathJoin(TestDir, "third.exe")},
			},
			ExpectedErrs: []error{
				// clunky
				fmt.Errorf("open %s: permission denied", pathJoin(TestDir, "first.exe")),
			},

			ParamSearch: TestDir,
			ParamDirs:   []PairPathPerm{},
			ParamFiles: []PairPathPerm{
				{pathJoin(TestDir, "first.exe"), 0222},
				{pathJoin(TestDir, "third.exe"), 0755},
			},
		},
		{
			Description: "nested file not accessible",
			Expected: []Exe{
				{Name: "first", Path: pathJoin(TestDir, "first.exe")},
				{Name: "third", Path: pathJoin(TestDir, "third.exe")},
			},
			ExpectedErrs: []error{
				// clunky
				fmt.Errorf("open %s: permission denied", pathJoin(TestDir, "extra/second.exe")),
			},

			ParamSearch: TestDir,
			ParamDirs: []PairPathPerm{
				{pathJoin(TestDir, "extra"), 0755},
			},
			ParamFiles: []PairPathPerm{
				{pathJoin(TestDir, "first.exe"), 0755},
				{pathJoin(TestDir, "extra/second.exe"), 0122},
				{pathJoin(TestDir, "third.exe"), 0755},
			},
		},
		{
			Description: "different and wrong extensions",
			Expected: []Exe{
				{Name: "alpha", Path: pathJoin(TestDir, "alpha.exe")},
			},
			ExpectedErrs: []error{},

			ParamSearch: TestDir,
			ParamDirs: []PairPathPerm{
				{pathJoin(TestDir, "ost"), 0755},
			},
			ParamFiles: []PairPathPerm{
				{pathJoin(TestDir, "alpha.exe"), 0755},
				{pathJoin(TestDir, "beta.xe"), 0755},
				{pathJoin(TestDir, "ost/zetta.mp3"), 0755},
			},
		},
	}

	// case cycling

	for _, testCase := range testTable {
		t.Run(testCase.Description, func(t *testing.T) {
			// create wanted directories
			for _, dirToMake := range testCase.ParamDirs {
				os.Mkdir(
					dirToMake.Path,
					fs.FileMode(dirToMake.Perm),
				)
				defer os.RemoveAll(dirToMake.Path)
			}

			// create wanted files
			for _, fileToMake := range testCase.ParamFiles {
				os.WriteFile(
					fileToMake.Path,
					// empty content
					[]byte(""),
					fs.FileMode(fileToMake.Perm),
				)
				defer os.RemoveAll(fileToMake.Path)
			}

			// run
			var gotten, gottenErrs = importFromScan(testCase.ParamSearch)

			// test

			if equalExeList(t, testCase.Expected, gotten) == false {
				errorExpGot(t, testCase.Expected, gotten, false)
			}

			if equalErrorList(t, testCase.ExpectedErrs, gottenErrs) == false {
				errorExpGot(t, testCase.ExpectedErrs, gottenErrs, true)
			}
		})
	}
}

func TestExportToFile(t *testing.T) {
	// make testing directory
	os.MkdirAll(TestDir, 0755)
	defer os.RemoveAll(TestDir)

	var testTable = []struct {
		Description string
		Expected    string
		ExpectedErr error

		ParamFile PairPathPerm
		ParamList []Exe
	}{
		{
			Description: "exporting a regular file",
			Expected:    "ck => ~/Games/ck/ck.exe\nfff => ~/Downloads/fff.exe\n",
			ExpectedErr: nil,

			ParamFile: PairPathPerm{
				Path: pathJoin(TestDir, "exportedFile"),
				Perm: 0755,
			},
			ParamList: []Exe{
				{"ck", "~/Games/ck/ck.exe"},
				{"fff", "~/Downloads/fff.exe"},
			},
		},
		{
			Description: "exporting to a path with no permission",
			Expected:    "",
			ExpectedErr: fmt.Errorf("open /root/exportedFile: permission denied"),

			ParamFile: PairPathPerm{
				Path: pathJoin("/root", "exportedFile"),
				Perm: 0755,
			},
			ParamList: []Exe{
				{"ck", "~/Games/ck/ck.exe"},
				{"fff", "~/Downloads/fff.exe"},
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Description, func(t *testing.T) {
			// run func
			var gottenErr = exportToFile(
				testCase.ParamFile.Path,
				testCase.ParamList,
			)

			// read file exported
			var data, _ = ioutil.ReadFile(testCase.ParamFile.Path)
			defer os.Remove(testCase.ParamFile.Path)

			// set gotten to it
			var gotten = string(data)

			// test

			if testCase.Expected != gotten {
				errorExpGot(t, testCase.Expected, gotten, false)
			}

			if equalErrorList(t, []error{testCase.ExpectedErr}, []error{gottenErr}) == false {
				errorExpGot(t, testCase.ExpectedErr, gottenErr, true)
			}
		})
	}
}
