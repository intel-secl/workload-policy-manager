/*
 * Copyright (C) 2019 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package containerimageflavor

/*
 *
 * @author arijitgh
 *
 */
import (
	"encoding/json"
	"fmt"
	cLog "intel/isecl/lib/common/v3/log"
	"intel/isecl/lib/flavor/v3"
	flavorUtil "intel/isecl/lib/flavor/v3/util"
	"intel/isecl/wpm/v3/consts"
	"intel/isecl/wpm/v3/pkg/util"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/pkg/errors"
)

const (
	DOCKER_CONTENT_TRUST_ENV_ENABLE        = "export DOCKER_CONTENT_TRUST=1"
	DOCKER_CONTENT_TRUST_ENV_CUSTOM_NOTARY = DOCKER_CONTENT_TRUST_ENV_ENABLE + "; export DOCKER_CONTENT_TRUST_SERVER="
	DEFAULT_NOTARY_SERVER_URL              = "https://notary.docker.io"
)

var (
	log    = cLog.GetDefaultLogger()
	secLog = cLog.GetSecurityLogger()
)

//CreateContainerImageFlavor is used to create flavor of a container image
func CreateContainerImageFlavor(imageName, tag, dockerFilePath, buildDir,
	keyID string, encryptionRequired, integrityEnforced bool, notaryServerURL, outputFlavorFilePath string) (string, error) {
	log.Trace("pkg/containerimageflavor/create_container_image_flavors.go:CreateContainerImageFlavor() Entering")
	defer log.Trace("pkg/containerimageflavor/create_container_image_flavors.go:CreateContainerImageFlavor() Leaving")

	var err error
	var wrappedKey []byte
	var keyUrlString string
	signedFlavor := ""

	// set logger fields
	log = log.WithFields(logrus.Fields{
		"imageName":            imageName,
		"encryptionRequired":   encryptionRequired,
		"integrityEnforced":    integrityEnforced,
		"dockerFilePath":       dockerFilePath,
		"outputFlavorFilePath": outputFlavorFilePath,
		"keyID":                keyID,
	})

	//Return usage if input params are provided incorrectly
	if len(strings.TrimSpace(imageName)) <= 0 {
		return signedFlavor, errors.New("Missing image name")
	}
	flavorLabel := imageName + ":" + tag
	if len(strings.TrimSpace(dockerFilePath)) > 0 || len(strings.TrimSpace(buildDir)) > 0 {

		//Error if Dockerfile specified doesn't exist
		_, err = os.Stat(dockerFilePath)
		if os.IsNotExist(err) {
			return signedFlavor, errors.Wrap(err, "Dockerfile does not exist")
		}

		//Error if build directory specified doesn't exist
		_, err = os.Stat(buildDir)
		if os.IsNotExist(err) {
			return signedFlavor, errors.Wrap(err, "Docker build directory does not exist")
		}

		//Encrypt the image with the key
		if encryptionRequired {
			wrappedKey, keyUrlString, err = util.FetchKey(keyID, "")
			if err != nil {
				return signedFlavor, errors.Wrap(err, "Error fetching KMS key")
			}
			// We infer the keyID from the keyUrlString
			if keyID == "" {
				keyUrl, _ := url.Parse(keyUrlString)
				re := regexp.MustCompile("(?i)([0-9A-F]{8}-[0-9A-F]{4}-4[0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12})")
				keyID = re.FindString(keyUrl.Path)
			}

			wrappedKeyFileName := "wrappedKey_" + keyID + "_"
			wrappedKeyFile, err := ioutil.TempFile("/tmp", wrappedKeyFileName)
			if err != nil {
				return signedFlavor, errors.Wrap(err, "Unable to create wrapped key file")
			}
			if _, err = wrappedKeyFile.Write(wrappedKey); err != nil {
				return signedFlavor, errors.Wrap(err, "Unable to write wrapped key to file")
			}
			defer os.Remove(wrappedKeyFile.Name())

			//Run docker build command to build encrypted image
			cmd := exec.Command("docker", "build", "--no-cache", "-t", imageName+":"+tag,
				"--imgcrypt-opt", "RequiresConfidentiality=true", "--imgcrypt-opt", "KeyFilePath="+wrappedKeyFile.Name(),
				"--imgcrypt-opt", "KeyType=key-type-kms", "-f", dockerFilePath, buildDir)

			log.Debugf("pkg/containerimageflavor/create_container_image_flavor.go:CreateContainerImageFlavor() Docker build command %s",
				fmt.Sprintf("docker build --no-cache -t %s:%s --imgcrypt-opt RequiresConfidentiality=true --imgcrypt-opt KeyFilePath=%s "+
					"--imgcrypt-opt KeyType=key-type-kms -f %s %s",
					imageName, tag, wrappedKeyFile.Name(), dockerFilePath, buildDir))
			_, err = cmd.CombinedOutput()
			if err != nil {
				if strings.Contains(fmt.Sprint(cmd.Stderr), "unknown flag: --imgcrypt-opt") {
					log.Errorf("Failed to build container image: %s", fmt.Sprint(cmd.Stderr) )
					return signedFlavor, errors.Wrap(errors.New("Secure Docker Daemon is not properly " +
						"installed. Check logs for more details"),  "Unable to build container image with " +
						"encryption")
				}
				return signedFlavor, errors.Wrap(err, "Unable to build container image with encryption")
			}

		} else {
			//Run docker build command to build plain image
			_, err = exec.Command("docker", "build", "--no-cache", "-t", imageName+":"+tag,
				"-f", dockerFilePath, buildDir).CombinedOutput()
			if err != nil {
				return signedFlavor, errors.Wrap(err, "Unable to build container image")
			}
		}
	} else {
		_, err = exec.Command("docker", "inspect", "--type=image", imageName+":"+tag).CombinedOutput()
		if err != nil {
			return signedFlavor, errors.Wrap(err, "Unable to find image with name: "+imageName+" and tag: "+tag+"\nImage should be present locally")
		}
	}

	if integrityEnforced && notaryServerURL == "" {
		//add public notary server url
		notaryServerURL = DEFAULT_NOTARY_SERVER_URL
		log.Infof("Using default notary server URL: %s", notaryServerURL)
	}

	//Create image flavor
	containerImageFlavor, err := flavor.GetContainerImageFlavor(flavorLabel, encryptionRequired, keyUrlString, integrityEnforced, notaryServerURL)
	if err != nil {
		return signedFlavor, errors.Wrap(err, "Error while creating image flavor: " + err.Error())
	}

	//Marshall the image flavor to a JSON string
	containerImageFlavorJSON, err := json.Marshal(containerImageFlavor)
	if err != nil {
		return signedFlavor, errors.Wrap(err, "Error while marshalling image flavor: " + err.Error())
	}

	signedFlavor, err = flavorUtil.GetSignedFlavor(string(containerImageFlavorJSON), consts.FlavorSigningKeyPath)
	if err != nil {
		return signedFlavor, errors.Wrap(err, "Error while signing image flavor: " + err.Error())
	}

	//If no output flavor file path was specified, return the marshalled image flavor
	if len(strings.TrimSpace(outputFlavorFilePath)) <= 0 {
		return signedFlavor, nil
	}

	//Otherwise, write it to the specified file
	err = ioutil.WriteFile(outputFlavorFilePath, []byte(signedFlavor), 0600)
	if err != nil {
		return signedFlavor, errors.Wrap(err, "Error writing image flavor to output file: " + err.Error())
	}
	return signedFlavor, nil
}
