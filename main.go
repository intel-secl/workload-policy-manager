
package main

import (
	c "intel/isecl/wpm/config"
	"os"
	"log"
	"encoding/json"
	"fmt"
)

//var config Configuration
func main() {
	filename :="configuration.json"
  	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&c.Configuration)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(c.Configuration.BaseURL)
}

