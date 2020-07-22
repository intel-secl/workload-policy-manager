/*
 * Copyright (C) 2019 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package util

import (
	"encoding/base64"
	"encoding/json"
	config "intel/isecl/wpm/v2/config"
	"intel/isecl/wpm/v2/consts"
	kmsc "intel/isecl/wpm/v2/pkg/kmsclient"
	"io/ioutil"
	"net/url"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

var (
	assetTagReg = regexp.MustCompile(`^[a-zA-Z0-9]+:[a-zA-Z0-9]+$`)
)

type keyInfo struct {
	KeyUrl string `json:"key_url"`
	Key    []byte `json:"key"`
}

//FetchKey from kms
func FetchKey(keyID string, assetTag string) ([]byte, string, error) {
	log.Trace("pkg/util/encrypt.go:FetchKey() Entering")
	defer log.Trace("pkg/util/encrypt.go:FetchKey() Leaving")

	var keyInfo kmsc.KeyInfo
	var keyUrlString string
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
		if assetTagReg.MatchString(strings.TrimSpace(assetTag)) {
			keyInfo.UsagePolicy = assetTag
		}
		log.Debug("pkg/util/fetch_key.go:FetchKey() Creating new key")
		key, err := kc.Keys().Create(keyInfo)
		if err != nil {
			return []byte(""), "", errors.Wrap(err, "pkg/util/fetch_key.go:FetchKey() Error creating the image encryption key")
		}
		keyID = key.KeyID
		log.Debugf("pkg/util/fetch_key.go:FetchKey() keyID: %s", keyID)
	}

	//Build the key URL, to be inserted later on when the image flavor is created
	keyUrl, err := url.Parse(config.Configuration.Kms.APIURL + "keys/" + keyID + "/transfer")
	if err != nil {
		return []byte(""), "", errors.Wrap(err, "Error building KMS key URL")
	}
	keyUrlString = keyUrl.String()
	log.Debugf("pkg/util/fetch_key.go:FetchKey() keyUrl: %s", keyUrlString)

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
	log.Debugf("pkg/util/fetch_key.go:FetchKey() %s | %s", base64.StdEncoding.EncodeToString(keyValue), keyUrlString)
	return keyValue, keyUrlString, nil
}

//FetchKeyForAssetTag is used to create flavor of an encrypted image
func FetchKeyForAssetTag(keyID string, assetTag string) ([]byte, error) {
	log.Trace("pkg/imageflavor/create_image_flavors.go:FetchKeyForAssetTag() Entering")
	defer log.Trace("pkg/imageflavor/create_image_flavors.go:FetchKeyForAssetTag() Leaving")

	var err error
	var wrappedKey []byte
	var keyUrlString string

	//Fetch the key
	wrappedKey, keyUrlString, err = FetchKey(keyID, assetTag)
	// unwrap
	key, err := UnwrapKey(wrappedKey, consts.EnvelopePrivatekeyLocation)

	var retrunkeyInfo = keyInfo{
		KeyUrl: keyUrlString,
		Key:    key,
	}

	//Marshall to a JSON string
	keyJSON, err := json.Marshal(retrunkeyInfo)
	if err != nil {
		return keyJSON, errors.Wrap(err, "Error while marshalling key info: ")
	}

	log.Info("pkg/util:FetchKeyForAssetTag() Successfully wrote image flavor to file")
	return keyJSON, nil
}
