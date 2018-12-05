package wpm
/*
 *
 * @author srege
 * @author hmgowda
 */
import(
	"log"
	"net/http" 
	"bytes"
	"io/ioutil"
	"fmt"
	"time"
	"crypto/tls"
	"encoding/json"
	"os"
	"strings"
	"io"
	"intel/isecl/lib/flavor"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	//"encoding/pem"
	"crypto/sha256"
	"encoding/base64"
	
)

type Client struct {
	Username string
	Password string
}
 
func NewBasicAuthClient(username string, password string) *Client {
	return &Client{
		Username: username,
		Password: password,
	}
}
type AuthToken struct{
	AuthorizationToken string `json:"authorization_token"` 
}

type KeyInfo struct{
	KeyId string `json:"id"`
	CipherMode string `json:"mode"`
	Algorithm string `json:"algorithm"`
	KeyLength string `json:"key_length"`
	PaddingMode string `json:"padding_mode"`
	TransferPolicy string `json:"transfer_policy"`
	TransferLink string `json:"transfer_link"`
	DigestAlgorithm string `json:"digest_algorithm"`
	
}

type KeyObj struct{
	Key string `json:"key"`
}

func (s *Client) CreateImageFlavor (imagePath string,encryptFilePath string,keyId string,encryptionRequired bool,integrityEnforced bool,outputFile string)(string,error) {
	var err error
	kms_ip := "10.105.168.214"
	kms_port :="443"
	var key string
	var keyUrl string

	//input validation
	if (len(strings.TrimSpace(imagePath)) <= 0){
        log.Fatal("image path not given")
	}
	 
	if (len(strings.TrimSpace(encryptFilePath)) <= 0){
		log.Fatal("encryption file path not given")
	}

	// check if image exists at the specified location
	_, err = os.Stat(imagePath)
	if !os.IsExist(err) {
		log.Fatal("image file does not exist")
	}

	// generate authentication token
	client := NewBasicAuthClient("username", "password")
	authToken,err := client.getAuthToken(kms_ip,kms_port)

	if err!=nil {
		log.Fatal("Error in generating token.",err)
	}
   
	//create key if keyId is not specified in input
	if len(strings.TrimSpace(keyId)) <= 0 {
		keyInformation := client.createKey(authToken.AuthorizationToken,kms_ip,kms_port)
		keyUrl = "https://" + kms_ip + kms_port + "/v1/keys/" + keyInformation.KeyId + "/transfer"  
		
		key = client.retrieveKey(authToken.AuthorizationToken,keyUrl)  
	} else {
        //retrieve key using keyid
		keyUrl = "https://" + kms_ip + kms_port + "/v1/keys/" + keyId + "/transfer"  
		key = client.retrieveKey(authToken.AuthorizationToken,keyUrl) 
	}

	// encrypt image using key 
	err,encryptedImage := encryptImage(imagePath, encryptFilePath, []byte(key))
	if err!=nil{
         log.Fatal("Error in encrypting image.",err)
	}
	
	//calculate SHA256 of the encrpted image
	digest := sha256.Sum256(encryptedImage)
	
	// create image flavor 
	imageFlavor,err := flavor.GetImageFlavor("label",encryptionRequired,keyUrl,base64.StdEncoding.EncodeToString(digest[:]))
	if err!=nil{
		log.Fatal("Error in creating image flavor.",err)
   }
    jsonFlavor, err := json.Marshal(imageFlavor)
	if len(strings.TrimSpace(outputFile)) <= 0 {
		 return string(jsonFlavor),nil
	}else {
		_, errorMsg := os.Stat(outputFile)
	if !os.IsExist(errorMsg) {
		log.Fatal("output file does not exist.")
	}
	
	 err = ioutil.WriteFile(outputFile,[]byte(jsonFlavor), 0644)
	  
	}
    
	
	return "",err
}

func (s *Client) sendRequest(req *http.Request) ([]byte, error) {

	tr := &http.Transport{
        TLSClientConfig: &tls.Config{
			//Needs to be changed to use secure tls cert 256
            InsecureSkipVerify: true,
        },
    }
   client := &http.Client{
        Transport: tr,
        Timeout:   3 * time.Second,
    }
	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Error in sending request.",err)
		return nil, err
    }
	defer response.Body.Close()

	//create byte array of HTTP response body
	body, err := ioutil.ReadAll(response.Body)
    
	if err != nil {
		return nil, err
	}
	
	if 200 != response.StatusCode {
		fmt.Println("Returned status code " + string(response.StatusCode),err)
		return nil, fmt.Errorf("%s", body)
	}
	
	return body, nil
}

func (s *Client) createKey(authToken string,kms_ip string,kms_port string) (KeyInfo) {
	var url string
	var requestBody []byte
	var keyObj KeyInfo
	

	url = "https://"+ kms_ip +  ":" + kms_port + "/v1/keys"
	requestBody = []byte(`{"algorithm": "AES","key_length": "256","mode": "GCM"}`)
		
	httpRequest,err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	httpRequest.Header.Set("Accept", "application/json")
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Authorization","Token " + authToken)
		
	if err != nil {
		log.Fatal("Error in request for key creation. ",err)
	}
	
	httpResponse, err := s.sendRequest(httpRequest)
	if err != nil {
		log.Fatal("Error in response from key creation request. ",err)
	} 
	
	_ = json.Unmarshal([]byte(httpResponse), &keyObj)
	return keyObj 
}

func (s *Client) getAuthToken(kms_ip string,kms_port string) (AuthToken,error) {
	var err error
	var requestBody []byte
	var authToken AuthToken	
	var url string
		
	url = "https://"+ kms_ip +  ":" + kms_port + "/v1/login"
	requestBody = []byte(`{"username": "admin","password": "password"}`)
	
    //POST request with Accept and Content-Type headers to generate token 
	httpRequest, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	httpRequest.Header.Set("Accept", "application/json")
	httpRequest.Header.Set("Content-Type", "application/json")
	
	if err != nil {
		return authToken,err
	}

  httpResponse, err := s.sendRequest(httpRequest)
	_ = json.Unmarshal([]byte(httpResponse), &authToken)
	if err != nil {
        return authToken,err
	}
	return authToken,nil
}

func (s *Client) retrieveKey(authToken string, keyUrl string) string {
    var keyValue KeyObj
	httpRequest, err := http.NewRequest("POST", keyUrl,nil)
	httpRequest.Header.Set("Accept", "application/json")
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Authorization","Token " + authToken)
	if err != nil {
	 log.Fatal(err)
	}
	httpResponse, err := s.sendRequest(httpRequest)
	if err != nil {
        panic(err.Error())
    }
	_ = json.Unmarshal([]byte(httpResponse), &keyValue)
    fmt.Println(keyValue.Key)
	return keyValue.Key
}

func encryptImage(imagePath string,encryptFilePath string, key []byte) (error,[]byte) {
	fmt.Println("Step 1")
	// reading the key file in Pem format
	data, err := ioutil.ReadFile(imagePath)
	fmt.Println("Step 1.2")
	if err != nil {
		fmt.Println("Error reading file")
		log.Fatal("Error reading the image file", err)
	}
    fmt.Println("Step 2")
	//decodedKey, _ := pem.Decode(key)
	  
	// creating a new cipher block of 128 bits
	block, err := aes.NewCipher(key)
	if err!= nil {
		log.Fatal("Error creating new cipher block", err)
	}

	fmt.Println("Step 3")
	gcm, err := cipher.NewGCM(block)
	fmt.Println("Step 4")
	if err != nil {
		log.Fatal("Error creating a cipher block", err)
	}
	// assigning a 12 byte empty array to store random value
	iv := make([]byte, 12)
	fmt.Println("Step 5")
	// reading random value into the byte array
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		log.Fatal("Error creating random IV value", err)
	}

	// encrypting the file(data) with IV value and appending
	// the IV to the first 12 bytes of the encrypted file
	ciphertext := gcm.Seal(iv, iv, data, nil)
	err = ioutil.WriteFile(encryptFilePath, ciphertext, 0644)
	return err,ciphertext
}


