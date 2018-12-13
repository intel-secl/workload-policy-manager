package flavor

/*
 *
 * @author srege
 * @author hmgowda
 */
import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go-workload-policy-manager/pkg/flavor"
	"io"
	"io/ioutil"
	flavor "lib-go-flavor"
	"log"
	"net/http"
	"os"
	"strings"
)

//AuthToken is a representation of token using for authentication
type AuthToken struct {
	AuthorizationToken string `json:"authorization_token"`
}

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

	//BaseURL := "https://10.105.168.214:443/v1/"
	configuration := getConfigurationVariables()
	//input validation
	if len(strings.TrimSpace(imagePath)) <= 0 {
		log.Fatal("image path not given")
	}

	if len(strings.TrimSpace(encryptFilePath)) <= 0 {
		log.Fatal("encryption file path not given")
	}

	// check if image exists at the specified location
	_, err = os.Stat(imagePath)
	if !os.IsExist(err) {
		log.Fatal("image file does not exist")
	}

	// generate authentication token
	authToken, err := GetAuthToken(&configuration)

	if err != nil {
		log.Fatal("Error in generating token.", err)
	}
	fmt.Println("Token generated")
	//create key if keyId is not specified in input
	if len(strings.TrimSpace(keyID)) <= 0 {
		keyInformation := client.createKey(authToken.AuthorizationToken, BaseURL)
		keyURL = BaseURL + "keys/" + keyInformation.KeyID + "/transfer"
		fmt.Println(keyURL)
		key = client.retrieveKey(authToken.AuthorizationToken, keyURL)
		fmt.Println(key)
	} else {
		//retrieve key using keyid
		keyURL = BaseURL + "keys/" + keyID + "/transfer"
		fmt.Println(keyUrl)
		key = client.retrieveKey(authToken.AuthorizationToken, keyURL)
		fmt.Println(key)
	}

	// encrypt image using key
	err, encryptedImage := encryptImage(imagePath, encryptFilePath, []byte(key))
	if err != nil {
		log.Fatal("Error in encrypting image.", err)
	}

	//calculate SHA256 of the encrpted image

	s := "Foo"
	digest := sha256.Sum256([]byte(s))

	// create image flavor
	imageFlavor, err := flavor.GetImageFlavor("label", encryptionRequired, keyUrl, base64.StdEncoding.EncodeToString(digest[:]))
	if err != nil {
		log.Fatal("Error in creating image flavor.", err)
	}
	jsonFlavor, err := json.Marshal(imageFlavor)
	if len(strings.TrimSpace(outputFile)) <= 0 {
		return string(jsonFlavor), nil
	}
	err = ioutil.WriteFile(outputFile, []byte(jsonFlavor), 0644)
	if err != nil {
		fmt.Println(err)
	}

	return "", err
}

func GetConfigurationVariables(filename string) Configuration {
	configuration := Configuration{}
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configuration)
	if err != nil {
		log.Fatal(err)
	}
	return configuration

}

func createKey(authToken string, config *Configuration) KeyInfo {
	var url string
	var requestBody []byte
	var keyObj KeyInfo

	url = config.BaseURL + "keys"
	requestBody = []byte(`{"algorithm": "AES","key_length": "256","mode": "GCM"}`)

	httpRequest, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	httpRequest.Header.Set("Accept", "application/json")
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Authorization", "Token "+authToken)

	if err != nil {
		log.Fatal("Error in request for key creation. ", err)
	}

	httpResponse, err := sendRequest(httpRequest)
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
	httpResponse, err := sendRequest(httpRequest)
	if err != nil {
		panic(err.Error())
	}
	_ = json.Unmarshal([]byte(httpResponse), &keyValue)
	fmt.Println(keyValue.Key)
	return keyValue.Key
}

func encryptImage(imagePath string, encryptFilePath string, key []byte) ([]byte, error) {
	fmt.Println("Step 1")
	// reading the key file in Pem format
	data, err := ioutil.ReadFile(imagePath)
	if err != nil {
		fmt.Println("Error reading file")
		log.Fatal("Error reading the image file", err)
	}
	//decodedKey, _ := pem.Decode(key)

	// creating a new cipher block of 128 bits
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal("Error creating new cipher block", err)
	}

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
	err = ioutil.WriteFile(encryptFilePath, ciphertext, 0644)
	return err, ciphertext
}
