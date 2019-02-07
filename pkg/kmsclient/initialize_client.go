package kmsclient

import (
	"encoding/hex"
	"errors"
	kms "intel/isecl/lib/kms-client"
	config "intel/isecl/wpm/config"
)

func InitializeClient() (*kms.Client, error) {
	var kc *kms.Client
	var certificateDigest [32]byte
	certDigestHex, err := hex.DecodeString(config.Configuration.Kms.TLSSha256)
	if err != nil {
		return kc, errors.New("error converting certificate digest to hex. " + err.Error())
	}
	copy(certificateDigest[:], certDigestHex)
	kc = &kms.Client{
		BaseURL:    config.Configuration.Kms.APIURL,
		Username:   config.Configuration.Kms.APIUsername,
		Password:   config.Configuration.Kms.APIPassword,
		CertSha256: &certificateDigest,
	}
	return kc, err
}
