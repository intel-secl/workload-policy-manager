package config

import (
	csetup "intel/isecl/lib/common/setup"
	"os"
	"time"

	logger "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
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
	EnvelopePublickeyLocation  string
	EnvelopePrivatekeyLocation string
}

const (
	KMS_API_URL      = "KMS_API_URL"
	KMS_API_USERNAME = "KMS_API_USERNAME"
	KMS_API_PASSWORD = "KMS_API_PASSWORD"
	KMS_TLS_SHA256   = "KMS_TLS_SHA256"
	ConfigFilePath   = "/opt/wpm/configuration/config.yml"
	LogDirPath       = "/opt/wpm/logs/"
	LogFileName      = "wpm.log"
)

// Save the configuration struct into configuration directory
func Save() error {
	file, err := os.OpenFile(ConfigFilePath, os.O_RDWR, 0)
	if err != nil {
		// we have an error
		if os.IsNotExist(err) {
			// error is that the config doesnt yet exist, create it
			file, err = os.Create(ConfigFilePath)
			if err != nil {
				return err
			}
		}
	}
	defer file.Close()
	return yaml.NewEncoder(file).Encode(Configuration)
}

// SaveConfiguration is used to save configurations that are provided in environment during setup tasks
// This is called when setup tasks are called
func SaveConfiguration(c csetup.Context) error {
	var err error
	Configuration.Kms.APIURL, err = c.GetenvString(KMS_API_URL, "Kms URL")
	if err != nil {
		return err
	}
	Configuration.Kms.APIUsername, err = c.GetenvString(KMS_API_USERNAME, "Kms Username")
	if err != nil {
		return err
	}
	Configuration.Kms.APIPassword, err = c.GetenvString(KMS_API_PASSWORD, "Kms Password")
	if err != nil {
		return err
	}
	Configuration.Kms.TLSSha256, err = c.GetenvString(KMS_TLS_SHA256, "Kms TLS SHA256")
	if err != nil {
		return err
	}
	return Save()
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
	// load from config
	/*file, err := os.Open(ConfigFilePath)
	if err == nil {
		defer file.Close()
		yaml.NewDecoder(file).Decode(&Configuration)
	}*/
	LogWriter = os.Stdout
}
