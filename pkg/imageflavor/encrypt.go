package imageflavor

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
	"io/ioutil"
)

func encrypt(imagePath string, privateKeyLocation string, encryptedFileLocation string, wrappedKey []byte) error {

	// reading image file
	image, err := ioutil.ReadFile(imagePath)
	if err != nil {
		return errors.New("Error reading the image file")
	}

	key, err := unwrapKey(wrappedKey, privateKeyLocation)
	if err != nil {
		return errors.New("Error while unwrapping the key")
	}
	// creating a new cipher block of 128 bits
	block, _ := aes.NewCipher(key)
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return errors.New("Error creating a cipher block")
	}
	// assigning a 12 byte empty array to store random value
	iv := make([]byte, 12)
	// reading random value into the byte array
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return errors.New("Error creating random IV value")
	}

	// encrypting the file(data) with IV value and appending
	// the IV to the first 12 bytes of the encrypted file
	ciphertext := gcm.Seal(iv, iv, image, nil)
	err = ioutil.WriteFile(encryptedFileLocation, ciphertext, 0600)
	if err != nil {
		return errors.New("Error during writing the encrypted image to file")
	}

	return nil
}

func unwrapKey(wrappedKey []byte, privateKeyLocation string) ([]byte, error) {

	var unwrappedKey []byte
	privateKey, err := ioutil.ReadFile(privateKeyLocation)
	if err != nil {
		return unwrappedKey, errors.New("Error reading the private key")
	}

	privateKeyBlock, _ := pem.Decode(privateKey)
	var pri *rsa.PrivateKey
	pri, err = x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		return unwrappedKey, errors.New("Error while parsing the private key")
	}

	decryptedKey, errDecrypt := rsa.DecryptOAEP(sha256.New(), rand.Reader, pri, wrappedKey, nil)
	if errDecrypt != nil {
		return unwrappedKey, errors.New("DecryptOEAP : Error while unwraping the key")
	}

	return decryptedKey, nil
}
