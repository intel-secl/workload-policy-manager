package config

import (
	"io/ioutil"
	"log"
	"strings"
)

/*
 *
 * @author srege
 *
 */
var Configuration struct {
	KmsAPIURL                  string
	KmsAPIUsername             string
	KmsAPIPassword             string
	KmsTLSSHA256               string
	EnvelopePublickeyLocation  string
	EnvelopePrivatekeyLocation string
}

var configFilePath = "/opt/wpm/configuration/wpm.properties"

// SetConfigValues receives a pointer to Foo so it can modify it.
func SetConfigValues() {
	fileContents, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Fatal("Error reading the config file")
	}

	configArray := strings.Split(string(fileContents), "\n")
	for i := 0; i < len(configArray)-1; i++ {
		tempConfig := strings.Split(configArray[i], "=")
		key := tempConfig[0]
		value := strings.Replace(tempConfig[1], "\"", "", -1)
		if strings.Contains(strings.ToLower(key), "url") {
			Configuration.KmsAPIURL = value
		} else if strings.Contains(strings.ToLower(key), "username") {
			Configuration.KmsAPIUsername = value
		} else if strings.Contains(strings.ToLower(key), "password") {
			Configuration.KmsAPIPassword = value
		} else if strings.Contains(strings.ToLower(key), "tls") {
			Configuration.KmsTLSSHA256 = value
		} else if strings.Contains(strings.ToLower(key), "private") {
			Configuration.EnvelopePrivatekeyLocation = value
		} else if strings.Contains(strings.ToLower(key), "public") {
			Configuration.EnvelopePublickeyLocation = value
		}
	}
}
