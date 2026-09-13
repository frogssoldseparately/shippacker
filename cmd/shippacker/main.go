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

// The entry point for shippacker.exe. Packages .*seq, .mmrs, and .seq+.meta pairs into an
// .o2r mod file with custom instrument banks.
func main() {
	failedStartup := false
	// Prompt for 2ship version
	cli.GetUserVersion()
	// Prompt for recursive directory search
	cli.GetRecurseDirectory()
	// Initialize global configs using selected version
	if err := globals.SetupByVersion(); err != nil {
		fmt.Printf("%s\n", err)
		failedStartup = true
	} else {
		// Get commandline arguments
		cli.ParseCommandLine()
		if err := cli.CheckInput(); err != nil {
			fmt.Printf("%s\n", err)
			failedStartup = true
		}
	}
	if !failedStartup {
		if err := TryO2RImport(); err != nil {
			fmt.Println(err)
		} else {
			globals.HasImportedO2R = true
		}
		err := shippacker.Pack(cli.MusicSrcPath, cli.O2ROutPath)
		if err != nil {
			fmt.Printf("%s\n", err)
		} else {
			fmt.Println("All Good")
			fmt.Printf(".o2r mod file in : %s\n", cli.O2ROutPath)
		}
	}
	fmt.Print("Press [ENTER] to close")
	fmt.Scanf(".")
}

func TryO2RImport() error {
	var filename string
	if globals.PortPlatform == "2S2H" {
		filename = "oot.o2r"
	} else {
		filename = "mm.o2r"
	}
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("If you would like to include translatable sequences, please place your %s file in the same directory as shippacker.exe\n", filename)
	}
	if err := ImportO2R(filename); err != nil {
		return fmt.Errorf("Could not unpack %s because\n\t%s\n", filename, err)
	}
	return nil
}

func ImportO2R(filename string) error {
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
