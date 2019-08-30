/*
 * Copyright (C) 2019 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package util

import (
	"errors"
	kms "intel/isecl/lib/kms-client"
	config "intel/isecl/wpm/config"
	"intel/isecl/wpm/consts"
	"intel/isecl/wpm/pkg/kmsclient"
	"net/url"
	"strings"
)

//fetch key from kms
func FetchKey(keyID string) ([]byte, string, error) {
	var keyInfo kms.KeyInfo
	var keyURLString string
	var keyValue []byte

	//Initialize the KMS client
	kc, err := kmsclient.InitializeClient()
	if err != nil {
		return []byte(""), "", errors.New("error initializing KMS client")
	}

	//If key ID is not specified, create a new key
	if len(strings.TrimSpace(keyID)) <= 0 {
		keyInfo.Algorithm = consts.KmsEncryptAlgo
		keyInfo.KeyLength = consts.KmsKeyLength
		keyInfo.CipherMode = consts.KmsCipherMode

		key, err := kc.Keys().Create(keyInfo)
		if err != nil {
			return []byte(""), "", errors.New("error creating the image encryption key: " + err.Error())
		}
		keyID = key.KeyID
	}

	//Build the key URL, to be inserted later on when the image flavor is created
	keyURL, err := url.Parse(config.Configuration.Kms.APIURL + "keys/" + keyID + "/transfer")
	if err != nil {
		return []byte(""), "", errors.New("error building KMS key URL: " + err.Error())
	}
	keyURLString = keyURL.String()

	//Retrieve key using key ID
	keyValue, err = kc.Key(keyID).Retrieve()
	if err != nil {
		return []byte(""), "", errors.New("error retrieving the image encryption key: " + err.Error())
	}
	return keyValue, keyURLString, err
}
