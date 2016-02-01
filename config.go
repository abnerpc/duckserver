package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

const configFilePath string = "./conf.json"

type Configuration struct {
	AdminKeys map[string]string
	UserKeys  map[string]string
}

// DefaultConfig return the default configuration. If the saveToFile parameter
// is true, than the configuration is saved to configFilePath file.
func DefaultConfig(saveToFile bool) (*Configuration, error) {
	config := &Configuration{}
	config.AdminKeys = map[string]string{"123": "admin"}
	if saveToFile {
		configJSON, err := json.Marshal(config)
		if err != nil {
			return nil, err
		}
		if err = ioutil.WriteFile(configFilePath, configJSON, 0644); err != nil {
			return nil, err
		}
	}
	return config, nil
}

// LoadConfiguration loads the configuration from configFilePath and returns a
// Configuration pointer
func LoadConfiguration() (*Configuration, error) {

	var config *Configuration

	file, err := os.Open(configFilePath)
	if err != nil {
		return DefaultConfig(true)
	}

	decoder := json.NewDecoder(file)
	config = &Configuration{}
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
