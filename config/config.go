/*
 * Copyright (C) 2019 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package config

import (
	"fmt"
	csetup "intel/isecl/lib/common/setup"
	"intel/isecl/wpm/consts"
	"io"
	"os"
	"strings"

	"github.com/pkg/errors"

	commLog "intel/isecl/lib/common/log"
	commLogInt "intel/isecl/lib/common/log/setup"

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
		CommonName   string
		Organization string
		Locality     string
		Province     string
		Country      string
	}
	Aas struct {
		APIURL string
	}
	Wpm struct {
		Username string
		Password string
	}
	LogLevel          logrus.Level
	LogEntryMaxLength int
	ConfigComplete    bool
}

var LogWriter io.Writer

func init() {
	// load from config
	file, err := os.Open(consts.ConfigFilePath)
	if err == nil {
		defer file.Close()
		yaml.NewDecoder(file).Decode(&Configuration)
	}
	LogWriter = os.Stdout
}

// Save the configuration struct into configuration directory
func Save() error {
	log.Trace("config/config:Save() Entering")
	defer log.Trace("config/config:Save() Leaving")

	file, err := os.OpenFile(consts.ConfigFilePath, os.O_RDWR, 0)
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

	servicePassword, err := c.GetenvString(consts.ServicePassword, "WPM AAS Password")
	if err == nil && servicePassword != "" {
		Configuration.Wpm.Password = servicePassword
		log.Infof("config/config.go:SaveConfiguration() %s config updated\n", consts.ServicePassword)
	} else if Configuration.Wpm.Password == "" {
		missingEnvVars = append(missingEnvVars, consts.ServicePassword)
		log.Errorf("config/config.go:SaveConfiguration() Environment variable %s required - but not set", consts.ServicePassword)
	}

	certCommonName, err := c.GetenvString(consts.WpmFlavorSignCertCommonNameEnv, "Common Name")
	if err == nil && certCommonName != "" {
		Configuration.Subject.CommonName = certCommonName
	} else if Configuration.Subject.CommonName == "" {
		log.Infof("config/config.go:SaveConfiguration() Using default value for %s\n", consts.WpmFlavorSignCertCommonNameEnv)
		Configuration.Subject.CommonName = consts.DefaultWpmFlavorSigningCn
	}

	certOrg, err := c.GetenvString(consts.WpmCertOrganizationEnv, "Organization")
	if err == nil && certOrg != "" {
		Configuration.Subject.Organization = certOrg
	} else if Configuration.Subject.Organization == "" {
		Configuration.Subject.Organization = consts.DefaultWpmOrganization
		log.Infof("config/config.go:SaveConfiguration() Using default value for %s\n", consts.WpmCertOrganizationEnv)
	}

	certCountry, err := c.GetenvString(consts.WpmCertCountryEnv, "Country")
	if err == nil && certCountry != "" {
		Configuration.Subject.Country = certCountry
	} else if Configuration.Subject.Country == "" {
		Configuration.Subject.Country = consts.DefaultWpmCountry
		log.Infof("config/config.go:SaveConfiguration() Using default value for %s\n", consts.WpmCertCountryEnv)
	}

	certProvince, err := c.GetenvString(consts.WpmCertProvinceEnv, "Province")
	if err == nil && certProvince != "" {
		Configuration.Subject.Province = certProvince
	} else if Configuration.Subject.Province == "" {
		Configuration.Subject.Province = consts.DefaultWpmProvince
		log.Infof("config/config.go:SaveConfiguration() Using default value for %s\n", consts.WpmCertProvinceEnv)
	}

	certLocality, err := c.GetenvString(consts.WpmCertLocalityEnv, "Locality")
	if err == nil && certLocality != "" {
		Configuration.Subject.Locality = certLocality
	} else if Configuration.Subject.Locality == "" {
		Configuration.Subject.Locality = consts.DefaultWpmLocality
		log.Infof("config/config.go:SaveConfiguration() Using default value for %s\n", consts.WpmCertLocalityEnv)
	}

	ll, err := c.GetenvString(consts.LogLevelEnvVar, "Logging Level")
	if err != nil {
		if Configuration.LogLevel.String() == "" {
			log.Infof("config/config:SaveConfiguration() %s not defined, using default log level: Info", consts.LogLevelEnvVar)
			Configuration.LogLevel = logrus.InfoLevel
		}
	} else {
		Configuration.LogLevel, err = logrus.ParseLevel(ll)
		if err != nil {
			log.Info("config/config:SaveConfiguration() Invalid log level specified in env, using default log level: Info")
			Configuration.LogLevel = logrus.InfoLevel
		} else {
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
	secLogFile, err := os.OpenFile(consts.SecLogFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		logrus.Fatal("Unable to open log file. " + err.Error())
	}
	defaultLogFile, err := os.OpenFile(consts.LogFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		logrus.Fatal("Unable to open log file. " + err.Error())
	}

	ioWriterDefault = defaultLogFile

	if stdOut && logFile {
		ioWriterDefault = io.MultiWriter(os.Stdout, defaultLogFile)
	}

	ioWriterSecurity := io.MultiWriter(ioWriterDefault, secLogFile)

	if Configuration.LogLevel.String() == "" {
		log.Infof("config/config:SaveConfiguration() %s not defined, using default log level: Info\n", consts.LogLevelEnvVar)
		Configuration.LogLevel = logrus.InfoLevel
	}

	commLogInt.SetLogger(commLog.DefaultLoggerName, Configuration.LogLevel, &commLog.LogFormatter{MaxLength: Configuration.LogEntryMaxLength}, ioWriterDefault, false)
	commLogInt.SetLogger(commLog.SecurityLoggerName, Configuration.LogLevel, &commLog.LogFormatter{MaxLength: Configuration.LogEntryMaxLength}, ioWriterSecurity, false)

	secLog.Trace("config/config:LogConfiguration() Security log initiated")
	log.Trace("config/config:LogConfiguration() Loggers setup finished")

	return nil
}
