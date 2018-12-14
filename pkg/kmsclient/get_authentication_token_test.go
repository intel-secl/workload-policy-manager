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
	authToken, err := GetAuthToken()
	fmt.Println(authToken)
	assert.Nil(t, err)
	assert.NotNil(t, authToken)
}
