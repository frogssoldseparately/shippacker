package cli

import (
	"fmt"

	"github.com/frogssoldseparately/shippacker/pkg/globals"
	"github.com/frogssoldseparately/shippacker/pkg/iohelper"
)

func DeduplicateSamplesPrompt() {
	msg := fmt.Sprintf("[DEDUPLICATION] Only copy vanilla samples that %s is missing for %s sequences (usually y)?", globals.PortPlatform, globals.TranslatableExtension)
	if globals.DeduplicateTranslationSamples = iohelper.PromptYesNo(msg); globals.DeduplicateTranslationSamples {
		fmt.Printf("shippacker will only copy missing samples\n\n")
	} else {
		fmt.Printf("shippacker will copy all requested samples\n\n")
	}
}

func AllowTranslationPrompt() {
	msg := fmt.Sprintf("[TRANSLATION] Translate %s sequences to %s sequences?", globals.TranslatableExtension, globals.NativeExtension)
	if globals.AllowTranslation = iohelper.PromptYesNo(msg); globals.AllowTranslation {
		fmt.Printf("shippacker will translate %s sequences for use in %s\n\n", globals.TranslatableExtension, globals.PortPlatform)
	} else {
		fmt.Printf("shippacker will only use %s sequences\n\n", globals.NativeExtension)
	}
}

func RecurseDirectoryPrompt() {
	msg := "[RECURSION] Read the music directory recursively?"
	if globals.RecurseSubdirectories = iohelper.PromptYesNo(msg); globals.RecurseSubdirectories {
		fmt.Printf("shippacker will read the music directory and its subdirectories\n\n")
	} else {
		fmt.Printf("shippacker will only read the music directory\n\n")
	}
}

func UserVersionPrompt() {
	promptMsg := "[VERSION] Version being built for:"
	opts := globals.SupportedVersions
	badMsg := "Invalid version number"
	globals.RomVersion = iohelper.PromptNumberedOption(promptMsg, opts, badMsg)
}
