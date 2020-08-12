/*
 * Copyright (C) 2019 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package main

import (
	"crypto/md5"
	"crypto/x509/pkix"
	base64 "encoding/base64"
	"flag"
	"fmt"
	csetup "intel/isecl/lib/common/v2/setup"
	"intel/isecl/lib/common/v2/validation"
	"intel/isecl/wpm/v2/config"
	"intel/isecl/wpm/v2/consts"

	"github.com/pkg/errors"

	containerImageFlavor "intel/isecl/wpm/v2/pkg/containerimageflavor"
	imageFlavor "intel/isecl/wpm/v2/pkg/imageflavor"
	"intel/isecl/wpm/v2/pkg/setup"
	"intel/isecl/wpm/v2/pkg/util"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"

	commLog "intel/isecl/lib/common/v2/log"
	commMsg "intel/isecl/lib/common/v2/log/message"

	"github.com/google/uuid"
)

var (
	// Version holds the version number for the WPM binary
	Version string = ""
	// Time holds the build timestamp for the WPM binary
	Time string = ""
	// Branch holds the git build branch for the WPM binary
	Branch string = ""
	// GitHash holds the commit hash for the WPM binary
	GitHash = ""
	// GitCommitDate holds the git commit date for the WPM binary
	GitCommitDate string = ""
	log                  = commLog.GetDefaultLogger()
	secLog               = commLog.GetSecurityLogger()
)

func printVersion() {
	fmt.Printf("Workload Policy Manager Version %s\nBuild %s at %s - %s\n", Version, Branch, Time, GitHash)
}

func main() {
	log.Trace("main:main() Entering")
	defer log.Trace("main:main() Leaving")

	var context csetup.Context
	var err error

	args := os.Args[1:]
	if len(args) <= 0 {
		usage()
		return
	}

	switch arg := strings.ToLower(args[0]); arg {
	case "setup":
		flags := args

		// if setup is provided without task name then print usage and quit
		if len(args) == 1 {
			usage()
			os.Exit(0)
		}

		// Set log configurations
		err = config.LogConfiguration(false, true)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error configuring logging.")
		}

		// Check if nosetup environment variable is true, if yes then skip the setup tasks
		if nosetup, err := strconv.ParseBool(os.Getenv(consts.WpmNosetupEnv)); err == nil && nosetup == true {
			fmt.Printf("%s is set to true. Skipping setup tasks.\n", consts.WpmNosetupEnv)
			log.Infof("main:main() %s is set to true. Skipping setup.\n", consts.WpmNosetupEnv)
			os.Exit(1)
		}

		if len(args) >= 2 &&
			args[1] != "createenvelopekey" &&
			args[1] != "download_ca_cert" &&
			args[1] != "download_cert" &&
			args[1] != "all" {
			fmt.Printf("Unrecognized command: %s %s\n", args[0], args[1])
			os.Exit(1)
		}

		if len(args) > 2 {
			// check if flavor-signing cert type was specified for download
			if strings.ToLower(args[1]) == "download_cert" {
				if strings.ToLower(args[2]) != "flavor-signing" {
					fmt.Println("Invalid cert type provided for download_cert setup task: Only flavor-signing cert type is supported. Aborting.")
					os.Exit(1)
				} else if len(args) > 3 {
					// flags will be post the flavor-signing arg
					flags = args[3:]
				}
			} else {
				// flags for arguments
				flags = args[2:]
			}
		}

		installRunner := &csetup.Runner{
			Tasks: []csetup.Task{
				setup.Configurer{},
			},
			AskInput: false,
		}
		err = installRunner.RunTasks("Configurer")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error validating WPM configuration: %s\n", err.Error())
			log.WithError(err).Errorf("%s Error validating configuration: %s", commMsg.AppRuntimeErr, err.Error())
			log.Tracef("%+v", err)
			os.Exit(1)
		}

		// save configuration from config.yml
		err = config.SaveConfiguration(context)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error updating configuration: %s\n.", err.Error())
			log.WithError(err).Errorf("main:main() %s - Error updating configuration: %s\n", commMsg.AppRuntimeErr, err.Error())
			log.Tracef("%+v", err)
			os.Exit(1)
		}

		// Run list of setup tasks one by one
		setupRunner := &csetup.Runner{
			Tasks: []csetup.Task{
				csetup.Download_Ca_Cert{
					Flags:                flags,
					CmsBaseURL:           config.Configuration.Cms.BaseURL,
					CaCertDirPath:        consts.TrustedCaCertsDir,
					TrustedTlsCertDigest: config.Configuration.CmsTLSCertDigest,
					ConsoleWriter:        os.Stdout,
				},
				csetup.Download_Cert{
					Flags:              flags,
					KeyFile:            config.Configuration.FlavorSigningKeyFile,
					CertFile:           config.Configuration.FlavorSigningCertFile,
					KeyAlgorithm:       consts.DefaultKeyAlgorithm,
					KeyAlgorithmLength: consts.DefaultKeyAlgorithmLength,
					CmsBaseURL:         config.Configuration.Cms.BaseURL,
					Subject: pkix.Name{
						CommonName: config.Configuration.Subject.CommonName,
					},
					SanList:       consts.DefaultWpmSan,
					CertType:      "Signing",
					CaCertsDir:    consts.TrustedCaCertsDir,
					BearerToken:   "",
					ConsoleWriter: os.Stdout,
				},
				setup.CreateEnvelopeKey{
					Flags: flags,
				},
			},
			AskInput: false,
		}

		// if "setup all" is passed we need to run all the tasks in order
		tasklist := []string{}
		if args[1] != "all" {
			tasklist = args[1:]
		}

		err = setupRunner.RunTasks(tasklist...)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error running setup tasks %s: %s", strings.Join(args[1:], " "), err.Error())
			log.WithError(err).Errorf("main:main() %s : Error running setup tasks %s: %s\n", commMsg.AppRuntimeErr, strings.Join(args[1:], " "), err.Error())
			log.Tracef("%+v", err)
			os.Exit(1)
		}

	case "create-image-flavor":
		// Set log configurations
		err := config.LogConfiguration(config.Configuration.LogEnableStdout, true)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error configuring logging.")
			log.Tracef("%+v", err)
		}

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
		flag.Usage = func() { imageFlavorUsage() }
		flag.CommandLine.Parse(os.Args[2:])

		if len(strings.TrimSpace(*flavorLabel)) <= 0 || len(strings.TrimSpace(*inputImageFilePath)) <= 0 {
			fmt.Fprintln(os.Stderr, "Error creating VM image flavor: Missing arguments Flavor label and image file path.")
			log.Errorf("main:main() %s : Error creating VM image flavor: Missing arguments Flavor label and image file path\n", commMsg.InvalidInputBadParam)
			log.Tracef("%+v", err)
			imageFlavorUsage()
			os.Exit(1)
		}

		// validate input strings
		inputArr := []string{*flavorLabel, *outputFlavorFilePath, *inputImageFilePath, *outputEncImageFilePath}
		if validationErr := validation.ValidateStrings(inputArr); validationErr != nil {
			fmt.Fprintln(os.Stderr, "Error creating VM image flavor: Invalid input args format")
			log.WithError(validationErr).Errorf("main:main() %s : Error creating VM image flavor. Parse error for input args: [ %s ] - %s\n", commMsg.InvalidInputBadParam, inputArr, validationErr.Error())
			log.Tracef("%+v", validationErr)
			imageFlavorUsage()
			os.Exit(1)
		}

		//If the key ID is specified, make sure it's a valid UUID
		if len(strings.TrimSpace(*keyID)) > 0 {
			if validatekeyIDErr := validation.ValidateUUIDv4(*keyID); validatekeyIDErr != nil {
				fmt.Fprintln(os.Stderr, "Error creating VM image flavor: Invalid Key UUID format")
				log.WithError(validatekeyIDErr).Errorf("main:main() %s : Error creating VM image flavor: Invalid UUID - %s\n", commMsg.InvalidInputBadParam, *keyID)
				log.Tracef("%+v", validatekeyIDErr)
				imageFlavorUsage()
				os.Exit(1)
			}
		}

		imageFlavor, err := imageFlavor.CreateImageFlavor(*flavorLabel, *outputFlavorFilePath, *inputImageFilePath,
			*outputEncImageFilePath, *keyID, false)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating VM image flavor: %s\n", err.Error())
			log.WithError(err).Errorf("main:main() %s - Error creating VM image flavor: %s\n", commMsg.AppRuntimeErr, err.Error())
			log.Tracef("%+v", err)
			os.Exit(1)
		}
		if len(imageFlavor) > 0 {
			fmt.Println(imageFlavor)
		}

	case "fetch-key":
		// Set log configurations
		err := config.LogConfiguration(config.Configuration.LogEnableStdout, true)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error configuring logging.")
			log.Tracef("%+v", err)
		}

		keyID := flag.String("k", "", "existing key ID")
		flag.StringVar(keyID, "key", "", "existing key ID")
		assetTag := flag.String("t", "", "asset tags associated with the new key")
		flag.StringVar(assetTag, "asset-tag", "", "asset tags associated with the new key")
		flag.Usage = func() { fetchKeyUsage() }
		flag.CommandLine.Parse(os.Args[2:])

		//If the key ID is specified, make sure it's a valid UUID
		if len(strings.TrimSpace(*keyID)) > 0 {
			if validatekeyIDErr := validation.ValidateUUIDv4(*keyID); validatekeyIDErr != nil {
				fmt.Fprintln(os.Stderr, "Error creating VM image flavor: Invalid Key UUID format")
				log.WithError(validatekeyIDErr).Errorf("main:main() %s : Error creating VM image flavor: Invalid UUID - %s\n", commMsg.InvalidInputBadParam, *keyID)
				log.Tracef("%+v", validatekeyIDErr)
				imageFlavorUsage()
				os.Exit(1)
			}
		}

		keyInfo, err := util.FetchKeyForAssetTag(*keyID, *assetTag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching key: %s\n", err.Error())
			log.WithError(err).Errorf("main:main() %s - Error fetching: %s\n", commMsg.AppRuntimeErr, err.Error())
			log.Tracef("%+v", err)
			os.Exit(1)
		}
		if len(keyInfo) > 0 {
			fmt.Println(string(keyInfo))
		}

	case "create-container-image-flavor":
		// Set log configurations
		err := config.LogConfiguration(config.Configuration.LogEnableStdout, true)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error configuring logging.")
		}

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
		flag.Usage = func() { containerFlavorUsage() }
		flag.CommandLine.Parse(os.Args[2:])

		if len(strings.TrimSpace(*imageName)) <= 0 {
			fmt.Println("Flavor label and image file path are required arguments.")
			containerFlavorUsage()
			os.Exit(1)
		}

		// validate input strings
		inputArr := []string{*imageName, *tagName, *dockerFilePath, *buildDir, *outputFlavorFilePath}
		if validationErr := validation.ValidateStrings(inputArr); validationErr != nil {
			fmt.Fprintln(os.Stderr, "Error Creating Container Flavor: Input strings contain invalid characters")
			log.WithError(validationErr).Errorf("main:main() %s : Error Creating Container Flavor: Validation error for input args: %s\n", commMsg.InvalidInputBadParam, inputArr)
			log.Tracef("%+v", validationErr)
			containerFlavorUsage()
			os.Exit(1)
		}

		//If the key ID is specified, make sure it's a valid UUID
		if len(strings.TrimSpace(*keyID)) > 0 {
			if validatekeyIDErr := validation.ValidateUUIDv4(*keyID); validatekeyIDErr != nil {
				fmt.Fprintln(os.Stderr, "Error Creating Container Flavor: Invalid UUID")
				log.WithError(validatekeyIDErr).Errorf("main:main() %s : Error Creating Container Flavor: %s\n", commMsg.InvalidInputBadParam, validatekeyIDErr.Error())
				log.Tracef("%+v", validatekeyIDErr)
				containerFlavorUsage()
				os.Exit(1)
			}
		}

		if *notaryServerURL != "" {
			notaryServerURIValue, _ := url.Parse(*notaryServerURL)
			protocol := make(map[string]byte)
			protocol["https"] = 0
			if validateURLErr := validation.ValidateURL(*notaryServerURL, protocol, notaryServerURIValue.RequestURI()); validateURLErr != nil {
				fmt.Fprintln(os.Stderr, "Error Creating Container Flavor: Invalid key URL format")
				log.WithError(validateURLErr).Errorf("main:main() %s : Error Creating Container Flavor: Invalid key URL format %s\n", commMsg.InvalidInputBadParam, validateURLErr.Error())
				log.Tracef("%+v", validateURLErr)
				containerFlavorUsage()
				os.Exit(1)
			}
		}

		containerImageFlavor, err := containerImageFlavor.CreateContainerImageFlavor(*imageName, *tagName, *dockerFilePath, *buildDir,
			*keyID, *encryptionRequired, *integrityEnforced, *notaryServerURL, *outputFlavorFilePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error Creating Container Flavor! %s\n", err.Error())
			log.WithError(err).Errorf("main:main() %s : Error Creating Container Flavor: %s\n", commMsg.AppRuntimeErr, err.Error())
			log.Tracef("%+v", err)
			os.Exit(1)
		}

		if len(containerImageFlavor) > 0 {
			log.Info("main:main() Successfully created container image flavor")
			fmt.Println(containerImageFlavor)
		}

	case "unwrap-key":
		err := config.LogConfiguration(config.Configuration.LogEnableStdout, true)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error configuring logging.")
		}

		wrappedKeyFilePath := flag.String("i", "", "wrapped key file path")
		flag.StringVar(wrappedKeyFilePath, "in", "", "wrapped key file path")
		flag.CommandLine.Parse(os.Args[2:])

		// validate input strings
		inputArr := []string{*wrappedKeyFilePath}
		if validationErr := validation.ValidateStrings(inputArr); validationErr != nil {
			fmt.Fprintf(os.Stderr, "Error unwrapping key: Invalid key file path %s\n", *wrappedKeyFilePath)
			log.WithError(err).Errorf("main:main() %s : Error unwrapping key: %s\n", commMsg.AppRuntimeErr, validationErr.Error())
			log.Tracef("%+v", err)
			os.Exit(1)
		}

		wrappedKey, err := ioutil.ReadFile(*wrappedKeyFilePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error unwrapping key: Unable to read from wrapped key file: %s\n", err.Error())
			log.WithError(err).Errorf("main:main() %s : Error unwrapping key: Unable to read from wrapped key file %s: %s\n", commMsg.AppRuntimeErr, *wrappedKeyFilePath, err.Error())
			log.Tracef("%+v", err)
			os.Exit(1)
		}

		unwrappedKey, err := util.UnwrapKey(wrappedKey, consts.EnvelopePrivatekeyLocation)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error unwrapping key %s - %s\n", *wrappedKeyFilePath, err.Error())
			log.WithError(err).Errorf("main:main() %s : Error unwrapping key: %s\n", commMsg.AppRuntimeErr, err.Error())
			log.Tracef("%+v", err)
			os.Exit(1)
		}
		log.Info("main:main() Successfully unwrapped key")
		fmt.Println(base64.StdEncoding.EncodeToString(unwrappedKey))

	case "get-container-image-id":
		err := config.LogConfiguration(config.Configuration.LogEnableStdout, true)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error configuring logging.")
		}

		if len(args[1:]) < 1 {
			fmt.Fprintf(os.Stderr, "Invalid number of parameters")
			os.Exit(1)
		}
		NameSpaceDNS := uuid.Must(uuid.Parse(consts.SampleUUID))
		imageUUID := uuid.NewHash(md5.New(), NameSpaceDNS, []byte(args[1]), 4)
		log.Infof("main:main() Successfully retrieved container image ID: %s\n", imageUUID)
		fmt.Println(imageUUID)

	case "uninstall":
		config.LogConfiguration(false, false)
		fmt.Println("Uninstalling WPM")

		_, err = exec.Command("ls", consts.OptDirPath+"secure-docker-daemon").Output()
		if err == nil {
			removeSecureDockerDaemon()
		}

		if len(args) > 1 && strings.ToLower(args[1]) == "--purge" {
			deleteFiles(consts.ConfigDirPath)
		}
		errorFiles, err := deleteFiles(consts.WpmSymLink, consts.OptDirPath, consts.LogDirPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error deleting files : %s\n", errorFiles)
			log.Tracef("%+v", err)
		}

	case "-h", "--help":
		usage()

	case "--version", "-v":
		printVersion()

	default:
		fmt.Printf("Unrecognized command: %s\n", arg)
	}
}

func usage() {
	log.Trace("main:usage() Entering")
	defer log.Trace("main:usage() Leaving")

	fmt.Fprintln(os.Stdout, "Usage:")
	fmt.Fprintf(os.Stdout, "    wpm <command> [arguments]\n 	")
	fmt.Fprintln(os.Stdout, "")
	fmt.Fprintln(os.Stdout, "Available Commands:")
	fmt.Fprintln(os.Stdout, "    -h|--help                        Show this help message")
	fmt.Fprintln(os.Stdout, "    -v|--version                     Print version/build information")
	fmt.Fprintln(os.Stdout, "    create-image-flavor              Create VM image flavors and encrypt the image")
	fmt.Fprintln(os.Stdout, "    create-container-image-flavor    Create container image flavors and encrypt the container image")
	fmt.Fprintln(os.Stdout, "    get-container-image-id           Fetch the container image ID given the sha256 digest of the image")
	fmt.Fprintln(os.Stdout, "    unwrap-key                       Unwraps the image encryption key fetched from KMS")
	fmt.Fprintln(os.Stdout, "    uninstall [--purge]              Uninstall wpm. --purge option needs to be applied to remove configuration and data files")
	fmt.Fprintf(os.Stdout, "    setup                            Run workload-policy-manager setup tasks\n")
	fmt.Fprintln(os.Stdout, "")
	imageFlavorUsage()
	fmt.Fprintln(os.Stdout, "    fetch-key                        Fetch key from KMS")
	fetchKeyUsage()
	fmt.Fprintln(os.Stdout, "")
	containerFlavorUsage()
	fmt.Fprintln(os.Stdout, "")
	fmt.Fprintf(os.Stdout, "usage: get-container-image-id [<sha256 digest of image>]\n")
	fmt.Fprintln(os.Stdout, "")
	fmt.Fprintf(os.Stdout, "usage: unwrap-key [-i |--in] <wrapped key file path>\n")
	fmt.Fprintln(os.Stdout, "")
	fmt.Fprintf(os.Stdout, "Setup command usage:     wpm setup [task] [--force]\n")
	fmt.Fprintln(os.Stdout, "")
	fmt.Fprintln(os.Stdout, "Available tasks for setup:")
	fmt.Fprintln(os.Stdout, "   all                                         Runs all setup tasks")
	fmt.Fprintln(os.Stdout, "                                               Required env variables:")
	fmt.Fprintln(os.Stdout, "                                                   - get required env variables from all the setup tasks")
	fmt.Fprintln(os.Stdout, "                                               Optional env variables:")
	fmt.Fprintf(os.Stdout, "                                                   - get optional env variables from all the setup tasks\n")
	fmt.Fprintln(os.Stdout, "")
	fmt.Fprintln(os.Stdout, "   download_ca_cert                            Download CMS root CA certificate")
	fmt.Fprintln(os.Stdout, "                                               - Option [--force] overwrites any existing files, and always downloads new root CA cert")
	fmt.Fprintln(os.Stdout, "                                               Required env variables if WPM_NOSETUP=true or variable not set in config.yml:")
	fmt.Fprintln(os.Stdout, "                                                   - KMS_API_URL=<url>                               : KMS API URL")
	fmt.Fprintln(os.Stdout, "                                                   - AAS_API_URL=<url>                               : AAS API URL")
	fmt.Fprintln(os.Stdout, "                                                   - WPM_SERVICE_USERNAME=<service username>         : WPM service username")
	fmt.Fprintln(os.Stdout, "                                                   - WPM_SERVICE_PASSWORD=<service password>         : WPM service password")
	fmt.Fprintln(os.Stdout, "                                               Required env variables specific to setup task are:")
	fmt.Fprintln(os.Stdout, "                                                   - CMS_BASE_URL=<url>                              : for CMS API url")
	fmt.Fprintf(os.Stdout, "                                                   - CMS_TLS_CERT_SHA384=<CMS TLS cert sha384 hash>  : to ensure that WPM is talking to the right CMS instance\n")
	fmt.Fprintln(os.Stdout, "")
	fmt.Fprintln(os.Stdout, "   download_cert flavor-signing                Generates Key pair and CSR, gets it signed from CMS")
	fmt.Fprintln(os.Stdout, "                                               - Option [--force] overwrites any existing files, and always downloads newly signed WPM Flavor Signing cert")
	fmt.Fprintln(os.Stdout, "                                               Required env variables if WPM_NOSETUP=true or variable not set in config.yml:")
	fmt.Fprintln(os.Stdout, "                                                   - CMS_TLS_CERT_SHA384=<CMS TLS cert sha384 hash>  : to ensure that WPM is talking to the right CMS instance")
	fmt.Fprintln(os.Stdout, "                                                   - KMS_API_URL=<url>                               : KMS API URL")
	fmt.Fprintln(os.Stdout, "                                                   - AAS_API_URL=<url>                               : AAS API URL")
	fmt.Fprintln(os.Stdout, "                                                   - WPM_SERVICE_USERNAME=<service username>         : WPM service username")
	fmt.Fprintln(os.Stdout, "                                                   - WPM_SERVICE_PASSWORD=<service password>         : WPM service password")
	fmt.Fprintln(os.Stdout, "                                               Required env variables specific to setup task are:")
	fmt.Fprintln(os.Stdout, "                                                   - CMS_BASE_URL=<url>                       : for CMS API url")
	fmt.Fprintln(os.Stdout, "                                                   - BEARER_TOKEN=<token>                     : for authenticating with CMS")
	fmt.Fprintln(os.Stdout, "                                               Optional env variables specific to setup task are:")
	fmt.Fprintln(os.Stdout, "                                                   - KEY_PATH=<key_path>                        : Path of file where Flavor-Signing key needs to be stored")
	fmt.Fprintln(os.Stdout, "                                                   - CERT_PATH=<cert_path>                      : Path of file/directory where Flavor-Signing certificate needs to be stored")
	fmt.Fprintf(os.Stdout, "                                                   - WPM_FLAVOR_SIGN_CERT_CN=<COMMON NAME>      : to override default specified in config\n")
	fmt.Fprintln(os.Stdout, "")
	fmt.Fprintln(os.Stdout, "   createenvelopekey                           Creates the key pair required to securely transfer key from KMS")
	fmt.Fprintln(os.Stdout, "                                               - Option [--force] overwrites existing envelope key pairs")
}

func deleteFiles(filePath ...string) (errorFiles []string, err error) {
	log.Trace("main:deleteFiles() Entering")
	defer log.Trace("main:deleteFiles() Leaving")

	for _, path := range filePath {
		log.Info("\n Deleting : ", path)
		fmt.Println("Deleting : ", path)
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

//Usage command line usage string
func imageFlavorUsage() {
	log.Trace("main:imageFlavorUsage() Entering")
	defer log.Trace("main:imageFlavorUsage() Leaving")

	fmt.Println("usage: wpm create-image-flavor [-l label] [-i in] [-o out] [-e encout] [-k key]\n" +
		"\t  -l, --label     image flavor label\n" +
		"\t  -i, --in        input image file path\n" +
		"\t  -o, --out       (optional) output image flavor file path\n" +
		"\t                  if not specified, will print to the console\n" +
		"\t  -e, --encout    (optional) output encrypted image file path\n" +
		"\t                  if not specified, encryption is skipped\n" +
		"\t  -k, --key       (optional) existing key ID\n" +
		"\t                  if not specified, a new key is generated")
}

//Usage command line usage string
func fetchKeyUsage() {
	log.Trace("main:fetchKeyUsage() Entering")
	defer log.Trace("main:fetchKeyUsage() Leaving")

	fmt.Println("usage: wpm fetch-key [-k key]\n" +
		"\t  -k, --key       (optional) existing key ID\n" +
		"\t                  if not specified, a new key is generated\n" +
		"\t  -t, --asset-tag (optional) asset tags associated with the new key\n" +
		"\t                  tags are key:value separated by comma\n" +
		"\t  -a, --asymmetric (optional) specify to use asymmetric encryption\n" +
		"\t                  currently only supports RSA")
}

//Usage command line usage string
func containerFlavorUsage() {
	log.Trace("main:containerFlavorUsage() Entering")
	defer log.Trace("main:containerFlavorUsage() Leaving")

	fmt.Println("usage: wpm create-container-image-flavor -i img-name [-t tag] [-f dockerFile] [-d build-dir] [-k keyId]\n" +
		"                            [-e] [-s] [-n notaryServer] [-o out-file]\n" +

		"\t  -i, --img-name                  container image name\n" +
		"\t  -t, --tag                       (optional) container image tag name\n" +
		"\t  -f, --docker-file               (optional) container file path\n" +
		"\t                                  to build the container image\n" +
		"\t  -d, --build-dir                 (optional) build directory to build the\n" +
		"\t                                  container image. To be provided when container\n" +
		"\t                                  file path [-f] is provided as parameter\n" +
		"\t  -k, --key-id                    (optional) existing key ID\n" +
		"\t                                  if not specified, a new key is generated\n" +
		"\t  -e, --encryption-required       (optional) boolean parameter specifies if\n" +
		"\t                                  container image needs to be encrypted\n" +
		"\t  -s, --integrity-enforced        (optional) boolean parameter specifies if\n" +
		"\t                                  container image should be signed\n" +
		"\t  -n, --notary-server             (optional) specify notary server url\n" +
		"\t  -o, --out-file                  (optional) specify output file path")

}

func removeSecureDockerDaemon() {
	fmt.Println("Uninstalling secure-docker-daemon")
	_, err := exec.Command(consts.OptDirPath + "secure-docker-daemon/uninstall-secure-docker-daemon.sh").Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to uninstall secure-docker-daemon Error %s:", err.Error())
	}
}
