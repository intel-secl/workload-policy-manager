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
	csetup "intel/isecl/lib/common/setup"
	"intel/isecl/lib/common/validation"
	"intel/isecl/wpm/config"
	"intel/isecl/wpm/consts"

	"github.com/pkg/errors"

	containerImageFlavor "intel/isecl/wpm/pkg/containerimageflavor"
	imageFlavor "intel/isecl/wpm/pkg/imageflavor"
	"intel/isecl/wpm/pkg/setup"
	"intel/isecl/wpm/pkg/util"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"

	commLog "intel/isecl/lib/common/log"

	"github.com/google/uuid"
)

var (
	// Version holds the version number for the WPM binary
	Version string = ""
	// Time holds the build timestamp for the WPM binary
	Time string = ""
	// Branch holds the git build branch for the WPM binary
	Branch string = ""
	log           = commLog.GetDefaultLogger()
	secLog        = commLog.GetSecurityLogger()
)

func printVersion() {
	fmt.Printf("Version %s\nBuild %s at %s\n", Version, Branch, Time)
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
		// Set log configurations
		err = config.LogConfiguration(false, true)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error configuring logging.")
		}

		// Check if nosetup environment variable is true, if yes then skip the setup tasks
		if nosetup, err := strconv.ParseBool(os.Getenv("WPM_NOSETUP")); err == nil && nosetup == true {
			fmt.Println("WPM_NOSETUP is set to true. Skipping setup tasks.")
			log.Info("main:main() WPM_NOSETUP is set to true. Skipping setup.")
			os.Exit(1)
		}

		// Everytime, we run setup, need to make sure that the configuration is complete
		// So lets run the Configurer as a seperate runner. We could have made a single runner
		// with the first task as the Configurer. However, the logic in the common setup task
		// runner runs only the tasks passed in the argument if there are 1 or more tasks.
		// This means that with current logic, if there are no specific tasks passed in the
		// argument, we will only run the confugurer but the intention was to run all of them

		// TODO : The right way to address this is to pass the arguments from the commandline
		// to a functon in the workload agent setup package and have it build a slice of tasks
		// to run.
		flags := args

		// if setup is provided without task name then print usage and quit
		if len(args) == 1 {
			usage()
			os.Exit(0)
		}

		if len(args) >= 2 &&
			args[1] != "CreateEnvelopeKey" &&
			args[1] != "download_ca_cert" &&
			args[1] != "download_cert" &&
			args[1] != "all" {
			fmt.Printf("Unrecognized command: %s %s\n", args[0], args[1])
			os.Exit(1)
		}

		if len(args) > 1 {
			flags = args[2:]
			if args[1] == "download_cert" && len(args) > 2 {
				flags = args[3:]
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
			fmt.Fprintln(os.Stderr, "Error loading configuration.")
			os.Exit(1)
		}

		// save configuration from config.yml
		err = config.SaveConfiguration(context)

		if err != nil {
			fmt.Fprintln(os.Stderr, "Error updating configuration. Refer logs for more information.")
			log.WithError(err).Errorf("main:main() Error updating configuration: %+v\n", err)
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
					KeyFile:            consts.FlavorSigningKeyPath,
					CertFile:           consts.FlavorSigningCertPath,
					KeyAlgorithm:       consts.DefaultKeyAlgorithm,
					KeyAlgorithmLength: consts.DefaultKeyAlgorithmLength,
					CmsBaseURL:         config.Configuration.Cms.BaseURL,
					Subject: pkix.Name{
						Country:      []string{config.Configuration.Subject.Country},
						Organization: []string{config.Configuration.Subject.Organization},
						Locality:     []string{config.Configuration.Subject.Locality},
						Province:     []string{config.Configuration.Subject.Province},
						CommonName:   config.Configuration.Subject.CommonName,
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
			log.WithError(err).Errorf("main:main() Error running setup tasks %s: %s\n", strings.Join(args[1:], " "), err.Error())
			log.Tracef("%+v", err)
			os.Exit(1)
		}

	case "create-image-flavor":
		// Set log configurations
		err := config.LogConfiguration(false, true)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error configuring logging.")
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
			imageFlavorUsage()
			os.Exit(1)
		}

		// validate input strings
		inputArr := []string{*flavorLabel, *outputFlavorFilePath, *inputImageFilePath, *outputEncImageFilePath}
		if validationErr := validation.ValidateStrings(inputArr); validationErr != nil {
			fmt.Fprintln(os.Stderr, "Error creating VM image flavor: Invalid string format")
			log.WithError(err).Errorf("main:main() Error creating VM image flavor: %+v\n", err)
			imageFlavorUsage()
			os.Exit(1)
		}

		//If the key ID is specified, make sure it's a valid UUID
		if len(strings.TrimSpace(*keyID)) > 0 {
			if validatekeyIDErr := validation.ValidateUUIDv4(*keyID); validatekeyIDErr != nil {
				fmt.Fprintln(os.Stderr, "Error creating VM image flavor: Invalid Key UUID format")
				log.WithError(validatekeyIDErr).Errorf("main:main() Error creating VM image flavor: %+v\n", validatekeyIDErr)
				imageFlavorUsage()
				os.Exit(1)
			}
		}

		imageFlavor, err := imageFlavor.CreateImageFlavor(*flavorLabel, *outputFlavorFilePath, *inputImageFilePath,
			*outputEncImageFilePath, *keyID, false)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error creating VM image flavor. Check logs for more information.")
			log.WithError(err).Errorf("main:main() Error creating VM image flavor: %+v\n", err)
			os.Exit(1)
		}
		if len(imageFlavor) > 0 {
			log.Info("main:main() Successfully created VM image flavor")
			fmt.Println(imageFlavor)
		}

	case "create-container-image-flavor":
		// Set log configurations
		err := config.LogConfiguration(false, true)
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
			fmt.Println("Flavor label and image file path should be given.")
			containerFlavorUsage()
			os.Exit(1)
		}

		// validate input strings
		inputArr := []string{*imageName, *tagName, *dockerFilePath, *buildDir, *outputFlavorFilePath}
		if validationErr := validation.ValidateStrings(inputArr); validationErr != nil {
			fmt.Fprintln(os.Stderr, "Error Creating Container Flavor: Invalid Input String Format")
			log.WithError(validationErr).Errorf("main:main() Error Creating Container Flavor: %+v\n", validationErr)
			containerFlavorUsage()
			os.Exit(1)
		}

		//If the key ID is specified, make sure it's a valid UUID
		if len(strings.TrimSpace(*keyID)) > 0 {
			if validatekeyIDErr := validation.ValidateUUIDv4(*keyID); validatekeyIDErr != nil {
				fmt.Fprintln(os.Stderr, "Error Creating Container Flavor: Invalid UUID")
				log.WithError(validatekeyIDErr).Errorf("main:main() Error Creating Container Flavor: %+v\n", validatekeyIDErr)
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
				log.WithError(validateURLErr).Errorf("Error Creating Container Flavor: %+v\n", validateURLErr)
				containerFlavorUsage()
				os.Exit(1)
			}
		}

		containerImageFlavor, err := containerImageFlavor.CreateContainerImageFlavor(*imageName, *tagName, *dockerFilePath, *buildDir,
			*keyID, *encryptionRequired, *integrityEnforced, *notaryServerURL, *outputFlavorFilePath)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error Creating Container Flavor! Check logs for more information.")
			log.WithError(err).Errorf("main:main() Error Creating Container Flavor: %+v", err)
			os.Exit(1)
		}

		if len(containerImageFlavor) > 0 {
			log.Info("main:main() Successfully created container image flavor")
			fmt.Println(containerImageFlavor)
		}

	case "unwrap-key":
		err := config.LogConfiguration(false, true)
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
			log.WithError(err).Errorf("main:main() Error unwrapping key: %+v\n", err)
			os.Exit(1)
		}

		wrappedKey, err := ioutil.ReadFile(*wrappedKeyFilePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error unwrapping key: Unable to read from wrapped key file: %s\n", *wrappedKeyFilePath)
			log.WithError(err).Errorf("main:main() Error unwrapping key: Unable to read from wrapped key file: %+v\n", err)
			os.Exit(1)
		}

		unwrappedKey, err := util.UnwrapKey(wrappedKey, consts.EnvelopePrivatekeyLocation)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error unwrapping key %s - check logs\n", *wrappedKeyFilePath)
			log.WithError(err).Errorf("main:main() Error unwrapping key: %+v\n" + err.Error())
			os.Exit(1)
		}
		log.Info("main:main() Successfully unwrapped key")
		fmt.Println(base64.StdEncoding.EncodeToString(unwrappedKey))

	case "get-container-image-id":
		err := config.LogConfiguration(false, true)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error configuring logging.")

		}

		if len(args[1:]) < 1 {
			fmt.Fprintf(os.Stderr, "Invalid number of parameters")
			os.Exit(1)
		}
		NameSpaceDNS := uuid.Must(uuid.Parse("6ba7b810-9dad-11d1-80b4-00c04fd430c8"))
		imageUUID := uuid.NewHash(md5.New(), NameSpaceDNS, []byte(args[1]), 4)
		log.Info("main:main() Successfully retrieved container image ID")
		fmt.Println(imageUUID)

	case "uninstall":
		config.LogConfiguration(false, false)
		fmt.Println("Uninstalling WPM")
		if len(args) > 1 && strings.ToLower(args[1]) == "--purge" {
			deleteFiles(consts.ConfigDirPath)
		}
		errorFiles, err := deleteFiles(consts.WpmSymLink, consts.OptDirPath, consts.ConfigDirPath, consts.LogDirPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error deleting files : %s\n", errorFiles)
		}

	case "help", "-help", "--help":
		usage()

	case "--version", "-v", "version", "-version":
		printVersion()

	case "create-software-flavor":
		fmt.Println("Not supported")

	default:
		fmt.Printf("Unrecognized command: %s\n", arg)
	}
}

func runCommand(cmd string, args []string) (string, error) {
	log.Trace("main:runCommand() Entering")
	defer log.Trace("main:runCommand() Leaving")

	out, err := exec.Command(cmd, args...).Output()
	return string(out), err
}

func usage() {
	log.Trace("main:usage() Entering")
	defer log.Trace("main:usage() Leaving")

	fmt.Printf("Workload Policy Manager\n")
	fmt.Printf("usage : %s <command> [<args>]\n\n", os.Args[0])
	fmt.Printf("Following are the list of commands\n")
	fmt.Printf("\tcreate-image-flavor|create-container-image-flavor|get-container-image-id|create-software-flavor|unwrap-key|uninstall|--help|--version\n\n")
	fmt.Printf("\tcreate-image-flavor - Used to create image flavors and encrypt the image\n")
	fmt.Printf("\t")
	imageFlavorUsage()
	fmt.Printf("\tcreate-container-image-flavor - Used to create container image flavors and encrypt the container image\n")
	fmt.Printf("\t")
	containerFlavorUsage()
	fmt.Printf("\n\t%s get-container-image-id [<sha256 of image id>] - Used to get the container image ID given the sha256 of the image\n", os.Args[0])
	fmt.Printf("\n\t%s unwrap-key [-i in]\n", os.Args[0])
	fmt.Printf("\t\t          -i, --in        wrapped key file path\n")
	fmt.Printf("\n\tuninstall          Uninstall wpm\n")
	fmt.Printf("\n\tuninstall --purge  Uninstalls wpm and deletes the existing configuration directory\n")
	fmt.Printf("\n\tusage : %s setup [<tasklist>]\n", os.Args[0])
	fmt.Printf("\t\t<tasklist>-space separated list of tasks\n")
	fmt.Printf("\t\t\t-Supported setup tasks - all download_ca_cert download_cert CreateEnvelopeKey\n")
	fmt.Printf("\tExample :-\n")
	fmt.Printf("\t\t%s setup all\n", os.Args[0])
	fmt.Printf("\t\t        - Runs all the setup tasks required in the right order\n")
	fmt.Printf("\t\t%s setup CreateEnvelopeKey\n", os.Args[0])
	fmt.Printf("\t\t        - Option [--force] overwrites any existing keypairs\n")
	fmt.Printf("\t\t%s setup download_ca_cert [--force]\n", os.Args[0])
	fmt.Printf("\t\t        - Download CMS root CA certificate\n")
	fmt.Printf("\t\t        - Option [--force] overwrites any existing files, and always downloads new root CA cert\n")
	fmt.Printf("\t\t       - Environment variable CMS_BASE_URL=<url> for CMS API URL\n")
	fmt.Printf("\t\t       - Environment variable CMS_TLS_CERT_SHA384=<sha384ForCMSTLSCert>\n")
	fmt.Printf("\t\t%s setup download_cert Flavor-Signing [--force]\n", os.Args[0])
	fmt.Printf("\t\t        - Generates Key pair and CSR, gets it signed from CMS\n")
	fmt.Printf("\t\t        - Option [--force] overwrites any existing files, and always downloads newly signed Flavor Signing cert\n")
	fmt.Printf("\t\t        - Environment variable CMS_BASE_URL=<url> for CMS API URL\n")
	fmt.Printf("\t\t        - Environment variable BEARER_TOKEN=<token> for downloading signed certificate from CMS\n")
	fmt.Printf("\t\t        - Environment variable CERT_PATH=<cert_path> to override default specified in config\n")
	fmt.Printf("\t\t        - Environment variable WPM_FLAVOR_SIGN_CERT_CN=<COMMON NAME> to override default specified in config\n")
	fmt.Printf("\t\t        - Environment variable WPM_CERT_ORG=<CERTIFICATE ORGANIZATION> to override default specified in config\n")
	fmt.Printf("\t\t        - Environment variable WPM_CERT_COUNTRY=<CERTIFICATE COUNTRY> to override default specified in config\n")
	fmt.Printf("\t\t        - Environment variable WPM_CERT_LOCALITY=<CERTIFICATE LOCALITY> to override default specified in config\n")
	fmt.Printf("\t\t        - Environment variable WPM_CERT_PROVINCE=<CERTIFICATE PROVINCE> to override default specified in config\n")
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
		"  -l, --label     image flavor label\n" +
		"  -i, --in        input image file path\n" +
		"  -o, --out       (optional) output image flavor file path\n" +
		"                  if not specified, will print to the console\n" +
		"  -e, --encout    (optional) output encrypted image file path\n" +
		"                  if not specified, encryption is skipped\n" +
		"  -k, --key       (optional) existing key ID\n" +
		"                  if not specified, a new key is generated\n")
}

//Usage command line usage string
func containerFlavorUsage() {
	log.Trace("main:containerFlavorUsage() Entering")
	defer log.Trace("main:containerFlavorUsage() Leaving")

	fmt.Println("usage: wpm create-container-image-flavor [-i img-name] [-t tag] [-f dockerFile] [-d build-dir] [-k keyId]\n" +
		"                            [-e] [-s] [-n notaryServer] [-o out-file]\n" +
		"  -i,       --img-name                     container image name\n" +
		"  -t,       --tag                          (optional)container image tag name\n" +
		"  -f,       --docker-file                  (optional) container file path\n" +
		"                                           to build the container image\n" +
		"  -d,       --build-dir                    (optional) build directory to\n" +
		"                                           build the container image\n" +
		"  -k,       --key-id                       (optional) existing key ID\n" +
		"                                           if not specified, a new key is generated\n" +
		"  -e,     --encryption-required            (optional) boolean parameter specifies if\n" +
		"                                           container image needs to be encrypted\n" +
		"  -s, 	   --integrity-enforced             (optional) boolean parameter specifies if\n" +
		"                                           container image should be signed\n" +
		"  -n,       --notary-server                (optional) specify notary server url\n" +
		"  -o,       --out-file                     (optional) specify output file path\n")

}
