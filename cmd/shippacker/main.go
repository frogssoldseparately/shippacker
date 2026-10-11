//go:build !wasm

package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/frogssoldseparately/shippacker/internal/cli"
	"github.com/frogssoldseparately/shippacker/pkg/globals"
	"github.com/frogssoldseparately/shippacker/pkg/sample"
	"github.com/frogssoldseparately/shippacker/pkg/shippacker"
	"github.com/frogssoldseparately/shippacker/pkg/soundfont"
	"github.com/frogssoldseparately/simpleseek/sreader"
)

// The entry point for shippacker.exe. Packages .*seq, .mmrs, .ootrs, and .seq+.meta pairs into an
// .o2r mod file with custom instrument banks.
func main() {
	if err := startup(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("All Good\n.o2r mod file in : %s\n", cli.O2ROutPath)
	}
	fmt.Print("Press [ENTER] to close")
	fmt.Scanf(".")
}

func startup() error {
	if err := getInputs(); err != nil {
		return err
	}
	if globals.HasImportedO2R {
		// Prompt for force copying of samples
		cli.DeduplicateSamplesPrompt()
	}
	return shippacker.Pack(cli.MusicSrcPath, cli.O2ROutPath)
}

func getInputs() error {
	// Prompt for 2ship version
	cli.UserVersionPrompt()
	// Initialize global configs using selected version
	if err := globals.SetupByVersion(); err != nil {
		return err
	}
	// Prompt for recursive directory search
	cli.RecurseDirectoryPrompt()
	// Prompt for translation between zelda64 games
	cli.AllowTranslationPrompt()
	// Get commandline arguments
	cli.ParseCommandLine()
	if err := cli.CheckInput(); err != nil {
		return err
	}
	if globals.AllowTranslation {
		if err := tryO2RImport(); err != nil {
			fmt.Println(err)
		} else {
			globals.HasImportedO2R = true
		}
	}
	return nil
}

func tryO2RImport() error {
	var filename string
	if globals.PortPlatform == "2S2H" {
		filename = "oot.o2r"
	} else {
		filename = "mm.o2r"
	}
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("If you would like to include %s sequences, please place your %s file in the same directory as shippacker.exe\n", globals.TranslatableExtension, filename)
	}
	if err := importO2R(filename); err != nil {
		return fmt.Errorf("Could not unpack %s because\n\t%s\n", filename, err)
	}
	return nil
}

func importO2R(filename string) error {
	if archive, err := sreader.OpenArchive(filename); err != nil {
		return err
	} else {
		if err := soundfont.RegisterSoundfonts(archive); err != nil {
			return err
		}
		sample.RegisterSamples(archive)
	}
	return nil
}
