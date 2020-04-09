/*
 * Copyright (C) 2019 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package util

import (
	"encoding/base64"
	config "intel/isecl/wpm/v2/config"
	"intel/isecl/wpm/v2/consts"
	kmsc "intel/isecl/wpm/v2/pkg/kmsclient"
	"io/ioutil"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

//fetch key from kms
func FetchKey(keyID string) ([]byte, string, error) {
	log.Trace("pkg/util/encrypt.go:FetchKey() Entering")
	defer log.Trace("pkg/util/encrypt.go:FetchKey() Leaving")

	var keyInfo kmsc.KeyInfo
	var keyURLString string
	var keyValue []byte

	//Initialize the KMS client
	kc, err := kmsc.InitializeKMSClient()
	if err != nil {
		return []byte(""), "", errors.Wrap(err, "pkg/util/fetch_key.go:FetchKey() Error initializing KMS client")
	}

	//If key ID is not specified, create a new key
	if len(strings.TrimSpace(keyID)) <= 0 {
		keyInfo.Algorithm = consts.KmsEncryptAlgo
		keyInfo.KeyLength = consts.KmsKeyLength
		keyInfo.CipherMode = consts.KmsCipherMode
		log.Debug("pkg/util/fetch_key.go:FetchKey() Creating new key")
		key, err := kc.Keys().Create(keyInfo)
		if err != nil {
			return []byte(""), "", errors.Wrap(err, "pkg/util/fetch_key.go:FetchKey() Error creating the image encryption key")
		}
		keyID = key.KeyID
		log.Debugf("pkg/util/fetch_key.go:FetchKey() keyID: %s", keyID)
	}

	//Build the key URL, to be inserted later on when the image flavor is created
	keyURL, err := url.Parse(config.Configuration.Kms.APIURL + "keys/" + keyID + "/transfer")
	if err != nil {
		return []byte(""), "", errors.Wrap(err, "Error building KMS key URL")
	}
	keyURLString = keyURL.String()
	log.Debugf("pkg/util/fetch_key.go:FetchKey() keyURL: %s", keyURLString)

	pubKey, err := ioutil.ReadFile(consts.EnvelopePublickeyLocation)
	if err != nil {
		return []byte(""), "", errors.Wrap(err, "pkg/util/fetch_key.go:FetchKey() Error reading envelop public key")
	}
	//Retrieve key using key ID
	keyValue, err = kc.Key(keyID).Retrieve(string(pubKey))
	if err != nil {
		return []byte(""), "", errors.Wrap(err, "pkg/util/fetch_key.go:FetchKey() Error retrieving the image encryption key")
	}
	log.Info("pkg/util/fetch_key.go:FetchKey() Successfully retrieved key")
	log.Debugf("pkg/util/fetch_key.go:FetchKey() %s | %s", base64.StdEncoding.EncodeToString(keyValue), keyURLString)
	return keyValue, keyURLString, nil
}
