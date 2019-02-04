package imageflavor

/*
 *
 * @author srege
 *
 */
import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	flavor "intel/isecl/lib/flavor"
	config "intel/isecl/wpm/config"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"intel/isecl/lib/kms-client"
	//"unsafe"
	logger "github.com/sirupsen/logrus"
)

//KeyInfo is a representation of key information
type KeyInfo struct {
	KeyID           string `json:"id"`
	CipherMode      string `json:"mode"`
	Algorithm       string `json:"algorithm"`
	KeyLength       int `json:"key_length"`
	PaddingMode     string `json:"padding_mode"`
	TransferPolicy  string `json:"transfer_policy"`
	TransferLink    string `json:"transfer_link"`
	DigestAlgorithm string `json:"digest_algorithm"`
}

//KeyObj is a represenation of the actual key
type KeyObj struct {
	Key []byte `json:"key"`
}

//CreateImageFlavor is used to create flavor of an encrypted image
func CreateImageFlavor(imagePath string, encryptFilePath string, keyID string, encryptionRequired bool, integrityEnforced bool, outputFile string) (string, error) {
	var err error
	var key []byte
	var keyURL string
	logger.Info("Creating image flavor")
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
	

	if err != nil {
		log.Fatal("Error in generating authentication token.", err)
	}

	//create key if keyId is not specified in input
	if len(strings.TrimSpace(keyID)) <= 0 {
		keyInformation := createKey()
		keyURL = config.Configuration.Kms.APIURL + "keys/" + keyInformation.KeyID + "/transfer"
		key = retrieveKey(keyURL)
	} else {
		//retrieve key using keyid
		keyURL = config.Configuration.Kms.APIURL + "keys/" + keyID + "/transfer"
		key = retrieveKey(keyURL)
	}

	// encrypt image using key
	err = encrypt(imagePath, config.Configuration.EnvelopePrivatekeyLocation, encryptFilePath, key)
	if err != nil {
		log.Fatal("Error in encrypting image.", err)
	}
	encryptedImage, err := ioutil.ReadFile(encryptFilePath)
	//calculate SHA256 of the encrpted image
	digest := sha256.Sum256([]byte(encryptedImage))

	// create image flavor
	imageFlavor, err := flavor.GetImageFlavor("label", encryptionRequired, keyURL, base64.StdEncoding.EncodeToString(digest[:]))
	if err != nil {
		log.Fatal("Error in creating image flavor.", err)
	}
	jsonFlavor, err := json.Marshal(imageFlavor)
	if len(strings.TrimSpace(outputFile)) <= 0 {
		return string(jsonFlavor), nil
	}
	//create outputFile for image flavor
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

func createKey() KeyInfo {
	var url string
	var keyObj KeyInfo
	kc := initializeClient()
	fmt.Println(kc.BaseURL)
	fmt.Println(kc.Username)
	fmt.Println(kc.Password)
	logger.Info("Creating transfer key")
	url = config.Configuration.Kms.APIURL + "keys"
	requestBody := []byte(`{"algorithm": "AES","key_length": "256","mode": "GCM"}`)
	// set POST request Accept, Content-Type and Authorization headers
	httpRequest, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	httpRequest.Header.Set("Accept", "application/json")
	httpRequest.Header.Set("Content-Type", "application/json")

	httpResponse, err := kc.DispatchRequest(httpRequest)
	if err != nil {
		fmt.Println(err)
		log.Fatal("Error in key creation. ", err)
	}
	fmt.Println(httpResponse)
	fmt.Println(httpResponse.Body)
	
	//deserialize the response to KeyInfo
	if err := json.NewDecoder(httpResponse.Body).Decode(&keyObj); err != nil {
		log.Fatal("Error in marshalling. ",err)
	}
	//err = json.NewDecoder(response).Decode(&keyObj)
	//fmt.Println(err)
    //err = json.Unmarshal(response, &keyObj)
	fmt.Println("Response returned")
	fmt.Println(err)
	fmt.Println(keyObj.KeyID)
	return keyObj
}

func retrieveKey(keyURL string) []byte {
	var keyValue KeyObj
	logger.Info("Retrieving transfer key")
	// set POST request Accept, Content-Type and Authorization headers
	kc := initializeClient()
	httpRequest, err := http.NewRequest("POST", keyURL, nil)
	httpRequest.Header.Set("Accept", "application/json")
	httpRequest.Header.Set("Content-Type", "application/json")
	
	if err != nil {
		log.Fatal(err)
	}
	httpResponse, err := kc.DispatchRequest(httpRequest)
	fmt.Println(httpResponse.Body)
	if err != nil {
		log.Fatal("Error in key retrieval. ", err)
	}
	fmt.Println(httpResponse)
	fmt.Println(httpResponse.Body)
	//deserialize the response to KeyInfo
	if err := json.NewDecoder(httpResponse.Body).Decode(&keyValue); err != nil {
		fmt.Println(err)
		log.Fatal("Error in marshalling. ",err)
	}
	fmt.Println(keyValue)
	//_ = json.Unmarshal(response, &keyValue)
	fmt.Println(keyValue.Key)
	return keyValue.Key
}

func initializeClient() (*kms.Client) {
	fmt.Println("Inside initializeClient")
	var certificateDigest [32]byte
	cert, err := hex.DecodeString(config.Configuration.Kms.TLSSha256)
	if err != nil {
		log.Fatal(err)
	}
	copy(certificateDigest[:], cert)
	kc := &kms.Client{
		BaseURL:  config.Configuration.Kms.APIURL,
		Username: config.Configuration.Kms.APIUsername,
		Password: config.Configuration.Kms.APIPassword,
        CertSha256 :&certificateDigest,
	}
	return kc
}
