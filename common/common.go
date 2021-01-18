package rk_common

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
)

// shutdown and print stack
func ShutdownWithError(err error) {
	if err == nil {
		err = errors.New("error is nil")
	}
	panic(err)
}

// read files with provided path, use working directory if given path is relative path
// shutdown process if any error occurs
// this is used for MUST SUCCESS scenario like reading config files
func MustReadFile(filePath string) []byte {
	if len(filePath) < 1 {
		ShutdownWithError(errors.New("empty file path"))
	}

	if !path.IsAbs(filePath) {
		wd, err := os.Getwd()

		if err != nil {
			ShutdownWithError(err)
		}
		filePath = path.Join(wd, filePath)
	}

	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		ShutdownWithError(err)
	}

	return bytes
}
