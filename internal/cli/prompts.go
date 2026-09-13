package cli

import (
	"fmt"

	"github.com/frogssoldseparately/shippacker/pkg/globals"
	"github.com/frogssoldseparately/shippacker/pkg/iohelper"
)

func GetRecurseDirectory() {
	msg := "Read the music directory recursively?"
	if globals.RecurseSubdirectories = iohelper.PromptYesNo(msg); globals.RecurseSubdirectories {
		fmt.Printf("shippacker will read the music directory and its subdirectories\n\n")
	} else {
		fmt.Printf("shippacker will only read the music directory\n\n")
	}
}

func GetUserVersion() {
	promptMsg := "Version being built for:"
	opts := globals.SupportedVersions
	badMsg := "Invalid version number"
	globals.RomVersion = iohelper.PromptNumberedOption(promptMsg, opts, badMsg)
}
