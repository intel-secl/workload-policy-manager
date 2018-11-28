
package main
/**
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
)

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

/*func createImageFlavor (imagePath string,encryptFilePath string,keyId string,encryptionRequired bool,integrityEnforced bool,outPutFile string)(ImageFlavor) {
	//input validation
	if (len(strings.TrimSpace(imagePath)) == 0 || len(strings.TrimSpace(encryptFilePath)) == 0 ){
        log.Fatal("Invalid input parameters") 
	}
	client := NewBasicAuthClient("username", "password")
	authToken,err := client.getAuthToken()

	if err!=nil {
		log.Fatal("Error in getting token")
	}

	if len(strings.TrimSpace(keyId)) == 0 {
		keyInformation = createKey(authToken)
		//hardcoded values need to be removed
		keyUrl := "https://10.1.70.56:443/v1/keys/" + keyInfo.KeyId + "/transfer"  
		
		retrieveTransferKey(authToken,keyUrl)  
	} else {
		keyUrl := "https://10.1.70.56:443/v1/keys/" + keyId + "/transfer"  
		retrieveTransferKey(authToken,keyUrl) 
	}

	//call to Encrypt image
    

	//GetImageFlavor constructs a new ImageFlavor with the specified label, encryption policy, KMS url, encryption IV, and digest of the encrypted payload

	imageFlavor := GetImageFlavor("label",encryptionRequired,keyUrl,)
	if len(strings.TrimSpace(imagePath)) == 0 {
		 return imageFlavor
	}else {
		fileOpen, err := os.Create(outPutFile)
       if err != nil {
        panic(err)
	  }
	  err = ioutil.WriteFile(fileOpen, imageFlavor, 0644)
	  
	}
    
	
}
*/
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


func main(){
    /* client := NewBasicAuthClient("username", "password")


 authToken,_ := client.GetAuth()
  fmt.Println(authToken.AuthorizationToken)
 keyInfo := client.AddTodo(authToken)
 fmt.Println(keyInfo)
  keyUrl := "https://10.1.70.56:443/v1/keys/" + keyInfo.KeyId + "/transfer"
  fmt.Println(keyUrl)
  client.retrieveTransferKey(authToken,keyUrl)
*/
 /*  fileOpen, err := os.Create("/root/salonee.txt")
       if err != nil {
        panic(err)
	  }*/
	  
	  err = ioutil.WriteFile("/root/salonee.txt", "I am awesome", 0644)
	  if err != nil{
		  fmt.Println(err)
	  }

}

func (s *Client) doRequest(req *http.Request) ([]byte, error) {
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
	fmt.Println(req)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error in sending request",err)
		return nil, err
    }
	fmt.Println("No error in sending request")
	defer resp.Body.Close()
	fmt.Println("The response is ",resp)
	body, err := ioutil.ReadAll(resp.Body)
    fmt.Println(body)
	if err != nil {
		return nil, err
	}
	fmt.Println("No error in creating body")
	if 200 != resp.StatusCode {
		fmt.Println("Error in status Code")
		fmt.Println(resp.StatusCode)
		return nil, fmt.Errorf("%s", body)
	}
	fmt.Println(resp.StatusCode)
	return body, nil
}

func (s *Client) createKey(authToken AuthToken) (KeyInfo) {
	url := "https://10.1.70.56:443/v1/keys"
	fmt.Println(url)
	var jsonStr = []byte(`{"algorithm": "AES","key_length": "256","cipher_mode": "GCM"}`)
	req,err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization","Token " + authToken.AuthorizationToken)
	fmt.Println("Inside AddtoDo")
	fmt.Println(authToken.AuthorizationToken)
	
	if err != nil {
		log.Fatal("Error in request for key creation. ",err)
	}
	fmt.Println("No error in request")
	fmt.Println(req.Header.Get("Authorization"))
	response, err := s.doRequest(req)
	if err != nil {
		log.Fatal("Error in response from key creation request. ",err)
	} 
	fmt.Println(response)
	keyObj := KeyInfo{}
	_ = json.Unmarshal([]byte(response), &keyObj)
	fmt.Println(keyObj)
	return keyObj 
}

func (s *Client) getAuthToken() (AuthToken,error) {
	authToken:= AuthToken{}
	url := "https://10.1.70.56/v1/login"
	fmt.Println(url)
	var jsonStr = []byte(`{"username": "admin","password": "password"}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Accept", "application/json")
    req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return authToken,err
	}
	fmt.Println("No error in request")
	response, err := s.doRequest(req)
	if err != nil {
        panic(err.Error())
    }

	errorMsg := json.Unmarshal([]byte(response), &authToken)

	if errorMsg != nil {
        fmt.Println("Error in seriliazing",errorMsg) 
	}
	fmt.Println(authToken)
	return authToken,err
}

func (s *Client) retrieveTransferKey(authToken AuthToken, keyUrl string)  {

	fmt.Println(keyUrl)
	
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
	
}


