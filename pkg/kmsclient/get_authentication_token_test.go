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
	c.Configuration.BaseURL = "https://10.105.168.214:443/v1/"
	c.Configuration.Username = "admin"
	c.Configuration.Password = "password"
	c.Configuration.KMSTlsCertSHA256 = "313f4798df8605b37bf89d68bef596e0a7ce338088a48dd389553d80bb512b76"
	authToken, err := GetAuthToken()
	fmt.Println(authToken)
	assert.Nil(t, err)
	assert.NotNil(t, authToken)
}
