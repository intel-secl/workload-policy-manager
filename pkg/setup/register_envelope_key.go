package setup

import (
	//"encoding/base64"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	log "github.com/sirupsen/logrus"
	csetup "intel/isecl/lib/common/setup"
	kms "intel/isecl/lib/kms-client"
	config "intel/isecl/wpm/config"
	"intel/isecl/wpm/consts"
	"intel/isecl/wpm/pkg/kmsclient"
	"io/ioutil"
	"strings"
)

type RegisterEnvelopeKey struct {
}
var UserInformation kms.UserInfo

// ValidateRegisterKey method is used to verify if the envelope key is registered with the KBS
func (re RegisterEnvelopeKey) Validate(c csetup.Context) error {
	var cert *x509.Certificate
	var wpmPublicKey *rsa.PublicKey

	log.Info("Validating register envelope key.")

	if len(strings.TrimSpace(UserInformation.TransferKeyPem)) <= 0 {
		return nil
	}

	publicKey, err := ioutil.ReadFile(consts.EnvelopePublickeyLocation)
	if err != nil {
		return errors.New("error reading envelop key from file.")
	}
    
	publicKeyDecoded, _ := pem.Decode(publicKey)

	parsedPublicKey, err := x509.ParsePKIXPublicKey(publicKeyDecoded.Bytes)
	if err != nil {
		return errors.New("could not parse public key from PEM decoded content")
	}

	wpmPublicKey, isRsaType := parsedPublicKey.(*rsa.PublicKey)
	if !isRsaType {
		return errors.New("public key not in RSA format")
	}

	block, _ := pem.Decode([]byte(UserInformation.TransferKeyPem))

	cert, err = x509.ParseCertificate(block.Bytes)
	if err != nil {
		return errors.New("could not parse public key from PEM decoded content")
	}

	kmsPublicKey := cert.PublicKey.(*rsa.PublicKey)

	if compareRSAPubKeys(wpmPublicKey, kmsPublicKey) {
		return errors.New("WPM and KMS public keys match")
	}
	return nil
}

//RegisterEnvelopeKey method is used to register the envelope public key with the KBS user
func (re RegisterEnvelopeKey) Run(c csetup.Context) error {
	log.Info("Registering envelope key")
    
	// save configuration from config.yml
	e := config.SaveConfiguration(c)
	if e != nil {
		log.Error("error saving configuration")
		return e
	}

	publicKey, err := ioutil.ReadFile(consts.EnvelopePublickeyLocation)
	if err != nil {
		return errors.New("error while reading the envelope public key")
	}

	kc := kmsclient.InitializeClient()

	UserInformation, err = getUserInfo()
	if err != nil {
		return errors.New("user does not exist in KMS")
	}

	err = kc.Keys().RegisterUserPubKey(publicKey, UserInformation.UserID)
	if err != nil {
		return errors.New("error while updating the KBS user with envelope public key")
	}

	log.Info("Envelop key registered successfully")
	return nil
}

func getUserInfo() (kms.UserInfo, error) {
	var userInfo kms.UserInfo
	kc := kmsclient.InitializeClient()

	userInfo, err := kc.Keys().GetKmsUser()
	if err != nil {
		return userInfo, errors.New("error while gettig the KMS user information")
	}

	return userInfo, nil
}
func compareRSAPubKeys(rsaPubKey1 *rsa.PublicKey, rsaPubKey2 *rsa.PublicKey) bool {
	return ((rsaPubKey1.N.Cmp(rsaPubKey2.N) == 0) && (rsaPubKey1.E == rsaPubKey2.E))
}
