package config

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
	csetup "intel/isecl/lib/common/setup"
	"intel/isecl/wpm/consts"
	"io"
	"os"
	"time"
)

/*
 *
 * @author srege
 *
 */

var Configuration struct {
	Kms struct {
		APIURL      string
		APIUsername string
		APIPassword string
		TLSSha256   string
	}
}

// Save the configuration struct into configuration directory
func Save() error {
	file, err := os.OpenFile(consts.ConfigFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModeAppend)
	defer file.Close()
	if err != nil {
		// we have an error
		if os.IsNotExist(err) {
			// error is that the config doesnt yet exist, create it
			file, err = os.Create(consts.ConfigFilePath)
			if err != nil {
				return err
			}
		}
	}
	return yaml.NewEncoder(file).Encode(consts.ConfigFilePath)
}

// SaveConfiguration is used to save configurations that are provided in environment during setup tasks
// This is called when setup tasks are called
func SaveConfiguration(c csetup.Context) error {
	var err error
	Configuration.Kms.APIURL, err = c.GetenvString(consts.KMS_API_URL, "Kms URL")
	if err != nil {
		return err
	}
	Configuration.Kms.APIUsername, err = c.GetenvString(consts.KMS_API_USERNAME, "Kms Username")
	if err != nil {
		return err
	}
	Configuration.Kms.APIPassword, err = c.GetenvString(consts.KMS_API_PASSWORD, "Kms Password")
	if err != nil {
		return err
	}
	Configuration.Kms.TLSSha256, err = c.GetenvString(consts.KMS_TLS_SHA256, "Kms TLS SHA256")
	if err != nil {
		return err
	}
	return Save()
}

// LogConfiguration is used to setup log configurations
func LogConfiguration() {
	// creating the log file if not preset
	LogFilePath := consts.LogDirPath + consts.LogFileName
	logFile, err := os.OpenFile(LogFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModeAppend)
	if err != nil {
		fmt.Println("unable to write file ", err)
		return
	}
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, TimestampFormat: time.RFC1123Z})
	logMultiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(logMultiWriter)
}

var LogWriter io.Writer

func init() {
	// load from config
	file, err := os.OpenFile(consts.ConfigFilePath,os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModeAppend)
	if err == nil {
		defer file.Close()
		yaml.NewDecoder(file).Decode(&Configuration)
	}
	LogWriter = os.Stdout
}
