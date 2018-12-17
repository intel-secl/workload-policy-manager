package setup

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type authToken struct {
	AuthorizationToken string `json:"authorization_token"`
}

func sendRequest(req *http.Request) ([]byte, error) {

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
		fmt.Println("Error in sending request.", err)
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	if 200 != response.StatusCode {
		fmt.Println("Returned status code "+string(response.StatusCode), err)
		return nil, fmt.Errorf("%s", body)
	}

	return body, nil
}

func getAuthToken(username, password string) (authToken, error) {

	var token authToken
	url := "https://10.1.70.56:443/v1/login"
	requestBody := []byte(`{"username": "admin","password": "password"}`)

	//POST request with Accept and Content-Type headers to generate token
	httpRequest, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatal("Error during initializing an http client")
	}
	httpRequest.Header.Set("Accept", "application/json")
	httpRequest.Header.Set("Content-Type", "application/json")

	response, err := sendRequest(httpRequest)
	_ = json.Unmarshal([]byte(response), &token)
	if err != nil {
		return token, err
	}
	return token, nil
}

func registerUserPubKey(publicKey []byte, token authToken) {
	requestURL := "https://10.1.70.56:443/v1/users/3aec609b-9226-420b-86f8-1121f5319773/transfer-key"
	httpRequest, err := http.NewRequest("PUT", requestURL, bytes.NewBuffer(publicKey))
	if err != nil {
		log.Fatal(err)
	}
	httpRequest.Header.Set("Content-Type", "application/x-pem-file")
	httpRequest.Header.Set("Authorization", "Token "+token.AuthorizationToken)

	_, _ = sendRequest(httpRequest)
}

func RegisterEnvelopeKey() {
	token, err := getAuthToken("admin", "password")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("token ", token.AuthorizationToken)

	publicKey, _ := ioutil.ReadFile("envelopeKey.pub")

	registerUserPubKey(publicKey, token)
}
