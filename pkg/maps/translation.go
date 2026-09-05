package maps

import (
	"fmt"

	"github.com/frogssoldseparately/shippacker/pkg/globals"
)

type TranslationMap map[string]uint32

func NewTranslationMaps() (*AssetMap, *TranslationMap, error) {
	key := "ackbar_delta_n64_ntsc_10"
	buf, ok := globals.SampleXmls[key]
	if !ok {
		return nil, nil, fmt.Errorf("Unknown xml key %s", key)
	}
	sampleEntries, err := ParseXML(buf)
	if err != nil {
		return nil, nil, err
	}
	am := AssetMap{}
	tm := TranslationMap{}

	for _, sample := range sampleEntries {
		currentOffset, _ := sample.GetOffset()
		am[currentOffset] = sample.Name + "_META"
		tm[sample.OriginalName+"_META"] = currentOffset
	}

	return &am, &tm, nil
}
