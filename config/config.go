package config

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	logger "github.com/sirupsen/logrus"
)

/*
 *
 * @author srege
 *
 */
const (
	LogDirPath  = "/opt/wpm/logs/"
	LogFileName = "wpm.log"
)

var Configuration struct {
	Kms struct {
		APIURL      string
		APIUsername string
		APIPassword string
		TLSSHA256   string
	}
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
			Configuration.Kms.APIURL = value
		} else if strings.Contains(strings.ToLower(key), "username") {
			Configuration.Kms.APIUsername = value
		} else if strings.Contains(strings.ToLower(key), "password") {
			Configuration.Kms.APIPassword = value
		} else if strings.Contains(strings.ToLower(key), "tls") {
			Configuration.Kms.TLSSHA256 = value
		} else if strings.Contains(strings.ToLower(key), "private") {
			Configuration.EnvelopePrivatekeyLocation = value
		} else if strings.Contains(strings.ToLower(key), "public") {
			Configuration.EnvelopePublickeyLocation = value
		}
	}
}

// LogConfiguration is used to setup log configurations
func LogConfiguration() {
	// creating the log file if not preset
	LogFilePath := LogDirPath + LogFileName
	logFile, err := os.OpenFile(LogFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModeAppend)
	if err != nil {
		fmt.Println("unable to write file ", err)
		return
	}
	logger.SetFormatter(&logger.TextFormatter{FullTimestamp: true, TimestampFormat: time.RFC1123Z})
	logMultiWriter := io.MultiWriter(os.Stdout, logFile)
	logger.SetOutput(logMultiWriter)
}

var LogWriter io.Writer

func init() {
	LogWriter = os.Stdout
}
