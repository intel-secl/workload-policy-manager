package kmsclient

import (
	"encoding/json"
	"errors"
	config "intel/isecl/wpm/config"
	"log"
	"net/http"
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
	requestURL := config.Configuration.KmsAPIURL + "users?usernameEqualTo=" + config.Configuration.KmsAPIUsername
	httpRequest, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	httpRequest.Header.Set("Accept", "application/json")
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Authorization", "Token "+token)
	var userInfo UserInfo
	var users Users

	httpResponse, err := SendRequest(httpRequest)
	if err != nil {
		return userInfo, errors.New("error while getting http response")
	}

	//deserialize the response to UserInfo response
	err = json.Unmarshal([]byte(httpResponse), &users)
	if err != nil {
		return userInfo, errors.New("error while unmarshalling the http response to the type users")
	}
	return users.Users[0], nil
}