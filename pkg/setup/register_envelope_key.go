/*
 * Copyright (C) 2019 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package setup

import (
	//"encoding/base64"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	csetup "intel/isecl/lib/common/setup"
	"intel/isecl/wpm/consts"
	kmsc "intel/isecl/wpm/pkg/kmsclient"

	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
)

type RegisterEnvelopeKey struct {
}

var UserInformation kmsc.UserInfo

// ValidateRegisterKey method is used to verify if the envelope key is registered with the KBS
func (re RegisterEnvelopeKey) Validate(c csetup.Context) error {
	log.Trace("pkg/setup/register_envelope_key.go:Validate() Entering")
	defer log.Trace("pkg/setup/register_envelope_key.go:Validate() Leaving")

	var cert *x509.Certificate
	var wpmPublicKey *rsa.PublicKey

	log.Info("pkg/setup/register_envelope_key.go:Validate() Validating register envelope key.")

	if len(strings.TrimSpace(UserInformation.TransferKeyPem)) <= 0 {
		return nil
	}

	publicKey, err := ioutil.ReadFile(consts.EnvelopePublickeyLocation)
	if err != nil {
		return errors.Wrap(err, "pkg/setup/register_envelope_key.go:Validate() Error reading envelop key from file.")
	}

	publicKeyDecoded, _ := pem.Decode(publicKey)

	parsedPublicKey, err := x509.ParsePKIXPublicKey(publicKeyDecoded.Bytes)
	if err != nil {
		return errors.Wrap(err, "pkg/setup/register_envelope_key.go:Validate() Could not parse public key from PEM decoded content. "+err.Error())
	}

	wpmPublicKey, isRsaType := parsedPublicKey.(*rsa.PublicKey)
	if !isRsaType {
		return errors.New("pkg/setup/register_envelope_key.go:Validate() Public key not in RSA format")
	}

	block, _ := pem.Decode([]byte(UserInformation.TransferKeyPem))

	cert, err = x509.ParseCertificate(block.Bytes)
	if err != nil {
		return errors.Wrap(err, "pkg/setup/register_envelope_key.go:Validate() Could not parse public key from PEM decoded content"+err.Error())
	}

	kmsPublicKey := cert.PublicKey.(*rsa.PublicKey)

	if compareRSAPubKeys(wpmPublicKey, kmsPublicKey) {
		return errors.Wrap(err, "pkg/setup/register_envelope_key.go:Validate() WPM and KMS public keys match")
	}
	return nil
}

//RegisterEnvelopeKey method is used to register the envelope public key with the KBS user
func (re RegisterEnvelopeKey) Run(c csetup.Context) error {
	log.Trace("pkg/setup/register_envelope_key.go:Run() Entering")
	defer log.Trace("pkg/setup/register_envelope_key.go:Run() Leaving")

	log.Info("Registering envelope key")

	publicKey, err := ioutil.ReadFile(consts.EnvelopePublickeyLocation)
	if err != nil {
		return errors.Wrap(err, "pkg/setup/register_envelope_key.go:Run() Error while reading the envelope public key")
	}

	log.Info("Registering envelope key")
	kc, err := kmsc.InitializeKMSClient()
	if err != nil {
		return errors.Wrap(err, "pkg/setup/register_envelope_key.go:Run() Failure to initialize KMS client")
	}

	UserInformation, err = getUserInfo()
	if err != nil {
		return errors.Wrap(err, "pkg/setup/register_envelope_key.go:Run() User does not exist in KMS")
	}

	err = kc.Keys().RegisterUserPubKey(publicKey, UserInformation.UserID)
	if err != nil {
		return errors.Wrap(err, "pkg/setup/register_envelope_key.go:Run() Error while updating the KBS user with envelope public key. "+err.Error())
	}

	log.Info("Envelop key registered successfully")
	return nil
}

func getUserInfo() (kmsc.UserInfo, error) {
	log.Trace("pkg/setup/register_envelope_key.go:getUserInfo() Entering")
	defer log.Trace("pkg/setup/register_envelope_key.go:getUserInfo() Leaving")

	var userInfo kmsc.UserInfo
	kc, err := kmsc.InitializeKMSClient()
	if err != nil {
		return userInfo, errors.Wrap(err, "pkg/setup/register_envelope_key.go:getUserInfo() Failure to initialize KMS client")
	}

	userInfo, err = kc.Keys().GetKmsUser()
	if err != nil {
		return userInfo, errors.New("pkg/setup/register_envelope_key.go:getUserInfo() Error while getting the KMS user information. " + err.Error())
	}

	return userInfo, nil
}

func compareRSAPubKeys(rsaPubKey1 *rsa.PublicKey, rsaPubKey2 *rsa.PublicKey) bool {
	log.Trace("pkg/setup/register_envelope_key.go:compareRSAPubKeys() Entering")
	defer log.Trace("pkg/setup/register_envelope_key.go:compareRSAPubKeys() Leaving")

	return ((rsaPubKey1.N.Cmp(rsaPubKey2.N) == 0) && (rsaPubKey1.E == rsaPubKey2.E))
}
