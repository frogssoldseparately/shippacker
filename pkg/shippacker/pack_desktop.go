//go:build !wasm

package shippacker

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/frogssoldseparately/shippacker/pkg/globals"
	"github.com/frogssoldseparately/shippacker/pkg/iohelper"
	"github.com/frogssoldseparately/shippacker/pkg/maps"
	"github.com/frogssoldseparately/shippacker/pkg/mmrs"
	"github.com/frogssoldseparately/shippacker/pkg/seq"
	"github.com/frogssoldseparately/simpleseek/swriter"
)

func Pack(musicSrcPath string, outPath string) error {
	endianness := binary.LittleEndian
	modWriter := swriter.NewEmptySimpleWriter(endianness)
	localW := swriter.NewEmptySimpleWriter(endianness)
	centralW := swriter.NewEmptySimpleWriter(endianness)
	assetMap, err := maps.NewAssetMap()
	swriter.TakeTimestamp()
	filesWritten, banksWritten, err := WriteModEntries(musicSrcPath, localW, centralW, assetMap, 0, globals.StartingBankIndex)
	if err != nil {
		return err
	}
	if filesWritten != 0 {
		modWriter.CopyFrom(localW)
		modWriter.CopyFrom(centralW)
		swriter.WriteCentralDirectoryEndRecord(modWriter, filesWritten, centralW.GetLength(), localW.GetLength())
		modFile, err := os.Create(filepath.Join(outPath, generateModFilename()))
		if err != nil {
			return err
		}
		defer modFile.Close()
		if _, err = modFile.Write(*modWriter.GetBuffer()); err != nil {
			return err
		}
	}
	fmt.Printf("Wrote %d songs and %d banks\n", filesWritten-uint16(banksWritten), banksWritten)
	return nil
}

func WriteModEntries(srcPath string, lw *swriter.SimpleWriter, cw *swriter.SimpleWriter, assetMap *maps.Assets, startFileCount uint16, startBankId uint64) (uint16, uint64, error) {
	files, err := os.ReadDir(srcPath)
	if err != nil {
		return 0, 0, err
	}
	bankId := startBankId
	includedFileCount := startFileCount
	for _, file := range files {
		if file.IsDir() {
			if globals.RecurseSubdirectories {
				newCount, newBanks, err := WriteModEntries(filepath.Join(srcPath, file.Name()), lw, cw, assetMap, includedFileCount, bankId)
				if err != nil {
					return 0, 0, err
				}
				includedFileCount += newCount
				bankId += newBanks
				if globals.EarlyExit {
					return includedFileCount - startFileCount, bankId - startBankId, nil
				}
			}
		} else {
			name := filepath.Base(file.Name())
			ext := filepath.Ext(name)
			if ext == ".mmrs" {
				if newCount, err := mmrs.RepackArchive(srcPath, file, lw, cw, assetMap, bankId); err != nil {
					fmt.Printf("Skipped %s because %s\n", name, err)
				} else {
					includedFileCount += newCount
					if newCount == 2 {
						bankId++
					}
				}
			} else if isSequenceExtension(ext) {
				if newCount, err := seq.RepackSequence(srcPath, file, lw, cw); err != nil {
					fmt.Printf("Skipped %s because %s\n", name, err)
				} else {
					includedFileCount += newCount
				}
			}
		}
		if globals.WarnOnTooManyBanks && bankId == globals.MaxBankCount {
			globals.WarnOnTooManyBanks = false
			switch iohelper.WarnPromptBanks() {
			case iohelper.EarlyExit:
				globals.EarlyExit = true
				return includedFileCount - startFileCount, bankId - startBankId, nil
			case iohelper.IgnoreOtherBanks:
				globals.AllowCustomBanks = false
			case iohelper.ContinueRunning:
				// do nothing
			case iohelper.HaltRunning:
				fallthrough
			default:
				return 0, 0, fmt.Errorf("Halt from excessive banks. No o2r file written")
			}
		}
		if globals.WarnOnTooManySongs && includedFileCount-(uint16(bankId-globals.StartingBankIndex)) == globals.MaxSongCount {
			globals.WarnOnTooManySongs = false
			switch iohelper.WarnPromptSongs() {
			case iohelper.ContinueRunning:
				// do nothing
			case iohelper.EarlyExit:
				globals.EarlyExit = true
				return includedFileCount - startFileCount, bankId - startBankId, nil
			default:
				return 0, 0, fmt.Errorf("Halt from excessive songs. No o2r file written")
			}
		}
	}
	return includedFileCount - startFileCount, bankId - startBankId, nil
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
