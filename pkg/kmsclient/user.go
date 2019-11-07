/*
 * Copyright (C) 2019 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package kmsclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"

	"intel/isecl/wpm/pkg/httpclient"
)

//UserInfo is a representation of key information
type UserInfo struct {
	UserID         string `json:"id"`
	Username       string `json:"username"`
	TransferKeyPem string `json:"transfer_key_pem"`
}

// Users is a collection of User
type Users struct {
	Users []UserInfo `json:"users"`
}

type User struct {
	client *Client
}

// GetKmsUser is used to get the kms user information
func (k *Keys) GetKmsUser() (UserInfo, error) {
	log.Trace("pkg/kmsclient/user.go:GetKmsUser() Entering")
	defer log.Trace("pkg/kmsclient/user.go:GetKmsUser() Leaving")

	var userInfo UserInfo
	var users Users
	log.Info("pkg/kmsclient/user.go:GetKmsUser Retrieving kms user information")

	baseURL, err := url.Parse(k.client.BaseURL)
	if err != nil {
		return userInfo, err
	}
	getUserURL, err := url.Parse(fmt.Sprintf("users?usernameEqualTo=%s", k.client.Username))
	if err != nil {
		log.Errorf("pkg/kmsclient/user.go:GetKmsUser Error parsing URL %s | %v", getUserURL.String(), err)
		return userInfo, errors.New("pkg/kmsclient/user.go:GetKmsUser Error parsing KMS URL")
	}
	reqURL := baseURL.ResolveReference(getUserURL)
	req, err := http.NewRequest("GET", reqURL.String(), nil)
	if err != nil {
		log.Errorf("pkg/kmsclient/user.go:GetKmsUser Error creating request for URL %s | %v", reqURL.String(), err)
		return userInfo, errors.New("pkg/kmsclient/user.go:GetKmsUser Error creating user transfer request")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	rsp, err := httpclient.SendRequest(req)
	if err != nil {
		log.Errorf("pkg/kmsclient/user.go:GetKmsUser Error response for GetUser  %s | %v", reqURL.String(), err)
		return userInfo, errors.New("pkg/kmsclient/user.go:GetKmsUser Error dispatching request")
	}
	err = json.Unmarshal(rsp, &users)
	if err != nil {
		log.Errorf("pkg/kmsclient/user.go:GetKmsUser Error unmarshal GetUser response  %s | %v", string(rsp), err)
		return userInfo, errors.Wrap(err, "pkg/kmsclient/user.go:GetKmsUser Error unmarshalling GetUser response")
	}
	log.Infof("pkg/kmsclient/user.go:GetKmsUser Successfully retrieved KMS user %v", users)
	return users.Users[0], nil
}

func (k *Keys) RegisterUserPubKey(publicKey []byte, userID string) error {
	log.Trace("pkg/kmsclient/user.go:RegisterUserPubKey() Entering")
	defer log.Trace("pkg/kmsclient/user.go:RegisterUserPubKey() Leaving")

	baseURL, err := url.Parse(k.client.BaseURL)
	if err != nil {
		return errors.Wrap(err, "Failed to parse KMS URL")
	}
	keyXferURL, err := url.Parse(fmt.Sprintf("users/%s/transfer-key", userID))
	if err != nil {
		log.Infof("pkg/kmsclient/user.go:RegisterUserPubKey Failed to parse key transfer URL %s | %v", keyXferURL.String(), err)
		return errors.Wrap(err, "Failed to parse Key Transfer URL")
	}
	reqURL := baseURL.ResolveReference(keyXferURL)
	req, err := http.NewRequest("PUT", reqURL.String(), bytes.NewBuffer(publicKey))
	if err != nil {
		log.Infof("pkg/kmsclient/user.go:RegisterUserPubKey Failed to parse URL %s | %v", reqURL.String(), err)
		return errors.Wrap(err, "Failed to register user public key")
	}
	req.Header.Set("Content-Type", "application/x-pem-file")
	_, err = httpclient.SendRequest(req)
	if err != nil {
		log.Infof("pkg/kmsclient/user.go:RegisterUserPubKey Failed to parse URL %s | %v", reqURL.String(), err)
		return errors.Wrap(err, "Error sending request")
	}

	return nil
}
