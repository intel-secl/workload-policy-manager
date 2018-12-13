package kbsclient

import (
	"bytes"
	"encoding/json"
	"net/http"
	s "syscall/js"
)

// add import and config
func GetAuthToken(config *Configuration) (AuthToken, error) {
	var err error
	var requestBody []byte
	var authToken AuthToken
	var url string

	//Add client here
	url = config.BaseURL + "login"
	//create the request body using some string builder or something
	//request := `"{"username":` + config.User
	requestBody = []byte(`{"username": " + config.Username,"password": "password"}`)

	//POST request with Accept and Content-Type headers to generate token
	httpRequest, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	httpRequest.Header.Set("Accept", "application/json")
	httpRequest.Header.Set("Content-Type", "application/json")

	httpResponse, err := s.sendRequest(httpRequest)
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal([]byte(httpResponse), &authToken)
	if err != nil {
		return authToken, err
	}
	return authToken, nil
}
