package main

import (
	"flag"
	"fmt"
	csetup "intel/isecl/lib/common/setup"
	config "intel/isecl/wpm/config"
	"intel/isecl/wpm/pkg"
	imageFlavor "intel/isecl/wpm/pkg/imageflavor"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	logger "github.com/sirupsen/logrus"
)

func main() {
	args := os.Args[1:]
	fmt.Println("Inside main")
	fmt.Println(args)
	if len(args) <= 0 {
		fmt.Println("Length")
		fmt.Println(len(args))
		usage()
		return
	}
	config.SetConfigValues()
	switch arg := strings.ToLower(args[0]); arg {
	case "setup":
		// Check if nosetup environment variable is true, if yes then skip the setup tasks
		if nosetup, err := strconv.ParseBool(os.Getenv("WPM_NOSETUP")); err != nil && nosetup == false {
			fmt.Println("Inside nosetup")
			// Run list of setup tasks one by one
			setupRunner := &csetup.Runner{
				Tasks: []csetup.Task{
					/*setup.CreateEnvelopeKey{
						T: t,
					},*/
					//setup.RegisterEnvelopeKey{},
					pkg.SaloneeInfo{},
				},
				AskInput: false,
			}
			fmt.Println("Before Runtasks")
			err = setupRunner.RunTasks(args[1:]...)
			fmt.Println("After Runtasks")
			if err != nil {
				fmt.Println("Error running setup: ", err)
				os.Exit(1)
			}
		} else {
			fmt.Println("WPM_NOSETUP is set, skipping setup")
			os.Exit(1)
		}
		fmt.Println("End of case")

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
			logger.Error("cannot create flavor")
		} else {
			logger.Info("Image flavor created successfully")
		}

	case "uninstall":
		logger.Info("Uninstalling WPM")
		deleteFile("/usr/local/bin/wpm")
		deleteFile("/opt/wpm/")

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
	fmt.Println("Usage: $0 uninstall|create-image-flavor|create-software-flavor")
	fmt.Println("Usage: $0 setup [--force|--noexec] [task1 task2 ...]")
	fmt.Println("Available setup tasks: CreateEnvelopKey and RegisterEnvelopeKeyWithKBS")
}

func isValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}

func deleteFile(path string) {
	log.Println("Deleting file: ", path)
	// delete file
	var err = os.RemoveAll(path)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("WPM uninstalled successfully")
}
