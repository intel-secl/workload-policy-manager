/*
 * Copyright (C) 2019 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package config

import (
	"fmt"
	csetup "intel/isecl/lib/common/v3/setup"
	"intel/isecl/wpm/v3/consts"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	commLog "intel/isecl/lib/common/v3/log"
	commLogInt "intel/isecl/lib/common/v3/log/setup"

	"github.com/sirupsen/logrus"

	yaml "gopkg.in/yaml.v2"
)

var (
	log    = commLog.GetDefaultLogger()
	secLog = commLog.GetSecurityLogger()
)

// Configuration holds the configuration required for WPM operations
var Configuration struct {
	CmsTLSCertDigest string
	Kms              struct {
		APIURL string
	}
	Cms struct {
		BaseURL string
	}
	Subject struct {
		CommonName string
	}
	Aas struct {
		APIURL string
	}
	Wpm struct {
		Username string
		Password string
	}
	FlavorSigningKeyFile       string
	FlavorSigningCertFile   string
	LogLevel          string
	LogEntryMaxLength int
	LogEnableStdout   bool
	ConfigComplete bool
}

var LogWriter io.Writer

func init() {
	// load from config
	file, err := os.Open(consts.ConfigFilePath)
	if err == nil {
		defer func() {
			derr := file.Close()
			if derr != nil {
				fmt.Fprintf(os.Stderr, "Error while closing file" + derr.Error())
	}
		}()
		err = yaml.NewDecoder(file).Decode(&Configuration)
	}
	LogWriter = os.Stdout
}

// Save the configuration struct into configuration directory
func Save() error {
	log.Trace("config/config:Save() Entering")
	defer log.Trace("config/config:Save() Leaving")

	file, err := os.OpenFile(consts.ConfigFilePath, os.O_RDWR, 0)
	defer func() {
		derr := file.Close()
		if derr != nil {
			fmt.Fprintf(os.Stderr, "Error while closing file" + derr.Error())
		}
	}()
	if err != nil {
		// we have an error
		if os.IsNotExist(err) {
			// error is that the config doesnt yet exist, create it
			file, err = os.OpenFile(consts.ConfigFilePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
			if err != nil {
				return errors.Wrapf(err, "Unable to write configuration file at %s", consts.ConfigFilePath)
			}
		}
	}
	return yaml.NewEncoder(file).Encode(Configuration)
}

// SaveConfiguration is used to save configurations that are provided in environment during setup tasks
// This is called when setup tasks are called
func SaveConfiguration(c csetup.Context) error {
	log.Trace("config/config:SaveConfiguration() Entering")
	defer log.Trace("config/config:SaveConfiguration() Leaving")

	var err error

	//clear the ConfigComplete flag and save the file. We will mark it complete on at the end.
	// we can use the ConfigComplete field to check if the configuration is complete before
	// running the other tasks.
	Configuration.ConfigComplete = false
	err = Save()
	if err != nil {
		return errors.Wrap(err, "Unable to save configuration file")
	}

	// we are going to check and set the required configuration variables
	// however, we do not want to error out after each one. We want to provide
	// entries in the log file indicating which ones are missing. At the
	// end of this section we will error out. Will use a flag to keep track
	missingEnvVars := []string{}

	cmsTLSCertDigest, err := c.GetenvString(consts.CmsTlsCertDigestEnv, "CMS TLS certificate SHA384 digest")
	if err == nil && cmsTLSCertDigest != "" {
		Configuration.CmsTLSCertDigest = cmsTLSCertDigest
		log.Infof("config/config.go:SaveConfiguration() %s config updated\n", consts.CmsTlsCertDigestEnv)
	} else if Configuration.CmsTLSCertDigest == "" {
		missingEnvVars = append(missingEnvVars, consts.CmsTlsCertDigestEnv)
		log.Errorf("config/config.go:SaveConfiguration() Environment variable %s required - but not set\n", consts.CmsTlsCertDigestEnv)
	}

	cmsBaseURL, err := c.GetenvString(consts.CmsBaseUrlEnv, "CMS Base URL")
	if err == nil && cmsBaseURL != "" {
		Configuration.Cms.BaseURL = cmsBaseURL
		log.Infof("config/config.go:SaveConfiguration() %s config updated\n", consts.CmsBaseUrlEnv)
	} else if Configuration.Cms.BaseURL == "" {
		missingEnvVars = append(missingEnvVars, consts.CmsBaseUrlEnv)
		log.Errorf("config/config.go:SaveConfiguration() Environment variable %s required - but not set", consts.CmsBaseUrlEnv)
	}

	kMSAPIURL, err := c.GetenvString(consts.KMSAPIURLEnv, "KMS API URL")
	if err == nil && kMSAPIURL != "" {
		Configuration.Kms.APIURL = kMSAPIURL
		log.Infof("config/config.go:SaveConfiguration() %s config updated\n", consts.KMSAPIURLEnv)
	} else if Configuration.Kms.APIURL == "" {
		missingEnvVars = append(missingEnvVars, consts.KMSAPIURLEnv)
		log.Errorf("config/config.go:SaveConfiguration() Environment variable %s required - but not set", consts.KMSAPIURLEnv)
	}

	aasAPIURLEnv, err := c.GetenvString(consts.AasAPIURLEnv, "AAS API URL")
	if err == nil && aasAPIURLEnv != "" {
		Configuration.Aas.APIURL = aasAPIURLEnv
		log.Infof("config/config.go:SaveConfiguration() %s config updated\n", consts.AasAPIURLEnv)
	} else if Configuration.Aas.APIURL == "" {
		missingEnvVars = append(missingEnvVars, consts.AasAPIURLEnv)
		log.Errorf("config/config.go:SaveConfiguration() Environment variable %s required - but not set", consts.AasAPIURLEnv)
	}

	serviceUserName, err := c.GetenvString(consts.ServiceUsername, "WPM AAS Username")
	if err == nil && serviceUserName != "" {
		Configuration.Wpm.Username = serviceUserName
		log.Infof("config/config.go:SaveConfiguration() %s config updated\n", consts.ServiceUsername)
	} else if Configuration.Wpm.Username == "" {
		missingEnvVars = append(missingEnvVars, consts.ServiceUsername)
		log.Errorf("config/config.go:SaveConfiguration() Environment variable %s required - but not set", consts.ServiceUsername)
	}

	servicePassword, err := c.GetenvSecret(consts.ServicePassword, "WPM AAS Password")
	if err == nil && servicePassword != "" {
		Configuration.Wpm.Password = servicePassword
		log.Infof("config/config.go:SaveConfiguration() %s config updated\n", consts.ServicePassword)
	} else if Configuration.Wpm.Password == "" {
		missingEnvVars = append(missingEnvVars, consts.ServicePassword)
		log.Errorf("config/config.go:SaveConfiguration() Environment variable %s required - but not set", consts.ServicePassword)
	}

	flavorSigningKeyPath, err := c.GetenvString("KEY_PATH", "Path of file where flavor signing key needs to be stored")
	if err == nil && flavorSigningKeyPath != "" {
		Configuration.FlavorSigningKeyFile = flavorSigningKeyPath
	} else if Configuration.FlavorSigningKeyFile == "" {
		Configuration.FlavorSigningKeyFile = consts.FlavorSigningKeyPath
	}

	flavorSigningCertPath, err := c.GetenvString("CERT_PATH", "Path of file/directory where flavor signing certificate needs to be stored")
	if err == nil && flavorSigningCertPath != "" {
		Configuration.FlavorSigningCertFile = flavorSigningCertPath
	} else if Configuration.FlavorSigningCertFile == "" {
		Configuration.FlavorSigningCertFile = consts.FlavorSigningCertPath
	}

	certCommonName, err := c.GetenvString(consts.WpmFlavorSignCertCommonNameEnv, "Common Name")
	if err == nil && certCommonName != "" {
		Configuration.Subject.CommonName = certCommonName
	} else if Configuration.Subject.CommonName == "" {
		log.Infof("config/config.go:SaveConfiguration() Using default value for %s\n", consts.WpmFlavorSignCertCommonNameEnv)
		Configuration.Subject.CommonName = consts.DefaultWpmFlavorSigningCn
	}

	ll, err := c.GetenvString(consts.LogLevelEnvVar, "Logging Level")
	if err != nil {
		if Configuration.LogLevel == "" {
			log.Infof("config/config:SaveConfiguration() %s not defined, using default log level: Info", consts.LogLevelEnvVar)
			Configuration.LogLevel = logrus.InfoLevel.String()
		}
	} else {
		llp, err := logrus.ParseLevel(ll)
		if err != nil {
			log.Info("config/config:SaveConfiguration() Invalid log level specified in env, using default log level: Info")
			Configuration.LogLevel = logrus.InfoLevel.String()
		} else {
			Configuration.LogLevel = llp.String()
			log.Infof("config/config:SaveConfiguration() Log level set %s\n", ll)
		}
	}

	logEntryMaxLength, err := c.GetenvInt(consts.LogEntryMaxlengthEnv, "Maximum length of each entry in a log")
	if err == nil && logEntryMaxLength >= 100 {
		Configuration.LogEntryMaxLength = logEntryMaxLength
	} else {
		log.Info("config/config:SaveConfiguration() Invalid Log Entry Max Length defined (should be > 100), " +
			"using default value")
		Configuration.LogEntryMaxLength = consts.DefaultLogEntryMaxlength
	}

	Configuration.LogEnableStdout = false
	logEnableStdout, err := c.GetenvString(consts.WPMConsoleEnableEnv, "Workload Policy Manager enable standard output")
	if err == nil && logEnableStdout != "" {
		Configuration.LogEnableStdout, err = strconv.ParseBool(logEnableStdout)
		if err != nil {
			log.Info("config/config:SaveConfiguration() Error while parsing the variable ", consts.WPMConsoleEnableEnv, " setting to default value false")
		}
	}

	if len(missingEnvVars) > 0 {
		return fmt.Errorf("Missing environment variables for setup not present: %s", strings.Join(missingEnvVars, " "))
	} else {
		Configuration.ConfigComplete = true
	}

	return Save()
}

// LogConfiguration is used to setup log configurations
func LogConfiguration(stdOut, logFile bool) error {
	log.Trace("config/config:LogConfiguration() Entering")
	defer log.Trace("config/config:LogConfiguration() Leaving")

	var ioWriterDefault io.Writer

	// creating the log file if not preset
	secLogFile, err := os.OpenFile(consts.SecLogFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		logrus.Fatal("Unable to open log file. " + err.Error())
	}
	defaultLogFile, err := os.OpenFile(consts.LogFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		logrus.Fatal("Unable to open log file. " + err.Error())
	}

	ioWriterDefault = defaultLogFile

	if stdOut && logFile {
		ioWriterDefault = io.MultiWriter(os.Stdout, defaultLogFile)
	}

	ioWriterSecurity := io.MultiWriter(ioWriterDefault, secLogFile)

	if Configuration.LogLevel == "" {
		log.Infof("config/config:SaveConfiguration() %s not defined, using default log level: Info\n", consts.LogLevelEnvVar)
		Configuration.LogLevel = logrus.InfoLevel.String()
	}

	llp, _ := logrus.ParseLevel(Configuration.LogLevel)
	commLogInt.SetLogger(commLog.DefaultLoggerName, llp, &commLog.LogFormatter{MaxLength: Configuration.LogEntryMaxLength}, ioWriterDefault, false)
	commLogInt.SetLogger(commLog.SecurityLoggerName, llp, &commLog.LogFormatter{MaxLength: Configuration.LogEntryMaxLength}, ioWriterSecurity, false)

	secLog.Trace("config/config:LogConfiguration() Security log initiated")
	log.Trace("config/config:LogConfiguration() Loggers setup finished")

	return nil
}
