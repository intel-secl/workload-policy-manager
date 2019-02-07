package setup

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	log "github.com/sirupsen/logrus"
	csetup "intel/isecl/lib/common/setup"
	"intel/isecl/wpm/config"
	"intel/isecl/wpm/consts"
	"os"
)

type CreateEnvelopeKey struct {
}

// ValidateCreateKey method is used to check if the envelope keys exists on disk
func (ek CreateEnvelopeKey) Validate(c csetup.Context) error {
	log.Info("Validating creating envelope key")

	_, err := os.Stat(consts.EnvelopePrivatekeyLocation)
	if os.IsNotExist(err) {
		return errors.New("private key does not exist")
	}

	_, err = os.Stat(consts.EnvelopePublickeyLocation)
	if os.IsNotExist(err) {
		return errors.New("public key does not exist")
	}
	return nil
}

func (ek CreateEnvelopeKey) Run(c csetup.Context) error {
	log.Info("Creating envelope key")

	// save configuration from config.yml
	e := config.SaveConfiguration(c)
	if e != nil {
		return errors.New("error saving configuration. " + e.Error())
	}
	bitSize := 2048
	keyPair, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return errors.New("error while generating a new RSA key pair")
	}

	// save private key
	privateKey := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(keyPair),
	}

	privateKeyFile, err := os.Create(consts.EnvelopePrivatekeyLocation)
	if err != nil {
		return errors.New("error while creating a new file. " + err.Error())
	}
	defer privateKeyFile.Close()
	err = pem.Encode(privateKeyFile, privateKey)
	if err != nil {
		return errors.New("error while encoding the private key. " + err.Error())
	}

	// save public key
	publicKey := &keyPair.PublicKey

	pubkeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return errors.New("error while marshalling the public key. " + err.Error())
	}
	var publicKeyInPem = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubkeyBytes,
	}

	publicKeyFile, err := os.Create(consts.EnvelopePublickeyLocation)
	if err != nil {
		return errors.New("error while creating a new file. " + err.Error())
	}
	defer publicKeyFile.Close()

	err = pem.Encode(publicKeyFile, publicKeyInPem)
	if err != nil {
		return errors.New("error while encoding the public key. " + err.Error())
	}
	return nil
}
