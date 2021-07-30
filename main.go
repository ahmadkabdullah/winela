package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	// only launch if given an argument
	if len(os.Args) > 1 {
		// let function handle import/export of runner
		var r = runnerInitMake()
		var returnedCode = launch(r, os.Args[1:])
		os.Exit(returnedCode)
	}

	fmt.Println(`winela [opts]
	-r   [num]   # run a program from the list
	-R   [num]   # run a program without forking the process
	-s   [dir]   # scan a directory to populate list with
	-l           # print out the list`)
}

// central function for usage of functions
// and handling of arguments
func launch(rnr Runner, args []string) int {
	switch args[0] {
	case "-r", "-R":
		// alert if no number given
		if len(args) == 1 {
			fmt.Printf("input error: give a number to launch\n")
			return 1
		}

		// convert given number
		var convertedInt, convErr = strconv.Atoi(args[1])
		if convErr != nil {
			fmt.Printf("conversion error: %v is not a number\n", args[1])
			return 2
		}

		// run
		switch args[0] {
		// if "r" then fork
		case "-r":
			var runErr = rnr.runFromList(convertedInt, true)
			if runErr != nil {
				fmt.Printf("run error: %s\n", runErr.Error())
				return 3
			}
			fmt.Printf("stat: number %v was run\n", convertedInt)
		// if "R" then don't fork
		case "-R":
			fmt.Printf("stat: number %v will be run\n", convertedInt)
			var runErr = rnr.runFromList(convertedInt, false)
			if runErr != nil {
				fmt.Printf("run error: %s\n", runErr.Error())
				return 3
			}
			fmt.Printf("stat: number %v finished running\n", convertedInt)
		}

	case "-s":
		// set target dir according to given value if any
		// otherwise use user home dir
		var dirToScan string
		switch len(args) {
		case 1:
			fmt.Printf("stat: no scan dir given - assuming home dir\n")

			var homedir, locateErr = os.UserHomeDir()
			if locateErr != nil {
				fmt.Printf("locating home dir error: %s\n", locateErr.Error())
				return 2
			}

			dirToScan = homedir
		case 2:
			dirToScan = args[1]
		}

		// do the scan
		var list, scanErr = importFromScan(dirToScan)
		// output all errors if they exist
		if scanErr != nil {
			for _, e := range scanErr {
				fmt.Printf("scanning error: %s\n", e.Error())
			}
			return 3
		}

		fmt.Printf("stat: dir %s was scanned\n", dirToScan)

		// export the scanned dir to wineladb
		var exportErr = exportToFile(rnr.ListFile, list)
		if exportErr != nil {
			fmt.Printf("exporting scanned list error: %s\n", exportErr.Error())
			return 1
		}

		fmt.Printf("stat: list exported to %s\n", rnr.ListFile)

	case "-l":
		// print every exe in list
		var toDisplay = rnr.displayList()
		fmt.Print(toDisplay)

		fmt.Printf("stat: list printed\n")

	default:
		fmt.Printf("input error: option %v is unusable\n", args[0])
		return 1
	}

	return 0
}
