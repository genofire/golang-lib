package file

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
)

// ReadConfigFile reads a config model from path of a yml file
func ReadTOML(path string, data interface{}) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	if err := toml.Unmarshal(file, data); err != nil {
		return err
	}

	return nil
}

// ReadJSON reads a config model from path of a yml file
func ReadJSON(path string, data interface{}) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	if err := json.NewDecoder(file).Decode(data); err != nil {
		return err
	}

	return nil
}

// SaveJSON to path
func SaveJSON(outputFile string, data interface{}) error {
	tmpFile := outputFile + ".tmp"

	file, err := os.OpenFile(tmpFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	if err := json.NewEncoder(file).Encode(data); err != nil {
		return err
	}

	file.Close()
	if err := os.Rename(tmpFile, outputFile); err != nil {
		return err
	}
	return nil
}
