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
	log "github.com/sirupsen/logrus"
	flavor "intel/isecl/lib/flavor"
	kms "intel/isecl/lib/kms-client"
	"intel/isecl/wpm/consts"
	"intel/isecl/wpm/pkg/kmsclient"
	"io/ioutil"
	"os"
	"strings"
)

//CreateImageFlavor is used to create flavor of an encrypted image
func CreateImageFlavor(label string, imagePath string, encryptFilePath string, keyID string, encryptionRequired bool, integrityEnforced bool, outputFile string) (string, error) {
	var err error
	var keyURL string
	var keyValue []byte
	var keyInfo kms.KeyInfo

	log.Info("Creating image flavor")
	//input validation
	if len(strings.TrimSpace(imagePath)) <= 0 {
		return "", errors.New("image path not given")
	}

	if len(strings.TrimSpace(label)) <= 0 {
		return "", errors.New("label for flavor not given")
	}

	if len(strings.TrimSpace(encryptFilePath)) <= 0 {
		return "", errors.New("encryption file path not given")
	}

	// check if image exists at the specified location
	_, err = os.Stat(imagePath)
	if os.IsNotExist(err) {
		log.Error(err)
		return "", errors.New("image file does not exist")
	}

	kc := kmsclient.InitializeClient()

	//create key if keyId is not specified in input
	if len(strings.TrimSpace(keyID)) <= 0 {
		//Iniliaze KeyInfo
		keyInfo.Algorithm = consts.KMS_ENCRYPTION_ALG
		keyInfo.KeyLength = consts.KMS_KEY_LENGTH
		keyInfo.CipherMode = consts.KMS_CIPHER_MODE

		key, err := kc.Keys().Create(keyInfo)
		if err != nil {
			log.Error(err)
			return "", errors.New("Error in creating transfer key")
		}
		keyID = key.KeyID
	}

	//retrieve key using keyid
	keyValue, err = kc.Key(keyID).Retrieve()
	if err != nil {
		log.Error(err)
		return "", errors.New("error in retrieving transfer key")
	}

	// encrypt image using key
	err = encrypt(imagePath, consts.EnvelopePrivatekeyLocation, encryptFilePath, keyValue)
	if err != nil {
		log.Error(err)
		return "", errors.New("error in encrypting image.")
	}
	encryptedImage, err := ioutil.ReadFile(encryptFilePath)
	if err != nil {
		log.Error(err)
		return "", errors.New("error reading from input file.")
	}
	//calculate SHA256 of the encrpted image
	digest := sha256.Sum256([]byte(encryptedImage))

	// create image flavor
	imageFlavor, err := flavor.GetImageFlavor(label, encryptionRequired, keyURL, base64.StdEncoding.EncodeToString(digest[:]))
	if err != nil {
		log.Error(err)
		return "", errors.New("error in creating image flavor.")
	}

	jsonFlavor, err := json.Marshal(imageFlavor)
	if len(strings.TrimSpace(outputFile)) <= 0 {
		return string(jsonFlavor), nil
	}
	//create outputFile for image flavor
	_, err = os.Create(outputFile)
	if err != nil {
		log.Error(err)
		return "", errors.New("error creating output file.")
	}

	err = ioutil.WriteFile(outputFile, []byte(jsonFlavor), 0600)
	if err != nil {
		log.Error(err)
		return "", errors.New("error writing image flavor to output file.")
	}
	return "", err
}
