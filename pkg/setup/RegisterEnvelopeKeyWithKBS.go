package setup

import (
	"bytes"
	"errors"
	config "intel/isecl/wpm/config"
	client "intel/isecl/wpm/pkg/kmsclient"
	"io/ioutil"
	"log"
	"net/http"
)

func validate() {

}

func registerUserPubKey(publicKey []byte, userID string, token string) error {
	requestURL := config.Configuration.KmsAPIURL + "users/" + userID + "/transfer-key"
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
func RegisterEnvelopeKey() error {
	token, err := client.GetAuthToken()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("token ", token)

	publicKey, _ := ioutil.ReadFile(config.Configuration.EnvelopePublickeyLocation)
	userInfo, err := client.GetKmsUser(token)
	if err != nil {
		return errors.New("Error while gettig the KMS user information")
	}
	err = registerUserPubKey(publicKey, userInfo.UserID, token)
	if err != nil {
		return errors.New("Error while updating the KBS user with envelope public key")
	}
	return nil
}
