package main

import (
	"fmt"
	config "intel/isecl/wpm/config"
	imageFlavor "intel/isecl/wpm/pkg/imageflavor"
	setup "intel/isecl/wpm/pkg/setup"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
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
		if len(os.Args[1:]) < 7 {
			usage()
			return
		}

		isEncryptionRequired, _ := strconv.ParseBool(os.Args[5])
		//isIntegrityRequired, _ := strconv.ParseBool(os.Args[6])
		var keyID string
		if isValidUUID(os.Args[4]) {
			keyID = os.Args[4]
		} else {
			keyID = ""
		}
		_, err := imageFlavor.CreateImageFlavor(os.Args[2], os.Args[3], keyID, isEncryptionRequired, false, os.Args[7])
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

	//remove wpm home directory
	args := []string{"-rf", wpmHomeDirectory}
	_, err := runCommand("rm", args)
	if err != nil {
		log.Fatal("Error trying to delete the WPM home directory")
	}

	//delete the wpm binary from installed location
	cmdArgs := []string{"-rf", "/usr/local/bin/wpm"}
	_, err = runCommand("rm", cmdArgs)
	if err != nil {
		log.Fatal("Error trying to delete the WPM binary")
	}
}

func runCommand(cmd string, args []string) (string, error) {
	out, err := exec.Command(cmd, args...).Output()
	return string(out), err
}

func usage() {
	fmt.Println("Usage: $0 uninstall|create-image-flavor|create-software-flavor")
	fmt.Println("Usage: $0 setup [--force|--noexec] [task1 task2 ...]")
	//fmt.Println("Usage: $0 export-config [outfile|--in=infile|--out=outfile|--stdout] [--env-password=PASSWORD_VAR]")
	fmt.Println("Available setup tasks:CreateEnvelopKey and RegisterEnvelopeKeyWithKBS")
}

func isValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}
