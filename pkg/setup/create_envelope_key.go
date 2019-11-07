/*
 * Copyright (C) 2019 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package setup

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	cLog "intel/isecl/lib/common/log"
	csetup "intel/isecl/lib/common/setup"
	"intel/isecl/wpm/consts"
	"os"

	"github.com/pkg/errors"
)

var (
	log    = cLog.GetDefaultLogger()
	secLog = cLog.GetSecurityLogger()
)

type CreateEnvelopeKey struct {
}

// ValidateCreateKey method is used to check if the envelope keys exists on disk
func (ek CreateEnvelopeKey) Validate(c csetup.Context) error {
	log.Trace("pkg/setup/create_envelope_key.go:Validate() Entering")
	defer log.Trace("pkg/setup/create_envelope_key.go:Validate() Leaving")

	log.Info("pkg/setup/create_envelope_key.go:Validate() Validating envelope key creation")

	_, err := os.Stat(consts.EnvelopePrivatekeyLocation)
	if os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Private key does not exist")
		return errors.Wrap(err, "pkg/setup/create_envelope_key.go:Validate() Private key does not exist")
	}

	_, err = os.Stat(consts.EnvelopePublickeyLocation)
	if os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Public key does not exist")
		return errors.Wrap(err, "pkg/setup/create_envelope_key.go:Validate() Public key does not exist")
	}
	return nil
}

func (ek CreateEnvelopeKey) Run(c csetup.Context) error {
	log.Trace("pkg/setup/create_envelope_key.go:Run() Entering")
	defer log.Trace("pkg/setup/create_envelope_key.go:Run() Leaving")

	log.Info("Creating envelope key")

	bitSize := consts.DefaultKeyAlgorithmLength
	keyPair, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while generating new RSA key pair")
		return errors.Wrap(err, "pkg/setup/create_envelope_key.go:Run() Error while generating a new RSA key pair")
	}

	// save private key
	privateKey := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(keyPair),
	}

	privateKeyFile, err := os.Create(consts.EnvelopePrivatekeyLocation)
	if err != nil {
		fmt.Fprintf(os.Stderr, "I/O error while saving private key file")
		return errors.Wrap(err, "pkg/setup/create_envelope_key.go:Run() I/O error while saving private key file")
	}
	defer privateKeyFile.Close()
	err = pem.Encode(privateKeyFile, privateKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "I/O error while encoding private key file")
		return errors.Wrap(err, "pkg/setup/create_envelope_key.go:Run() Error while encoding the private key.")
	}

	// save public key
	publicKey := &keyPair.PublicKey

	pubkeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "I/O error while encoding private key file")
		return errors.Wrap(err, "pkg/setup/create_envelope_key.go:Run() Error while marshalling the public key.")
	}
	var publicKeyInPem = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubkeyBytes,
	}

	publicKeyFile, err := os.Create(consts.EnvelopePublickeyLocation)
	if err != nil {
		fmt.Fprintf(os.Stderr, "I/O error while encoding public envelope key file")
		return errors.Wrap(err, "pkg/setup/create_envelope_key.go:Run() Error while creating a new file. ")
	}
	defer publicKeyFile.Close()

	err = pem.Encode(publicKeyFile, publicKeyInPem)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while encoding the public envelope key")
		return errors.Wrap(err, "pkg/setup/create_envelope_key.go:Run() Error while encoding the public key.")
	}
	return nil
}
