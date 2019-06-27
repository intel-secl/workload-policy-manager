package containerimageflavor

/*
 *
 * @author arijitgh
 *
 */
import (
	"encoding/json"
	"errors"
	"intel/isecl/lib/flavor"
	flavorUtil "intel/isecl/lib/flavor/util"
	"intel/isecl/wpm/config"
	"intel/isecl/wpm/consts"
	"intel/isecl/wpm/pkg/util"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

const (
	DOCKER_CONTENT_TRUST_ENV_ENABLE        = "export DOCKER_CONTENT_TRUST=1"
	DOCKER_CONTENT_TRUST_ENV_CUSTOM_NOTARY = DOCKER_CONTENT_TRUST_ENV_ENABLE + "; export DOCKER_CONTENT_TRUST_SERVER="
	DEFAULT_NOTARY_SERVER_URL              = "https://notary.docker.io"
)

//CreateContainerImageFlavor is used to create flavor of a container image
func CreateContainerImageFlavor(imageName, tag, dockerFilePath, buildDir,
	keyID string, encryptionRequired, integrityEnforced bool, notaryServerURL, outputFlavorFilePath string) (string, error) {
	var err error
	var wrappedKey []byte
	var keyURLString string

	//Return usage if input params are provided incorrectly
	if len(strings.TrimSpace(imageName)) <= 0 {
		return "", errors.New("missing image name")
	}
	flavorLabel := imageName + ":" + tag
	if len(strings.TrimSpace(dockerFilePath)) > 0 || len(strings.TrimSpace(buildDir)) > 0 {

		//Error if Dockerfile specified doesn't exist
		_, err = os.Stat(dockerFilePath)
		if os.IsNotExist(err) {
			return "", errors.New("Dockerfile does not exist")
		}

		//Error if build directory specified doesn't exist
		_, err = os.Stat(buildDir)
		if os.IsNotExist(err) {
			return "", errors.New("docker build directory does not exist")
		}

		//Encrypt the image with the key
		if encryptionRequired {
			wrappedKey, keyURLString, err = util.FetchKey(keyID)
			if keyID == "" {
				keyID = strings.Split(strings.Split(keyURLString, "/transfer")[0], config.Configuration.Kms.APIURL+"keys/")[1]
			}

			wrappedKeyFilePath := "/tmp/wrappedKey_" + keyID
			os.Create(wrappedKeyFilePath)
			err = ioutil.WriteFile(wrappedKeyFilePath, wrappedKey, 0600)

			//Run docker build command to build encrypted image
			cmd := exec.Command("docker", "build", "--no-cache", "-t", imageName+":"+tag,
				"--imgcrypt-opt", "RequiresConfidentiality=true", "--imgcrypt-opt", "KeyFilePath="+wrappedKeyFile.Name(),
				"--imgcrypt-opt", "KeyType=key-type-kms", "-f", dockerFilePath, buildDir)

			_, err = cmd.CombinedOutput()
			if err != nil {
				return "", errors.New("could not build container image" + err.Error())
			}

		} else {
			//Run docker build command to build plain image
			_, err = exec.Command("docker", "build", "--no-cache", "-t", imageName+":"+tag,
				"-f", dockerFilePath, buildDir).CombinedOutput()
			if err != nil {
				return "", errors.New("could not build container image")
			}
		}
	} else {
		_, err = exec.Command("docker", "inspect", "--type=image", imageName+":"+tag).CombinedOutput()
		if err != nil {
			return "", errors.New("Could not find image with name:" + imageName + " and tag:" + tag + "\nImage should be present locally")
		}
	}

	if integrityEnforced && notaryServerURL == "" {
		//add public notary server url
		notaryServerURL = DEFAULT_NOTARY_SERVER_URL
	}

	//Create image flavor
	containerImageFlavor, err := flavor.GetContainerImageFlavor(flavorLabel, encryptionRequired, keyURLString, integrityEnforced, notaryServerURL)
	if err != nil {
		return "", errors.New("error creating image flavor:" + err.Error())
	}

	//Marshall the image flavor to a JSON string
	containerImageFlavorJSON, err := json.Marshal(containerImageFlavor)
	if err != nil {
		return "", errors.New("error while marshalling image flavor:" + err.Error())
	}

	signedFlavor, err := flavorUtil.GetSignedFlavor(string(containerImageFlavorJSON), consts.FlavorSigningKeyPath)

	//If no output flavor file path was specified, return the marshalled image flavor
	if len(strings.TrimSpace(outputFlavorFilePath)) <= 0 {
		return signedFlavor, nil
	}

	//Otherwise, write it to the specified file
	err = ioutil.WriteFile(outputFlavorFilePath, []byte(signedFlavor), 0600)
	if err != nil {
		return "", errors.New("error writing image flavor to output file")
	}
	return "", err
}
