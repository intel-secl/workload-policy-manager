/*
 * Copyright (C) 2019 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package consts

const (
	KMSAPIURLEnv                   = "KMS_API_URL"
	KmsEncryptAlgo                 = "AES"
	KmsKeyLength                   = 256
	KmsCipherMode                  = "GCM"
	CmsBaseUrlEnv                  = "CMS_BASE_URL"
	WpmFlavorSignCertCommonNameEnv = "WPM_FLAVOR_SIGN_CERT_CN"
	LogLevelEnvVar                 = "WPM_LOG_LEVEL"
	OptDirPath                     = "/opt/workload-policy-manager/"
	ConfigDirPath                  = "/etc/workload-policy-manager/"
	ConfigFilePath                 = ConfigDirPath + "config.yml"
	FlavorSigningCertPath          = ConfigDirPath + "flavor-signing-cert.pem"
	FlavorSigningKeyPath           = ConfigDirPath + "flavor-signing-key.pem"
	TrustedCaCertsDir              = ConfigDirPath + "cacerts/"
	TrustedJWTSigningCertsDir      = ConfigDirPath + "jwt/"
	DefaultKeyAlgorithm            = "rsa"
	DefaultKeyAlgorithmLength      = 3072
	CertApproverGroupName          = "CertApprover"
	DefaultWpmFlavorSigningCn      = "WPM Flavor Signing Certificate"
	DefaultWpmSan                  = "127.0.0.1,localhost"
	LogDirPath                     = "/var/log/workload-policy-manager/"
	LogEntryMaxlengthEnv           = "LOG_ENTRY_MAXLENGTH"
	DefaultLogEntryMaxlength       = 300
	LogFileName                    = LogDirPath + "wpm.log"
	SecLogFileName                 = LogDirPath + "wpm_security.log"
	EnvelopePublickeyLocation      = ConfigDirPath + "envelopePublicKey.pub"
	EnvelopePrivatekeyLocation     = ConfigDirPath + "envelopePrivateKey.pem"
	WpmSymLink                     = "/usr/local/bin/wpm"
	AasAPIURLEnv                   = "AAS_API_URL"
	BearerTokenEnv                 = "BEARER_TOKEN"
	KMSKeyRetrievalGroupName       = "KeyCRUD"
	ServiceName                    = "WPM"
	ServiceUsername                = "WPM_SERVICE_USERNAME"
	ServicePassword                = "WPM_SERVICE_PASSWORD"
	CmsTlsCertDigestEnv            = "CMS_TLS_CERT_SHA384"
	WPMConsoleEnableEnv            = "WPM_ENABLE_CONSOLE_LOG"
	SampleUUID                     = "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
	WpmNosetupEnv                  = "WPM_NOSETUP"
)
