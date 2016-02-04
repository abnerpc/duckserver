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

var Config *Configuration

// LoadDefaultConfig return the default configuration. If the saveToFile parameter
// is true, than the configuration is saved to configFilePath file.
func LoadDefaultConfig(saveToFile bool) error {

	Config = &Configuration{}
	Config.AdminKeys = map[string]string{"123": "admin"}
	Config.UserKeys = make(map[string]string)

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

func (c *Configuration) ChangeUserPassword(user string, password string) (string, bool) {
	currentPassword := ""
	for k, v := range c.UserKeys {
		if v == user {
			currentPassword = k
			break
		}
	}
	if currentPassword == "" {
		return "User not found", false
	}
	delete(c.UserKeys, currentPassword)
	c.UserKeys[password] = user
	err := WriteConfiguration()
	if err != nil {
		return "Problem to save configuration", false
	}
	return "Success", true
}

func (c *Configuration) AddUser(user string, password string) (string, bool) {
	_, exists := c.UserKeys[password]
	if exists {
		return "Password in use", false
	}
	c.UserKeys[password] = user
	err := WriteConfiguration()
	if err != nil {
		return "Problem to save configuration", false
	}
	return "Success", true
}
