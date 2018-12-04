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
	var err error
	fmt.Println("Inside test for encrypted Image")
	imagePath := "/tmp/testEncryption/go-workload-policy-manager/cirros-x86.qcow2"
	fmt.Println(imagePath)
	_, errorMsg := os.Stat(imagePath)
	if !os.IsExist(errorMsg) {
		log.Fatal("image file does not exist.")
	}
	key := "2Hc3ZDesuyk2IrBEnO6Hqpfo+dYhwxHAoGrQzuB7G6U=PDR6P+iJYrxJ8DB6HfuTzRzAJUhwWmSs2o/Ixg9Yx4dtB4jrHU5x9WnVdNZ4T255WFz17MRmyevC+Ih6pBRGemy61Izd1iU7orZUy4d4G21rtnJlsNDdRDGnKKsSja8fe3xh1OshgpaW3vodKBirTirrNM7LYrZkiOY9T8gq3IjvLcYsMd3V3ylByFTBYS0BufFDz0miWixdMH8LeJ0I2jAiDyENKOt5azkh0TaV5tEqvovY1Dblm5LNhiYHOuM21xz/5yuwPlYqhTLH7ZX9xg1u7y/T6gd3qr4d9K2Tx2BGJpap0zYBT3gqyzh+5otVkcJME1ACCY3fDfVunHQLkw==PDR6P+iJYrxJ8DB6HfuTzRzAJUhwWmSs2o/Ixg9Yx4dtB4jrHU5x9WnVdNZ4T255WFz17MRmyevC+Ih6pBRGemy61Izd1iU7orZUy4d4G21rtnJlsNDdRDGnKKsSja8fe3xh1OshgpaW3vodKBirTirrNM7LYrZkiOY9T8gq3IjvLcYsMd3V3ylByFTBYS0BufFDz0miWixdMH8LeJ0I2jAiDyENKOt5azkh0TaV5tEqvovY1Dblm5LNhiYHOuM21xz/5yuwPlYqhTLH7ZX9xg1u7y/T6gd3qr4d9K2Tx2BGJpap0zYBT3gqyzh+5otVkcJME1ACCY3fDfVunHQLkw=="
	encryptedPath := "/tmp/testEncryption/go-workload-policy-manager/cirros-x86.qcow2_enc"
	err = encryptImage(imagePath,encryptedPath,[]byte(key))
	assert.Nil(t, err)
	assert.NotNil(t, key)
	fmt.Println(key)
}