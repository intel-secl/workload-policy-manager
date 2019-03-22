package imageflavor

/*
 *
 * @author srege
 *
 */
import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"intel/isecl/lib/flavor"
	"intel/isecl/wpm/consts"
	"intel/isecl/wpm/pkg/util"
	"intel/isecl/lib/common/validation"
	"io/ioutil"
	"os"
	"strings"
)

//CreateImageFlavor is used to create flavor of an encrypted image
func CreateImageFlavor(flavorLabel string, outputFlavorFilePath string, inputImageFilePath string, outputEncImageFilePath string,
	keyID string, integrityRequired bool) (string, error) {

	var err error
	var wrappedKey []byte
	var keyURLString string
	encRequired := true
	imageFilePath := inputImageFilePath

	//Return usage if input params are provided incorrectly
	if len(strings.TrimSpace(flavorLabel)) <= 0 || len(strings.TrimSpace(inputImageFilePath)) <= 0 {
		return "", errors.New(Usage())
	}

	inputArr := []string{flavorLabel, outputFlavorFilePath, inputImageFilePath, outputEncImageFilePath}
	if validationErr := validation.ValidateStrings(inputArr); validationErr != nil {
		return "", validationErr
	}

	//Determine if encryption is required
	outputEncImageFilePath = strings.TrimSpace(outputEncImageFilePath)
	if len(outputEncImageFilePath) <= 0 {
		encRequired = false
	}

	//Error if image specified doesn't exist
	_, err = os.Stat(inputImageFilePath)
	if os.IsNotExist(err) {
		return "", errors.New("image file does not exist")
	}

	//Encrypt the image with the key
	if encRequired {
		//If the key ID is specified, make sure it's a valid UUID		
		if len(strings.TrimSpace(keyID)) > 0 {
			if validatekeyIDErr := validation.ValidateUUID(keyID); validatekeyIDErr != nil {
				return "", errors.New("incorrectly formatted key ID")
			}
		}
		wrappedKey, keyURLString, err = util.FetchKey(keyID)

		err = util.Encrypt(inputImageFilePath, consts.EnvelopePrivatekeyLocation, outputEncImageFilePath, wrappedKey)
		if err != nil {
			return "", errors.New("error encrypting image: " + err.Error())
		}
		imageFilePath = outputEncImageFilePath
	}

	//Check the encrypted image output file
	imageFile, err := ioutil.ReadFile(imageFilePath)
	if err != nil {
		return "", errors.New("error while reading the image file")
	}

	//Take the digest of the encrypted image
	digest := sha256.Sum256([]byte(imageFile))

	//Create image flavor
	imageFlavor, err := flavor.GetImageFlavor(flavorLabel, encRequired, keyURLString, base64.StdEncoding.EncodeToString(digest[:]))
	if err != nil {
		return "", errors.New("error creating image flavor:" + err.Error())
	}

	//Marshall the image flavor to a JSON string
	imageFlavorJSON, err := json.Marshal(imageFlavor)
	if err != nil {
		return "", errors.New("error while marshalling image flavor:" + err.Error())
	}

	//If no output flavor file path was specified, return the marshalled image flavor
	if len(strings.TrimSpace(outputFlavorFilePath)) <= 0 {
		return string(imageFlavorJSON), nil
	}

	//Otherwise, write it to the specified file
	err = ioutil.WriteFile(outputFlavorFilePath, imageFlavorJSON, 0600)
	if err != nil {
		return "", errors.New("error writing image flavor to output file")
	}
	return "", err
}

//Usage command line usage string
func Usage() string {
	return "usage: wpm create-image-flavor [-l label] [-i in] [-o out] [-e encout] [-k key]\n" +
		"  -l, --label     image flavor label\n" +
		"  -i, --in        input image file path\n" +
		"  -o, --out       (optional) output image flavor file path\n" +
		"                  if not specified, will print to the console\n" +
		"  -e, --encout    (optional) output encrypted image file path\n" +
		"                  if not specified, encryption is skipped\n" +
		"  -k, --key       (optional) existing key ID\n" +
		"                  if not specified, a new key is generated\n"
}
