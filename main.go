package main

import (
	"crypto/md5"
	"errors"
	"flag"
	"fmt"
	csetup "intel/isecl/lib/common/setup"
	"intel/isecl/wpm/config"
	"intel/isecl/wpm/consts"
	containerImageFlavor "intel/isecl/wpm/pkg/containerimageflavor"
	imageFlavor "intel/isecl/wpm/pkg/imageflavor"
	"intel/isecl/wpm/pkg/setup"
	"intel/isecl/wpm/pkg/util"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

var (
	Version string = ""
	Time    string = ""
	Branch  string = ""
)

func printVersion() {
	fmt.Printf("Version %s\nBuild %s at %s\n", Version, Branch, Time)
}

func main() {
	args := os.Args[1:]
	if len(args) <= 0 {
		usage()
		return
	}

	// Save log configurations
	err := config.LogConfiguration()
	if err != nil {
		log.Error("error in configuring logs.")
	}

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
		flavorLabel := flag.String("l", "", "flavor label")
		flag.StringVar(flavorLabel, "label", "", "flavor label")
		inputImageFilePath := flag.String("i", "", "input image file path")
		flag.StringVar(inputImageFilePath, "in", "", "input image file path")
		outputFlavorFilePath := flag.String("o", "", "output flavor file path")
		flag.StringVar(outputFlavorFilePath, "out", "", "output flavor file path")
		outputEncImageFilePath := flag.String("e", "", "output encrypted image file path")
		flag.StringVar(outputEncImageFilePath, "encout", "", "output encrypted image file path")
		keyID := flag.String("k", "", "existing key ID")
		flag.StringVar(keyID, "key", "", "existing key ID")
		flag.Usage = func() { fmt.Println(imageFlavor.Usage()) }
		flag.CommandLine.Parse(os.Args[2:])

		imageFlavor, err := imageFlavor.CreateImageFlavor(*flavorLabel, *outputFlavorFilePath, *inputImageFilePath,
			*outputEncImageFilePath, *keyID, false)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
			os.Exit(1)
		}
		if len(imageFlavor) > 0 {
			fmt.Println(imageFlavor)
		}

	case "create-container-image-flavor":
		imageName := flag.String("i", "", "docker image name")
		flag.StringVar(imageName, "img-name", "", "docker image name")
		tagName := flag.String("t", "latest", "docker image tag")
		flag.StringVar(tagName, "tag", "latest", "docker image tag")
		dockerFilePath := flag.String("f", "", "Dockerfile path")
		flag.StringVar(dockerFilePath, "docker-file", "", "Dockerfile path")
		buildDir := flag.String("d", "", "build directory path containing source to build the docker image")
		flag.StringVar(buildDir, "build-dir", "", "build directory path containing source to build the docker image")
		keyID := flag.String("k", "", "key ID of key used for encrypting the image")
		flag.StringVar(keyID, "key-id", "", "key ID of key used for encrypting the image")
		encryptionRequired := flag.Bool("e", false, "specifies if image needs to be encrypted")
		flag.BoolVar(encryptionRequired, "encryption-required", false, "specifies if image needs to be encrypted")
		integrityEnforced := flag.Bool("s", false, "specifies if container image should be signed")
		flag.BoolVar(integrityEnforced, "integrity-enforced", false, "specifies if container image needs to be signed")
		notaryServerURL := flag.String("n", "", "notary server url to pull signed images")
		flag.StringVar(notaryServerURL, "notary-server", "", "notary server url to pull signed images")
		outputFlavorFilePath := flag.String("o", "", "output flavor file path")
		flag.StringVar(outputFlavorFilePath, "out-file", "", "output flavor file path")
		flag.Usage = func() { fmt.Println(containerImageFlavor.Usage()) }
		flag.CommandLine.Parse(os.Args[2:])

		containerImageFlavor, err := containerImageFlavor.CreateContainerImageFlavor(*imageName, *tagName, *dockerFilePath, *buildDir,
			*keyID, *encryptionRequired, *integrityEnforced, *notaryServerURL, *outputFlavorFilePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
			os.Exit(1)
		}
		if len(containerImageFlavor) > 0 {
			fmt.Println(containerImageFlavor)
		}

	case "unwrap-key":
		wrappedKeyFilePath := flag.String("i", "", "wrapped key file path")
		flag.StringVar(wrappedKeyFilePath, "in", "", "wrapped key file path")
		flag.CommandLine.Parse(os.Args[2:])

		wrappedKey, err := ioutil.ReadFile(*wrappedKeyFilePath)
		if err != nil {
			fmt.Println("Cannot read from file: " + err.Error())
			os.Exit(1)
		}
		unwrappedKey, err := util.UnwrapKey(wrappedKey, consts.EnvelopePrivatekeyLocation)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Println(unwrappedKey)

	case "get-container-image-id":
		if len(args[1:]) < 1 {
			fmt.Println("Invalid number of parameters")
			os.Exit(1)
		}
		NameSpaceDNS := uuid.Must(uuid.Parse("6ba7b810-9dad-11d1-80b4-00c04fd430c8"))
		imageUUID := uuid.NewHash(md5.New(), NameSpaceDNS, []byte(args[1]), 4)
		fmt.Println(imageUUID)

	case "uninstall":
		fmt.Println("Uninstalling WPM")
		if len(args) > 1 && strings.ToLower(args[1]) == "--purge" {
			deleteFiles(consts.ConfigDirPath)
		}
		errorFiles, err := deleteFiles(consts.WpmSymLink, consts.OptDirPath, consts.ConfigDirPath, consts.LogDirPath)
		if err != nil {
			fmt.Printf("Error deleting files : %s", errorFiles)
		}

	case "help", "-help", "--help":
		usage()

	case "--version", "-v", "version", "-version":
		printVersion()

	case "create-software-flavor":
		fmt.Println("Not supported")

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
	fmt.Printf("Workload Policy Manager\n")
	fmt.Printf("usage : %s <command> [<args>]\n\n", os.Args[0])
	fmt.Printf("Following are the list of commands\n")
	fmt.Printf("\tcreate-image-flavor|create-container-image-flavor|get-container-image-id|create-software-flavor|uninstall|--help|--version\n\n")
	fmt.Printf("\tusage : %s setup [<tasklist>]\n", os.Args[0])
	fmt.Printf("\t\t<tasklist>-space separated list of tasks\n")
	fmt.Printf("\t\t\t-Supported tasks - CreateEnvelopeKey and RegisterEnvelopeKey\n")
	fmt.Printf("\tExample :-\n")
	fmt.Printf("\t\t%s setup\n", os.Args[0])
	fmt.Printf("\t\t%s setup CreateEnvelopeKey\n", os.Args[0])
}

func deleteFiles(filePath ...string) (errorFiles []string, err error) {
	for _, path := range filePath {
		log.Info("\n Deleting : ", path)
		err := os.RemoveAll(path)
		if err != nil {
			errorFiles = append(errorFiles, path)
		}
	}
	if len(errorFiles) > 0 {
		return errorFiles, errors.New("error deleting files")
	}
	return nil, nil
}
