package util

import "os"

//
// DirectoryExists return true if `filename` exists and is a directory
//
func DirectoryExists(filename string) bool {

	stat, err := os.Stat(filename)

	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		panic(err)
	} else {
		return stat.Mode().IsDir()
	}
}

//
// FileExists return true if `filename` exists and is a file
//
func FileExists(filename string) bool {

	stat, err := os.Stat(filename)

	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		panic(err)
	} else {
		return stat.Mode().IsRegular()
	}

}

//
// PathExists return true if `filename` exists
//
func PathExists(filename string) bool {

	_, err := os.Stat(filename)

	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		panic(err)
	} else {
		return true
	}

}

//
// returns true if toFind is in `list`
//
func StringInArray(toFind string, list []string) bool {
	for _, curr := range list {
		if curr == toFind {
			return true
		}
	}
	return false
}
