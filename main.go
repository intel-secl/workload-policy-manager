package main

import (
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

	log "github.com/sirupsen/logrus"
)


const WPM_HOME="/opt/wpm"

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
		imageName := flag.String("n", "", "container image name")
		flag.StringVar(imageName, "img-name", "", "container image name")
		tagName := flag.String("t", "latest", "container image tag")
		flag.StringVar(tagName, "tag-name", "latest", "container image tag")
		dockerFilePath := flag.String("f", "", "Dockerfile path")
		flag.StringVar(dockerFilePath, "docker-file", "", "Dockerfile path")
		buildDir := flag.String("d", "", "build directory path containing source to build the container image")
		flag.StringVar(buildDir, "build-dir", "", "build directory path containing source to build the container image")
		keyID := flag.String("k", "", "key ID of key used for encrypting the container image")
		flag.StringVar(keyID, "key-id", "", "key ID of key used for encrypting the container image")
		encryptionRequired := flag.Bool("enc", false, "specifies if image needs to be encrypted")
		flag.BoolVar(encryptionRequired, "encryption-required", false, "specifies if image needs to be encrypted")
		integrityEnforced := flag.Bool("enforce", false, "specifies if workload flavor should be enforced on image during launch")
		flag.BoolVar(integrityEnforced, "integrity-enforced", false, "specifies if workload flavor should be enforced on image during launch")
		notaryServerUrl := flag.String("s", "", "notary server url to pull signed images")
		flag.StringVar(notaryServerUrl, "notary-server", "", "notary server url to pull signed images")
		outputFlavorFilePath := flag.String("o", "", "output flavor file path")
		flag.StringVar(outputFlavorFilePath, "out-file", "", "output flavor file path")
		flag.Usage = func() { fmt.Println(containerImageFlavor.Usage()) }
		flag.CommandLine.Parse(os.Args[2:])

		containerImageFlavor, err := containerImageFlavor.CreateContainerImageFlavor(*imageName, *tagName, *dockerFilePath, *buildDir,
			*keyID, *encryptionRequired, *integrityEnforced, *notaryServerUrl, *outputFlavorFilePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
			os.Exit(1)
		}
		if len(containerImageFlavor) > 0 {
			fmt.Println(containerImageFlavor)
		}

	case "unwrap-key":
		wrappedKeyFilePath := flag.String("f", "", "wrapped key file path")
		flag.StringVar(wrappedKeyFilePath, "key-file-path", "", "wrapped key file path")
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

	case "uninstall":
		fmt.Println("Uninstalling WPM")
	        _, err = exec.Command("ls", consts.OptDirPath+"/secure-docker-daemon").Output()
                if err == nil {
                   removeSecureDockerDaemon()
                }

		errorFiles, err := deleteFiles("/usr/local/bin/wpm", consts.OptDirPath, consts.ConfigDirPath, consts.LogDirPath)
		if err != nil {
			fmt.Println(err)
			fmt.Println(errorFiles)
		}
                
	case "help", "-help", "--help":
		usage()

	default:
		fmt.Printf("Unrecognized option : %s\n", arg)
		usage()
	}
}

func removeSecureDockerDaemon(){
         _, err := exec.Command(WPM_HOME+"/secure-docker-daemon/uninstall-secure-docker-daemon.sh").Output()
         if err != nil {
                 fmt.Println(err)
         }
}

func runCommand(cmd string, args []string) (string, error) {
	out, err := exec.Command(cmd, args...).Output()
	return string(out), err
}

func usage() {
	fmt.Printf("Usage: $0 uninstall|create-image-flavor|create-container-image-flavor|create-software-flavor\n")
	fmt.Printf("Usage: $0 setup [--force|--noexec] [task1 task2 ...]\n")
	fmt.Printf("Available setup tasks: CreateEnvelopKey and RegisterEnvelopeKeyWithKBS\n")
}

func deleteFiles(filePath ...string) (errorFiles []string, err error) {
	for _, path := range filePath {
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
