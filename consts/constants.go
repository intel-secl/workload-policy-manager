package consts

const (
	KMS_API_URL                = "KMS_API_URL"
	KMS_API_USERNAME           = "KMS_API_USERNAME"
	KMS_API_PASSWORD           = "KMS_API_PASSWORD"
	KMS_TLS_SHA256             = "KMS_TLS_SHA256"
	KmsEncryptAlgo             = "AES"
	KmsKeyLength               = 256
	KmsCipherMode              = "GCM"
	OptDirPath                 = "/opt/workload-policy-manager/"
	ConfigDirPath              = "/etc/workload-policy-manager/"
	ConfigFilePath             = ConfigDirPath + "config.yml"
	LogDirPath                 = "/var/log/workload-policy-manager/"
	LogFileName                = LogDirPath + "wpm.log"
	EnvelopePublickeyLocation  = ConfigDirPath + "envelopePublicKey.pub"
	EnvelopePrivatekeyLocation = ConfigDirPath + "envelopePrivateKey.pem"
	WpmSymLink                 = "/usr/local/bin/wpm"
)
