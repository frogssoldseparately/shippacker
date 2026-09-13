package globals

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/frogssoldseparately/simpleseek/swriter"
)

const OsPlatform string = runtime.GOOS

// Set by internal/cli
var RecurseSubdirectories bool = true

// Modified by SetupByVersion()
var RomVersion string = "Keiichi_Charlie_2S2H"
var RomPlatform string = "N64_US"
var TranslatableRomVersion string = "Ackbar_Delta_SoH"
var TranslatableRomPlatform string = "N64_NTSC_10"
var PortPlatform string = "2S2H"
var StartingBankIndex uint64 = 41
var MaxBankCount uint64
var MaxSongCount uint16
var UseCRC64Encoding bool
var UseNumericCategories bool

// Control flow
var WarnOnTooManyBanks bool = true
var WarnOnTooManySongs bool = true
var AllowCustomBanks bool = true
var EarlyExit bool = false
var HasImportedO2R = false

func GetNativeAudioXmlKey() string {
	return strings.ToLower(fmt.Sprintf("%s_%s", RomVersion, RomPlatform))
}

func GetTranslatedAudioXmlKey() string {
	return strings.ToLower(fmt.Sprintf("%s_%s", TranslatableRomVersion, TranslatableRomPlatform))
}

func GetCurrentBank(zipWriter *swriter.SimpleZipWriter) uint64 {
	return uint64(zipWriter.GetTypedFileCount("Soundfont")) + StartingBankIndex
}
