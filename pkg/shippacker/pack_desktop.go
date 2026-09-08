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
	"github.com/frogssoldseparately/shippacker/pkg/ootrs"
	"github.com/frogssoldseparately/shippacker/pkg/seq"
	"github.com/frogssoldseparately/simpleseek/swriter"
)

func Pack(musicSrcPath string, outPath string) error {
	zipWriter := swriter.NewEmptyZipWriter(binary.LittleEndian)
	sampleMap, err := maps.NewSampleMap()
	if err != nil {
		return err
	}
	if globals.HasOotO2r {
		if err := ootrs.InjectOotSamples(zipWriter, sampleMap.OcarinaOfTime); err != nil {
			fmt.Printf("Could not inject oot samples because %s\n", err)
		}
	}
	if err := WriteModEntries(musicSrcPath, zipWriter, sampleMap); err != nil {
		return err
	}
	sequencesWritten := zipWriter.GetTypedFileCount("Sequence")
	soundfontsWritten := zipWriter.GetTypedFileCount("Soundfont")
	samplesWritten := zipWriter.GetTypedFileCount("Sample")
	if sequencesWritten != 0 {
		fO2R, err := os.Create(filepath.Join(outPath, generateModFilename()))
		if err != nil {
			return err
		}
		defer fO2R.Close()
		o2rWriter := zipWriter.Finish()
		if _, err = fO2R.Write(*o2rWriter.GetBuffer()); err != nil {
			return err
		}
	} else {
		soundfontsWritten = 0
		samplesWritten = 0
	}
	fmt.Printf("Wrote %d sequences", sequencesWritten)
	if sequencesWritten != 0 {
		if soundfontsWritten != 0 {
			fmt.Print(", ")
			if samplesWritten == 0 {
				fmt.Printf("and ")
			}
			fmt.Printf("%d soundfonts", soundfontsWritten)
			if samplesWritten != 0 {
				fmt.Printf(", and %d samples", samplesWritten)
			}
		}
	}
	fmt.Print("\n")
	return nil
}

func WriteModEntries(srcPath string, zipWriter *swriter.SimpleZipWriter, sampleMap *maps.SampleMap) error {
	files, err := os.ReadDir(srcPath)
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.IsDir() {
			if globals.RecurseSubdirectories {
				if err := WriteModEntries(filepath.Join(srcPath, file.Name()), zipWriter, sampleMap); err != nil {
					return err
				}
				if globals.EarlyExit {
					return nil
				}
			}
		} else {
			name := filepath.Base(file.Name())
			ext := filepath.Ext(name)
			if ext == ".mmrs" {
				if err := mmrs.RepackArchive(srcPath, file, zipWriter, sampleMap.MajorasMask); err != nil {
					fmt.Printf("Skipped %s because %s\n", name, err)
				}
			} else if ext == ".ootrs" {
				if err := ootrs.RepackArchive(srcPath, file, zipWriter, sampleMap.OcarinaOfTime); err != nil {
					fmt.Printf("Skipped %s because %s\n", name, err)
				}
			} else if isSequenceExtension(ext) {
				if err := seq.RepackSequence(srcPath, file, zipWriter); err != nil {
					fmt.Printf("Skipped %s because %s\n", name, err)
				}
			}
		}
		if globals.WarnOnTooManyBanks && globals.GetCurrentBank(zipWriter) == globals.MaxBankCount {
			globals.WarnOnTooManyBanks = false
			switch iohelper.WarnPromptBanks() {
			case iohelper.EarlyExit:
				globals.EarlyExit = true
				return nil
			case iohelper.IgnoreOtherBanks:
				globals.AllowCustomBanks = false
			case iohelper.ContinueRunning:
				// do nothing
			case iohelper.HaltRunning:
				fallthrough
			default:
				return fmt.Errorf("Halt from excessive banks. No o2r file written")
			}
		}
		if globals.WarnOnTooManySongs && zipWriter.GetTypedFileCount("Sequence") == globals.MaxSongCount {
			globals.WarnOnTooManySongs = false
			switch iohelper.WarnPromptSongs() {
			case iohelper.ContinueRunning:
				// do nothing
			case iohelper.EarlyExit:
				globals.EarlyExit = true
				return nil
			default:
				return fmt.Errorf("Halt from excessive songs. No o2r file written")
			}
		}
	}
	return nil
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
