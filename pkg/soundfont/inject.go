package soundfont

import (
	"archive/zip"
	"fmt"
	maps0 "maps"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/frogssoldseparately/shippacker/pkg/maps"
	"github.com/frogssoldseparately/simpleseek/sreader"
	"github.com/frogssoldseparately/simpleseek/swriter"
)

var storedOotSoundfonts = map[uint64]*zip.File{}
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
	for _, soundfontEntry := range archive.GetAllByPrefix("audio/fonts/") {
		soundfontName := filepath.Base(soundfontEntry.Name)
		bankString := soundfontName[0:strings.Index(soundfontName, "_")]
		bankId, err := strconv.ParseUint(bankString, 10, 32)
		if err != nil {
			return err
		}
		storedOotSoundfonts[bankId] = soundfontEntry
	}
	return nil
}

func InjectSoundfont(zipWriter *swriter.SimpleZipWriter, metaBank string, currentBankId *uint64, name string, gsm *maps.GameSampleMap) error {
	var parsedBank uint64
	var err error
	if len(metaBank) >= 2 && metaBank[0:2] == "0x" {
		parsedBank, err = strconv.ParseUint(metaBank[2:], 16, 32)
		if err != nil {
			return err
		}
	} else {
		parsedBank, err = strconv.ParseUint(metaBank, 16, 32)
		if err != nil {
			return err
		}
	}
	if usedBankId, ok := includedBanks[parsedBank]; ok {
		*currentBankId = usedBankId
		return nil
	}
	soundfontEntry, ok := storedOotSoundfonts[parsedBank]
	if !ok {
		return fmt.Errorf("could not find OoT bank with id %s\n", metaBank)
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

func QueueSoundfont(ootBankId uint64, realBankId uint64) {
	(*soundfontQueue)[ootBankId] = realBankId
}
