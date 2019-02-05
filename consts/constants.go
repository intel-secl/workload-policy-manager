package consts

const (
	KMS_API_URL                = "KMS_API_URL"
	KMS_API_USERNAME           = "KMS_API_USERNAME"
	KMS_API_PASSWORD           = "KMS_API_PASSWORD"
	KMS_TLS_SHA256             = "KMS_TLS_SHA256"
	KMS_ENCRYPTION_ALG         = "AES"
	KMS_KEY_LENGTH             = 256
	KMS_CIPHER_MODE            = "GCM"
	ConfigFilePath             = "/etc/wpm/config.yml"
	ConfigDirPath              = "/etc/wpm/"
	LogDirPath                 = "/var/log/wpm/"
	LogFileName                = "wpm.log"
	EnvelopePublickeyLocation  = "/etc/wpm/envelopePublicKey.pub"
	EnvelopePrivatekeyLocation = "/etc/wpm/envelopePrivateKey.pem"
)
