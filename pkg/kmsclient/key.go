/*
 * Copyright (C) 2019 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package kmsclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"intel/isecl/wpm/pkg/httpclient"
	"io"
	"net/http"
	"net/url"
	"strings"

	logger "github.com/sirupsen/logrus"

	"github.com/pkg/errors"
)

// KeyID represents a single key id on the KMS, equating to /keys/id
type KeyID struct {
	client *Client
	ID     string
}

// KeyObj is a represenation of the actual key
type KeyObj struct {
	Key []byte `json:"key,omitempty"`
}

// KeyInfo is a representation of key information
type KeyInfo struct {
	KeyID           string `json:"id,omitempty"`
	CipherMode      string `json:"mode,omitempty"`
	Algorithm       string `json:"algorithm,omitempty"`
	KeyLength       int    `json:"key_length,omitempty"`
	PaddingMode     string `json:"padding_mode,omitempty"`
	TransferPolicy  string `json:"transfer_policy,omitempty"`
	TransferLink    string `json:"transfer_link,omitempty"`
	DigestAlgorithm string `json:"digest_algorithm,omitempty"`
}

// Error is a error struct that contains error information thrown by the actual KMS
type Error struct {
	StatusCode int
	Message    string
}

func (k Error) Error() string {
	return fmt.Sprintf("kms-client: failed (HTTP Status Code: %d)\nMessage: %s", k.StatusCode, k.Message)
}

// Transfer performs a POST to /key/{id}/transfer to retrieve the actual key data from the KMS
func (k *KeyID) Transfer(saml []byte) ([]byte, error) {
	logger.Trace("pkg/kmsclient/key.go:Transfer() Entering")
	defer logger.Trace("pkg/kmsclient/key.go:Transfer() Leaving")

	logger.Info("pkg/kmsclient/key.go:GetKmsUser Retrieving kms user information")

	var samlBuf io.Reader
	if saml != nil {
		samlBuf = bytes.NewBuffer(saml)
	}
	keyXferURL, err := url.Parse(k.client.BaseURL)
	if err != nil {
		logger.Errorf("pkg/kmsclient/key.go:GetKmsUser Error in key transfer: %s", err.Error())
		return nil, errors.Wrap(err, "pkg/kmsclient/key.go:Transfer() Invalid key transfer URL passed")
	}
	req, err := http.NewRequest("POST", keyXferURL.String(), samlBuf)
	if err != nil {
		logger.Errorf("pkg/kmsclient/key.go:GetKmsUser Error in key transfer: %s", err.Error())
		return nil, errors.Wrap(err, "pkg/kmsclient/key.go:Transfer() Error creating key transfer request")
	}
	req.Header.Set("Accept", "application/octet-stream")
	req.Header.Set("Content-Type", "application/samlassertion+xml")
	rsp, err := httpclient.SendRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "pkg/kmsclient/key.go:Transfer() Error response to key transfer request")
	}
	return rsp, nil
}

// Keys represents the resource collection of Keys on the KMS
type Keys struct {
	client *Client
}

// Create sends a POST to /keys to create a new Key with the specified parameters
func (k *Keys) Create(key KeyInfo) (*KeyInfo, error) {
	logger.Trace("pkg/kmsclient/key.go:Create() Entering")
	defer logger.Trace("pkg/kmsclient/key.go:Create() Leaving")

	kiJSON, err := json.Marshal(&key)
	if err != nil {
		return nil, errors.Wrap(err, "pkg/kmsclient/key.go:CreateKey() Error marshalling key creation request")
	}
	baseURL, err := url.Parse(k.client.BaseURL)
	if err != nil {
		return nil, errors.Wrap(err, "pkg/kmsclient/key.go:CreateKey() Error parsing key creation URL")
	}
	keysURL, _ := url.Parse("keys")
	reqURL := baseURL.ResolveReference(keysURL)
	req, err := http.NewRequest("POST", reqURL.String(), bytes.NewBuffer(kiJSON))
	if err != nil {
		return nil, errors.Wrap(err, "pkg/kmsclient/key.go:CreateKey() Error creating key creation request")
	}

	// Set the request headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	rsp, err := httpclient.SendRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "pkg/kmsclient/key.go:CreateKey() Error response from key creation")
	}

	// Parse response
	var kiOut KeyInfo
	err = json.Unmarshal(rsp, &kiOut)
	if err != nil {
		return nil, errors.Wrap(err, "pkg/kmsclient/key.go:CreateKey() Response unmarshal failure")
	}
	return &kiOut, nil
}

// Retrieve performs a POST to /key/{id}/transfer to retrieve the actual key data from the KMS
func (k *KeyID) Retrieve(pubKey string) ([]byte, error) {
	logger.Trace("pkg/kmsclient/key.go:Retrieve() Entering")
	defer logger.Trace("pkg/kmsclient/key.go:Retrieve() Leaving")

	var keyValue KeyObj
	baseURL, err := url.Parse(k.client.BaseURL)
	if err != nil {
		return nil, errors.Wrap(err, "pkg/kmsclient/key.go:Retrieve() Failed parsing KMS base URL")
	}
	keyXferURL, err := url.Parse(fmt.Sprintf("keys/%s/transfer", k.ID))
	if err != nil {
		return nil, errors.Wrap(err, "pkg/kmsclient/key.go:Retrieve() Failed parsing key retrieve URL")
	}
	reqURL := baseURL.ResolveReference(keyXferURL)

	req, err := http.NewRequest("POST", reqURL.String(), strings.NewReader(pubKey))
	if err != nil {
		return nil, errors.Wrap(err, "pkg/kmsclient/key.go:Retrieve() Error creating key retrieval request")
	}

	// Set request headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "text/plain")


	rsp, err := httpclient.SendRequest(req)
	log.Debugf("pkg/kmsclient/key.go:Retrieve() HTTP response %s", string(rsp))
	if err != nil {
		return nil, errors.Wrap(err, "pkg/kmsclient/key.go:Retrieve() Error response from key retrieve request")
	}
	err = json.Unmarshal(rsp, &keyValue)
	if err != nil {
		return nil, errors.Wrap(err, "pkg/kmsclient/key.go:Retrieve() Error unmarshaling key retrieve response")
	}
	log.Debugf("pkg/kmsclient/key.go:Retrieve() After unmarshalling Key %v", keyValue)
	return keyValue.Key, nil
}
