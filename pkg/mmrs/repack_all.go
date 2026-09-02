package mmrs

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
	"github.com/frogssoldseparately/shippacker/pkg/sample"
	"github.com/frogssoldseparately/shippacker/pkg/seq"
	"github.com/frogssoldseparately/shippacker/pkg/seqcat"
	"github.com/frogssoldseparately/shippacker/pkg/soundfont"
	"github.com/frogssoldseparately/simpleseek/sreader"
	"github.com/frogssoldseparately/simpleseek/swriter"
)

func RepackArchiveFromZipReader(archive *sreader.SimpleZipReader, lw *swriter.SimpleWriter, cw *swriter.SimpleWriter, am *maps.Assets, bankId uint64) (uint16, error) {
	archiveFilename := filepath.Base(archive.Name())
	archiveExtension := filepath.Ext(archiveFilename)
	archiveBasename := archiveFilename[0 : len(archiveFilename)-len(archiveExtension)]
	sequenceSuffix := "bgm"
	fontCount := uint32(1)
	seqEntry, ok := archive.GetFirstByAnyExt([]string{".seq", ".zseq", ".aseq"})
	if !ok {
		return 0, fmt.Errorf("it did not have a valid sequence file.\n")
	}
	var seqName string
	{
		seqFull := seqEntry.Name
		seqExt := filepath.Ext(seqFull)
		seqName = seqFull[0 : len(seqFull)-len(seqExt)]
	}
	isFanfare := false
	if catEntry, ok := archive.GetFile("categories.txt"); ok {
		if categories := getCategoriesFromArchive(catEntry); categories != nil {
			isFanfare = seqcat.HasFanfareCategories(*categories)
			if globals.UseNumericCategories {
				sequenceSuffix = strings.Join(*categories, "-")
			} else if isFanfare {
				sequenceSuffix = "fanfare"
			}
		}
	}

	bufferedLW := swriter.NewEmptySimpleWriter(binary.LittleEndian)
	bufferedCW := swriter.NewEmptySimpleWriter(binary.LittleEndian)
	filesWritten := uint16(0)
	if bankEntry, ok := archive.GetFile(seqName + ".zbank"); ok {
		if !globals.AllowCustomBanks {
			return 0, fmt.Errorf("it has a custom bank\n")
		}
		metaEntry, ok := archive.GetFile(seqName + ".bankmeta")
		if !ok {
			return 0, fmt.Errorf("it is missing a .bankmeta file\n")
		}
		fBank, err := bankEntry.Open()
		if err != nil {
			return 0, fmt.Errorf("its .zbank file could not be opened\n")
		}
		fMeta, err := metaEntry.Open()
		// Determine fontname and index by global settings
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
		// Get custom samples
		customSamples := []*sample.Sample{}
		if zsoundEntries, ok := archive.GetAllByExt(".zsound"); ok {
			for _, zsoundEntry := range zsoundEntries {
				instName := zsoundEntry.Name
				baseName := instName[0:strings.LastIndex(instName, "_")]
				sampleName := fmt.Sprintf("%s_%s_META", baseName, stamp)
				addrHex := instName[len(baseName)+1 : strings.LastIndex(instName, ".")]
				addr, err := strconv.ParseUint(addrHex, 16, 32)
				if err != nil {
					return 0, err
				}
				fInst, err := zsoundEntry.Open()
				if err != nil {
					return 0, err
				}
				customSample, err := sample.NewSampleFromStream(fInst, uint32(addr), sampleName, am)
				if err != nil {
					return 0, err
				}
				customSamples = append(customSamples, customSample)
				// So this sample can be referenced in .zbank files
				(*am)[customSample.Addr] = customSample.Name
			}
		}
		// Generate zippable soundfont container
		sf, err := soundfont.NewSoundfontFromBankStreams(fBank, fMeta, fontName, am)
		if err != nil {
			return 0, fmt.Errorf("its soundfont could not be generated\n")
		}
		// Should this always be 1?
		if isFanfare {
			sf.Meta.CachePolicy = int8(0x1)
		} else {
			sf.Meta.CachePolicy = int8(0x2)
		}
		// Write zippable custom samples
		for _, customSample := range customSamples {
			loopPtr, ok := (*sf.LoopMap)[customSample.Addr]
			if !ok {
				return 0, fmt.Errorf("could not find AdpcmLoop for instrument")
			}
			customSample.Loop = loopPtr
			bookPtr, ok := (*sf.BookMap)[customSample.Addr]
			if !ok {
				return 0, fmt.Errorf("could not find AdpcmBook for instrument")
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
		// No custom instrument bank
		seqFilename := seqEntry.Name
		seqExt := filepath.Ext(seqFilename)
		seqBasename := seqFilename[0 : len(seqFilename)-len(seqExt)]
		newBank, err := strconv.ParseUint(seqBasename, 16, 16)
		bankId = uint64(newBank)
		if err != nil {
			return 0, fmt.Errorf("its bank could not be parsed")
		}
	}
	fSeq, err := seqEntry.Open()
	if err != nil {
		return 0, err
	}
	sequenceName := strings.ReplaceAll(archiveBasename, "_", " ")
	sequenceName += "_" + sequenceSuffix
	banks := makeFontIdArray(bankId, fontCount)
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

func makeFontIdArray(id uint64, len uint32) *[]byte {
	// clamp array length between 1 and 8
	out := make([]byte, min32(max32(len, 1), 8))
	for i := range out {
		out[i] = byte(id & 0xFF)
		id = id >> 8
	}
	return &out
}

// No standard library min/max for integers
func max32(a uint32, b uint32) uint32 {
	if a > b {
		return a
	}
	return b
}

// No standard library min/max for integers
func min32(a uint32, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}

// Parse contents of an .mmrs file's categories.txt file.
func getCategoriesFromArchive(src *zip.File) *[]string {
	fSrc, err := src.Open()
	if err != nil {
		return nil
	}
	defer fSrc.Close()
	catText := readIntoString(fSrc, src.FileInfo().Size())
	// prevent newlines in category.txt from messing with suffixes
	catText = strings.Replace(catText, "\r", "", 1)
	newlineIndex := strings.Index(catText, "\n")
	if newlineIndex != -1 {
		catText = catText[0:newlineIndex]
	}
	return seqcat.GetCategoriesFromString(catText)
}

func readIntoString(r io.Reader, len int64) string {
	buf := make([]byte, len)
	r.Read(buf)
	return string(buf[:])
}
