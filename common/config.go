package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const ConfigFilePath string = "./conf.json"

type Configuration struct {
	AdminKeys map[string]string
	UserKeys  map[string]string
}

// WriteNewConfig writes default configuration and returns a Configuration pointer
func WriteNewConfig() (*Configuration, error) {
	config := &Configuration{}
	config.AdminKeys = map[string]string{"123": "admin"}
	configJSON, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	if err = ioutil.WriteFile(ConfigFilePath, configJSON, 0644); err != nil {
		return nil, err
	}
	return config, nil
}

// LoadConfiguration loads the configuration fron ConfigFilePath and returns a Configuration pointer
func LoadConfiguration() (*Configuration, error) {

	var config *Configuration

	file, err := os.Open(ConfigFilePath)
	if err != nil {
		config, err = WriteNewConfig()
		if err != nil {
			fmt.Println("Can't create new config file")
			return nil, err
		}
		return config, nil
	}

	decoder := json.NewDecoder(file)
	config = &Configuration{}
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
