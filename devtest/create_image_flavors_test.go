package imageflavor

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	config "intel/isecl/wpm/config"
	"testing"
	i "intel/isecl/wpm/pkg/imageflavor"
)

func TestCreateImageFlavor(t *testing.T) {
	config.Configuration.Kms.APIURL = "https://10.105.168.214:443/v1/"
	config.Configuration.Kms.APIUsername = "kms-admin"
	config.Configuration.Kms.APIPassword = "password"
	config.Configuration.EnvelopePrivatekeyLocation = "admin-privatekey.pem"
	config.Configuration.Kms.TLSSha256 = "313f4798df8605b37bf89d68bef596e0a7ce338088a48dd389553d80bb512b76"
	imageFlavor, err := i.CreateImageFlavor("cirros-x86.qcow2", "cirros-x86.qcow2_enc", "", true, false, "")
	assert.Nil(t, err)
	assert.NotNil(t, imageFlavor)
}

func TestCreateImageFlavorToFile(t *testing.T) {
	config.Configuration.Kms.APIURL = "https://10.105.168.214:443/v1/"
	config.Configuration.Kms.APIUsername = "kms-admin"
	config.Configuration.Kms.APIPassword = "password"
	config.Configuration.EnvelopePrivatekeyLocation = "admin-privatekey.pem"
	config.Configuration.Kms.TLSSha256 = "313f4798df8605b37bf89d68bef596e0a7ce338088a48dd389553d80bb512b76"
	imageFlavor, err := i.CreateImageFlavor("cirros-x86.qcow2", "cirros-x86.qcow2_enc", "", true, false, "image_flavor.txt")
	assert.Nil(t, err)
	assert.NotNil(t, imageFlavor)
}
