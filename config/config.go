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
	LogLevel       string
	ConfigComplete bool
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
		return errors.Wrap(err, "config/config:SaveConfiguration() unable to save configuration file")
	}

	// we are going to check and set the required configuration variables
	// however, we do not want to error out after each one. We want to provide
	// entries in the log file indicating which ones are missing. At the
	// end of this section we will error out. Will use a flag to keep track

	requiredConfigsPresent := true

	requiredConfigs := [...]csetup.EnvVars{
		{
			Name:        consts.CmsTlsCertDigestEnv,
			ConfigVar:   &Configuration.CmsTLSCertDigest,
			Description: "CMS TLS certificate SHA384 digest",
			EmptyOkay:   false,
		},
		{
			Name:        consts.CmsBaseUrlEnv,
			ConfigVar:   &Configuration.Cms.BaseURL,
			Description: "CMS Base URL",
			EmptyOkay:   false,
		},
		{
			Name:        consts.KMSAPIURLEnv,
			ConfigVar:   &Configuration.Kms.APIURL,
			Description: "KMS URL",
			EmptyOkay:   false,
		},
		{
			Name:        consts.AasAPIURLEnv,
			ConfigVar:   &Configuration.Aas.APIURL,
			Description: "AAS API URL",
			EmptyOkay:   false,
		},
		{
			Name:        consts.ServiceUsername,
			ConfigVar:   &Configuration.Wpm.Username,
			Description: "WPM Service Username",
			EmptyOkay:   false,
		},
		{
			Name:        consts.ServicePassword,
			ConfigVar:   &Configuration.Wpm.Password,
			Description: "WPM Service Password",
			EmptyOkay:   false,
		},
		{
			Name:        consts.WpmFlavorSignCertCommonNameEnv,
			ConfigVar:   &Configuration.Subject.CommonName,
			Description: "WPM Signing Certificate Common Name",
			EmptyOkay:   true,
		},
		{
			Name:        consts.WpmCertOrganizationEnv,
			ConfigVar:   &Configuration.Subject.Organization,
			Description: "WPM Signing Certificate Organization",
			EmptyOkay:   true,
		},
		{
			Name:        consts.WpmCertCountryEnv,
			ConfigVar:   &Configuration.Subject.Country,
			Description: "WPM Signing Certificate Country",
			EmptyOkay:   true,
		},
		{
			Name:        consts.WpmCertLocalityEnv,
			ConfigVar:   &Configuration.Subject.Locality,
			Description: "WPM Signing Certificate Locality",
			EmptyOkay:   true,
		},
		{
			Name:        consts.LogLevelEnvVar,
			ConfigVar:   &Configuration.LogLevel,
			Description: "WPM Logging Level",
			EmptyOkay:   true,
		},
	}

	for _, cv := range requiredConfigs {
		_, _, err = c.OverrideValueFromEnvVar(cv.Name, cv.ConfigVar, cv.Description, cv.EmptyOkay)
		if err != nil {
			requiredConfigsPresent = false
			fmt.Fprintf(os.Stderr, "Environment variable %s required - but not set\n", cv.Name)
			fmt.Fprintln(os.Stderr, err)
		}
	}

	certCommonName, err := c.GetenvString(consts.WpmFlavorSignCertCommonNameEnv, "Common Name")
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

	// check if the logging level is provided in env
	logLevel, err := c.GetenvString(consts.LogLevelEnvVar, "Logging Level")
	if err != nil && logLevel == "" {
		fmt.Fprintln(os.Stderr, "No logging level specified, using default logging level: Info")
		Configuration.LogLevel = logrus.InfoLevel.String()
	} else {
		// check if the provided log level is valid
		_, err = logrus.ParseLevel(logLevel)
		if err != nil {
			// fall back to the default logging level
			fmt.Fprintln(os.Stderr, "Invalid logging level specified, using default logging level: Info")
			Configuration.LogLevel = logrus.InfoLevel.String()
		} else {
			if Configuration.LogLevel == "" {
				// update the logging level
				ll, _ := logrus.ParseLevel(logLevel)
				Configuration.LogLevel = ll.String()
			}
		}
	}

	if requiredConfigsPresent {
		Configuration.ConfigComplete = true
		return Save()
	}

	return errors.New("config/config one or more required environment variables for setup not present. log file has details")
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

	ll, _ := logrus.ParseLevel(Configuration.LogLevel)

	commLogInt.SetLogger(commLog.DefaultLoggerName, ll, nil, ioWriterDefault, false)
	commLogInt.SetLogger(commLog.SecurityLoggerName, ll, nil, ioWriterSecurity, false)
	secLog.Trace("config/config:LogConfiguration() Security log initiated")
	log.Trace("config/config:LogConfiguration() Loggers setup finished")

	return nil
}
