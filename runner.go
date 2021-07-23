package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type Runner struct {
	Program     string
	ProgramArgs string
	List        []exe

	ConfigFile string
	ListFile   string
}

// see if there is a configuration stored in configuration dir
// if not then create it by exporting default config and creating files
// if yes then import the configuration and set it as running configuration
func RunnerInitMake() (ret Runner) {
	// fix: deal with confdir error
	var confDir, _ = os.UserConfigDir()
	var progDir = PathJoin(confDir, "winel")

	// set defaults
	ret.Program = "wine"
	ret.ProgramArgs = ""
	ret.List = []exe{}
	ret.ConfigFile = PathJoin(progDir, "winelrc")
	ret.ListFile = PathJoin(progDir, "wineldb")

	// try import and go from there

	// read program config dir
	var _, readDirErr = ioutil.ReadDir(progDir)
	switch os.IsNotExist(readDirErr) {
	case true:
		// if it doesn't exist then create it
		os.MkdirAll(progDir, os.FileMode(0755))
		// then return as nothing can be imported
		return
	}

	// read config file
	var _, readRCErr = ioutil.ReadFile(ret.ConfigFile)
	switch os.IsNotExist(readRCErr) {
	case true:
		ret.RunnerWriteConfig()
	case false:
		ret.RunnerReadConfig()
	}

	// read list file
	var _, readDBErr = ioutil.ReadFile(ret.ListFile)
	switch os.IsNotExist(readDBErr) {
	case true:
		// if it doesn't exist create it but don't populate it
		ioutil.WriteFile(ret.ListFile, []byte{}, os.FileMode(0755))
	case false:
		// if it exists then import the list from it
		var importedList, importErr = ImportFromFile(ret.ListFile)
		// and if no errors occured set it to returning list
		if importErr == nil {
			ret.List = importedList
		}
	}

	return
}

// read the configuration from previously saved file
// into the current runner (startup, init)
func (r *Runner) RunnerReadConfig() {
	var readData, _ = ioutil.ReadFile(r.ConfigFile)
	var strData = string(readData)
	var lines = strings.Split(strData, "\n")
	for _, line := range lines {
		var pair = strings.Split(line, ":")

		if len(pair) != 2 {
			continue
		}

		var left = strings.TrimSpace(pair[0])
		var right = strings.TrimSpace(pair[1])

		switch left {
		case "Program":
			r.Program = right
		case "Args":
			r.ProgramArgs = right
		}
	}
}

// write the configuration in the runner into a file
func (r *Runner) RunnerWriteConfig() {
	var leftList = []string{
		"Program", "Args",
	}
	var rightList = []string{
		r.Program, r.ProgramArgs,
	}

	var strList string

	for i := range leftList {
		strList += fmt.Sprintf("%v = %v\n", leftList[i], rightList[i])
	}

	ioutil.WriteFile(
		r.ConfigFile,
		[]byte(strList),
		os.FileMode(0755),
	)
}

// run specified exe from the list of exes
// choosing whether to fork the process or not
func (r Runner) RunFromList(elementNumber int, shouldFork bool) error {
	// vars for storing target exe and its found stat
	var targetExe exe
	var foundTarget bool

	// see if target exe is in the list
	for _, exeEntry := range r.List {
		if exeEntry.Number == elementNumber {
			targetExe = exeEntry
			foundTarget = true
			break
		}
	}

	// if target was not found return not found error
	if !foundTarget {
		return fmt.Errorf("exe number %d: not in list", elementNumber)
	}

	// make up command with or without arguments
	var commandToRun *exec.Cmd
	switch r.ProgramArgs {
	case "":
		commandToRun = exec.Command(r.Program, targetExe.Path)
	default:
		commandToRun = exec.Command(r.Program, r.ProgramArgs, targetExe.Path)
	}

	if shouldFork {
		// start and letgo
		var execErr = commandToRun.Start()
		if execErr != nil {
			return fmt.Errorf("could not execute %s: %s", r.Program, execErr.Error())
		}
	} else {
		// for non-repetition
		var scanPipe = func(prefix string, scanner *bufio.Scanner, channel chan bool) {
			for scanner.Scan() {
				fmt.Println(prefix, scanner.Text())
			}
			channel <- true
		}

		// pipe errors to a reader and scan it in a go routine
		var errReader, _ = commandToRun.StderrPipe()
		var errScanner = bufio.NewScanner(errReader)
		var finishedError = make(chan bool)
		go scanPipe("ERR:", errScanner, finishedError)

		// pipe output to a reader and scan it in a go routine
		var outReader, _ = commandToRun.StdoutPipe()
		var outScanner = bufio.NewScanner(outReader)
		var finishedOutput = make(chan bool)
		go scanPipe("OUT:", outScanner, finishedOutput)

		// run the command
		var execErr = commandToRun.Start()
		if execErr != nil {
			return fmt.Errorf("could not execute %s: %s", r.Program, execErr.Error())
		}

		// wait until channels send finish bool
		<-finishedError
		<-finishedOutput
	}

	return nil
}

// return the list as a numbered string
func (r Runner) DisplayList() (ret string) {
	for _, entry := range r.List {
		ret += fmt.Sprintf("%v %v\n", entry.Number, entry.Name)
	}
	return
}
