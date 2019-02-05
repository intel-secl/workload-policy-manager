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

func (ek CreateEnvelopeKey) Run(c csetup.Context) error {

	/*if ek.Validate(c) != nil {
		fmt.Println("Envelope key already created. Skipping this setup task.")
		return nil
	}*/
	// save configuration from config.yml
	e := config.SaveConfiguration(c)
	if e != nil {
		log.Error(e.Error())
		return e
	}
	bitSize := 2048
	keyPair, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return errors.New("Error while generating a new RSA key pair")
	}

	// save private key
	privateKey := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(keyPair),
	}

	privateKeyFile, err := os.Create(consts.EnvelopePrivatekeyLocation)
	if err != nil {
		return errors.New("Error while creating a new file")
	}
	defer privateKeyFile.Close()
	err = pem.Encode(privateKeyFile, privateKey)
	if err != nil {
		return errors.New("Error while encoding the private key")
	}

	// save public key
	publicKey := &keyPair.PublicKey

	pubkeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return errors.New("Error while getting the public key from private key")
	}

	var publicKeyInPem = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubkeyBytes,
	}

	publicKeyFile, err := os.Create(consts.EnvelopePublickeyLocation)
	if err != nil {
		return errors.New("Error while creating a new file")
	}
	defer publicKeyFile.Close()

	err = pem.Encode(publicKeyFile, publicKeyInPem)
	if err != nil {
		return errors.New("Error while encoding the public key")
	}
	return nil
}

// ValidateCreateKey method is used to check if the envelope keys exists on disk
func (ek CreateEnvelopeKey) Validate(c csetup.Context) error {
	_, err := os.Stat(consts.EnvelopePrivatekeyLocation)
	if os.IsNotExist(err) {
		return errors.New("Private key exists")
	}

	_, err = os.Stat(consts.EnvelopePublickeyLocation)
	if os.IsNotExist(err) {
		return errors.New("Public key exists")
	}
	return nil
}
