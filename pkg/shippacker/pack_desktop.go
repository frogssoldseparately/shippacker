//go:build !wasm

package shippacker

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"slices"
	"strings"
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
	songList, err := FindSongs(musicSrcPath)
	if err != nil {
		return err
	}
	rand.Shuffle(len(songList), func(i, j int) {
		songList[i], songList[j] = songList[j], songList[i]
	})
	if err := WriteModEntries(songList, zipWriter, sampleMap); err != nil {
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

type CustomSong struct {
	Path string
	File os.DirEntry
}

func FindSongs(root string) ([]CustomSong, error) {
	out := []CustomSong{}
	files, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		ext := filepath.Ext(file.Name())
		if file.IsDir() {
			if globals.RecurseSubdirectories && !strings.HasPrefix(file.Name(), "_") {
				newEntries, err := FindSongs(filepath.Join(root, file.Name()))
				if err != nil {
					return nil, err
				}
				out = slices.Concat(out, newEntries)
			}
		} else if ext == ".mmrs" || ext == ".ootrs" || isSequenceExtension(ext) {
			out = append(out, CustomSong{
				Path: filepath.Join(root, file.Name()),
				File: file,
			})
		}
	}
	return out, nil
}

func WriteModEntries(customSongs []CustomSong, zipWriter *swriter.SimpleZipWriter, sampleMap *maps.SampleMap) error {
	for _, customSong := range customSongs {
		file := customSong.File
		path := customSong.Path
		dir := filepath.Dir(path)
		name := filepath.Base(path)
		ext := filepath.Ext(name)
		if ext == ".mmrs" {
			if err := mmrs.RepackArchive(dir, file, zipWriter, sampleMap.MajorasMask); err != nil {
				fmt.Printf("\tSkipped \"%s\"\n\tbecause %s\n", name, err)
			}
		} else if ext == ".ootrs" {
			if err := ootrs.RepackArchive(dir, file, zipWriter, sampleMap.OcarinaOfTime); err != nil {
				fmt.Printf("\tSkipped \"%s\"\n\tbecause %s\n", name, err)
			}
		} else if isSequenceExtension(ext) {
			if err := seq.RepackSequence(dir, file, zipWriter); err != nil {
				fmt.Printf("\tSkipped \"%s\"\n\tbecause %s\n", name, err)
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
