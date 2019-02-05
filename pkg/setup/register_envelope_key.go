package setup

import (
	"encoding/base64"
	"errors"
	log "github.com/sirupsen/logrus"
	csetup "intel/isecl/lib/common/setup"
	"intel/isecl/lib/kms-client"
	config "intel/isecl/wpm/config"
	"intel/isecl/wpm/consts"
	"intel/isecl/wpm/pkg/kmsclient"
	"io/ioutil"
	"strings"
)

type RegisterEnvelopeKey struct {
}

// ValidateRegisterKey method is used to verify if the envelope key is registered with the KBS
func (re RegisterEnvelopeKey) Validate(c csetup.Context) error {
	log.Info("Validating register envelope key.")
	userInfo, err := getUserInfo()
	if err != nil {
		return errors.New("user does not exist in KMS")
	}
	publicKey, err := ioutil.ReadFile(consts.EnvelopePublickeyLocation)
	if err != nil {
		return errors.New("error reading envelop key")
	}
	encodedPublicKey := base64.StdEncoding.EncodeToString(publicKey)
	if strings.EqualFold(userInfo.TransferKeyPem, encodedPublicKey) {
		return errors.New("validation failed. Certificates from WPM and KMS do not match")
	}
	return nil

}

//RegisterEnvelopeKey method is used to register the envelope public key with the KBS user
func (re RegisterEnvelopeKey) Run(c csetup.Context) error {
	err := re.Validate(c)
	if err != nil {
		return errors.New("Envelope public key is already registered on KBS. Skipping this setup task....")
	}

	log.Info("Registering envelope key")

	userInfo, err := getUserInfo()
	if err != nil {
		return errors.New("Error while gettig the KMS user information")
	}

	// save configuration from config.yml
	e := config.SaveConfiguration(c)
	if e != nil {
		log.Error("Error saving configuration.")
		return e
	}
	publicKey, err := ioutil.ReadFile(consts.EnvelopePublickeyLocation)
	if err != nil {
		return errors.New("Error while reading the envelope public key")
	}
	kc := kmsclient.InitializeClient()
	err = kc.Keys().RegisterUserPubKey(publicKey, userInfo.UserID)
	if err != nil {
		return errors.New("Error while updating the KBS user with envelope public key")
	}
	log.Info("Envelop key registered successfully")
	return nil
}

func getUserInfo() (kms.UserInfo, error) {
	var userInfo kms.UserInfo
	kc := kmsclient.InitializeClient()
	userInfo, err := kc.Keys().GetKmsUser()
	if err != nil {
		return userInfo, errors.New("Error while gettig the KMS user information")
	}
	return userInfo, nil
}
