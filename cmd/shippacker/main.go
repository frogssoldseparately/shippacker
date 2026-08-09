package main

import (
	"fmt"

	"github.com/frogssoldseparately/shippacker/internal/cli"
	"github.com/frogssoldseparately/shippacker/pkg/maps"
	"github.com/frogssoldseparately/shippacker/pkg/shippacker"
)

// The entry point for shippacker.exe. Packages .*seq, .mmrs, and .seq+.meta pairs into an
// .o2r mod file with custom instrument banks.
func main() {
	// Get user input
	cli.ParseCommandLine()
	if err := cli.CheckInput(); err != nil {
		fmt.Printf("%s\n", err)
	} else {
		paths := maps.BundlePaths(cli.MusicSrcPath, cli.O2ROutPath)
		// Get to packing
		err := shippacker.Pack(paths)
		if err != nil {
			fmt.Printf("%s\n", err)
		} else {
			fmt.Println("All Good")
			fmt.Printf(".o2r mod file in : %s\n", paths.OOut)
		}
	}
	fmt.Print("Press [ENTER] to close")
	fmt.Scanf(".")
}
