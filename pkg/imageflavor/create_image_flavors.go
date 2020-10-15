/*
 * Copyright (C) 2019 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package imageflavor

/*
 *
 * @author srege
 *
 */
import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	cLog "intel/isecl/lib/common/v3/log"
	"intel/isecl/lib/flavor/v3"
	flavorUtil "intel/isecl/lib/flavor/v3/util"
	"intel/isecl/wpm/v3/consts"
	"intel/isecl/wpm/v3/pkg/util"

	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	log    = cLog.GetDefaultLogger()
	secLog = cLog.GetSecurityLogger()
)

//CreateImageFlavor is used to create flavor of an encrypted image
func CreateImageFlavor(flavorLabel string, outputFlavorFilePath string, inputImageFilePath string, outputEncImageFilePath string,
	keyID string, integrityRequired bool) (string, error) {
	log.Trace("pkg/imageflavor/create_image_flavors.go:CreateImageFlavor() Entering")
	defer log.Trace("pkg/imageflavor/create_image_flavors.go:CreateImageFlavor() Leaving")

	var err error
	var wrappedKey []byte
	var keyUrlString string
	encRequired := true
	imageFilePath := inputImageFilePath

	//Determine if encryption is required
	outputEncImageFilePath = strings.TrimSpace(outputEncImageFilePath)
	if len(outputEncImageFilePath) <= 0 {
		encRequired = false
	}

	// set logger fields
	log = log.WithFields(logrus.Fields{
		"flavorLabel":            flavorLabel,
		"encryptionRequired":     encRequired,
		"integrityrequired":      integrityRequired,
		"inputImageFilePath":     inputImageFilePath,
		"outputFlavorFilePath":   outputFlavorFilePath,
		"outputEncImageFilePath": keyID,
	})

	//Error if image specified doesn't exist
	_, err = os.Stat(inputImageFilePath)
	if os.IsNotExist(err) {
		return "", errors.Wrap(err, "I/O error reading image file: "+err.Error())
	}

	//Encrypt the image with the key
	if encRequired {
		// fetch the key to encrypt the image
		wrappedKey, keyUrlString, err = util.FetchKey(keyID, "")
		if err != nil {
			return "", errors.Wrap(err, "Fetch key failed: "+err.Error())
		}
		// encrypt the image with key retrieved from KBS
		err = util.Encrypt(inputImageFilePath, consts.EnvelopePrivatekeyLocation, outputEncImageFilePath, wrappedKey)
		if err != nil {
			return "", errors.Wrap(err, "Image encryption failed: "+err.Error())
		}
		imageFilePath = outputEncImageFilePath
	}

	//Check the encrypted image output file
	imageFile, err := ioutil.ReadFile(imageFilePath)
	if err != nil {
		return "", errors.Wrap(err, "I/O Error creating encrypted image file: "+err.Error())
	}

	//Take the digest of the encrypted image
	digest := sha512.Sum384([]byte(imageFile))

	//Create image flavor
	imageFlavor, err := flavor.GetImageFlavor(flavorLabel, encRequired, keyUrlString, base64.StdEncoding.EncodeToString(digest[:]))
	if err != nil {
		return "", errors.Wrap(err, "Error creating image flavor: "+err.Error())
	}

	//Marshall the image flavor to a JSON string
	imageFlavorJSON, err := json.Marshal(imageFlavor)
	if err != nil {
		return "", errors.Wrap(err, "Error while marshalling image flavor: "+err.Error())
	}

	signedFlavor, err := flavorUtil.GetSignedFlavor(string(imageFlavorJSON), consts.FlavorSigningKeyPath)
	if err != nil {
		return "", errors.Wrap(err, "Error signing flavor for image: "+err.Error())
	}

	log.Info("pkg/imageflavor/create_image_flavors.go:CreateImageFlavor() Successfully created image flavor")
	log.Debugf("pkg/imageflavor/create_image_flavors.go:CreateImageFlavor() Successfully created image flavor %s", signedFlavor)

	//If no output flavor file path was specified, return the marshalled image flavor
	if len(strings.TrimSpace(outputFlavorFilePath)) <= 0 {
		return signedFlavor, nil
	}

	//Otherwise, write it to the specified file
	err = ioutil.WriteFile(outputFlavorFilePath, []byte(signedFlavor), 0600)
	if err != nil {
		return "", errors.Wrapf(err, "I/O Error writing image flavor to output file %s", outputFlavorFilePath)
	}

	log.Info("pkg/imageflavor/create_image_flavors.go:CreateImageFlavor() Successfully wrote image flavor to file")
	return "", nil
}
