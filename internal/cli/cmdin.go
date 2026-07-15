package cli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var MusicSrcPath string
var O2RSrcPath string
var O2ROutPath string
var AudioXMLPath string

func ParseCommandLine() {
	flag.StringVar(&MusicSrcPath, "msrc", "./music", "Directory containing unmodified custom sequences")
	flag.StringVar(&O2RSrcPath, "o2r", "./mm.o2r", "Path to mm.o2r file")
	flag.StringVar(&O2ROutPath, "out", "./mods", "Directory to write .o2r mod file in")
	flag.StringVar(&AudioXMLPath, "xml", "./assets/xml/N64_US/audio/Audio.xml", "Path to Audio.xml file that comes with 2ship2harkinian")

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
