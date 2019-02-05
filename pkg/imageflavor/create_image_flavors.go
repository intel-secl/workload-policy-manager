package imageflavor

/*
 *
 * @author srege
 *
 */
import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	t "intel/isecl/lib/common/tls"
	"fmt"
	flavor "intel/isecl/lib/flavor"
	config "intel/isecl/wpm/config"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"intel/isecl/lib/kms-client"
	"crypto/tls"
	"net/http"
	logger "github.com/sirupsen/logrus"
)

//CreateImageFlavor is used to create flavor of an encrypted image
func CreateImageFlavor(imagePath string, encryptFilePath string, keyID string, encryptionRequired bool, integrityEnforced bool, outputFile string) (string, error) {
	var err error
	var keyURL string
	var keyValue []byte

	logger.Info("Creating image flavor")
	//input validation
	if len(strings.TrimSpace(imagePath)) <= 0 {
		log.Fatal("image path not given")
	}
	

	if len(strings.TrimSpace(encryptFilePath)) <= 0 {
		log.Fatal("encryption file path not given")
	}
  

	// check if image exists at the specified location
	_, err = os.Stat(imagePath)
	if os.IsNotExist(err) {
		log.Fatal("image file does not exist")
	}
	
	//create key if keyId is not specified in input
	if len(strings.TrimSpace(keyID)) <= 0 {
		keyInformation,err := createKey()
		if err!=nil{
			
			log.Fatal("Error in creating transfer key")
            return "" ,err
		}
	   	keyValue,err = retrieveKey(keyInformation.KeyID)
		if err!=nil{
			log.Fatal("Error in retrieving transfer key")
            return "" ,err
		}
	} else {
		//retrieve key using keyid
		keyValue,err = retrieveKey(keyID)
		if err!=nil{
			log.Fatal("Error in retrieving transfer key")
            return "" ,err
		}
	}
   
	// encrypt image using key
	err = encrypt(imagePath, config.Configuration.EnvelopePrivatekeyLocation, encryptFilePath, keyValue)
	if err != nil {
		log.Fatal("Error in encrypting image.", err)
	}
	encryptedImage, err := ioutil.ReadFile(encryptFilePath)
	//calculate SHA256 of the encrpted image
	digest := sha256.Sum256([]byte(encryptedImage))

	// create image flavor
	imageFlavor, err := flavor.GetImageFlavor("label", encryptionRequired, keyURL, base64.StdEncoding.EncodeToString(digest[:]))
	if err != nil {
		log.Fatal("Error in creating image flavor.", err)
	}
	jsonFlavor, err := json.Marshal(imageFlavor)
	if len(strings.TrimSpace(outputFile)) <= 0 {
		return string(jsonFlavor), nil
	}
	//create outputFile for image flavor
	_, err = os.Create(outputFile)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	err = ioutil.WriteFile(outputFile, []byte(jsonFlavor), 0600)
	if err != nil {
		fmt.Println(err)
	}

	return "", err

}

func createKey() (*kms.KeyInfo,error) {
	var keyInfo kms.KeyInfo	
	kc := initializeClient()
	keyInfo.Algorithm = "AES"	
	keyInfo.KeyLength =  256
	keyInfo.CipherMode ="GCM"
	key, err := kc.Keys().Create(keyInfo)
	if err!=nil{
		logger.Error("Error creating key")
		return key , errors.New("Error creating key.")
	}
	return key,nil
	
}

func retrieveKey(keyID string) ([]byte,error) {
	kc := initializeClient()
	key, err := kc.Key(keyID).Retrieve()
	if err!=nil{
		logger.Error("Error retrieving key")
		return key , errors.New("Error retrieving key.")
	}
	return key,nil
}

func initializeClient() (*kms.Client) {
	var certificateDigest [32]byte
	cert, err := hex.DecodeString(config.Configuration.Kms.TLSSha256)
	if err != nil {
		log.Fatal(err)
	}
	copy(certificateDigest[:], cert)
	client := &http.Client{
					Transport: &http.Transport{
						TLSClientConfig: &tls.Config{
							InsecureSkipVerify: true,
							VerifyPeerCertificate: t.VerifyCertBySha256(certificateDigest),
						},
					},
				}
	kc := &kms.Client{
		BaseURL:  config.Configuration.Kms.APIURL,
		Username: config.Configuration.Kms.APIUsername,
		Password: config.Configuration.Kms.APIPassword,
		CertSha256 :&certificateDigest,
        HTTPClient: client,
	}
	return kc
}
