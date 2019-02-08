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
	config "intel/isecl/wpm/config"
	"intel/isecl/wpm/consts"
	"intel/isecl/wpm/pkg/kmsclient"
	"io/ioutil"
	"os"
	"strings"
	"net/url"
)

//CreateImageFlavor is used to create flavor of an encrypted image
func CreateImageFlavor(label string, imagePath string, encryptFilePath string, keyID string, encryptionRequired bool, integrityEnforced bool, outputFile string) (string, error) {
	var err error
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
		return "", errors.New("image file does not exist")
	}

	kc, err := kmsclient.InitializeClient()
	if err != nil {
		return "", errors.New("error initializing KMS client")
	}

	//create key if keyId is not specified in input
	if len(strings.TrimSpace(keyID)) <= 0 {
		//Initialize KeyInfo
		keyInfo.Algorithm = consts.KMS_ENCRYPTION_ALG
		keyInfo.KeyLength = consts.KMS_KEY_LENGTH
		keyInfo.CipherMode = consts.KMS_CIPHER_MODE

		key, err := kc.Keys().Create(keyInfo)
		if err != nil {
			return "", errors.New("error in creating transfer key: " + err.Error())
		}
		keyID = key.KeyID
	}

	// set key URL
	keyURL, err := url.Parse(config.Configuration.Kms.APIURL + "keys/" + keyID + "/transfer")
	if err != nil {
		return "", errors.New("error building KMS key URL: " + err.Error())
	}

	//retrieve key using keyid
	keyValue, err = kc.Key(keyID).Retrieve()
	if err != nil {
		return "", errors.New("error in retrieving transfer key: " + err.Error())
	}

	// encrypt image using key
	err = encrypt(imagePath, consts.EnvelopePrivatekeyLocation, encryptFilePath, keyValue)
	if err != nil {
		return "", errors.New("error in encrypting image: " + err.Error())
	}
	encryptedImage, err := ioutil.ReadFile(encryptFilePath)
	if err != nil {
		return "", errors.New("error reading from input file")
	}
	//calculate SHA256 of the encrpted image
	digest := sha256.Sum256([]byte(encryptedImage))

	// create image flavor
	imageFlavor, err := flavor.GetImageFlavor(label, encryptionRequired, keyURL.String(), base64.StdEncoding.EncodeToString(digest[:]))
	if err != nil {
		return "", errors.New("error in creating image flavor:" + err.Error())
	}

	jsonFlavor, err := json.Marshal(imageFlavor)

	if len(strings.TrimSpace(outputFile)) <= 0 {
		return string(jsonFlavor), nil
	}
	
	err = ioutil.WriteFile(outputFile, jsonFlavor, 0600)
	if err != nil {
		return "", errors.New("error writing image flavor to output file")
	}
	return "", err
}
