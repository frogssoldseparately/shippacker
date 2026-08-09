package shippacker

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/frogssoldseparately/shippacker/pkg/maps"
	"github.com/frogssoldseparately/shippacker/pkg/mmrs"
	"github.com/frogssoldseparately/shippacker/pkg/seq"
	"github.com/frogssoldseparately/simpleseek/swriter"
)

func Pack(paths maps.Paths) error {
	endianness := binary.LittleEndian
	modWriter := swriter.NewEmptySimpleWriter(endianness)
	localW := swriter.NewEmptySimpleWriter(endianness)
	centralW := swriter.NewEmptySimpleWriter(endianness)
	assetMap, err := maps.NewAssetMap()
	swriter.TakeTimestamp()
	filesWritten, err := WriteModEntries(paths, localW, centralW, assetMap)
	if err != nil {
		return err
	}
	modWriter.CopyFrom(localW)
	modWriter.CopyFrom(centralW)
	swriter.WriteCentralDirectoryEndRecord(modWriter, filesWritten, centralW.GetLength(), localW.GetLength())
	modFile, err := os.Create(filepath.Join(paths.OOut, generateModFilename()))
	if err != nil {
		return err
	}
	defer modFile.Close()
	if _, err = modFile.Write(*modWriter.GetBuffer()); err != nil {
		return err
	}
	return nil
}

func WriteModEntries(paths maps.Paths, lw *swriter.SimpleWriter, cw *swriter.SimpleWriter, assetMap *maps.Assets) (uint16, error) {
	bankId := uint64(41)
	files, err := os.ReadDir(paths.MSrc)
	if err != nil {
		return 0, err
	}
	includedFileCount := uint16(0)
	for _, file := range files {
		name := filepath.Base(file.Name())
		ext := filepath.Ext(name)
		if ext == ".mmrs" {
			if newCount, err := mmrs.RepackArchive(paths, file, lw, cw, assetMap, bankId); err != nil {
				fmt.Printf("Could not include %s because %s\n", name, err)
			} else {
				includedFileCount += newCount
				if newCount == 2 {
					bankId++
				}
			}
		} else if isSequenceExtension(ext) {
			if newCount, err := seq.RepackSequence(paths, file, lw, cw); err != nil {
				fmt.Printf("Could not include %s because %s\n", name, err)
			} else {
				includedFileCount += newCount
			}
		}
	}
	return includedFileCount, nil
}

func generateModFilename() string {
	stamp := time.Now().Unix()
	return fmt.Sprintf("%d.o2r", stamp)
}

// Any file that ends with .*seq will be treated as a sequence file.
func isSequenceExtension(ext string) bool {
	if len(ext) < 4 {
		return false
	}
	ending := ext[len(ext)-3:]
	return ending == "seq"
}
