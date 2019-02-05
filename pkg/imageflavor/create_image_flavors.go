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
	"fmt"
	log "github.com/sirupsen/logrus"
	flavor "intel/isecl/lib/flavor"
	"intel/isecl/wpm/consts"
	"intel/isecl/wpm/pkg/kmsclient"
	kms "intel/isecl/lib/kms-client"
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
		log.Error("image path not given")
	}
	if len(strings.TrimSpace(label)) <= 0 {
		log.Error("label for flavor not given")
	}
	if len(strings.TrimSpace(encryptFilePath)) <= 0 {
		log.Error("encryption file path not given")
	}

	// check if image exists at the specified location
	_, err = os.Stat(imagePath)
	if os.IsNotExist(err) {
		log.Error("image file does not exist")
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
			log.Error("Error in creating transfer key")
			return "", err
		}
		keyID = key.KeyID
		fmt.Println("Inside if for create key keyID")
		fmt.Println(keyID)
	}

	//retrieve key using keyid
	keyValue, err = kc.Key(keyID).Retrieve()
	if err != nil {
		log.Error("Error in retrieving transfer key")
		return "", err
	}

	// encrypt image using key
	err = encrypt(imagePath, consts.EnvelopePrivatekeyLocation, encryptFilePath, keyValue)
	if err != nil {
		log.Error("Error in encrypting image.", err)
	}
	encryptedImage, err := ioutil.ReadFile(encryptFilePath)
	//calculate SHA256 of the encrpted image
	digest := sha256.Sum256([]byte(encryptedImage))

	// create image flavor
	imageFlavor, err := flavor.GetImageFlavor(label, encryptionRequired, keyURL, base64.StdEncoding.EncodeToString(digest[:]))
	if err != nil {
		log.Error("Error in creating image flavor.", err)
	}
	jsonFlavor, err := json.Marshal(imageFlavor)
	if len(strings.TrimSpace(outputFile)) <= 0 {
		return string(jsonFlavor), nil
	}
	//create outputFile for image flavor
	_, _ = os.Create(outputFile)

	_ = ioutil.WriteFile(outputFile, []byte(jsonFlavor), 0600)
	if err != nil {
		log.Error("Error writing image flavor to output file")
	}
	return "", err

}

/*
func createKey() (*kms.KeyInfo, error) {
	var keyInfo kms.KeyInfo
	keyInfo.Algorithm = "AES"
	keyInfo.KeyLength = 256
	keyInfo.CipherMode = "GCM"
	key, err := kc.Keys().Create(keyInfo)
	if err != nil {
		logger.Error("Error creating key")
		return key, errors.New("Error creating key.")
	}
	return key, nil

}

func retrieveKey(keyID string) ([]byte, error) {
	kc := initializeClient()
	key, err := kc.Key(keyID).Retrieve()
	if err != nil {
		logger.Error("Error retrieving key")
		return key, errors.New("Error retrieving key.")
	}
	return key, nil
}*/
