//go:build wasm

package shippacker

import (
	"encoding/binary"
	"fmt"

	"github.com/frogssoldseparately/shippacker/pkg/globals"
	"github.com/frogssoldseparately/shippacker/pkg/maps"
	"github.com/frogssoldseparately/shippacker/pkg/mmrs"
	"github.com/frogssoldseparately/simpleseek/swriter"
)

func Pack(srcPaths []string) []byte {
	endianness := binary.LittleEndian
	modWriter := swriter.NewEmptySimpleWriter(endianness)
	localW := swriter.NewEmptySimpleWriter(endianness)
	centralW := swriter.NewEmptySimpleWriter(endianness)
	assetMap, err := maps.NewAssetMap()
	swriter.TakeTimestamp()
	filesWritten, _, err := WriteModEntries(srcPaths, localW, centralW, assetMap)
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

func WriteModEntries(srcPaths []string, lw *swriter.SimpleWriter, cw *swriter.SimpleWriter, assetMap *maps.Assets) (uint16, uint64, error) {
	bankId := globals.StartingBankIndex
	includedFileCount := uint16(0)
	for _, path := range srcPaths {
		if newCount, err := mmrs.RepackArchive(path, lw, cw, assetMap, bankId); err != nil {
			fmt.Printf("Skipped %s because %s\n", path, err)
		} else {
			includedFileCount += newCount
			if newCount == 2 {
				bankId++
			}
		}
	}
	return includedFileCount, bankId - globals.StartingBankIndex, nil
}
