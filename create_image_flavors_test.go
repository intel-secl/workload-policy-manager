package wpm

import (
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGetAuthToken(t *testing.T) {
	    var authToken AuthToken
	    client := NewBasicAuthClient("username", "password")
		authToken,err := client.getAuthToken("10.1.70.56","443")
		fmt.Println(authToken.AuthorizationToken)
		
		assert.Nil(t, err)
	    assert.NotNil(t, authToken)
}

func TestCreateKey(t *testing.T){
	var authToken AuthToken
	var keyInfo KeyInfo
	client := NewBasicAuthClient("username", "password")
	authToken,err := client.getAuthToken("10.1.70.56","443")
	keyInfo = client.createKey(authToken.AuthorizationToken,"10.1.70.56","443")
	
	assert.Nil(t, err)
	assert.NotNil(t, keyInfo)
	fmt.Println(keyInfo)
}

func Test(t *testing.T) {
	var authToken AuthToken
	client := NewBasicAuthClient("username", "password")
	authToken,err := client.getAuthToken("10.1.70.56","443")
	fmt.Println(authToken.AuthorizationToken)
	
	assert.Nil(t, err)
	assert.NotNil(t, authToken)
}