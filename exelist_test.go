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
		Expected     []exe
		ExpectedErrs []error

		ParamSearchDir string
		ParamWriteFile string
		ParamDirs      []PairPathPerm
		ParamFiles     []PairPathPerm
	}{
		{
			Description: "scan a dir and export result to a file then import back from exported file",
			Expected: []exe{
				{1, "flap", PathJoin(TestDir, "games/flap.exe")},
				{2, "paint", PathJoin(TestDir, "ms/paint.exe")},
				{3, "pt", PathJoin(TestDir, "pt.exe")},
			},
			ExpectedErrs: []error{},

			ParamSearchDir: TestDir,
			ParamWriteFile: PathJoin(TestDir, "expoFile"),
			ParamDirs: []PairPathPerm{
				{PathJoin(TestDir, "ms"), 0755},
				{PathJoin(TestDir, "games"), 0755},
			},
			ParamFiles: []PairPathPerm{
				{PathJoin(TestDir, "ms/paint.exe"), 0755},
				{PathJoin(TestDir, "games/flap.exe"), 0755},
				{PathJoin(TestDir, "pt.exe"), 0755},
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
			var listFromScan, scanErrs = ImportFromScan(testCase.ParamSearchDir)
			if scanErrs != nil {
				gottenErrs = append(gottenErrs, scanErrs...)
			}

			// step 2 export into a file
			var writeErr = ExportToFile(testCase.ParamWriteFile, listFromScan)
			defer os.RemoveAll(testCase.ParamWriteFile)
			if writeErr != nil {
				gottenErrs = append(gottenErrs, writeErr)
			}

			// step 3 read from file
			var listFromFile, readErr = ImportFromFile(testCase.ParamWriteFile)
			if readErr != nil {
				gottenErrs = append(gottenErrs, readErr)
			}

			// assign step 3 results to gotten results
			var gotten = listFromFile

			// test

			if EqualExeList(t, testCase.Expected, gotten) == false {
				ErrorExpGot(t, testCase.Expected, gotten, false)
			}

			if EqualErrorList(t, testCase.ExpectedErrs, gottenErrs) == false {
				ErrorExpGot(t, testCase.ExpectedErrs, gottenErrs, true)
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
	var testFileName = PathJoin(TestDir, "testFile")

	var testTable = []struct {
		Description string
		Expected    []exe
		ExpectedErr error

		ParamContent string
		ParamFile    PairPathPerm
	}{
		{
			Description: "import regular file",
			Expected:    []exe{{1, "okay", "~/Downloads/okay.exe"}},
			ExpectedErr: nil,

			ParamContent: "okay => ~/Downloads/okay.exe\n",
			ParamFile:    PairPathPerm{Path: testFileName, Perm: 0755},
		},
		{
			Description: "file can't be read",
			Expected:    []exe{},
			ExpectedErr: fmt.Errorf("open %s: permission denied", testFileName),

			ParamContent: "okay => ~/Downloads/okay.exe\n" + "yes=> ~/go/bin/yes.exe\n",
			ParamFile:    PairPathPerm{Path: testFileName, Perm: 0111},
		},
		{
			Description: "multiple separators in one line in file",
			Expected: []exe{
				{1, "okay", "~/Downloads/okay.exe"},
				{2, "yes", "~/go/bin/yes.exe"},
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
			var gotten, gottenErr = ImportFromFile(testFileName)

			// test data
			if EqualExeList(t, testCase.Expected, gotten) == false {
				ErrorExpGot(t, testCase.Expected, gotten, false)
			}

			// test error
			if EqualErrorList(t, []error{testCase.ExpectedErr}, []error{gottenErr}) == false {
				ErrorExpGot(t, testCase.ExpectedErr, gottenErr, true)
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
		Expected     []exe
		ExpectedErrs []error

		ParamSearch string
		ParamDirs   []PairPathPerm
		ParamFiles  []PairPathPerm
	}{
		{
			Description: "scanning a regular dir",
			Expected: []exe{
				{Name: "second", Path: PathJoin(TestDir, "extra/second.exe")},
				{Name: "first", Path: PathJoin(TestDir, "first.exe")},
				{Name: "third", Path: PathJoin(TestDir, "third.exe")},
			},
			ExpectedErrs: []error{},

			ParamSearch: TestDir,
			ParamDirs: []PairPathPerm{
				{PathJoin(TestDir, "extra"), 0755},
			},
			ParamFiles: []PairPathPerm{
				{PathJoin(TestDir, "first.exe"), 0755},
				{PathJoin(TestDir, "extra/second.exe"), 0755},
				{PathJoin(TestDir, "third.exe"), 0755},
			},
		},
		{
			Description: "dir not accessible",
			Expected: []exe{
				{Name: "first", Path: PathJoin(TestDir, "first.exe")},
				{Name: "third", Path: PathJoin(TestDir, "third.exe")},
			},
			ExpectedErrs: []error{
				// clunky
				fmt.Errorf("open %s: permission denied", PathJoin(TestDir, "extra")),
			},

			ParamSearch: TestDir,
			ParamDirs: []PairPathPerm{
				{PathJoin(TestDir, "extra"), 0111},
			},
			ParamFiles: []PairPathPerm{
				{PathJoin(TestDir, "first.exe"), 0755},
				{PathJoin(TestDir, "extra/second.exe"), 0755},
				{PathJoin(TestDir, "third.exe"), 0755},
			},
		},
		{
			Description: "file not accessible",
			Expected: []exe{
				{Name: "third", Path: PathJoin(TestDir, "third.exe")},
			},
			ExpectedErrs: []error{
				// clunky
				fmt.Errorf("open %s: permission denied", PathJoin(TestDir, "first.exe")),
			},

			ParamSearch: TestDir,
			ParamDirs:   []PairPathPerm{},
			ParamFiles: []PairPathPerm{
				{PathJoin(TestDir, "first.exe"), 0222},
				{PathJoin(TestDir, "third.exe"), 0755},
			},
		},
		{
			Description: "nested file not accessible",
			Expected: []exe{
				{Name: "first", Path: PathJoin(TestDir, "first.exe")},
				{Name: "third", Path: PathJoin(TestDir, "third.exe")},
			},
			ExpectedErrs: []error{
				// clunky
				fmt.Errorf("open %s: permission denied", PathJoin(TestDir, "extra/second.exe")),
			},

			ParamSearch: TestDir,
			ParamDirs: []PairPathPerm{
				{PathJoin(TestDir, "extra"), 0755},
			},
			ParamFiles: []PairPathPerm{
				{PathJoin(TestDir, "first.exe"), 0755},
				{PathJoin(TestDir, "extra/second.exe"), 0122},
				{PathJoin(TestDir, "third.exe"), 0755},
			},
		},
		{
			Description: "different and wrong extensions",
			Expected: []exe{
				{Name: "alpha", Path: PathJoin(TestDir, "alpha.exe")},
			},
			ExpectedErrs: []error{},

			ParamSearch: TestDir,
			ParamDirs: []PairPathPerm{
				{PathJoin(TestDir, "ost"), 0755},
			},
			ParamFiles: []PairPathPerm{
				{PathJoin(TestDir, "alpha.exe"), 0755},
				{PathJoin(TestDir, "beta.xe"), 0755},
				{PathJoin(TestDir, "ost/zetta.mp3"), 0755},
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
			var gotten, gottenErrs = ImportFromScan(testCase.ParamSearch)

			// test

			if EqualExeList(t, testCase.Expected, gotten) == false {
				ErrorExpGot(t, testCase.Expected, gotten, false)
			}

			if EqualErrorList(t, testCase.ExpectedErrs, gottenErrs) == false {
				ErrorExpGot(t, testCase.ExpectedErrs, gottenErrs, true)
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
		ParamList []exe
	}{
		{
			Description: "exporting a regular file",
			Expected:    "ck => ~/Games/ck/ck.exe\nfff => ~/Downloads/fff.exe\n",
			ExpectedErr: nil,

			ParamFile: PairPathPerm{
				Path: PathJoin(TestDir, "exportedFile"),
				Perm: 0755,
			},
			ParamList: []exe{
				{1, "ck", "~/Games/ck/ck.exe"},
				{2, "fff", "~/Downloads/fff.exe"},
			},
		},
		{
			Description: "exporting to a path with no permission",
			Expected:    "",
			ExpectedErr: fmt.Errorf("open /root/exportedFile: permission denied"),

			ParamFile: PairPathPerm{
				Path: PathJoin("/root", "exportedFile"),
				Perm: 0755,
			},
			ParamList: []exe{
				{1, "ck", "~/Games/ck/ck.exe"},
				{2, "fff", "~/Downloads/fff.exe"},
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Description, func(t *testing.T) {
			// run func
			var gottenErr = ExportToFile(
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
				ErrorExpGot(t, testCase.Expected, gotten, false)
			}

			if EqualErrorList(t, []error{testCase.ExpectedErr}, []error{gottenErr}) == false {
				ErrorExpGot(t, testCase.ExpectedErr, gottenErr, true)
			}
		})
	}
}
