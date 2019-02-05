package setup

import (
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
	"crypto/tls"
	t "intel/isecl/lib/common/tls"

	logger "github.com/sirupsen/logrus"
)

type RegisterEnvelopeKey struct {
}

// ValidateRegisterKey method is used to verify if the envelope key is registered with the KBS
func (re RegisterEnvelopeKey) Validate(c csetup.Context) error {
	return nil
}

// ValidateRegisterKey method is used to verify if the envelope key is registered with the KBS
func ValidateRegisterKey() (string, bool) {

	userInfo, err := getUserInfo()
	if err != nil {
		return "",  true
	}
	if len(strings.TrimSpace(userInfo.TransferKeyPem)) < 0 {
		return userInfo.UserID, true
	} else {
		publicKey, err := ioutil.ReadFile(config.Configuration.EnvelopePublickeyLocation)
		if err != nil {
			return userInfo.UserID, true
		}
		encodedPublicKey := base64.StdEncoding.EncodeToString(publicKey)
		log.Println("encoded public key : ", encodedPublicKey)
		log.Println("user pub key : ", userInfo.TransferKeyPem)
		if strings.EqualFold(userInfo.TransferKeyPem, encodedPublicKey) {
			return userInfo.UserID, false
		}
		return userInfo.UserID,  true
	}

}

//RegisterEnvelopeKey method is used to register the envelope public key with the KBS user
func (re RegisterEnvelopeKey) Run(c csetup.Context) error {
    userID, isValidated := ValidateRegisterKey()
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
    kc := initializeClient()
	err = kc.Keys().RegisterUserPubKey(publicKey, userID)
	if err != nil {
		return errors.New("Error while updating the KBS user with envelope public key")
	}
	logger.Println("Envelop key registered successfully")
	return nil
}

func getUserInfo() (kms.UserInfo, error) {
    var userInfo kms.UserInfo
	kc := initializeClient()
	userInfo, err := kc.Keys().GetKmsUser()
	if err != nil {
		return userInfo, errors.New("Error while gettig the KMS user information")
	}
	return userInfo, nil
}
func initializeClient() (*kms.Client) {
	fmt.Println("Inside initializeClient")
	var certificateDigest [32]byte
	cert, err := hex.DecodeString(config.Configuration.Kms.TLSSha256)
	if err != nil {
		log.Fatal(err)
	}
	copy(certificateDigest[:], cert)
	client := &http.Client{
					Transport: &http.Transport{
						TLSClientConfig: &tls.Config{
							InsecureSkipVerify: true,
							VerifyPeerCertificate: t.VerifyCertBySha256(certificateDigest),
						},
					},
				}
	kc := &kms.Client{
		BaseURL:  config.Configuration.Kms.APIURL,
		Username: config.Configuration.Kms.APIUsername,
		Password: config.Configuration.Kms.APIPassword,
		CertSha256 :&certificateDigest,
        HTTPClient: client,
	}
	return kc
}
