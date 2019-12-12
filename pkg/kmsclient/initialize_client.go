/*
 * Copyright (C) 2019 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package kmsclient

import (
	"intel/isecl/lib/clients"
	config "intel/isecl/wpm/config"
	"intel/isecl/wpm/consts"

	cLog "intel/isecl/lib/common/log"
)

var (
	log    = cLog.GetDefaultLogger()
	secLog = cLog.GetSecurityLogger()
)

func InitializeKMSClient() (*Client, error) {
	log.Trace("pkg/kmsclient/initialize_client.go:InitializeKMSClient() Entering")
	defer log.Trace("pkg/kmsclient/initialize_client.go:InitializeKMSClient() Leaving")

	var kc *Client
	hc, err := clients.HTTPClientWithCADir(consts.TrustedCaCertsDir)
	kc = &Client{
		BaseURL:    config.Configuration.Kms.APIURL,
		HTTPClient: hc,
	}

	if err != nil {
		return nil, err
	}
	return kc, nil
}
