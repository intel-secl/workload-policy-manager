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
		return userInfo, errors.Wrapf(err, "pkg/kmsclient/user.go:GetKmsUser Error parsing URL %s", getUserURL.String())
	}
	reqURL := baseURL.ResolveReference(getUserURL)
	req, err := http.NewRequest("GET", reqURL.String(), nil)
	if err != nil {
		return userInfo, errors.Wrapf(err, "pkg/kmsclient/user.go:GetKmsUser Error creating user transfer request for URL %s", reqURL.String())
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	rsp, err := httpclient.SendRequest(req)
	if err != nil {
		return userInfo, errors.Wrapf(err, "pkg/kmsclient/user.go:GetKmsUser Error response for GetUser URL %s", reqURL.String())
	}
	err = json.Unmarshal(rsp, &users)
	if err != nil {
		return userInfo, errors.Wrapf(err, "pkg/kmsclient/user.go:GetKmsUser Error unmarshal GetUser response")
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
		return errors.Wrapf(err, "pkg/kmsclient/user.go:RegisterUserPubKey Failed to parse user envelope key register URL %s", keyXferURL.String())
	}
	reqURL := baseURL.ResolveReference(keyXferURL)
	log.Debugf("pkg/kmsclient/user.go:RegisterUserPubKey Envelope key register URL %s", reqURL.String())
	req, err := http.NewRequest("PUT", reqURL.String(), bytes.NewBuffer(publicKey))
	if err != nil {
		return errors.Wrap(err, "pkg/kmsclient/user.go:RegisterUserPubKey Failed to create envelope key register request")
	}
	req.Header.Set("Content-Type", "application/x-pem-file")
	_, err = httpclient.SendRequest(req)
	if err != nil {
		return errors.Wrap(err, "pkg/kmsclient/user.go:RegisterUserPubKey Register Envelope Key failure response")
	}
	log.Info("pkg/kmsclient/user.go:GetKmsUser Successfully registered envelope key for KMS user")
	return nil
}
