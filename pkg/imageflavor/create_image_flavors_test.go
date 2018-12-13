package flavor

import (
	"fmt"
	"testing"
)

//"os"
//"log"
//"io/ioutil"

/*func TestGetAuthToken(t *testing.T) {
	client := NewBasicAuthClient("username", "password")
	authToken, err := client.getAuthToken("10.105.168.214", "443")
	fmt.Println(authToken.AuthorizationToken)

	assert.Nil(t, err)
	assert.NotNil(t, authToken)
}

func TestCreateKey(t *testing.T) {
	client := NewBasicAuthClient("username", "password")
	authToken, err := client.getAuthToken("10.105.168.214", "443")
	keyInfo := client.createKey(authToken.AuthorizationToken, "10.105.168.214", "443")
	assert.Nil(t, err)
	assert.NotNil(t, keyInfo)
	fmt.Println(keyInfo.KeyId)
}

func TestRetrieveTransferKey(t *testing.T) {
	var authToken AuthToken
	client := NewBasicAuthClient("username", "password")
	authToken, err := client.getAuthToken("10.105.168.214", "443")
	fmt.Println(authToken.AuthorizationToken)
	key := client.retrieveKey(authToken.AuthorizationToken, "https://10.105.168.214:443/v1/keys/756aa06b-14d1-4fa0-8d6d-363eb292e8c5/transfer")
	fmt.Println(len(key))
	assert.Nil(t, err)
	assert.NotNil(t, key)
	fmt.Println(key)
}

func TestEncryptImage(t *testing.T) {
	fmt.Println("Inside test for encrypted Image")
	imagePath := "salonee.txt"
	fmt.Println(imagePath)
	_,errorMsg := os.Stat(imagePath)
		if(os.IsNotExist(errorMsg)) {
			log.Fatal("image file does not exist.")
	    }
	key := "QHmVKX/25rRoHODlFJ5/RfQXxfYBL4yVysWDsx3hTEdndtML1Gf7sn/Q07+sEpaXzdR9DURjGFbSZON1Wt82Jm1oYNV3pSf6M73JBxOvQ3o0RAY3BJNFjbL3IW0i9VqvmPeWFIUdqRCuvpgy+fASoWvzRZ58w02QHLW2uIaVLn1uKBZYVtXOhUG5DLFTEv+Ju/edxn60tdvaEiKZ/jxu/euVFUXQyij55YRMpzFLYz6fW0P3piipAB5eij/ax6ShVUqLZbk3HtMyF+RTe6hLBkUnyNG0JOeN4JetbkAulpnWIMT7+b/cAqGStZD7IzsbDBQm/mEAcFZowvG0oJ4+sw=="
	err := ioutil.WriteFile("key.txt", []byte(key), 0600)
	if err != nil {
		fmt.Println("error during writing to file")
	}
	keyBytes := []byte(key)
	cmd := exec.Command("openssl","enc","-d","-aes-256-cbc")
	encryptedPath := "salonee.txt_enc"
	fmt.Println("Key length : ")
	fmt.Println(len(key))
	keyByte := ([]byte(key))
	fmt.Println(len(keyByte))
	err, encryptedImage := encryptImage(imagePath, encryptedPath, []byte(key))
	assert.Nil(t, err)
	assert.NotNil(t, encryptedImage)
	fmt.Println(encryptedImage)

}

func TestCreateImageFlavor(t *testing.T) {
	imageFlavor, err := CreateImageFlavor("cirros-x86.qcow2", "cirros-x86.qcow2_enc", "", true, false, "")
	assert.Nil(t, err)
	assert.NotNil(t, imageFlavor)
	fmt.Println(imageFlavor)
}

func TestCreateImageFlavorToFile(t *testing.T) {
	imageFlavor, err := CreateImageFlavor("cirros-x86.qcow2", "cirros-x86.qcow2_enc", "", true, false, "image_flavor.text")
	assert.Nil(t, err)
	assert.NotNil(t, imageFlavor)
	fmt.Println(imageFlavor)
}
*/
func TestGetConfigurationVariables(t *testing.T) {
	configuration := GetConfigurationVariables("configuration.json")
	fmt.Println(configuration.BaseURL)
}
