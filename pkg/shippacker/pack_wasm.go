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
	endianness := binary.LittleEndian
	modWriter := swriter.NewEmptySimpleWriter(endianness)
	localW := swriter.NewEmptySimpleWriter(endianness)
	centralW := swriter.NewEmptySimpleWriter(endianness)
	mmrsAssetMap, err := maps.NewAssetMap()
	if err != nil {
		return nil
	}
	ootrsAssetMap, ootrsTranslationMap, err := maps.NewTranslationMaps()
	swriter.TakeTimestamp()
	filesWritten, _, _, err := WriteModEntries(srcPaths, localW, centralW, mmrsAssetMap, ootrsAssetMap, ootrsTranslationMap)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if filesWritten != 0 {
		modWriter.CopyFrom(localW)
		modWriter.CopyFrom(centralW)
		swriter.WriteCentralDirectoryEndRecord(modWriter, filesWritten, centralW.GetLength(), localW.GetLength())
		return *modWriter.GetBuffer()
	}
	fmt.Println("Nothing to write")
	return nil
}

func WriteModEntries(srcPaths []string, lw *swriter.SimpleWriter, cw *swriter.SimpleWriter, mmrsAssetMap *maps.AssetMap, ootrsAssetMap *maps.AssetMap, ootrsTranslationMap *maps.TranslationMap) (uint16, uint16, uint64, error) {
	bankId := globals.StartingBankIndex
	includedFileCount := uint16(0)
	songCount := uint16(0)
	for _, path := range srcPaths {
		switch filepath.Ext(path) {
		case ".mmrs":
			if newCount, err := mmrs.RepackArchive(path, lw, cw, mmrsAssetMap, bankId); err != nil {
				fmt.Printf("Skipped %s because %s\n", path, err)
			} else {
				includedFileCount += newCount
				songCount++
				if newCount >= 2 {
					bankId++
				}
			}
		case ".ootrs":
			if newCount, err := ootrs.RepackArchive(path, lw, cw, ootrsAssetMap, ootrsTranslationMap, bankId); err != nil {
				fmt.Printf("Skipped %s because %s\n", path, err)
			} else {
				includedFileCount += newCount
				songCount++
				if newCount >= 2 {
					bankId++
				}
			}
		default:
			// do nothing
		}
	}
	return includedFileCount, songCount, bankId - globals.StartingBankIndex, nil
}
