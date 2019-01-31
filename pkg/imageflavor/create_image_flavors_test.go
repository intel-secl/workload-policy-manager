package imageflavor

import (
	"fmt"
	config "intel/isecl/wpm/config"
	client "intel/isecl/wpm/pkg/kmsclient"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateKey(t *testing.T) {
	config.Configuration.Kms.APIURL = "https://10.105.168.214:443/v1/"
	config.Configuration.Kms.APIUsername = "admin"
	config.Configuration.Kms.APIPassword = "password"
	config.Configuration.Kms.TLSSha256 = "313f4798df8605b37bf89d68bef596e0a7ce338088a48dd389553d80bb512b76"
	authToken, err := client.GetAuthToken()
	keyInfo := createKey(authToken)
	assert.Nil(t, err)
	assert.NotNil(t, keyInfo)
	fmt.Println(keyInfo.KeyID)
}

func TestRetrieveTransferKey(t *testing.T) {
	config.Configuration.Kms.APIURL = "https://10.105.168.214:443/v1/"
	config.Configuration.Kms.APIUsername = "admin"
	config.Configuration.Kms.APIPassword = "password"
	config.Configuration.Kms.TLSSha256 = "313f4798df8605b37bf89d68bef596e0a7ce338088a48dd389553d80bb512b76"
	authToken, err := client.GetAuthToken()
	keyID := "d10220b7-4398-48d7-8843-ccd9675f0d16"
	keyURL := config.Configuration.Kms.APIURL + "keys/" + keyID + "/transfer"
	key := retrieveKey(authToken, keyURL)
	assert.Nil(t, err)
	assert.NotNil(t, key)
}
func TestCreateImageFlavor(t *testing.T) {
	config.Configuration.Kms.APIURL = "https://10.105.168.214:443/v1/"
	config.Configuration.Kms.APIUsername = "admin"
	config.Configuration.Kms.APIPassword = "password"
	config.Configuration.EnvelopePrivatekeyLocation = "admin-privatekey.pem"
	config.Configuration.Kms.TLSSha256 = "313f4798df8605b37bf89d68bef596e0a7ce338088a48dd389553d80bb512b76"
	imageFlavor, err := CreateImageFlavor("cirros-x86.qcow2", "cirros-x86.qcow2_enc", "", true, false, "")
	assert.Nil(t, err)
	assert.NotNil(t, imageFlavor)
}

func TestCreateImageFlavorToFile(t *testing.T) {
	config.Configuration.Kms.APIURL = "https://10.105.168.214:443/v1/"
	config.Configuration.Kms.APIUsername = "admin"
	config.Configuration.Kms.APIPassword = "password"
	config.Configuration.EnvelopePrivatekeyLocation = "admin-privatekey.pem"
	config.Configuration.Kms.TLSSha256 = "313f4798df8605b37bf89d68bef596e0a7ce338088a48dd389553d80bb512b76"
	imageFlavor, err := CreateImageFlavor("cirros-x86.qcow2", "cirros-x86.qcow2_enc", "", true, false, "image_flavor.txt")
	assert.Nil(t, err)
	assert.NotNil(t, imageFlavor)
}
