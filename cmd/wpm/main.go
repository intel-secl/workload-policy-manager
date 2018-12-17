package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {

	if len(os.Args[0:]) < 1 {
		log.Fatal("Usage :  wpm <task list>")
	}
	var task = os.Args[1]

	switch arg := strings.ToLower(task); arg {
	case "setup":
		log.Println("harshitha")

	case "create-image-flavor":
		log.Println("create image flavor")

	case "unintall":
		log.Println("Uninstall")

	case "help", "-help", "--help":
		usage()

	default:
		fmt.Printf("Unrecognized option : %s\n", arg)
		usage()
	}
}

func usage() {
	fmt.Println("Usage: $0 uninstall|create-image-flavor|create-software-flavor")
	fmt.Println("Usage: $0 setup [--force|--noexec] [task1 task2 ...]")
	fmt.Println("Usage: $0 export-config [outfile|--in=infile|--out=outfile|--stdout] [--env-password=PASSWORD_VAR]")
	fmt.Println("Available setup tasks:CreateEnvelopKey and RegisterEnvelopeKeyWithKBS")
}
