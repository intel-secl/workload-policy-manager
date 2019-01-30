package setup

import (
	"bytes"
	"encoding/base64"
	"errors"
	config "intel/isecl/wpm/config"
	client "intel/isecl/wpm/pkg/kmsclient"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// ValidateRegisterKey method is used to verify if the envelope key is registered with the KBS
func ValidateRegisterKey() (string, string, bool) {

	userInfo, token, err := getUserInfo()
	if err != nil {
		return "", "", true
	}
	if len(strings.TrimSpace(userInfo.TransferKeyPem)) < 0 {
		return userInfo.UserID, token, true
	} else {
		publicKey, err := ioutil.ReadFile(config.Configuration.EnvelopePublickeyLocation)
		if err != nil {
			return userInfo.UserID, token, true
		}
		encodedPublicKey := base64.StdEncoding.EncodeToString(publicKey)
		log.Println("encoded public key : ", encodedPublicKey)
		log.Println("user pub key : ", userInfo.TransferKeyPem)
		if strings.EqualFold(userInfo.TransferKeyPem, encodedPublicKey) {
			return userInfo.UserID, token, false
		}
		return userInfo.UserID, token, true
	}

}

func registerUserPubKey(publicKey []byte, userID string, token string) error {
	requestURL := config.Configuration.Kms.APIURL + "users/" + userID + "/transfer-key"
	httpRequest, err := http.NewRequest("PUT", requestURL, bytes.NewBuffer(publicKey))
	if err != nil {
		return errors.New("Error while creating a http request object")
	}
	httpRequest.Header.Set("Content-Type", "application/x-pem-file")
	httpRequest.Header.Set("Authorization", "Token "+token)

	_, err = client.SendRequest(httpRequest)
	if err != nil {
		return errors.New("Error while sending a PUT request with envelope public key")
	}
	return nil
}

// RegisterEnvelopeKey method is used to register the envelope public key with the KBS user
func RegisterEnvelopeKey(userID, token string) error {

	publicKey, err := ioutil.ReadFile(config.Configuration.EnvelopePublickeyLocation)
	if err != nil {
		return errors.New("Error while reading the envelope public key")
	}

	err = registerUserPubKey(publicKey, userID, token)
	if err != nil {
		return errors.New("Error while updating the KBS user with envelope public key")
	}
	return nil
}

func getUserInfo() (client.UserInfo, string, error) {

	var userInfo client.UserInfo
	var token string

	token, err := client.GetAuthToken()
	if err != nil {
		return userInfo, token, errors.New("Error while getting authentication token")
	}

	userInfo, err = client.GetKmsUser(token)
	if err != nil {
		return userInfo, token, errors.New("Error while gettig the KMS user information")
	}

	return userInfo, token, nil
}
