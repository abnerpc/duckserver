package main

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
)

const configFilePath string = "./access.json"

// Admin is the value for identify admin users
const Admin byte = 0

// User is the value for identify common (non admin) users
const User byte = 1

// Configuration is loaded from the configFilePath and has the current
// users configuration
type Configuration struct {
	Users      map[string]string
	AccessKeys map[string]byte
}

// Config is the global current configuration
var Config *Configuration

// WriteConfiguration saves the current Configuration state to configFilePath
func WriteConfiguration() error {
	configJSON, err := json.Marshal(Config)
	if err == nil {
		err = ioutil.WriteFile(configFilePath, configJSON, 0644)
	}
	return err
}

// LoadDefaultConfig returns the default configuration. If the saveToFile parameter
// is true, than the configuration is saved to configFilePath file.
func LoadDefaultConfig(saveToFile bool) error {

	Config = &Configuration{make(map[string]string), make(map[string]byte)}
	accessKey := basicAuth("admin", "123")
	Config.Users["admin"] = accessKey
	Config.AccessKeys[accessKey] = Admin

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

// See 2 (end of page 4) http://www.ietf.org/rfc/rfc2617.txt
// "To receive authorization, the client sends the userid and password,
// separated by a single colon (":") character, within a base64
// encoded string in the credentials."
// It is not meant to be urlencoded.
func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// parseBasicAuth parses an HTTP Basic Authentication string.
// "Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==" returns ("Aladdin", "open sesame", true).
func parseBasicAuth(auth string) (username, password string, ok bool) {
	const prefix = "Basic "
	if !strings.HasPrefix(auth, prefix) {
		return
	}
	c, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return
	}
	cs := string(c)
	s := strings.IndexByte(cs, ':')
	if s < 0 {
		return
	}
	return cs[:s], cs[s+1:], true
}

func (c *Configuration) changePassword(userName, newPassword string) (string, bool) {

	keyUserName := ""
	accessKey, ok := c.Users[userName]
	if ok {
		keyUserName, _, ok = parseBasicAuth("Basic " + accessKey)
	}
	if !ok || userName != keyUserName {
		return "User is invalid.", false
	}
	userType, ok := c.AccessKeys[accessKey]
	if !ok {
		return "User is invalid.", false
	}

	delete(c.Users, userName)
	delete(c.AccessKeys, accessKey)

	newAccessKey := basicAuth(userName, newPassword)
	c.Users[userName] = newAccessKey
	c.AccessKeys[newAccessKey] = userType

	err := WriteConfiguration()
	if err != nil {
		return "Problem to save configuration.", false
	}

	return "Success", true
}

func (c *Configuration) addUser(userName, password string, userType byte) (string, bool) {

	_, ok := c.Users[userName]
	if ok {
		return "User already exists.", false
	}

	accessKey := basicAuth(userName, password)
	c.Users[userName] = accessKey
	c.AccessKeys[accessKey] = userType

	err := WriteConfiguration()
	if err != nil {
		return "Problem to save configuration.", false
	}

	return "Success", true
}

func (c *Configuration) deleteUser(userName string) (string, bool) {

	accessKey, ok := c.Users[userName]
	if !ok {
		return "User not found.", false
	}
	delete(c.Users, userName)
	delete(c.AccessKeys, accessKey)

	err := WriteConfiguration()
	if err != nil {
		return "Problem to save configuration.", false
	}

	return "Success", true
}
