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
	"intel/isecl/wpm/config"
	"intel/isecl/wpm/pkg/util"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

const (
	DOCKER_CONTENT_TRUST_ENV_ENABLE        = "export DOCKER_CONTENT_TRUST=1"
	DOCKER_CONTENT_TRUST_ENV_CUSTOM_NOTARY = DOCKER_CONTENT_TRUST_ENV_ENABLE + "; export DOCKER_CONTENT_TRUST_SERVER="
	DEFAULT_NOTARY_SERVER_URL              = "https://notary.docker.io"
)

//CreateContainerImageFlavor is used to create flavor of a container image
func CreateContainerImageFlavor(imageName, tagName, dockerFilePath, buildDir,
	keyID string, encryptionRequired, integrityEnforced bool, notaryServerURL, outputFlavorFilePath string) (string, error) {
	var err error
	var wrappedKey []byte
	var keyURLString string

	//Return usage if input params are provided incorrectly
	if len(strings.TrimSpace(imageName)) <= 0 {
		return "", errors.New(Usage())
	}

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
			//If the key ID is specified, make sure it's a valid UUID
			if len(strings.TrimSpace(keyID)) > 0 && !isValidUUID(keyID) {
				return "", errors.New("incorrectly formatted key ID")
			}

			wrappedKey, keyURLString, err = util.FetchKey(keyID)
			if keyID == "" {
				keyID = strings.TrimLeft(strings.TrimRight(keyURLString, "/transfer"), config.Configuration.Kms.APIURL+"keys/")
			}

			wrappedKeyFileName := "wrappedKey_" + keyID + "_"
			wrappedKeyFile, err := ioutil.TempFile("/tmp", wrappedKeyFileName)
			if err != nil {
				return "", errors.New("could not create wrapped key file")
			}
			if _, err = wrappedKeyFile.Write(wrappedKey); err != nil {
				return "", errors.New("could not write the wrapped key in to the file")
			}
			defer os.Remove(wrappedKeyFile.Name())

			//Run docker build command to build encrypted image
			cmd := exec.Command("docker", "build", "--no-cache", "-t", imageName+":"+tagName,
				"--storage-opt", "RequiresConfidentiality=true", "--storage-opt", "KeyFilePath="+wrappedKeyFile.Name(),
				"--squash", "-f", dockerFilePath, buildDir)

			_, err = cmd.CombinedOutput()
			if err != nil {
				return "", errors.New("could not build container image with encryption" + err.Error())
			}

		} else {
			//Run docker build command to build plain image
			_, err = exec.Command("docker", "build", "--no-cache", "-t", imageName+":"+tagName,
				"-f", dockerFilePath, buildDir).CombinedOutput()
			if err != nil {
				return "", errors.New("could not build container image")
			}
		}
	} else {
		if integrityEnforced {
			if notaryServerURL == "" {
				//add public notary server url
				notaryServerURL = DEFAULT_NOTARY_SERVER_URL
			}
			//Pull signed image
			_, err = exec.Command(DOCKER_CONTENT_TRUST_ENV_CUSTOM_NOTARY+notaryServerURL, ";",
				"docker", "pull", imageName+":"+tagName).CombinedOutput()
			if err != nil {
				return "", errors.New("could not pull signed docker image:" + err.Error())
			}
		} else {
			//Pull plain image
			_, err = exec.Command("docker", "pull", imageName+":"+tagName).CombinedOutput()
			if err != nil {
				return "", errors.New("could not pull docker image:" + err.Error())
			}
		}
	}

	flavorLabel := imageName + ":" + tagName

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

	//If no output flavor file path was specified, return the marshalled image flavor
	if len(strings.TrimSpace(outputFlavorFilePath)) <= 0 {
		return string(containerImageFlavorJSON), nil
	}

	//Otherwise, write it to the specified file
	err = ioutil.WriteFile(outputFlavorFilePath, containerImageFlavorJSON, 0600)
	if err != nil {
		return "", errors.New("error writing image flavor to output file")
	}
	return "", err
}

//Regex match to determine if string is valid UUID
//TODO: move to common lib validation package
func isValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}

//Usage command line usage string
func Usage() string {
	return "usage: wpm create-container-image-flavor [-n img-name] [-t tag] [-f dockerFile] [-d build-dir] [-k keyId]\n" +
		"                            [-enc] [-enforce] [-s notaryServer] [-o out-file]\n" +
		"  -i,       --img-name                     container image name\n" +
		"  -t,       --tag                          (optional)container image tag name\n" +
		"  -f,       --docker-file                  (optional) container file path\n" +
		"                                           to build the container image\n" +
		"  -d,       --build-dir                    (optional) build directory to\n" +
		"                                           build the container image\n" +
		"  -k,       --key-id                       (optional) existing key ID\n" +
		"                                           if not specified, a new key is generated\n" +
		"  -e,     --encryption-required            (optional) boolean parameter specifies if\n" +
		"                                           container image needs to be encrypted\n" +
		"  -s, 	   --integrity-enforced             (optional) boolean parameter specifies if\n" +
		"                                           container image should be signed\n" +
		"  -n,       --notary-server                (optional) specify notary server url\n" +
		"  -o,       --out-file                     (optional) specify output file path\n"

}
