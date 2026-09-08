package maps

import (
	"fmt"

	"github.com/frogssoldseparately/shippacker/pkg/globals"
)

type SampleMap struct {
	MajorasMask   *GameSampleMap
	OcarinaOfTime *GameSampleMap
}

type GameSampleMap struct {
	ByAddress *SampleAddressMap
	ByName    *SampleNameMap
	Unmapped  *SampleNameMap
}

type SampleAddressMap map[uint32]string
type SampleNameMap map[string]uint32

func NewSampleMap() (*SampleMap, error) {
	mmrsMap, err := NewMmrsMap()
	if err != nil {
		return nil, err
	}
	ootrsMap, err := NewOotrsMap()
	if err != nil {
		return nil, err
	}
	return &SampleMap{
		MajorasMask:   mmrsMap,
		OcarinaOfTime: ootrsMap,
	}, nil
}

// Generate an asset map using mm.o2r and Audio.xml for .zbank -> Soundfont conversion.
func NewMmrsMap() (*GameSampleMap, error) {
	key := globals.GetAudioXmlKey()
	buf, ok := globals.SampleXmls[key]
	if !ok {
		return nil, fmt.Errorf("Unknown xml key %s", key)
	}
	sampleEntries, err := ParseXML(buf)
	if err != nil {
		return nil, err
	}

	byAddress := SampleAddressMap{}
	byName := SampleNameMap{}
	unmapped := SampleNameMap{}

	for _, sample := range sampleEntries {
		newName := sample.Name + "_META"
		currentOffset, _ := sample.GetOffset()
		byAddress[currentOffset] = newName
		byName[newName] = currentOffset
	}
	return &GameSampleMap{
		ByAddress: &byAddress,
		ByName:    &byName,
		Unmapped:  &unmapped,
	}, nil
}

func NewOotrsMap() (*GameSampleMap, error) {
	key := "ackbar_delta_n64_ntsc_10"
	buf, ok := globals.SampleXmls[key]
	if !ok {
		return nil, fmt.Errorf("Unknown xml key %s", key)
	}
	sampleEntries, err := ParseXML(buf)
	if err != nil {
		return nil, err
	}

	byAddress := SampleAddressMap{}
	byName := SampleNameMap{}
	unmapped := SampleNameMap{}

	for _, sample := range sampleEntries {
		oldName := sample.OriginalName + "_META"
		newName := sample.Name + "_META"
		currentOffset, err := sample.GetOffset()
		if err != nil {
			return nil, err
		}
		if newName == "unknown" {
			unmapped[oldName] = currentOffset
		} else {
			byAddress[currentOffset] = newName
			byName[oldName] = currentOffset
		}
	}
	return &GameSampleMap{&byAddress, &byName, &unmapped}, nil
}
