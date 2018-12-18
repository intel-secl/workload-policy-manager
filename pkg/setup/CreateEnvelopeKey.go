package setup

import (
	"intel/isecl/wpm/config"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

// Validate method is used to check if the envelope keys exists on disk
func Validate() bool {
	isPubKeyExists := false
	isPriKeyExists := false
	privateKey := "/opt/wpm/configuration/envelopePrivateKey.pem"
	publicKey := "/opt/wpm/configuration/envelopePublicKey.pub"

	_, err := os.Stat(privateKey)
	if !os.IsNotExist(err) {
		isPriKeyExists = true
	}

	_, err = os.Stat(publicKey)
	if !os.IsNotExist(err) {
		isPubKeyExists = true
	}

	if isPriKeyExists && isPubKeyExists {
		return false
	} else {
		return true
	}

}

// CreateEnvelopeKey method is used t create the envelope key
func CreateEnvelopeKey() error {
	savePrivateFileTo := "/opt/wpm/configuration/envelopePrivateKey.pem"
	savePublicFileTo := "/opt/wpm/configuration/envelopePublicKey.pub"
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

	privateKeyFile, err := os.Create(savePrivateFileTo)
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

	publicKeyFile, err := os.Create(savePublicFileTo)
	if err != nil {
		return errors.New("Error while creating a new file")
	}
	defer publicKeyFile.Close()

	err = pem.Encode(publicKeyFile, publicKeyInPem)
	if err != nil {
		return errors.New("Error while encoding the public key")
	}

	config.Configuration.EnvelopePrivatekeyLocation = savePrivateFileTo
	config.Configuration.EnvelopePublickeyLocation = savePublicFileTo
	
	return nil
}
