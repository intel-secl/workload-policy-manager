package kmsclient

/*
 *
 * @author srege
 *
 */
import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

//add imports
func SendRequest(req *http.Request) ([]byte, error) {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			//add tls cert 256
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

	//create byte array of HTTP response body
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
