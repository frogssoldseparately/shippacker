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
	ByAddress         *SampleAddressMap
	ByName            *SampleNameMap
	UnmappedByAddress *SampleAddressMap
	UnmappedByName    *SampleNameMap
}

type SampleAddressMap map[uint32]string
type SampleNameMap map[string]uint32

func NewSampleMap() (*SampleMap, error) {
	nativeMap, err := NewNativeMap()
	if err != nil {
		return nil, err
	}
	translationMap, err := NewTranslationMap()
	if err != nil {
		return nil, err
	}
	if globals.PortPlatform == "2S2H" {
		return &SampleMap{
			MajorasMask:   nativeMap,
			OcarinaOfTime: translationMap,
		}, nil
	} else {
		return &SampleMap{
			MajorasMask:   translationMap,
			OcarinaOfTime: nativeMap,
		}, nil
	}
}

func NewNativeMap() (*GameSampleMap, error) {
	key := globals.GetNativeAudioXmlKey()
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
	unmappedByAddress := SampleAddressMap{}
	unmappedByName := SampleNameMap{}

	for _, sample := range sampleEntries {
		newName := sample.OriginalName + "_META"
		currentOffset, _ := sample.GetOffset()
		byAddress[currentOffset] = newName
		byName[newName] = currentOffset
	}
	return &GameSampleMap{
		ByAddress:         &byAddress,
		ByName:            &byName,
		UnmappedByAddress: &unmappedByAddress,
		UnmappedByName:    &unmappedByName,
	}, nil
}

func NewTranslationMap() (*GameSampleMap, error) {
	key := globals.GetTranslatedAudioXmlKey()
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
	unmappedByAddress := SampleAddressMap{}
	unmappedByName := SampleNameMap{}

	for _, sample := range sampleEntries {
		oldName := sample.OriginalName + "_META"
		newName := sample.Name + "_META"
		currentOffset, err := sample.GetOffset()
		if err != nil {
			return nil, err
		}
		if newName == "unknown_META" {
			unmappedByAddress[currentOffset] = oldName
			unmappedByName[oldName] = currentOffset
		} else {
			byAddress[currentOffset] = newName
			byName[oldName] = currentOffset
		}
	}
	return &GameSampleMap{
		ByAddress:         &byAddress,
		ByName:            &byName,
		UnmappedByAddress: &unmappedByAddress,
		UnmappedByName:    &unmappedByName,
	}, nil
}
