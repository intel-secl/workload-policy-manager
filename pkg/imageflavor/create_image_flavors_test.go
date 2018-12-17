package imageflavor

import (
	"fmt"
	config "intel/isecl/wpm/config"
	client "intel/isecl/wpm/pkg/kmsclient"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateKey(t *testing.T) {
	config.Configuration.KmsAPIURL = "https://10.105.168.214:443/v1/"
	config.Configuration.KmsAPIUsername = "admin"
	config.Configuration.KmsAPIPassword = "password"
	authToken, err := client.GetAuthToken()
	keyInfo := createKey(authToken)
	assert.Nil(t, err)
	assert.NotNil(t, keyInfo)
	fmt.Println(keyInfo.KeyID)
}

func TestRetrieveTransferKey(t *testing.T) {
	config.Configuration.KmsAPIURL = "https://10.105.168.214:443/v1/"
	config.Configuration.KmsAPIUsername = "admin"
	config.Configuration.KmsAPIPassword = "password"
	authToken, err := client.GetAuthToken()
	keyID := "69010ca8-462d-42b9-a1a7-f0426121831d"
	keyURL := c.Configuration.BaseURL + "keys/" + keyID + "/transfer"
	fmt.Println(keyURL)
	key := retrieveKey(authToken, keyURL)
	fmt.Println(len(key))
	assert.Nil(t, err)
	assert.NotNil(t, key)
	fmt.Println(key)
}
func TestCreateImageFlavor(t *testing.T) {
	config.Configuration.KmsAPIURL = "https://10.105.168.214:443/v1/"
	config.Configuration.KmsAPIUsername = "admin"
	config.Configuration.KmsAPIPassword = "password"
	config.Configuration.EnvelopeKeyLocation = "admin-privatekey.pem"
	imageFlavor, err := CreateImageFlavor("cirros-x86.qcow2", "cirros-x86.qcow2_enc", "", true, false, "")
	assert.Nil(t, err)
	assert.NotNil(t, imageFlavor)
}

func TestCreateImageFlavorToFile(t *testing.T) {
	config.Configuration.KmsAPIURL = "https://10.105.168.214:443/v1/"
	config.Configuration.KmsAPIUsername = "admin"
	config.Configuration.KmsAPIPassword = "password"
	imageFlavor, err := CreateImageFlavor("cirros-x86.qcow2", "cirros-x86.qcow2_enc", "", true, false, "image_flavor.text")
	assert.Nil(t, err)
	assert.NotNil(t, imageFlavor)
}
