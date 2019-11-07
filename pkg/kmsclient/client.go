/*
 * Copyright (C) 2019 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package kmsclient

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type ISO8601Time struct {
	time.Time
}

const ISO8601Layout = "2006-01-02T15:04:05-0700"

func (t *ISO8601Time) MarshalJSON() ([]byte, error) {
	log.Trace("pkg/kmsclient/client.go:MarshalJSON() Entering")
	defer log.Trace("pkg/kmsclient/client.go:MarshalJSON() Leaving")

	tstr := t.Format(ISO8601Layout)
	return []byte(strconv.Quote(tstr)), nil
}

func (t *ISO8601Time) UnmarshalJSON(b []byte) (err error) {
	log.Trace("pkg/kmsclient/client.go:UnmarshalJSON() Entering")
	defer log.Trace("pkg/kmsclient/client.go:UnmarshalJSON() Leaving")

	t.Time, err = time.Parse(ISO8601Layout, strings.Trim(string(b), "\""))
	return errors.Wrap(err, "pkg/kmsclient/client.go:UnmarshalJSON Invalid timestamp")
}

// AuthToken issued by KMS
type AuthToken struct {
	AuthorizationToken string        `json:"authorization_token"`
	AuthorizationDate  ISO8601Time   `json:"authorization_date"`
	NotAfter           ISO8601Time   `json:"not_after"`
	Faults             []interface{} `json:"faults"`
}

// A Client is defines parameters to connect and Authenticate with a KMS
type Client struct {
	// BaseURL specifies the URL base for the KMS, for example https://keymanagement.server/v1
	BaseURL string
	// HttpClient used for outbound requests to KMS
	HTTPClient *http.Client
	// Username for KMS client account
	Username string
	// Password for KMS client account
	Password string
}

// Key returns a reference to a KeyID on the KMS. It is a reference only, and does not immediately contain any Key information.
func (c *Client) Key(uuid string) *KeyID {
	log.Trace("pkg/kmsclient/client.go:Key() Entering")
	defer log.Trace("pkg/kmsclient/client.go:Key() Leaving")

	return &KeyID{client: c, ID: uuid}
}

// Keys returns a sub client that operates on KMS /keys endpoints, such as creating a new key
func (c *Client) Keys() *Keys {
	log.Trace("pkg/kmsclient/client.go:Keys() Entering")
	defer log.Trace("pkg/kmsclient/client.go:Keys() Leaving")

	return &Keys{client: c}
}
