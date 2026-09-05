//go:build !wasm

package main

import (
	"archive/zip"
	"fmt"
	"strings"

	"github.com/frogssoldseparately/shippacker/internal/cli"
	"github.com/frogssoldseparately/shippacker/pkg/globals"
	"github.com/frogssoldseparately/shippacker/pkg/ootrs"
	"github.com/frogssoldseparately/shippacker/pkg/shippacker"
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
		if ootO2RArchive, err := sreader.OpenArchive("oot.o2r"); err == nil {
			soundfontEntries := []*zip.File{}
			for _, entry := range ootO2RArchive.GetFiles() {
				if strings.Contains(entry.Name, "audio/fonts/") {
					soundfontEntries = append(soundfontEntries, entry)
				}
			}
			if err := ootrs.PrepareOotSoundfonts(&soundfontEntries); err == nil {
				globals.HasOotO2r = true
			} else {
				fmt.Printf("Could not unpack oot.o2r because %s\n", err)
			}
		}
		// Get packing!
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
