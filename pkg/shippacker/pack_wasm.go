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
	sampleMap, err := maps.NewSampleMap()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if globals.HasOotO2r {
		if err := ootrs.InjectOotSamples(zipWriter, sampleMap.OcarinaOfTime); err != nil {
			fmt.Printf("Could not inject oot samples because %s\n", err)
		}
	}
	if err := WriteModEntries(srcPaths, zipWriter, sampleMap); err != nil {
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

func WriteModEntries(srcPaths []string, zipWriter *swriter.SimpleZipWriter, sampleMap *maps.SampleMap) error {
	for _, path := range srcPaths {
		switch filepath.Ext(path) {
		case ".mmrs":
			if err := mmrs.RepackArchive(path, zipWriter, sampleMap.MajorasMask); err != nil {
				fmt.Printf("Skipped %s because %s\n", path, err)
			}
		case ".ootrs":
			if err := ootrs.RepackArchive(path, zipWriter, sampleMap.OcarinaOfTime); err != nil {
				fmt.Printf("Skipped %s because %s\n", path, err)
			}
		default:
			// do nothing
		}
	}
	return nil
}
