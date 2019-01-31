package kmsclient

/*
 *
 * @author srege
 *
 */
import (
	"fmt"
	c "intel/isecl/wpm/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAuthToken(t *testing.T) {
	c.Configuration.Kms.APIURL = "https://10.105.168.214:443/v1/"
	c.Configuration.Kms.APIUsername = "admin"
	c.Configuration.Kms.APIPassword = "password"
	c.Configuration.Kms.TLSSHA256 = "313f4798df8605b37bf89d68bef596e0a7ce338088a48dd389553d80bb512b76"
	authToken, err := GetAuthToken()
	fmt.Println(authToken)
	assert.Nil(t, err)
	assert.NotNil(t, authToken)
}
