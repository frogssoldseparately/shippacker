package cli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var MusicSrcPath string
var O2ROutPath string

func ParseCommandLine() {
	flag.StringVar(&MusicSrcPath, "msrc", "./music", "Directory containing unmodified custom sequences")
	flag.StringVar(&O2ROutPath, "out", "./mods", "Directory to write .o2r mod file in")

	flag.Parse()

	srcOverride := flag.Arg(0)
	if len(srcOverride) != 0 {
		MusicSrcPath = srcOverride
		cwp, err := os.Executable()
		if err == nil {
			cwd := filepath.Dir(cwp)
			fmt.Printf("Changing working directory to: %s\n", cwd)
			os.Chdir(cwd)
		}
	}
}
