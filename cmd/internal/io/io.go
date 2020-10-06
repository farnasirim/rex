package io

import (
	"io/ioutil"

	log "github.com/sirupsen/logrus"
)

// ReadFileOrFatal reads a file and returns its contents. Will exit if an
// error is encountered in the process.
func ReadFileOrFatal(filepath string) []byte {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatalf("Failed to read %s: %v\n", filepath, err)
	}
	return content
}
