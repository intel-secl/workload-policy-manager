package main

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	csetup "intel/isecl/lib/common/setup"
	"intel/isecl/wpm/config"
	"intel/isecl/wpm/consts"
	imageFlavor "intel/isecl/wpm/pkg/imageflavor"
	"intel/isecl/wpm/pkg/setup"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func main() {
	args := os.Args[1:]
	if len(args) <= 0 {
		usage()
		return
	}

	// Save log configurations
	config.LogConfiguration()

	switch arg := strings.ToLower(args[0]); arg {
	case "setup":
		// Check if nosetup environment variable is true, if yes then skip the setup tasks
		if nosetup, err := strconv.ParseBool(os.Getenv("WPM_NOSETUP")); err != nil && nosetup == false {
			// Run list of setup tasks one by one
			setupRunner := &csetup.Runner{
				Tasks: []csetup.Task{
					setup.CreateEnvelopeKey{},
					setup.RegisterEnvelopeKey{},
				},
				AskInput: false,
			}

			err = setupRunner.RunTasks(args[1:]...)
			if err != nil {
				fmt.Println("Error running setup: ", err)
				os.Exit(1)
			}
		} else {
			fmt.Println("WPM_NOSETUP is set, skipping setup")
			os.Exit(1)
		}

	case "create-image-flavor":
		label := flag.String("l", "", "Label for flavor not given.")
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
		if *label == "" {
			fmt.Printf("Please provide the input image file path using -l option. It is a " +
				"required parameter.\n")
			os.Exit(1)
		}
		fmt.Printf("label:%s, inputImageFilePath: %s, outputEncryptedImageFilePath: %s, outputFlavorFilePath:"+
			" %s, keyID: %s, isEncryRequired: %t\n", *label, *inputImageFilePath, *outputEncryptedImageFilePath,
			*outputFlavorFilePath, *inputKeyID, *isEncryRequired)

		var keyID string
		if isValidUUID(*inputKeyID) {
			keyID = *inputKeyID
		} else {
			keyID = ""
		}
		_, err := imageFlavor.CreateImageFlavor(*label, *inputImageFilePath, *outputEncryptedImageFilePath,
			keyID, *isEncryRequired, false, *outputFlavorFilePath)
		if err != nil {
			log.Error("cannot create flavor")
		} else {
			log.Info("Image flavor created successfully")
		}

	case "uninstall":
		log.Info("Uninstalling WPM")
		deleteFile("/usr/local/bin/wpm")
		deleteFile("/opt/wpm/")
		deleteFile(consts.ConfigDirPath)
		deleteFile(consts.LogDirPath)
		log.Info("WPM uninstalled successfully")

	case "help", "-help", "--help":
		usage()

	default:
		fmt.Printf("Unrecognized option : %s\n", arg)
		usage()
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
	fmt.Printf("Usage: $0 uninstall|create-image-flavor|create-software-flavor")
	fmt.Printf("Usage: $0 setup [--force|--noexec] [task1 task2 ...]")
	fmt.Printf("Available setup tasks: CreateEnvelopKey and RegisterEnvelopeKeyWithKBS")
}

func isValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}

func deleteFile(path string) {
	log.Info("Deleting file: ", path)
	// delete file
	var err = os.RemoveAll(path)
	if err != nil {
		log.Error(err)
	}
}
