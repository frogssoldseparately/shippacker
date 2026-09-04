package ootrs

import (
	"archive/zip"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/frogssoldseparately/shippacker/pkg/crc64"
	"github.com/frogssoldseparately/shippacker/pkg/globals"
	"github.com/frogssoldseparately/shippacker/pkg/maps"
	"github.com/frogssoldseparately/shippacker/pkg/mmrs"
	"github.com/frogssoldseparately/shippacker/pkg/sample"
	"github.com/frogssoldseparately/shippacker/pkg/seq"
	"github.com/frogssoldseparately/shippacker/pkg/soundfont"
	"github.com/frogssoldseparately/simpleseek/sreader"
	"github.com/frogssoldseparately/simpleseek/swriter"
)

func RepackArchiveFromZipReader(archive *sreader.SimpleZipReader, lw *swriter.SimpleWriter, cw *swriter.SimpleWriter, am *maps.AssetMap, tm *maps.TranslationMap, bankId uint64) (uint16, error) {
	archiveFilename := filepath.Base(archive.Name())
	archiveExtension := filepath.Ext(archiveFilename)
	archiveBasename := archiveFilename[0 : len(archiveFilename)-len(archiveExtension)]
	fontCount := uint32(1)
	metaEntry, ok := archive.GetFirstByExt(".meta")
	if !ok {
		return 0, fmt.Errorf("it did not have a valid meta file.\n")
	}
	metadata, err := processOotrsMeta(metaEntry)
	if err != nil {
		return 0, err
	}
	sequenceSuffix := metadata.Type
	isFanfare := sequenceSuffix == "fanfare"
	// TODO: convert ootrs categories to mmrs categories
	// if globals.UseNumericCategories {
	// 	sequenceSuffix = strings.Join(mapOotrsCategories(metadata.Categories), "-")
	// }
	seqEntry, ok := archive.GetFirstByAnyExt([]string{".seq", ".zseq", ".aseq"})
	if !ok {
		return 0, fmt.Errorf("it did not have a valid sequence file.\n")
	}
	var fontName string
	stamp := fmt.Sprintf("%d%d%d", os.Getpid(), swriter.GetTimestamp(), bankId)
	if globals.UseCRC64Encoding {
		// Generate a hash to prevent (or minimize) collisions between soundfonts
		stampArr := []byte(stamp)
		stampHash := crc32.ChecksumIEEE(stampArr)
		fontName = fmt.Sprintf("custom/fonts/Soundfont_%d", stampHash)
		fontNameArr := []byte(fontName)
		// Makes 2ship find the correct soundfont by crc instead of index
		bankId = crc64.CRC64(&fontNameArr)
		fontCount = 0xFFFFFFFF
	} else {
		fontName = fmt.Sprintf("audio/fonts/Soundfont_%d", bankId)
	}

	bufferedLW := swriter.NewEmptySimpleWriter(binary.LittleEndian)
	bufferedCW := swriter.NewEmptySimpleWriter(binary.LittleEndian)
	filesWritten := uint16(0)
	if bankEntry, ok := archive.GetFirstByExt(".zbank"); ok {
		if !globals.AllowCustomBanks {
			return 0, fmt.Errorf("it has a custom bank\n")
		}
		var bankName string
		{
			bankBase := bankEntry.Name
			bankExt := ".zbank"
			bankName = bankBase[0 : len(bankBase)-len(bankExt)]
		}
		bankmetaEntry, ok := archive.GetFile(bankName + ".bankmeta")
		if !ok {
			return 0, fmt.Errorf("it is missing a .bankmeta file\n")
		}
		fBank, err := bankEntry.Open()
		if err != nil {
			return 0, fmt.Errorf("its .zbank file could not be opened\n")
		}
		fBankmeta, err := bankmetaEntry.Open()
		if err != nil {
			return 0, fmt.Errorf("its .bankmeta file could not be opened\n")
		}
		// Get custom samples
		customSamples := []*sample.Sample{}
		for _, zsoundInfo := range metadata.CustomSamples {
			parts := strings.Split(zsoundInfo, ":")
			if len(parts) < 3 {
				return 0, fmt.Errorf("bad custom sample entry")
			}
			sourceName := parts[1]
			sourceExt := filepath.Ext(sourceName)
			baseName := sourceName[0 : len(sourceName)-len(sourceExt)]
			sampleName := fmt.Sprintf("%s_%s_META", baseName, stamp)
			addrHex := parts[2]
			addr, err := strconv.ParseUint(addrHex, 16, 32)
			if err != nil {
				return 0, err
			}
			zsoundEntry, ok := archive.GetFile(sourceName)
			if !ok {
				return 0, fmt.Errorf("could not find %s in archive\n", sourceName)
			}
			fSample, err := zsoundEntry.Open()
			if err != nil {
				return 0, err
			}
			customSample, err := sample.NewSampleFromStream(fSample, uint32(addr), sampleName, am)
			if err != nil {
				return 0, err
			}
			customSamples = append(customSamples, customSample)
			(*am)[customSample.Addr] = customSample.Name
		}
		// Generate zippable soundfont container
		sf, err := soundfont.NewSoundfontFromBankStreams(fBank, fBankmeta, fontName, am)
		if err != nil {
			return 0, fmt.Errorf("its soundfont could not be generated\n")
		}
		if isFanfare {
			sf.Meta.CachePolicy = 0x1
		} else {
			sf.Meta.CachePolicy = 0x2
		}
		// Write zippable custom samples
		for _, customSample := range customSamples {
			loopPtr, ok := (*sf.LoopMap)[customSample.Addr]
			if !ok {
				return 0, fmt.Errorf("could not find AdpcmLoop for sample")
			}
			customSample.Loop = loopPtr
			bookPtr, ok := (*sf.BookMap)[customSample.Addr]
			if !ok {
				return 0, fmt.Errorf("could not find AdpcmBook for custom sample")
			}
			customSample.Book = bookPtr
			if err := swriter.WriteZipEntry(customSample, bufferedLW, bufferedCW, lw.GetLength()); err != nil {
				return 0, err
			}
			filesWritten++
		}
		// Write zippable soundfont
		if err := swriter.WriteZipEntry(sf, bufferedLW, bufferedCW, lw.GetLength()); err != nil {
			return 0, err
		}
		filesWritten++
	} else {
		if !globals.HasOotO2r {
			return 0, fmt.Errorf("oot.o2r was not provided\n")
		}
		var parsedBank uint64
		if len(metadata.Bank) >= 2 && metadata.Bank[0:2] == "0x" {
			parsedBank, err = strconv.ParseUint(metadata.Bank[2:], 16, 32)
			if err != nil {
				return 0, err
			}
		} else {
			parsedBank, err = strconv.ParseUint(metadata.Bank, 16, 32)
			if err != nil {
				return 0, err
			}
		}
		if usedBankId, ok := includedBanks[parsedBank]; ok {
			bankId = usedBankId
		} else {
			soundfontEntry, ok := ootSoundFonts[parsedBank]
			if !ok {
				return 0, fmt.Errorf("could not find OoT bank with id %s\n", metadata.Bank)
			}
			fSoundfont, err := soundfontEntry.Open()
			if err != nil {
				return 0, err
			}
			sf, err := soundfont.ReadSoundfont(fSoundfont, fontName, am, tm)
			if err != nil {
				return 0, err
			}
			fmt.Println(fontName)
			if err := swriter.WriteZipEntry(sf, bufferedLW, bufferedCW, lw.GetLength()); err != nil {
				return 0, err
			}
			filesWritten++

			includedBanks[parsedBank] = bankId
		}
	}
	fSeq, err := seqEntry.Open()
	if err != nil {
		return 0, err
	}
	// TODO: convert ootrs categories to mmrs categories
	// sequenceName := metadata.Name + "_" + sequenceSuffix
	sequenceName := archiveBasename + "_" + sequenceSuffix
	banks := mmrs.MakeFontIdArray(bankId, fontCount)
	seq, err := seq.NewSequenceFromStream(fSeq, sequenceName, banks)
	seq.NumFonts = fontCount
	if err := swriter.WriteZipEntry(seq, bufferedLW, bufferedCW, lw.GetLength()); err != nil {
		return 0, err
	}
	filesWritten++
	lw.CopyFrom(bufferedLW)
	cw.CopyFrom(bufferedCW)
	return filesWritten, nil
}

type OotrsMeta struct {
	Name          string
	Bank          string
	Type          string
	Categories    []string
	CustomSamples []string
}

func processOotrsMeta(f *zip.File) (*OotrsMeta, error) {
	r, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer r.Close()
	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	rawLines := strings.Split(string(buf), "\n")
	lines := []string{}
	for _, rawLine := range rawLines {
		line := strings.Trim(rawLine, " \r\t")
		if len(line) > 0 {
			lines = append(lines, line)
		}
	}
	if len(lines) == 2 {
		lines = append(lines, "bgm")
	}
	if len(lines) == 3 {
		lines = append(lines, "all")
	}
	return &OotrsMeta{
		lines[0],
		lines[1],
		lines[2],
		strings.FieldsFunc(lines[3], func(r rune) bool { return r == ',' || r == '-' }),
		lines[4:],
	}, nil
}

// TODO: convert ootrs categories to mmrs categories
// func mapOotrsCategories(original []string) []string {
// 	return original
// }

var includedBanks = map[uint64]uint64{}

var ootSoundFonts = map[uint64]*zip.File{}

func PrepareOotSoundfonts(soundfontEntries *[]*zip.File) {
	for _, soundfontEntry := range *soundfontEntries {
		base := filepath.Base(soundfontEntry.Name)
		bankDec := base[0:strings.Index(base, "_")]
		if bankNum, err := strconv.ParseUint(bankDec, 10, 32); err != nil {
			fmt.Println(err)
		} else {
			ootSoundFonts[bankNum] = soundfontEntry
		}
	}
}
