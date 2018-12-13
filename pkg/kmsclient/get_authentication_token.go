package kmsclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	c "intel/isecl/wpm/config"
	"net/http"
)

//AuthToken is a representation of token using for authentication
type AuthToken struct {
	AuthorizationToken string `json:"authorization_token"`
}

// add import and config
func GetAuthToken() (string, error) {
	var err error
	var authToken AuthToken
	var url string
	var requestBody bytes.Buffer

	//Add client here
	url = c.Configuration.BaseURL + "login"
	fmt.Println(c.Configuration.BaseURL)
	//requestBody = []byte(`{"username": " + config.Username,"password": "password"}`)

	//build request body using username and password from config
	requestBody.WriteString(`{"username":"`)
	requestBody.WriteString(c.Configuration.Username)
	requestBody.WriteString(`","password":"`)
	requestBody.WriteString(c.Configuration.Password)
	requestBody.WriteString(`"}`)

	//Construct POST request with Accept and Content-Type headers to generate token
	httpRequest, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(requestBody.String())))
	httpRequest.Header.Set("Accept", "application/json")
	httpRequest.Header.Set("Content-Type", "application/json")

	httpResponse, err := SendRequest(httpRequest)
	if err != nil {
		return "", err
	}

	_ = json.Unmarshal([]byte(httpResponse), &authToken)
	if err != nil {
		return "", err
	}
	return authToken.AuthorizationToken, nil
}
