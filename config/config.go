package config

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
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

// LogConfiguration is used to setup log rotation configurations
func LogConfiguration() {
	// creating the log file if not preset
	LogFilePath := LogDirPath + LogFileName
	_, err := os.Stat(LogFilePath)
	if os.IsNotExist(err) {
		logger.Debug("Log file does not exist. Creating the file.")
		_, touchErr := exec.Command("touch", LogFilePath).Output()
		if touchErr != nil {
			fmt.Println("Error while creating the log file.", touchErr)
			return
		}
	}
	logFile, err := os.OpenFile(LogFilePath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Printf("unable to write file on filehook %v\n", err)
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
