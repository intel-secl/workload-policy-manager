package kmsclient

/*
 *
 * @author srege
 *
 */
import (
	"bytes"
	"encoding/json"
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
	url = c.Configuration.KmsAPIURL + "login"

	//build request body using username and password from config
	requestBody.WriteString(`{"username":"`)
	requestBody.WriteString(c.Configuration.KmsAPIUsername)
	requestBody.WriteString(`","password":"`)
	requestBody.WriteString(c.Configuration.KmsAPIPassword)
	requestBody.WriteString(`"}`)

	// set POST request Accept and Content-Type headers
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
