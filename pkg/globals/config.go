package globals

import (
	"fmt"
	"runtime"
	"strings"
)

const Platform string = runtime.GOOS

// Set by internal/cli
var Version string = "Keiichi_Charlie" // set by internal/cli
var RecurseSubdirectories bool = true  // set by internal/cli

// Currently unchanged
var RomPlatform string = "N64_US"
var StartingBankIndex uint64 = 41

// Initialized by SetupByVersion()
var MaxBankCount uint64
var MaxSongCount uint16
var UseCRC64Encoding bool
var UseNumericCategories bool

// Control flow
var WarnOnTooManyBanks bool = true
var WarnOnTooManySongs bool = true
var AllowCustomBanks bool = true
var EarlyExit bool = false
var HasOotO2r = false

func GetAudioXmlKey() string {
	return strings.ToLower(fmt.Sprintf("%s_%s", Version, RomPlatform))
}
