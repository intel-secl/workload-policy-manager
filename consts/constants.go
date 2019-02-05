package consts

const (
	KMS_API_URL                = "KMS_API_URL"
	KMS_API_USERNAME           = "KMS_API_USERNAME"
	KMS_API_PASSWORD           = "KMS_API_PASSWORD"
	KMS_TLS_SHA256             = "KMS_TLS_SHA256"
	KMS_ENCRYPTION_ALG         = "AES"
	KMS_KEY_LENGTH             = 256
	KMS_CIPHER_MODE            = "GCM"
	ConfigFilePath             = "/etc/wpm/configuration/config.yml"
	LogDirPath                 = "/var/log/wpm/"
	LogFileName                = "wpm.log"
	EnvelopePublickeyLocation  = "/etc/wpm/configuration/envelopePublicKey.pub"
	EnvelopePrivatekeyLocation = "/etc/wpm/configuration/envelopePrivateKey.pem"
)
