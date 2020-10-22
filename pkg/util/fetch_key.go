/*
 * Copyright (C) 2019 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package util

import (
	"encoding/json"
	"intel/isecl/wpm/v3/config"
	"intel/isecl/wpm/v3/consts"
	"io/ioutil"
	"net/url"
	"regexp"
	"strings"

	kmsc "github.com/intel-secl/intel-secl/v3/pkg/clients/kbs"
	"github.com/intel-secl/intel-secl/v3/pkg/lib/common/crypt"
	"github.com/intel-secl/intel-secl/v3/pkg/model/kbs"
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

	aasUrl, err := url.Parse(config.Configuration.Aas.APIURL)
	if err != nil {
		return nil, "", errors.Wrap(err, "pkg/util/fetch_key.go:FetchKey() Error parsing AAS url")
	}

	kmsUrl, err := url.Parse(config.Configuration.Kms.APIURL)
	if err != nil {
		return nil, "", errors.Wrap(err, "pkg/util/fetch_key.go:FetchKey() Error parsing KMS url")
	}

	//Load trusted CA certificates
	caCerts, err := crypt.GetCertsFromDir(consts.TrustedCaCertsDir)
	if err != nil {
		return nil, "", errors.Wrap(err, "pkg/util/fetch_key.go:FetchKey() Error loading CA certificates")
	}

	//Initialize the KMS client
	kc := kmsc.NewKBSClient(aasUrl, kmsUrl, config.Configuration.Wpm.Username, config.Configuration.Wpm.Password, caCerts)

	var keyUrlString string
	//If key ID is not specified, create a new key
	if len(strings.TrimSpace(keyID)) <= 0 {
		var keyInfo kbs.KeyInformation
		var keyRequest kbs.KeyRequest

		keyInfo.Algorithm = consts.KmsEncryptAlgo
		keyInfo.KeyLength = consts.KmsKeyLength
		keyRequest.KeyInformation = &keyInfo
		if assetTagReg.MatchString(strings.TrimSpace(assetTag)) {
			keyRequest.Usage = assetTag
		} else {
			log.Warn("pkg/util/fetch_key.go:FetchKey() Asset Tags provided are not in valid format. Skipping associating usage policy")
		}
		log.Debug("pkg/util/fetch_key.go:FetchKey() Creating new key")
		keyResponse, err := kc.CreateKey(&keyRequest)
		if err != nil {
			return nil, "", errors.Wrap(err, "pkg/util/fetch_key.go:FetchKey() Error creating the image encryption key")
		}

		keyID = keyResponse.KeyInformation.ID.String()
		log.Debugf("pkg/util/fetch_key.go:FetchKey() keyID: %s", keyID)
		keyUrlString = keyResponse.TransferLink

	} else {
		//Build the key URL, to be inserted later on when the image flavor is created
		keyUrl, err := url.Parse(config.Configuration.Kms.APIURL + "keys/" + keyID + "/transfer")
		if err != nil {
			return nil, "", errors.Wrap(err, "Error building KMS key URL")
		}
		keyUrlString = keyUrl.String()
	}

	log.Debugf("pkg/util/fetch_key.go:FetchKey() keyUrl: %s", keyUrlString)

	pubKey, err := ioutil.ReadFile(consts.EnvelopePublickeyLocation)
	if err != nil {
		return nil, "", errors.Wrap(err, "pkg/util/fetch_key.go:FetchKey() Error reading envelop public key")
	}
	//Retrieve key using key ID
	keyValue, err := kc.TransferKey(keyID, string(pubKey))
	if err != nil {
		return nil, "", errors.Wrap(err, "pkg/util/fetch_key.go:FetchKey() Error retrieving the image encryption key")
	}
	log.Info("pkg/util/fetch_key.go:FetchKey() Successfully retrieved key")
	log.Debugf("pkg/util/fetch_key.go:FetchKey() %s", keyUrlString)
	return []byte(keyValue.KeyData), keyUrlString, nil
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
	if err != nil {
		return nil, errors.Wrap(err, "Error while fetching the key")
	}

	// unwrap
	key, err := UnwrapKey(wrappedKey, consts.EnvelopePrivatekeyLocation)
	if err != nil {
		return nil, errors.Wrap(err, "Error while unwrapping the key")
	}

	var returnKeyInfo = keyInfo{
		KeyUrl: keyUrlString,
		Key:    key,
	}

	//Marshall to a JSON string
	keyJSON, err := json.Marshal(returnKeyInfo)
	if err != nil {
		return nil, errors.Wrap(err, "Error while marshalling key info")
	}

	log.Info("pkg/util:FetchKeyForAssetTag() Successfully received encryption key from kms")
	return keyJSON, nil
}
