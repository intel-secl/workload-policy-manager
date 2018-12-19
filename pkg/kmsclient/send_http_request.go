package kmsclient

/*
 *
 * @author srege
 *
 */
import (
	"crypto/tls"
	"encoding/hex"
	"fmt"
	t "intel/isecl/lib/common/tls"
	"intel/isecl/wpm/config"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

//SendRequest method is used to create an http client object and send the request to the server
func SendRequest(req *http.Request) ([]byte, error) {

	cert, err := hex.DecodeString(config.Configuration.KMSTlsCertSHA256)
	if err != nil {
		log.Fatal(err)
	}

	var certificateDigest [32]byte
	copy(certificateDigest[:], cert)

	tlsConfig := tls.Config{
		InsecureSkipVerify:    true,
		VerifyPeerCertificate: t.VerifyCertBySha256(certificateDigest),
	}
	transport := http.Transport{
		TLSClientConfig: &tlsConfig,
	}
	client := &http.Client{
		Transport: &transport,
		Timeout:   3 * time.Second,
	}

	/*tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify:    true,
			VerifyPeerCertificate: t.VerifyCertBySha256(certificateDigest),
		},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   3 * time.Second,
	}*/
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
