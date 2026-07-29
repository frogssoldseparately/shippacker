package globals

import (
	"fmt"
	"runtime"
	"strings"
)

const Platform string = runtime.GOOS

var Version string = "Keiichi_Charlie"
var RomPlatform string = "N64_US"

func GetAudioXmlKey() string {
	return strings.ToLower(fmt.Sprintf("%s_%s", Version, RomPlatform))
}
