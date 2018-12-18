package kmsclient

/*
 *
 * @author srege
 *
 */
import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

//SendRequest method is used to create an http client object and send the request to the server
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
		log.Println("Error in sending request.", err)
		return nil, err
	}
	defer response.Body.Close()

	//create byte array of HTTP response body
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	log.Println("status code returned : ", strconv.Itoa(response.StatusCode))
	return body, nil
}
