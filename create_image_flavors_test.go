package wpm

import (
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
	"os"
	"log"
)

func TestGetAuthToken(t *testing.T) {
	    client := NewBasicAuthClient("username", "password")
		authToken,err := client.getAuthToken("10.105.168.214","443")
		fmt.Println(authToken.AuthorizationToken)
		
		assert.Nil(t, err)
	    assert.NotNil(t, authToken)
}

func TestCreateKey(t *testing.T){
	client := NewBasicAuthClient("username", "password")
	authToken,err := client.getAuthToken("10.105.168.214","443")
	keyInfo := client.createKey(authToken.AuthorizationToken,"10.105.168.214","443")
	assert.Nil(t, err)
	assert.NotNil(t, keyInfo)
	fmt.Println(keyInfo.KeyId)
}

func TestRetrieveTransferKey(t *testing.T) {
	var authToken AuthToken
	client := NewBasicAuthClient("username", "password")
	authToken,err := client.getAuthToken("10.105.168.214","443")
	fmt.Println(authToken.AuthorizationToken)
	key := client.retrieveKey(authToken.AuthorizationToken,"https://10.105.168.214:443/v1/keys/f4634e81-613b-4677-8e4c-95caf61a3d65/transfer")
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
	key := "UyUGqlszfE5WQrf2PaCrxd0PbfqredR98P5jcCrwPhQ1b34n2SNvX1E0R7hdRdwfVVyH5T6aeEbGe4IJbwQ8zDG06g53Tl8koI8/WxmJEboaVc7brCUywdW2+/4TmeU+NfKvvlyfOWkg1XcVjt7YLaNMWxr2oV6U9enukq9qXE1uv6MdtZ6RDQ0712fhYR6QnQzOFw7Iv0YwMb/Fj12K3LrrxWpwCbtSmTKbGkCc0nBXU8CJ2xQNCnc4FyKuwqbwQWKaiN3vMNIPOmHduOVtA3HlsmSJAycal4GcAg9av7ZK6AOuoqYRmRvp+sfddMI6+wqZCyqoyhjCUD8ubhrzsg=="
	encryptedPath := "salonee.txt_enc"
	err,encryptedImage := encryptImage(imagePath,encryptedPath,[]byte(key))
	assert.Nil(t, err)
	assert.NotNil(t, encryptedImage)
	fmt.Println(encryptedImage)
}