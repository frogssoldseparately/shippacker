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
var TargetOotSamples = map[uint32]string{
	0x10EFF0: "Small Cymbal_META",
	0x29FE30: "Ghostly Wind_META",
	0x2F8690: "Mechanical Rumbling_META",
	0x343870: "Church Organ_META",
	0x35C900: "Clap_META",
	0x35D100: "Sample_590_META",
	0x365DB0: "Sample_591_META",
	0x36C840: "Islamic Chant_META",
	0x38D230: "Off Hi-Hat_META",
	0x34D670: "Electric Organ_META",
	0x391FD0: "Drum and Sitar_META",
	0x395710: "Percussive Resonance_META",
	0x3D1980: "Flute and Voice_META",
	0x3D2840: "Off-Kilter Flute_META",
	0x3D37A0: "Voice_META",
	0x3FA9E0: "sample_2_003FA9E0_META",
	0x4006B0: "sample_3_004006B0_META",
	0x409270: "sample_3_00409270_META",
	0x40BA60: "sample_3_0040BA60_META",
	0x416B30: "sample_3_00416B30_META",
	0x4377E0: "sample_6_004377E0_META",
	0x4428B0: "sample_6_004428B0_META",
	0x447BC0: "sample_6_00447BC0_META",
}

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
