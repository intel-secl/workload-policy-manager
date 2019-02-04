package setup

import (
	"encoding/json"
	"errors"
	config "intel/isecl/wpm/config"
	"log"
	"net/http"
    "encoding/hex"
	logger "github.com/sirupsen/logrus"
	"intel/isecl/lib/kms-client"
	"fmt"
)

//UserInfo is a representation of key information
type UserInfo struct {
	UserID         string `json:"id"`
	Username       string `json:"username"`
	TransferKeyPem string `json:"transfer_key_pem"`
}

//Users is a representation of key information
type Users struct {
	Users []UserInfo `json:"users"`
}

// GetKmsUser is used to get the kms user information
func GetKmsUser(token string) (UserInfo, error) {
	logger.Info("Retrieving kms user information")
	fmt.Println("Inside Get Kms user")
	fmt.Println(token)
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

	requestURL := config.Configuration.Kms.APIURL + "users?usernameEqualTo=" + config.Configuration.Kms.APIUsername
	httpRequest, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	httpRequest.Header.Set("Accept", "application/json")
	httpRequest.Header.Set("Content-Type", "application/json")
	var userInfo UserInfo
	var users Users

	httpResponse, err := kc.DispatchRequest(httpRequest)
	fmt.Println(httpResponse)
	if err != nil {
		fmt.Println(err)
		return userInfo, errors.New("error while getting http response")
	}
   	//deserialize the response to UserInfo response
	if err := json.NewDecoder(httpResponse.Body).Decode(&users); err != nil {
		return userInfo, errors.New("error while unmarshalling the http response to the type users")
	}
	fmt.Println("Inside get kms user")
	fmt.Println(users.Users[0])
	return users.Users[0], nil
}
