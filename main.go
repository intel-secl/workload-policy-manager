package main

import (
	"flag"
	"fmt"
	config "intel/isecl/wpm/config"
	imageFlavor "intel/isecl/wpm/pkg/imageflavor"
	setup "intel/isecl/wpm/pkg/setup"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func main() {
	config.SetConfigValues()
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
		inputImageFilePath := flag.String("i", "", "Input image file path.")
		outputEncryptedImageFilePath := flag.String("p", "", "Output encrypted image file path.")
		outputFlavorFilePath := flag.String("o", "", "Output flavor file path. If not specified, the "+
			"command will output on console by default.")
		inputKeyID := flag.String("id", "", "Specify Key ID to get the image encryption key. If not "+
			"specified, it will create a key by default.")
		isEncryRequired := flag.Bool("enc", false, "Boolean parameter to specify if image has to "+
			"be encrypted on the host when it is downloaded from the cloud orchestrator.")
		flag.CommandLine.Parse(os.Args[2:])

		if *inputImageFilePath == "" {
			fmt.Printf("Please provide the input image file path using -i option. It is a " +
				"required parameter.\n")
			os.Exit(1)
		}
		if *outputEncryptedImageFilePath == "" {
			fmt.Printf("Please provide the output encrypted image file path using -p option. " +
				"It is a required parameter.\n")
			os.Exit(1)
		}
		fmt.Printf("inputImageFilePath: %s, outputEncryptedImageFilePath: %s, outputFlavorFilePath:"+
			" %s, keyID: %s, isEncryRequired: %t\n", *inputImageFilePath, *outputEncryptedImageFilePath,
			*outputFlavorFilePath, *inputKeyID, *isEncryRequired)

		var keyID string
		if isValidUUID(*inputKeyID) {
			keyID = *inputKeyID
		} else {
			keyID = ""
		}
		_, err := imageFlavor.CreateImageFlavor(*inputImageFilePath, *outputEncryptedImageFilePath,
			keyID, *isEncryRequired, false, *outputFlavorFilePath)
		if err != nil {
			log.Fatal("cannot create flavor")
		} else {
			log.Println("Image flavor created successfully")
		}

	case "uninstall":
		uninstall()

	case "help", "-help", "--help":
		usage()

	default:
		fmt.Printf("Unrecognized option : %s\n", arg)
		usage()
	}
}

func createKeys() {
	if setup.ValidateCreateKey() {
		err := setup.CreateEnvelopeKey()
		if err != nil {
			log.Fatal("Error creating the envelope key")
		} else {
			log.Println("Envelope key created successfully")
		}
	} else {
		log.Println("Envelope keys are already created by WPM. Skipping this setup task....")
		return
	}
}

func registerKeys() {
	userID, token, isValidated := setup.ValidateRegisterKey()
	if isValidated {
		err := setup.RegisterEnvelopeKey(userID, token)
		if err != nil {
			log.Fatal("Error registering the envelope key")
		} else {
			log.Println("Envelope key registered successfully")
		}
	} else {
		log.Println("Envelope public key is already registered on KBS. Skipping this setup task....")
		return
	}

}

func uninstall() {
	var wpmHomeDirectory = "/opt/wpm/"
	var wpmBinFile = "/usr/local/bin/wpm"

	//remove wpm home directory
	args := []string{"-rf", wpmHomeDirectory}
	_, err := runCommand("rm", args)
	if err != nil {
		log.Fatal("Error trying to delete the WPM home directory")
	}
	log.Println("Deleting file: ", wpmHomeDirectory)

	//delete the wpm binary from installed location
	cmdArgs := []string{"-rf", wpmBinFile}
	_, err = runCommand("rm", cmdArgs)
	if err != nil {
		log.Fatal("Error trying to delete the WPM binary")
	}
	log.Println("Deleting file: ", wpmBinFile)
	log.Println("WPM uninstalled.")
}

func runCommand(cmd string, args []string) (string, error) {
	out, err := exec.Command(cmd, args...).Output()
	return string(out), err
}

func usage() {
	fmt.Println("Usage: $0 uninstall|create-image-flavor|create-software-flavor")
	fmt.Println("Usage: $0 setup [--force|--noexec] [task1 task2 ...]")
	fmt.Println("Available setup tasks: CreateEnvelopKey and RegisterEnvelopeKeyWithKBS")
}

func isValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}
