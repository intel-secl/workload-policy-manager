/*
 * Copyright (C) 2019 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package config

import (
	"errors"
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
		TLSSha384   string
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
	var err error

	kmsApiUrl, err := c.GetenvString(consts.KMS_API_URL, "Kms URL")
	if err == nil && kmsApiUrl != "" {
		Configuration.Kms.APIURL = kmsApiUrl
	} else if Configuration.Kms.APIURL == "" {
		return errors.New("KMS API URL is not defined in environment or config file")
	}

	kmsApiUsername, err := c.GetenvString(consts.KMS_API_USERNAME, "Kms Username")
	if err == nil && kmsApiUsername!= "" {
		Configuration.Kms.APIUsername = kmsApiUsername
	} else if Configuration.Kms.APIUsername == "" {
		return errors.New("KMS API Username is not defined in environment or config file")
	}

	kmsApiPassword, err := c.GetenvString(consts.KMS_API_PASSWORD, "Kms Password")
	if err == nil && kmsApiPassword != ""{
		Configuration.Kms.APIPassword = kmsApiPassword
	} else if Configuration.Kms.APIPassword == "" {
		return errors.New("KMS API Password is not defined in environment or config file")
	}

	kmsTlsSha384, err := c.GetenvString(consts.KMS_TLS_SHA384, "Kms TLS SHA384")
	if err == nil && kmsTlsSha384 != ""{
		Configuration.Kms.TLSSha384 = kmsTlsSha384
	} else if Configuration.Kms.TLSSha384 == "" {
		return errors.New("KMS TLS is not defined in environment or config file")
	}

	cmsBaseUrl, err := c.GetenvString(consts.CmsBaseUrlEnv, "CMS Base URL")
	if err == nil && cmsBaseUrl != "" {
		Configuration.Cms.BaseUrl = cmsBaseUrl
	} else if Configuration.Cms.BaseUrl == "" {
		return errors.New("CMS Base URL is not defined in environment or config file")
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
	return Save()
}

// LogConfiguration is used to setup log configurations
func LogConfiguration() error {
	// creating the log file if not preset
	logFile, err := os.OpenFile(consts.LogFileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModeAppend)
	if err != nil {
		return errors.New("unable to write file. " + err.Error())
	}
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, TimestampFormat: time.RFC1123Z})
	logMultiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(logMultiWriter)

	return nil
}
