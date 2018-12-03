package wpm

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/pem"
	"io"
	"io/ioutil"
	"log"
)

func encrypt(data []byte, keyLocation string) []byte {
	// reading the key file in Pem format
	keyFile, err := ioutil.ReadFile(keyLocation)
	if err != nil {
		log.Fatal("Error reading the key file", err)
	}
	// decoding the key file
	key, _ := pem.Decode(keyFile)
	// creating a new cipher block of 128 bits
	block, _ := aes.NewCipher(key.Bytes)
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatal("Error creating a cipher block", err)
	}
	// assigning a 12 byte empty array to store random value
	iv := make([]byte, 12)
	// reading random value into the byte array
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		log.Fatal("Error creating random IV value", err)
	}

	// encrypting the file(data) with IV value and appending
	// the IV to the first 12 bytes of the encrypted file
	ciphertext := gcm.Seal(iv, iv, data, nil)
	return ciphertext
}
