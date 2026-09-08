package ootrs

import (
	"archive/zip"
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

func RepackArchiveFromZipReader(archive *sreader.SimpleZipReader, zipWriter *swriter.SimpleZipWriter, am *maps.AssetMap, tm *maps.TranslationMap) error {
	bankId := globals.GetCurrentBank(zipWriter)
	bufferedWriter := zipWriter.NewBuffer()
	fontCount := uint32(1)
	metaEntry, ok := archive.GetFirstByExt(".meta")
	if !ok {
		return fmt.Errorf("it did not have a valid meta file.\n")
	}
	metadata, err := processOotrsMeta(metaEntry)
	if err != nil {
		return err
	}
	sequenceSuffix := metadata.Type
	isFanfare := sequenceSuffix == "fanfare"
	// TODO: convert ootrs categories to mmrs categories
	// if globals.UseNumericCategories {
	// 	sequenceSuffix = strings.Join(mapOotrsCategories(metadata.Categories), "-")
	// }
	seqEntry, ok := archive.GetFirstByAnyExt([]string{".seq", ".zseq", ".aseq"})
	if !ok {
		return fmt.Errorf("it did not have a valid sequence file.\n")
	}
	var fontName string
	stamp := fmt.Sprintf("%d%d%d", os.Getpid(), zipWriter.GetTimestamp(), bankId)
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

	if bankEntry, ok := archive.GetFirstByExt(".zbank"); ok {
		if !globals.AllowCustomBanks {
			return fmt.Errorf("it has a custom bank\n")
		}
		var bankName string
		{
			bankBase := bankEntry.Name
			bankExt := ".zbank"
			bankName = bankBase[0 : len(bankBase)-len(bankExt)]
		}
		bankmetaEntry, ok := archive.GetFile(bankName + ".bankmeta")
		if !ok {
			return fmt.Errorf("it is missing a .bankmeta file\n")
		}
		fBank, err := bankEntry.Open()
		if err != nil {
			return fmt.Errorf("its .zbank file could not be opened\n")
		}
		fBankmeta, err := bankmetaEntry.Open()
		if err != nil {
			return fmt.Errorf("its .bankmeta file could not be opened\n")
		}
		// Get custom samples
		customSamples := []*sample.Sample{}
		for _, zsoundInfo := range metadata.CustomSamples {
			parts := strings.Split(zsoundInfo, ":")
			if len(parts) < 3 {
				return fmt.Errorf("bad custom sample entry\n")
			}
			sourceName := parts[1]
			sourceExt := filepath.Ext(sourceName)
			baseName := sourceName[0 : len(sourceName)-len(sourceExt)]
			sampleName := fmt.Sprintf("%s_%s_META", baseName, stamp)
			addrHex := parts[2]
			addr, err := strconv.ParseUint(addrHex, 16, 32)
			if err != nil {
				return err
			}
			zsoundEntry, ok := archive.GetFile(sourceName)
			if !ok {
				return fmt.Errorf("could not find %s in archive\n", sourceName)
			}
			fSample, err := zsoundEntry.Open()
			if err != nil {
				return err
			}
			customSample, err := sample.NewSampleFromStream(fSample, uint32(addr), sampleName, am)
			if err != nil {
				return err
			}
			customSamples = append(customSamples, customSample)
			(*am)[customSample.Addr] = customSample.Name
		}
		// Generate zippable soundfont container
		sf, err := soundfont.NewSoundfontFromBankStreams(fBank, fBankmeta, fontName, am)
		if err != nil {
			return fmt.Errorf("its soundfont could not be generated\n")
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
				return fmt.Errorf("could not find AdpcmLoop for sample\n")
			}
			customSample.Loop = loopPtr
			bookPtr, ok := (*sf.BookMap)[customSample.Addr]
			if !ok {
				return fmt.Errorf("could not find AdpcmBook for custom sample\n")
			}
			customSample.Book = bookPtr
			if err := bufferedWriter.WriteEntry(customSample); err != nil {
				return err
			}
		}
		// Write zippable soundfont
		if err := bufferedWriter.WriteEntry(sf); err != nil {
			return err
		}
	} else {
		if !globals.HasOotO2r {
			return fmt.Errorf("oot.o2r was not provided\n")
		}
		var parsedBank uint64
		if len(metadata.Bank) >= 2 && metadata.Bank[0:2] == "0x" {
			parsedBank, err = strconv.ParseUint(metadata.Bank[2:], 16, 32)
			if err != nil {
				return err
			}
		} else {
			parsedBank, err = strconv.ParseUint(metadata.Bank, 16, 32)
			if err != nil {
				return err
			}
		}
		if usedBankId, ok := includedBanks[parsedBank]; ok {
			bankId = usedBankId
		} else {
			if parsedBank < 3 {
				return fmt.Errorf("bank %d is skipped\n", parsedBank)
			}
			soundfontEntry, ok := ootSoundFonts[parsedBank]
			if !ok {
				return fmt.Errorf("could not find OoT bank with id %s\n", metadata.Bank)
			}
			fSoundfont, err := soundfontEntry.Open()
			if err != nil {
				return err
			}
			sf, err := soundfont.ReadSoundfont(fSoundfont, fontName, am, tm)
			if err != nil {
				return err
			}
			if err := bufferedWriter.WriteEntry(sf); err != nil {
				return err
			}
			includedBanks[parsedBank] = bankId
		}
	}
	fSeq, err := seqEntry.Open()
	if err != nil {
		return err
	}
	// TODO: convert ootrs categories to mmrs categories
	sequenceName := strings.ReplaceAll(metadata.Name, "/", "-") + "_" + sequenceSuffix
	banks := mmrs.MakeFontIdArray(bankId, fontCount)
	seq, err := seq.NewSequenceFromStream(fSeq, sequenceName, banks)
	seq.NumFonts = fontCount
	if err := bufferedWriter.WriteEntry(seq); err != nil {
		return err
	}
	zipWriter.ConsumeBuffer()
	return nil
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
	zsoundLines := []string{}
	for _, rawLine := range rawLines {
		line := strings.Trim(rawLine, " \r\t")
		if len(line) > 0 {
			if strings.Contains(line, "ZSOUND") {
				zsoundLines = append(zsoundLines, line)
			} else if len(lines) == 4 {
				lines[3] = line
			} else {
				lines = append(lines, line)
			}
		}
	}
	if len(lines) == 2 {
		lines = append(lines, "bgm")
	}
	if len(lines) == 3 {
		lines = append(lines, "bgm")
	}
	return &OotrsMeta{
		lines[0],
		lines[1],
		lines[2],
		strings.FieldsFunc(lines[3], func(r rune) bool { return r == ',' || r == '-' }),
		zsoundLines,
	}, nil
}

// TODO: convert ootrs categories to mmrs categories
// func mapOotrsCategories(original []string) []string {
// 	return original
// }

var includedBanks = map[uint64]uint64{}

var ootSoundFonts = map[uint64]*zip.File{}

var ootSamples = []*sample.StubbedSample{}

func PrepareOotSoundfonts(soundfontEntries *[]*zip.File) error {
	for _, soundfontEntry := range *soundfontEntries {
		base := filepath.Base(soundfontEntry.Name)
		bankDec := base[0:strings.Index(base, "_")]
		if bankNum, err := strconv.ParseUint(bankDec, 10, 32); err != nil {
			return err
		} else {
			ootSoundFonts[bankNum] = soundfontEntry
		}
	}
	return nil
}

func PrepareOotSamples(sampleEntries *map[uint32]*zip.File) error {
	for addr, sampleEntry := range *sampleEntries {
		sampleName := filepath.Base(sampleEntry.Name)
		fSample, err := sampleEntry.Open()
		if err != nil {
			return err
		}
		sample, err := sample.ReadShipSample(fSample, addr, sampleName)
		if err != nil {
			return err
		}
		ootSamples = append(ootSamples, sample)
	}
	return nil
}

func InjectOotSamples(zipWriter *swriter.SimpleZipWriter, am *maps.AssetMap, tm *maps.TranslationMap) error {
	buffered := zipWriter.NewBuffer()
	for _, entry := range ootSamples {
		if err := buffered.WriteEntry(entry); err != nil {
			return err
		}
	}
	// Second for loop so the asset maps don't have pointers to samples that weren't injected
	for _, entry := range ootSamples {
		(*am)[entry.Addr] = entry.Name
		(*tm)[entry.Name] = entry.Addr
	}
	zipWriter.ConsumeBuffer()
	return nil
}
