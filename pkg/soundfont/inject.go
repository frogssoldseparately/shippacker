package soundfont

import (
	"archive/zip"
	"fmt"
	maps0 "maps"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/frogssoldseparately/shippacker/pkg/globals"
	"github.com/frogssoldseparately/shippacker/pkg/maps"
	"github.com/frogssoldseparately/simpleseek/sreader"
	"github.com/frogssoldseparately/simpleseek/swriter"
)

var storedSoundfonts = map[uint64]*zip.File{}
var includedBanks = map[uint64]uint64{}
var soundfontQueue *map[uint64]uint64

func NewSoundfontQueue() {
	soundfontQueue = &map[uint64]uint64{}
}

func AcceptQueuedSoundfonts() {
	if soundfontQueue != nil {
		maps0.Copy(includedBanks, *soundfontQueue)
	}
	soundfontQueue = nil
}

func RegisterSoundfonts(archive *sreader.SimpleZipReader) error {
	// A little messy, but this ensures wasm builds don't store two different o2rs
	storedSoundfonts = map[uint64]*zip.File{}
	for _, soundfontEntry := range archive.GetAllByPrefix("audio/fonts/") {
		soundfontName := filepath.Base(soundfontEntry.Name)
		var bankStr string
		if globals.PortPlatform == "2S2H" {
			bankStr = soundfontName[0:strings.Index(soundfontName, "_")]
		} else {
			bankStr = soundfontName[strings.Index(soundfontName, "_")+1:]
		}
		bankId, err := strconv.ParseUint(bankStr, 10, 32)
		if err != nil {
			return err
		}
		storedSoundfonts[bankId] = soundfontEntry
	}
	return nil
}

func InjectSoundfont(zipWriter *swriter.SimpleZipWriter, metaBank string, currentBankId *uint64, name string, gsm *maps.GameSampleMap) error {
	parsedBank, err := strconv.ParseUint(strings.TrimPrefix(metaBank, "0x"), 16, 32)
	if err != nil {
		return err
	}
	if usedBankId, ok := includedBanks[parsedBank]; ok {
		*currentBankId = usedBankId
		return nil
	}
	if !globals.AllowCustomBanks {
		return fmt.Errorf("it has a custom bank\n")
	}
	soundfontEntry, ok := storedSoundfonts[parsedBank]
	if !ok {
		return fmt.Errorf("could not find translatable soundfont with id %s\n", metaBank)
	}
	fSoundfont, err := soundfontEntry.Open()
	if err != nil {
		return err
	}
	sf, err := ReadSoundfont(fSoundfont, name, zipWriter, gsm)
	if err != nil {
		return err
	}
	if err := zipWriter.WriteEntry(sf); err != nil {
		return err
	}
	QueueSoundfont(parsedBank, *currentBankId)
	return nil
}

func QueueSoundfont(vanillaBankId uint64, realBankId uint64) {
	(*soundfontQueue)[vanillaBankId] = realBankId
}
