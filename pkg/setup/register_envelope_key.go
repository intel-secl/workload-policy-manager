package setup

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	csetup "intel/isecl/lib/common/setup"
	config "intel/isecl/wpm/config"
	//client "intel/isecl/wpm/pkg/kmsclient"
	"intel/isecl/lib/kms-client"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"encoding/hex"

	logger "github.com/sirupsen/logrus"
)

type RegisterEnvelopeKey struct {
}

// ValidateRegisterKey method is used to verify if the envelope key is registered with the KBS
func (re RegisterEnvelopeKey) Validate(c csetup.Context) error {
	return nil
}

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

//RegisterEnvelopeKey method is used to register the envelope public key with the KBS user
func (re RegisterEnvelopeKey) Run(c csetup.Context) error {
    fmt.Println("Before call to validateRegisterKey in Run")
	userID, token, isValidated := ValidateRegisterKey()
	if !isValidated {
		return errors.New("Envelope public key is already registered on KBS. Skipping this setup task....")
	}
	// save configuration from config.yml
	e := config.SaveConfiguration(c)
	if e != nil {
		logger.Error(e.Error())
		return e
	}
	publicKey, err := ioutil.ReadFile(config.Configuration.EnvelopePublickeyLocation)
	if err != nil {
		return errors.New("Error while reading the envelope public key")
	}

	err = registerUserPubKey(publicKey, userID, token)
	if err != nil {
		return errors.New("Error while updating the KBS user with envelope public key")
	}
	fmt.Println("Envelop key registered successfully")
	return nil
}

func registerUserPubKey(publicKey []byte, userID string, token string) error {
	var certificateDigest [32]byte
	cert, err := hex.DecodeString(config.Configuration.Kms.TLSSha256)
	if err != nil {
		log.Fatal(err)
	}
	copy(certificateDigest[:], cert)
	kc := &kms.Client{
		BaseURL:  config.Configuration.Kms.APIURL,
		Username: config.Configuration.Kms.APIUsername,
		Password: config.Configuration.Kms.APIPassword,
        CertSha256 :&certificateDigest,
	}

	requestURL := config.Configuration.Kms.APIURL + "users/" + userID + "/transfer-key"
	httpRequest, err := http.NewRequest("PUT", requestURL, bytes.NewBuffer(publicKey))
	if err != nil {
		fmt.Println(err)
		return errors.New("Error while creating a http request object")
	}
   	httpRequest.Header.Set("Content-Type", "application/x-pem-file")
	_, err = kc.DispatchRequest(httpRequest)
	if err != nil {
		fmt.Println(err)
		return errors.New("Error while sending a PUT request with envelope public key")
	}
	
	return nil
}

func getUserInfo() (UserInfo, string, error) {
    fmt.Println("Inside getUserInfo")
	var userInfo UserInfo
	var token string
	var certificateDigest [32]byte
	cert, err := hex.DecodeString(config.Configuration.Kms.TLSSha256)
	if err != nil {
		log.Fatal(err)
	}
	copy(certificateDigest[:], cert)
	kc := &kms.Client{
		BaseURL:  config.Configuration.Kms.APIURL,
		Username: config.Configuration.Kms.APIUsername,
		Password: config.Configuration.Kms.APIPassword,
		CertSha256 :&certificateDigest,
	}
	fmt.Println(kc.BaseURL)
    fmt.Println(kc.Username)
	fmt.Println(kc.Password)
	err = kc.RefreshAuthToken()
	if err != nil {
		fmt.Println(err)
		return userInfo, token, errors.New("Error while getting authentication token")
	}
    fmt.Println("Token inside register Envelope key")
    fmt.Println(kc.AuthenticationToken.AuthorizationToken)
	userInfo, err = GetKmsUser(kc.AuthenticationToken.AuthorizationToken)
	if err != nil {
		return userInfo, token, errors.New("Error while gettig the KMS user information")
	}

	return userInfo, token, nil
}
