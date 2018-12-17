package setup

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
)

func Validate() bool{
	isPubKeyExists bool := false
	isPriKeyExists bool := false
	privateKey := "/opt/wpm/configuration/envelopePrivateKey.pem"
	publicKey := "/opt/wpm/configuration/envelopePublicKey.pub"

	_, err = os.Stat(privateKey)
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

func CreateEnvelopeKey() {
	savePrivateFileTo := "/opt/wpm/configuration/envelopeKey.pem"
	savePublicFileTo := "/opt/wpm/configuration/envelopeKey.pub"
	bitSize := 2048

	keyPair, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		log.Fatal(err)
	}

	// save private key
	privateKey := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(keyPair),
	}

	privateKeyFile, err := os.Create(savePrivateFileTo)
	if err != nil {
		log.Fatal(err)
	}
	defer privateKeyFile.Close()
	err = pem.Encode(privateKeyFile, privateKey)
	if err != nil {
		log.Fatal(err)
	}

	// save public key
	publicKey := &keyPair.PublicKey

	pubkeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		log.Fatal(err)
	}
	
	var publicKeyInPem = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubkeyBytes,
	}

	publicKeyFile, err := os.Create(savePublicFileTo)
	if err != nil {
		log.Fatal(err)
	}
	defer publicKeyFile.Close()

	err = pem.Encode(publicKeyFile, publicKeyInPem)
	if err != nil {
		log.Fatal(err)
	}
}
