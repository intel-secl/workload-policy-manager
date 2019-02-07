package kmsclient

import (
	"crypto/tls"
	"encoding/hex"
	"errors"
	t "intel/isecl/lib/common/tls"
	kms "intel/isecl/lib/kms-client"
	config "intel/isecl/wpm/config"
	"net/http"
)

func InitializeClient() (*kms.Client, error) {
	var kc *kms.Client
	var certificateDigest [32]byte
	certDigestHex, err := hex.DecodeString(config.Configuration.Kms.TLSSha256)
	if err != nil {
		return kc, errors.New("error converting certificate digest to hex. " + err.Error())
	}
	copy(certificateDigest[:], certDigestHex)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify:    true,
				VerifyPeerCertificate: t.VerifyCertBySha256(certificateDigest),
			},
		},
	}
	kc = &kms.Client{
		BaseURL:    config.Configuration.Kms.APIURL,
		Username:   config.Configuration.Kms.APIUsername,
		Password:   config.Configuration.Kms.APIPassword,
		CertSha256: &certificateDigest,
		HTTPClient: client,
	}
	return kc, err
}
