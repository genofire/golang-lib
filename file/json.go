// Package file provides functionality to load and save marshal files
package file

import (
	"encoding/json"
	"os"
)

// ReadJSON reads a config model from path of a yml file
func ReadJSON(path string, data interface{}) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	return json.NewDecoder(file).Decode(data)
}

// SaveJSON to path
func SaveJSON(outputFile string, data interface{}) error {
	tmpFile := outputFile + ".tmp"

	file, err := os.OpenFile(tmpFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	err = json.NewEncoder(file).Encode(data)
	if err != nil {
		return err
	}

	file.Close()
	return os.Rename(tmpFile, outputFile)
}
