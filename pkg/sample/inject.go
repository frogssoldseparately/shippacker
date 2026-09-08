package sample

import (
	"archive/zip"
	"fmt"
	"path/filepath"

	"github.com/frogssoldseparately/shippacker/pkg/maps"
	"github.com/frogssoldseparately/simpleseek/sreader"
	"github.com/frogssoldseparately/simpleseek/swriter"
)

var storedOotSamples = map[string]*zip.File{}
var sampleQueueByAddress *map[uint32]string
var sampleQueueByName *map[string]uint32

func NewSampleQueue() {
	sampleQueueByAddress = &map[uint32]string{}
	sampleQueueByName = &map[string]uint32{}
}

func GetQueuedByAddress(assetAddr uint32) (string, bool) {
	assetName, ok := (*sampleQueueByAddress)[assetAddr]
	return assetName, ok
}

func AcceptQueuedSamples(gsm *maps.GameSampleMap) {
	for assetName, assetAddr := range *sampleQueueByName {
		(*gsm.ByAddress)[assetAddr] = assetName
		(*gsm.ByName)[assetName] = assetAddr
	}
}

func RegisterSamples(archive *sreader.SimpleZipReader) {
	for _, sampleEntry := range archive.GetAllByPrefix("audio/samples/") {
		assetName := filepath.Base(sampleEntry.Name)
		storedOotSamples[assetName] = sampleEntry
	}
}

func InjectSampleByAddress(zipWriter *swriter.SimpleZipWriter, assetAddr uint32, gsm *maps.GameSampleMap) (string, error) {
	if assetName, ok := (*sampleQueueByAddress)[assetAddr]; ok {
		return assetName, nil
	}
	assetName, ok := (*gsm.UnmappedByAddress)[assetAddr]
	if !ok {
		return "", fmt.Errorf("Ship of Harkinian sample of address %X could not be found\n", assetAddr)
	}
	sampleEntry, ok := storedOotSamples[assetName]
	if !ok {
		return "", fmt.Errorf("Ship of Harkinian sample of address %X could not be found\n", assetAddr)
	}
	fSample, err := sampleEntry.Open()
	if err != nil {
		return "", err
	}
	sample, err := ReadShipSample(fSample, assetAddr, assetName)
	if err != nil {
		return "", err
	}
	if err := zipWriter.WriteEntry(sample); err != nil {
		return "", err
	}
	QueueSample(assetName, assetAddr)
	return assetName, nil
}

func InjectSampleByName(zipWriter *swriter.SimpleZipWriter, assetName string, gsm *maps.GameSampleMap) (uint32, error) {
	if assetAddr, ok := (*sampleQueueByName)[assetName]; ok {
		return assetAddr, nil
	}
	assetAddr, ok := (*gsm.UnmappedByName)[assetName]
	if !ok {
		return 0, fmt.Errorf("Ship of Harkinian sample \"%s\" could not be found\n", assetName)
	}
	sampleEntry, ok := storedOotSamples[assetName]
	if !ok {
		return 0, fmt.Errorf("Ship of Harkinian sample \"%s\" could not be found\n", assetName)
	}
	fSample, err := sampleEntry.Open()
	if err != nil {
		return 0, err
	}
	sample, err := ReadShipSample(fSample, assetAddr, assetName)
	if err != nil {
		return 0, err
	}
	if err := zipWriter.WriteEntry(sample); err != nil {
		return 0, err
	}
	QueueSample(assetName, assetAddr)
	return assetAddr, nil
}

func QueueSample(assetName string, assetAddr uint32) {
	(*sampleQueueByName)[assetName] = assetAddr
	(*sampleQueueByAddress)[assetAddr] = assetName
}
