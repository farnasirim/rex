package io

import (
	"io/ioutil"

	log "github.com/sirupsen/logrus"
)

func ReadFileOrFatal(filepath string) []byte {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatalf("Failed to read %s: %v\n", filepath, err)
	}
	return content
}
