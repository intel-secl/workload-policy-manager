package main

import (
	"fmt"
	imageFlavor "intel/isecl/wpm/pkg/imageflavor"
	setup "intel/isecl/wpm/pkg/setup"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {

	if len(os.Args[0:]) <= 1 {
		usage()
		return
	}
	var task = os.Args[1]

	switch arg := strings.ToLower(task); arg {
	case "setup":
		switch setupTask := strings.ToLower(os.Args[2]); setupTask {
		case "create-envelope-key":
			createKeys()

		case "register-envelope-key":
			registerKeys()

		case "--all":
			createKeys()
			registerKeys()

		default:
			usage()
			return
		}

	case "create-image-flavor":
		if len(os.Args[1:]) < 6 {
			usage()
			return
		}

		isEncryptionRequired, _ := strconv.ParseBool(os.Args[5])
		_, err := imageFlavor.CreateImageFlavor(os.Args[2], os.Args[3], os.Args[4], isEncryptionRequired, false, os.Args[6])
		if err != nil {
			log.Fatal("cannot create flavor")
		} else {
			log.Println("Image flavor created successfully")
		}

	case "uninstall":
		log.Println("Uninstall")

	case "help", "-help", "--help":
		usage()

	default:
		fmt.Printf("Unrecognized option : %s\n", arg)
		usage()
	}
}

func createKeys() {
	err := setup.CreateEnvelopeKey()
	if err != nil {
		log.Fatal("Error creating the envelope key")
	} else {
		log.Println("Envelope key created successfully")
	}
}

func registerKeys() {
	err := setup.RegisterEnvelopeKey()
	if err != nil {
		log.Fatal("Error creating the envelope key")
	} else {
		log.Println("Envelope key created successfully")
	}
}

func usage() {
	fmt.Println("Usage: $0 uninstall|create-image-flavor|create-software-flavor")
	fmt.Println("Usage: $0 setup [--force|--noexec] [task1 task2 ...]")
	//fmt.Println("Usage: $0 export-config [outfile|--in=infile|--out=outfile|--stdout] [--env-password=PASSWORD_VAR]")
	fmt.Println("Available setup tasks:CreateEnvelopKey and RegisterEnvelopeKeyWithKBS")
}
