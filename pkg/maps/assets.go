package maps

import (
	"fmt"

	"github.com/frogssoldseparately/shippacker/pkg/globals"
)

// A map of sample pointers (relative addresses within then sample banks on ROM) to the asset names used by 2Ship2Harkinian.
type AssetMap map[uint32]string

// Generate an asset map using mm.o2r and Audio.xml for .zbank -> Soundfont conversion.
func NewAssetMap() (*AssetMap, error) {
	key := globals.GetAudioXmlKey()
	buf, ok := globals.SampleXmls[key]
	if !ok {
		return nil, fmt.Errorf("Unknown xml key %s", key)
	}
	sampleEntries, err := ParseXML(buf)
	if err != nil {
		return nil, err
	}
	out := AssetMap{}

	for _, sample := range sampleEntries {
		currentOffset, _ := sample.GetOffset()
		out[currentOffset] = sample.Name + "_META"
	}
	return &out, nil
}
