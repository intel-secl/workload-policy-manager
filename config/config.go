/*
 * Copyright (C) 2019 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package config

import (
	"errors"
	"fmt"
	csetup "intel/isecl/lib/common/setup"
	"intel/isecl/wpm/consts"
	"io"
	"os"

	commLog "intel/isecl/lib/common/log"
	commLogInt "intel/isecl/lib/common/log/setup"

	"github.com/sirupsen/logrus"

	yaml "gopkg.in/yaml.v2"
)

var (
	log    = commLog.GetDefaultLogger()
	secLog = commLog.GetSecurityLogger()
)

var Configuration struct {
	Kms struct {
		APIURL   string
		Username string
		Password string
	}
	Cms struct {
		BaseUrl string
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
	LogLevel logrus.Level
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

	kmsApiUrl, err := c.GetenvString(consts.KMSAPIURLEnv, "KMS URL")
	if err == nil && kmsApiUrl != "" {
		Configuration.Kms.APIURL = kmsApiUrl
	} else if Configuration.Kms.APIURL == "" {
		return errors.New("KMS API URL is not defined in environment or config file")
	}

	kmsUsername, err := c.GetenvString(consts.KMSUsernameEnv, "KMS Username")
	if err == nil && kmsUsername != "" {
		Configuration.Kms.Username = kmsUsername
	} else if Configuration.Kms.Username == "" {
		return errors.New("KMS Username is not defined in environment or config file")
	}

	kmsPassword, err := c.GetenvString(consts.KMSPasswordEnv, "KMS Password")
	if err == nil && kmsPassword != "" {
		Configuration.Kms.Password = kmsPassword
	} else if Configuration.Kms.Password == "" {
		return errors.New("KMS Password is not defined in environment or config file")
	}

	cmsBaseUrl, err := c.GetenvString(consts.CmsBaseUrlEnv, "CMS Base URL")
	if err == nil && cmsBaseUrl != "" {
		Configuration.Cms.BaseUrl = cmsBaseUrl
	} else if Configuration.Cms.BaseUrl == "" {
		return errors.New("CMS Base URL is not defined in environment or config file")
	}

	aasAPIURL, err := c.GetenvString(consts.AasAPIURLEnv, "AAS API URL")
	if err == nil && aasAPIURL != "" {
		Configuration.Aas.APIURL = aasAPIURL
	} else if Configuration.Aas.APIURL == "" {
		return errors.New("AAS API URL is not defined in environment or config file")
	}

	wpmAASUsername, err := c.GetenvString(consts.ServiceUsername, "AAS API Username")
	if err == nil && wpmAASUsername != "" {
		Configuration.Wpm.Username = wpmAASUsername
	} else if Configuration.Wpm.Username == "" {
		return errors.New("WPM AAS Username is not defined in environment or config file")
	}

	wpmAASPassword, err := c.GetenvString(consts.ServicePassword, "AAS API Password")
	if err == nil && wpmAASPassword != "" {
		Configuration.Wpm.Password = wpmAASPassword
	} else if Configuration.Wpm.Password == "" {
		return errors.New("WPM AAS Password is not defined in environment or config file")
	}

	certCommonName, err := c.GetenvString(consts.WpmFlavorSignCertCommonNameEnv, "Common name")
	if err == nil && certCommonName != "" {
		Configuration.Subject.CommonName = certCommonName
	} else if Configuration.Subject.CommonName == "" {
		Configuration.Subject.CommonName = consts.DefaultWpmFlavorSigningCn
	}

	certOrg, err := c.GetenvString(consts.WpmCertOrganizationEnv, "Organization")
	if err == nil && certOrg != "" {
		Configuration.Subject.Organization = certOrg
	} else if Configuration.Subject.Organization == "" {
		Configuration.Subject.Organization = consts.DefaultWpmOrganization
	}

	certCountry, err := c.GetenvString(consts.WpmCertCountryEnv, "Country")
	if err == nil && certCountry != "" {
		Configuration.Subject.Country = certCountry
	} else if Configuration.Subject.Country == "" {
		Configuration.Subject.Country = consts.DefaultWpmCountry
	}

	certProvince, err := c.GetenvString(consts.WpmCertProvinceEnv, "Province")
	if err == nil && certProvince != "" {
		Configuration.Subject.Province = certProvince
	} else if Configuration.Subject.Province == "" {
		Configuration.Subject.Province = consts.DefaultWpmProvince
	}

	certLocality, err := c.GetenvString(consts.WpmCertLocalityEnv, "Locality")
	if err == nil && certLocality != "" {
		Configuration.Subject.Locality = certLocality
	} else if Configuration.Subject.Locality == "" {
		Configuration.Subject.Locality = consts.DefaultWpmLocality
	}

	logLevel, err := c.GetenvString(consts.LogLevelEnvVar, "Logging Level")
	if err != nil && logLevel == "" {
		fmt.Fprintln(os.Stderr, "No logging level specified, using default logging level: Error")
		Configuration.LogLevel = logrus.ErrorLevel
	}
	Configuration.LogLevel, err = logrus.ParseLevel(logLevel)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Invalid logging level specified, using default logging level: Error")
		Configuration.LogLevel = logrus.ErrorLevel
	}

	return Save()
}

// LogConfiguration is used to setup log configurations
func LogConfiguration() error {
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

	ioWriterSecurity := io.MultiWriter(ioWriterDefault, secLogFile)

	commLogInt.SetLogger(commLog.DefaultLoggerName, Configuration.LogLevel, nil, ioWriterDefault, false)
	commLogInt.SetLogger(commLog.SecurityLoggerName, Configuration.LogLevel, nil, ioWriterSecurity, false)
	secLog.Trace("config/config:LogConfiguration() Security log initiated")
	log.Trace("config/config:LogConfiguration() Loggers setup finished")

	return nil
}
