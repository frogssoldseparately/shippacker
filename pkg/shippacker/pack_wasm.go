//go:build wasm

package shippacker

import (
	"encoding/binary"
	"fmt"
	"path/filepath"

	"github.com/frogssoldseparately/shippacker/pkg/globals"
	"github.com/frogssoldseparately/shippacker/pkg/maps"
	"github.com/frogssoldseparately/shippacker/pkg/mmrs"
	"github.com/frogssoldseparately/shippacker/pkg/ootrs"
	"github.com/frogssoldseparately/simpleseek/swriter"
)

func Pack(srcPaths []string) []byte {
	zipWriter := swriter.NewEmptyZipWriter(binary.LittleEndian)
	mmrsAssetMap, err := maps.NewAssetMap()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	ootrsAssetMap, ootrsTranslationMap, err := maps.NewTranslationMaps()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if globals.HasOotO2r {
		if err := ootrs.InjectOotSamples(zipWriter, ootrsAssetMap, ootrsTranslationMap); err != nil {
			fmt.Printf("Could not inject oot samples because %s\n", err)
		}
	}
	if err := WriteModEntries(srcPaths, zipWriter, mmrsAssetMap, ootrsAssetMap, ootrsTranslationMap); err != nil {
		fmt.Println(err)
		return nil
	}
	if zipWriter.GetTypedFileCount("Sequence") != 0 {
		o2rWriter := zipWriter.Finish()
		return *o2rWriter.GetBuffer()
	}
	fmt.Println("Nothing to write")
	return nil
}

func WriteModEntries(srcPaths []string, zipWriter *swriter.SimpleZipWriter, mmrsAssetMap *maps.AssetMap, ootrsAssetMap *maps.AssetMap, ootrsTranslationMap *maps.TranslationMap) error {
	for _, path := range srcPaths {
		switch filepath.Ext(path) {
		case ".mmrs":
			if err := mmrs.RepackArchive(path, zipWriter, mmrsAssetMap); err != nil {
				fmt.Printf("Skipped %s because %s\n", path, err)
			}
		case ".ootrs":
			if err := ootrs.RepackArchive(path, zipWriter, ootrsAssetMap, ootrsTranslationMap); err != nil {
				fmt.Printf("Skipped %s because %s\n", path, err)
			}
		default:
			// do nothing
		}
	}
	return nil
}
