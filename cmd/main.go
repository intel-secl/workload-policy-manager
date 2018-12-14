package main

import (
	"encoding/json"
	"fmt"
	c "intel/isecl/wpm/config"
	f "intel/isecl/wpm/pkg/imageflavor"
	"log"
	"os"
)

func main() {

	//read values from configuration file
	filename := "configuration.json"
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&c.Configuration)
	if err != nil {
		log.Fatal(err)
	}
	outPutFilePath := "image_flavor.txt"
	_, err = f.CreateImageFlavor("", "cirros-x86.qcow2_enc", "", true, false, outPutFilePath)
	if err != nil {
		log.Fatal("cannot create flavor")
	} else {
		fmt.Println("Image flavor created successfully at " + outPutFilePath)
	}
}
