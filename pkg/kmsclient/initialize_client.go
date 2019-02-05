package kmsclient

import (
	log "github.com/sirupsen/logrus"
	"encoding/hex"
	kms "intel/isecl/lib/kms-client"
	config "intel/isecl/wpm/config"
	"net/http"
	"crypto/tls"
	t "intel/isecl/lib/common/tls"
)

func InitializeClient() *kms.Client {
	var certificateDigest [32]byte
	certDigestHex, err := hex.DecodeString(config.Configuration.Kms.TLSSha256)
	if err != nil {
		log.Error("Error converting certificate digest to hex")
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
	kc := &kms.Client{
		BaseURL:    config.Configuration.Kms.APIURL,
		Username:   config.Configuration.Kms.APIUsername,
		Password:   config.Configuration.Kms.APIPassword,
		CertSha256: &certificateDigest,
		HTTPClient: client,
	}
	return kc
}
