/*
 * Copyright (C) 2019 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package consts

const (
	KMSAPIURLEnv                   = "KMS_API_URL"
	KMSUsernameEnv                 = "KMS_USERNAME"
	KMSPasswordEnv                 = "KMS_PASSWORD"
	KmsEncryptAlgo                 = "AES"
	KmsKeyLength                   = 256
	KmsCipherMode                  = "GCM"
	CmsBaseUrlEnv                  = "CMS_BASE_URL"
	WpmFlavorSignCertCommonNameEnv = "WPM_FLAVOR_SIGN_CERT_CN"
	WpmCertOrganizationEnv         = "WPM_CERT_ORG"
	WpmCertCountryEnv              = "WPM_CERT_COUNTRY"
	WpmCertLocalityEnv             = "WPM_CERT_LOCALITY"
	WpmCertProvinceEnv             = "WPM_CERT_PROVINCE"
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
	DefaultWpmOrganization         = "INTEL"
	DefaultWpmCountry              = "US"
	DefaultWpmProvince             = "SF"
	DefaultWpmLocality             = "SC"
	DefaultWpmSan                  = "127.0.0.1,localhost"
	LogDirPath                     = "/var/log/workload-policy-manager/"
	LogFileName                    = LogDirPath + "wpm.log"
	SecLogFileName                 = LogDirPath + "wpm_security.log"
	EnvelopePublickeyLocation      = ConfigDirPath + "envelopePublicKey.pub"
	EnvelopePrivatekeyLocation     = ConfigDirPath + "envelopePrivateKey.pem"
	WpmSymLink                     = "/usr/local/bin/wpm"
	AasAPIURLEnv                   = "AAS_API_URL"
	BearerTokenEnv                 = "BEARER_TOKEN"
	KMSKeyRetrievalGroupName       = "KeyCRUD"
	ServiceName                    = "WPM"
	ServiceUsername                = "WPM_USERNAME"
	ServicePassword                = "WPM_PASSWORD"
	CmsTlsCertDigestEnv            = "CMS_TLS_CERT_SHA384"
)
