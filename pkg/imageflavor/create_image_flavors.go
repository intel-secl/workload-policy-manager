package imageflavor

/*
 *
 * @author srege
 * @author hmgowda
 */

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	flavor "intel/isecl/lib/flavor"
	c "intel/isecl/wpm/config"
	kms "intel/isecl/wpm/pkg/kmsclient"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

//KeyInfo is a representation of key information
type KeyInfo struct {
	KeyID           string `json:"id"`
	CipherMode      string `json:"mode"`
	Algorithm       string `json:"algorithm"`
	KeyLength       string `json:"key_length"`
	PaddingMode     string `json:"padding_mode"`
	TransferPolicy  string `json:"transfer_policy"`
	TransferLink    string `json:"transfer_link"`
	DigestAlgorithm string `json:"digest_algorithm"`
}

//KeyObj is a represenation of the actual key
type KeyObj struct {
	Key string `json:"key"`
}

//CreateImageFlavor is used to create flavor of an encrypted image
func CreateImageFlavor(imagePath string, encryptFilePath string, keyID string, encryptionRequired bool, integrityEnforced bool, outputFile string) (string, error) {
	var err error
	var key string
	var keyURL string

	//input validation
	if len(strings.TrimSpace(imagePath)) <= 0 {
		log.Fatal("image path not given")
	}

	if len(strings.TrimSpace(encryptFilePath)) <= 0 {
		log.Fatal("encryption file path not given")
	}

	// check if image exists at the specified location
	_, err = os.Stat(imagePath)
	if os.IsNotExist(err) {
		log.Fatal("image file does not exist")
	}
	// generate authentication token
	authToken, err := kms.GetAuthToken()

	if err != nil {
		log.Fatal("Error in generating token.", err)
	}

	//create key if keyId is not specified in input
	if len(strings.TrimSpace(keyID)) <= 0 {
		keyInformation := createKey(authToken)
		keyURL = c.Configuration.BaseURL + "keys/" + keyInformation.KeyID + "/transfer"
		fmt.Println(keyURL)
		key = retrieveKey(authToken, keyURL)
		fmt.Println(key)
	} else {
		//retrieve key using keyid
		keyURL = c.Configuration.BaseURL + "keys/" + keyID + "/transfer"
		fmt.Println(keyURL)
		key = retrieveKey(authToken, keyURL)
		fmt.Println(key)
	}

	// encrypt image using key
	//encryptedImage := encrypt(imagePath, encryptFilePath, key)

	//calculate SHA256 of the encrpted image
	s := "foo"
	digest := sha256.Sum256([]byte(s))

	// create image flavor
	imageFlavor, err := flavor.GetImageFlavor("label", encryptionRequired, keyURL, base64.StdEncoding.EncodeToString(digest[:]))
	if err != nil {
		log.Fatal("Error in creating image flavor.", err)
	}
	jsonFlavor, err := json.Marshal(imageFlavor)
	if len(strings.TrimSpace(outputFile)) <= 0 {
		return string(jsonFlavor), nil
	}
	//create outputFile
	_, err = os.Create(outputFile)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	err = ioutil.WriteFile(outputFile, []byte(jsonFlavor), 0600)
	if err != nil {
		fmt.Println(err)
	}

	return "", err

}

func createKey(authToken string) KeyInfo {
	var url string
	var requestBody []byte
	var keyObj KeyInfo

	url = c.Configuration.BaseURL + "keys"
	requestBody = []byte(`{"algorithm": "AES","key_length": "256","mode": "GCM"}`)

	httpRequest, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	httpRequest.Header.Set("Accept", "application/json")
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Authorization", "Token "+authToken)

	if err != nil {
		log.Fatal("Error in request for key creation. ", err)
	}

	httpResponse, err := kms.SendRequest(httpRequest)
	if err != nil {
		log.Fatal("Error in response from key creation request. ", err)
	}

	_ = json.Unmarshal([]byte(httpResponse), &keyObj)
	return keyObj
}

func retrieveKey(authToken string, keyURL string) string {
	var keyValue KeyObj
	httpRequest, err := http.NewRequest("POST", keyURL, nil)
	httpRequest.Header.Set("Accept", "application/json")
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Authorization", "Token "+authToken)
	if err != nil {
		log.Fatal(err)
	}
	httpResponse, err := kms.SendRequest(httpRequest)
	if err != nil {
		panic(err.Error())
	}
	_ = json.Unmarshal([]byte(httpResponse), &keyValue)
	fmt.Println(keyValue.Key)
	return keyValue.Key
}
