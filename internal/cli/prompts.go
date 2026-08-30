package cli

import (
	"fmt"
	"strings"

	"github.com/frogssoldseparately/shippacker/pkg/globals"
)

func GetRecurseDirectory() {
	fmt.Printf("Read the music directory recursively? (y/n) ")
	var response string
	fmt.Scan(&response)
	fmt.Scanln() // discard newline?
	globals.RecurseSubdirectories = strings.ToLower(response[0:1]) == "y"
	if globals.RecurseSubdirectories {
		fmt.Println("Will read music directory and subdirectories")
	} else {
		fmt.Println("Will only read music directory")
	}
}

func GetUserVersion() {
	svs := globals.SupportedVersions
	svCount := len(svs)
	fmt.Println("Please pick the version you are building for:")
	userInput := -1
	for i, version := range globals.SupportedVersions {
		fmt.Printf("%d) %s\n", i, version)
	}
	// Continue prompting until a valid number has been provided
	for true {
		fmt.Scan(&userInput)
		fmt.Scanln() // discard newline?
		if userInput >= 0 && userInput < svCount {
			break
		}
		fmt.Println("Invalid version number")
	}
	globals.Version = svs[userInput]
}
