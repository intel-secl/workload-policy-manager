/*
 * Copyright (C) 2019 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package httpclient

import (
	"io/ioutil"
	"net/http"
	"sync"

	"intel/isecl/wpm/v3/config"
	"intel/isecl/wpm/v3/consts"

	"intel/isecl/lib/clients/v3"
	"intel/isecl/lib/clients/v3/aas"

	"github.com/pkg/errors"

	cLog "intel/isecl/lib/common/v3/log"
)

var (
	log       = cLog.GetDefaultLogger()
	secLog    = cLog.GetSecurityLogger()
	aasClient = aas.NewJWTClient(config.Configuration.Aas.APIURL)
	aasRWLock = sync.RWMutex{}
)

// init sets up a single shared AAS client for use
func init() {
	log.Trace("pkg/httpclient/send_http_request.go:init() Entering")
	defer log.Trace("pkg/httpclient/send_http_request.go:init() Leaving")

	aasRWLock.Lock()
	if aasClient.HTTPClient == nil {
		c, err := clients.HTTPClientWithCADir(consts.TrustedCaCertsDir)
		if err != nil {
			return
		}
		aasClient.HTTPClient = c
	}
	aasRWLock.Unlock()
}

// addJWTToken fetches and adds JWT token for the AAS client
func addJWTToken(req *http.Request) error {
	log.Trace("pkg/httpclient/send_http_request.go:addJWTToken() Entering")
	defer log.Trace("pkg/httpclient/send_http_request.go:addJWTToken() Leaving")

	if aasClient.BaseURL == "" {
		aasClient = aas.NewJWTClient(config.Configuration.Aas.APIURL)
		if aasClient.HTTPClient == nil {
			c, err := clients.HTTPClientWithCADir(consts.TrustedCaCertsDir)
			if err != nil {
				return errors.Wrap(err, "clients/send_http_request.go:addJWTToken() Error initializing http client")
			}
			aasClient.HTTPClient = c
		}
	}

	aasRWLock.RLock()
	jwtToken, err := aasClient.GetUserToken(config.Configuration.Wpm.Username)
	aasRWLock.RUnlock()
	// something wrong
	if err != nil {
		// lock aas with w lock
		aasRWLock.Lock()
		// check if other thread fix it already
		jwtToken, err = aasClient.GetUserToken(config.Configuration.Wpm.Username)
		// it is not fixed
		if err != nil {
			// these operation cannot be done in init() because it is not sure
			// if config.Configuration is loaded at that time
			aasClient.AddUser(config.Configuration.Wpm.Username, config.Configuration.Wpm.Password)
			aasClient.FetchTokenForUser(config.Configuration.Wpm.Username)
			jwtToken, err = aasClient.GetUserToken(config.Configuration.Wpm.Username)
			if err != nil {
				return errors.Wrap(err, "pkg/httpclient/send_http_request:addJWTToken() Failed to fetch JWT token")
			}
		}
		aasRWLock.Unlock()
	}
	req.Header.Set("Authorization", "Bearer "+string(jwtToken))
	log.Info("pkg/httpclient/send_http_request:addJWTToken() Successfully fetched the JWT token for the user")
	return nil
}

//SendRequest method is used to create an cert-chained http client and send the request with JWT-auth token
func SendRequest(req *http.Request) ([]byte, error) {
	log.Trace("pkg/httpclient/send_http_request.go:SendRequest() Entering")
	defer log.Trace("pkg/httpclient/send_http_request.go:SendRequest() Leaving")

	// Fetch a HTTP client that has validated cert chain
	client, err := clients.HTTPClientWithCADir(consts.TrustedCaCertsDir)
	if err != nil {
		return nil, errors.Wrap(err, "pkg/httpclient/send_http_request:SendRequest() Error intitializing HTTP client")
	}

	// Add the JWT to the auth header in the request
	err = addJWTToken(req)
	if err != nil {
		return nil, errors.Wrap(err, "pkg/httpclient/send_http_request:SendRequest() Error fetching JWT auth token")
	}

	// Dispatch request
	response, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "pkg/httpclient/send_http_request:SendRequest() HTTP Response failure")
	}
	defer response.Body.Close()

	// Reauthenticate to fetch a fresh token
	if response.StatusCode == http.StatusUnauthorized {
		// fetch token and try again
		aasRWLock.Lock()
		aasClient.FetchAllTokens()
		aasRWLock.Unlock()
		err = addJWTToken(req)
		if err != nil {
			return nil, errors.Wrap(err, "pkg/httpclient/send_http_request:SendRequest() Failed fetching JWT auth token")
		}
		response, err = client.Do(req)
		if err != nil {
			return nil, errors.Wrap(err, "pkg/httpclient/send_http_request:SendRequest() Error from response")
		}
	}

	// 404 -- key not found, URL invalid - how do we interpret this?
	if response.StatusCode == http.StatusNotFound {
		return nil, errors.Wrap(err, "pkg/httpclient/send_http_request:SendRequest() Resource not found")
	}

	//create byte array of HTTP response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "pkg/httpclient/send_http_request:SendRequest() Error reading response body")
	}

	log.Info("pkg/httpclient/send_http_request.go:SendRequest() Successfully processed http request")
	return body, nil
}

//SendBasicAuthRequest method is used to create an cert-chained http client and send the request with basic-auth credentials
func SendBasicAuthRequest(req *http.Request) ([]byte, error) {
	log.Trace("pkg/httpclient/send_http_request.go:SendBasicAuthRequest() Entering")
	defer log.Trace("pkg/httpclient/send_http_request.go:SendBasicAuthRequest() Leaving")

	// Fetch a HTTP client that has validated cert chain
	client, err := clients.HTTPClientWithCADir(consts.TrustedCaCertsDir)
	if err != nil {
		return nil, errors.Wrap(err, "pkg/httpclient/send_http_request:SendBasicAuthRequest() Error intitializing HTTP client")
	}

	// Dispatch request
	response, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "pkg/httpclient/send_http_request:SendBasicAuthRequest() HTTP Response failure")
	}
	defer response.Body.Close()

	//create byte array of HTTP response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "pkg/httpclient/send_http_request:SendBasicAuthRequest() Error reading response body")
	}

	log.Info("pkg/httpclient/send_http_request.go:SendBasicAuthRequest() Successfully processed http request")
	return body, nil
}
