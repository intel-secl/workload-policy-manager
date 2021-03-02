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
	TrustedCaCertsDir              = ConfigDirPath + "certs/trustedca/"
	FlavorSigningCertPath          = ConfigDirPath + "certs/flavorsign/flavor-signing-cert.pem"
	FlavorSigningKeyPath           = ConfigDirPath + "certs/flavorsign/flavor-signing-key.pem"
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
	EnvelopePublickeyLocation      = ConfigDirPath + "certs/kbs/envelopePublicKey.pub"
	EnvelopePrivatekeyLocation     = ConfigDirPath + "certs/kbs/envelopePrivateKey.pem"
	WpmSymLink                     = "/usr/local/bin/wpm"
	AasAPIURLEnv                   = "AAS_API_URL"
	BearerTokenEnv                 = "BEARER_TOKEN"
	KMSKeyRetrievalGroupName       = "KeyCRUD"
	ServiceName                    = "WPM"
	ExplicitServiceName            = "Workload Policy Manager"
	ServiceUsername                = "WPM_SERVICE_USERNAME"
	ServicePassword                = "WPM_SERVICE_PASSWORD"
	CmsTlsCertDigestEnv            = "CMS_TLS_CERT_SHA384"
	WPMConsoleEnableEnv            = "WPM_ENABLE_CONSOLE_LOG"
	SampleUUID                     = "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
	WpmNosetupEnv                  = "WPM_NOSETUP"
)
