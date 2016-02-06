package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

const configFilePath string = "./conf.json"

const initialKey string = "123"
const Admin string = "admin"
const User string = "user"

type Configuration struct {
	AccessKeys map[string]string
}

var Config *Configuration

// LoadDefaultConfig return the default configuration. If the saveToFile parameter
// is true, than the configuration is saved to configFilePath file.
func LoadDefaultConfig(saveToFile bool) error {

	Config = &Configuration{}
	Config.AccessKeys = make(map[string]string)
	Config.AccessKeys[initialKey] = Admin

	if saveToFile {
		return WriteConfiguration()
	}
	return nil
}

// LoadConfiguration loads the configuration from configFilePath and returns a
// Configuration pointer
func LoadConfiguration() error {

	file, err := os.Open(configFilePath)
	if err != nil {
		return LoadDefaultConfig(true)
	}

	decoder := json.NewDecoder(file)
	Config = &Configuration{}
	return decoder.Decode(Config)
}

// WriteConfiguration saves the current Configuration state to configFilePath
func WriteConfiguration() error {
	configJSON, err := json.Marshal(Config)
	if err == nil {
		err = ioutil.WriteFile(configFilePath, configJSON, 0644)
	}
	return err
}

func (c *Configuration) ChangeAccessKey(oldKey, newKey string) (string, bool) {

	userType, exists := c.AccessKeys[oldKey]
	if !exists {
		return "Access Key not found", false
	}
	delete(c.AccessKeys, oldKey)
	c.AccessKeys[newKey] = userType
	err := WriteConfiguration()
	if err != nil {
		return "Problem to save configuration", false
	}
	return "Success", true
}

func (c *Configuration) AddAccessKey(key, userType string) (string, bool) {
	if userType != Admin && userType != User {
		return "Invalid user type", false
	}
	_, exists := c.AccessKeys[key]
	if exists {
		return "Access Key in use", false
	}
	c.AccessKeys[key] = userType
	err := WriteConfiguration()
	if err != nil {
		return "Problem to save configuration", false
	}
	return "Success", true
}

func (c *Configuration) DeleteAccessKey(key string) (string, bool) {

	_, exists := c.AccessKeys[key]
	if !exists {
		return "Access Key not found", false
	}
	delete(c.AccessKeys, key)
	err := WriteConfiguration()
	if err != nil {
		return "Problem to save configuration", false
	}
	return "Success", true
}

func (c *Configuration) ListAccessKeys() (map[string]string, error) {
	return c.AccessKeys, nil
}
