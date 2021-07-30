package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type exe struct {
	Name string
	Path string
}

// read a file with a specific (exelist) format and get the list in it
func ImportFromFile(fileName string) (retList []exe, retErr error) {
	// read the file
	data, retErr := ioutil.ReadFile(fileName)

	if retErr != nil {
		return
	}

	// split file into string slice, each line a string
	var entriesInFile = strings.Split(string(data), "\n")

	// loop through the strings
	for _, entryFull := range entriesInFile {

		// split each line into two parts
		var entryInTwo = strings.Split(entryFull, "=>")

		// check if split was done right (has two parts)
		// and only proceed then
		if len(entryInTwo) != 2 {
			continue
		}

		// append the two parts each to a field in an exe struct
		var tempExe = exe{
			Name: strings.TrimSpace(entryInTwo[0]),
			Path: strings.TrimSpace(entryInTwo[1]),
		}

		// append the exe struct to the returned list
		retList = append(
			retList,
			tempExe,
		)
	}

	return
}

// scan a directory and get a list of exe files in it (recursively)
func ImportFromScan(dirName string) (retList []exe, retErr []error) {
	// read the dir
	var dirEntryList, readErr = ioutil.ReadDir(dirName)

	// if had any errors
	if readErr != nil {
		// add error to list
		retErr = append(retErr, readErr)
		// exit
		return
	}

	// go through dir entries
	for _, dirEntry := range dirEntryList {

		// create vars for the entry
		var (
			dirEntryName         = dirEntry.Name()
			dirEntryPath         = PathJoin(dirName, dirEntry.Name())
			dirEntryNameNoSuffix = strings.TrimSuffix(dirEntryName, ".exe")
		)

		// act depending on it being a dir or file
		if dirEntry.IsDir() {
			// recursive call to read dirs
			var recurList, recurErr = ImportFromScan(dirEntryPath)
			// assign the recursive err to return one
			retErr = recurErr
			// add result to caller
			retList = append(retList, recurList...)

		} else {
			// if file has wrong extension just skip
			if filepath.Ext(dirEntry.Name()) != ".exe" {
				continue
			}

			// try opening the file
			var readFile, readErr = os.Open(dirEntryPath)
			defer readFile.Close()

			// don't add files that can't be read
			if readErr != nil {
				// add caught error to list and go to next file
				retErr = append(retErr, readErr)
				continue
			}

			// then add to list
			retList = append(retList,
				exe{
					Name: dirEntryNameNoSuffix,
					Path: dirEntryPath,
				},
			)
		}
	}

	return
}

func ExportToFile(fileName string, listToWrite []exe) (retErr error) {
	// read out the whole struct into a string
	// with proper formatting
	var dataAsString string
	for _, element := range listToWrite {
		dataAsString += fmt.Sprintf("%s => %s\n", element.Name, element.Path)
	}

	// write the acquired string to the specidied file
	var writeErr = ioutil.WriteFile(fileName, []byte(dataAsString), os.FileMode(0755))

	// if there was an error in writing return it
	if writeErr != nil {
		retErr = writeErr
		return
	}

	return
}
