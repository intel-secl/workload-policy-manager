package wpm
/*
 *
 * @author srege
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
	//"os"
	//"strings"
	//"errors"
	//"intel/isecl/lib/flavor"
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
	AuthorizationDate string  `json:"authorization_date"`
    NotAfter string  `json:"not_after"`
	Faults []string  `json:"faults"`
}

type KeyInfo struct{
	KeyId string `json:"id"`
	CipherMode string `json:"cipher_mode"`
	Algorithm string `json:"algorithm"`
	KeyLength string `json:"key_length"`
	PaddingMode string `json:"padding_mode"`
	TransferPolicy string `json:"transfer_policy"`
	TransferLink string `json:"transfer_link"`
	DigestAlgorithm string `json:"digest_algorithm"`
	
}

func createImageFlavor (imagePath string,encryptFilePath string,keyId string,encryptionRequired bool,integrityEnforced bool,outPutFile string)(error) {
	/*//input validation
	var err error
	if (len(strings.TrimSpace(imagePath)) <= 0){
        return errors.New("image path not given")
	}
	 
	if (len(strings.TrimSpace(encryptFilePath)) <= 0){
        return errors.New("encryption file path not given")
	}

	// check if image exists at the specified location
	_, err = os.Stat(imagePath)
	if !os.IsExist(err) {
		return errors.New("image file does nt exist")
	}

	client := NewBasicAuthClient("username", "password")
	authToken,err := client.getAuthToken()

	if err!=nil {
		log.Fatal("Error in generating token.",err)
	}
    */
	/*if len(strings.TrimSpace(keyId)) <= 0 {
		keyInformation = createKey(authToken)

		//hardcoded values need to be removed
		keyUrl := "https://10.1.70.56:443/v1/keys/" + keyInfo.KeyId + "/transfer"  
		
		encryptedKey := retrieveTransferKey(authToken,keyUrl)  
	} else {
		keyUrl := "https://10.1.70.56:443/v1/keys/" + keyId + "/transfer"  
		encryptedKey := retrieveTransferKey(authToken,keyUrl) 
	}

	//call to Encrypt image. 
    //digest := encrypt(imagePath,encryptFilePath,encryptedKey)

	//GetImageFlavor constructs a new ImageFlavor with the specified label, encryption policy, KMS url, encryption IV, and digest of the encrypted payload

	imageFlavor := GetImageFlavor("label",encryptionRequired,keyUrl,digest)
	if len(strings.TrimSpace(imagePath)) == 0 {
		 return imageFlavor
	}else {
		fileOpen, err := os.Create(outPutFile)
       if err != nil {
        panic(err)
	  }
	  err = ioutil.WriteFile(fileOpen, imageFlavor, 0644)
	  
	}
    
	*/
	return nil
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
	var keyObj KeyInfo
	url = "https://"+ kms_ip +  ":" + kms_port + "/v1/keys"
	
	var jsonStr = []byte(`{"algorithm": "AES","key_length": "256","cipher_mode": "GCM"}`)
	
	
	httpRequest,err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
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
	var jsonStr = `{"username": "admin","password": "password"}`
	requestBody = []byte(jsonStr)

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

/*func (s *Client) retrieveTransferKey(authToken AuthToken, keyUrl string) byte[]  {

	req, err := http.NewRequest("POST", keyUrl,nil)
	req.Header.Set("Accept", "application/x-pem-file")
	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("Authorization","Token " + authToken.AuthorizationToken)
	if err != nil {
	 log.Fatal(err)
	}
	fmt.Println("No error in request")
	fmt.Println(req)
	response, err := s.doRequest(req)
	if err != nil {
        panic(err.Error())
    }
	fmt.Println(response)
	return response
}*/



