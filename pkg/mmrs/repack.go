package mmrs

import (
	"archive/zip"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"strconv"
	"strings"

	"github.com/frogssoldseparately/shippacker/pkg/maps"
	"github.com/frogssoldseparately/shippacker/pkg/seq"
	"github.com/frogssoldseparately/shippacker/pkg/seqcat"
	"github.com/frogssoldseparately/shippacker/pkg/soundfont"
	"github.com/frogssoldseparately/shippacker/pkg/swriter"
)

func RepackArchive(paths maps.Paths, file os.DirEntry, lw *swriter.SimpleWriter, cw *swriter.SimpleWriter, am *maps.Assets, bankId uint32) (uint16, error) {
	archiveFilename := file.Name()
	archiveExtension := filepath.Ext(archiveFilename)
	archiveBasename := archiveFilename[0 : len(archiveFilename)-len(archiveExtension)]
	sequenceSuffix := "bgm"
	archive, err := zip.OpenReader(filepath.Join(paths.MSrc, archiveFilename))
	if err != nil {
		return 0, err
	}
	defer archive.Close()

	entries := map[string]*zip.File{}

	for _, file := range archive.File {
		ext := filepath.Ext(file.Name)
		if strings.Contains(ext, "seq") {
			ext = ".seq"
		} else if strings.Contains(ext, "txt") {
			if file.Name != "categories.txt" {
				ext = ""
			}
		}
		entries[ext] = file
	}
	if _, ok := entries[".zsound"]; ok {
		return 0, fmt.Errorf("%s relies on custom instruments. Skipping.\n", archiveFilename)
	}
	seqEntry, ok := entries[".seq"]
	if !ok {
		return 0, fmt.Errorf("%s did not have a valid sequence file. Skipping.\n", archiveFilename)
	}
	isFanfare := false
	if catEntry, ok := entries[".txt"]; ok {
		if categories := getCategoriesFromArchive(catEntry); categories != nil {
			isFanfare = seqcat.HasFanfareCategories(*categories)
			sequenceSuffix = strings.Join(*categories, "-")
		}
	}

	bufferedLW := swriter.NewEmptySimpleWriter(binary.LittleEndian)
	bufferedCW := swriter.NewEmptySimpleWriter(binary.LittleEndian)
	filesWritten := uint16(0)
	if bankEntry, ok := entries[".zbank"]; ok {
		metaEntry, ok := entries[".bankmeta"]
		if !ok {
			return 0, fmt.Errorf("%s is missing a .bankmeta file. Skipping\n", archiveFilename)
		}
		fBank, err := bankEntry.Open()
		if err != nil {
			return 0, fmt.Errorf("Couldn't open %s's .zbank file. Skipping\n", archiveFilename)
		}
		fMeta, err := metaEntry.Open()
		sf, err := soundfont.NewSoundfontFromBankStreams(fBank, fMeta, fmt.Sprintf("Soundfont_%d", bankId), am)
		// Should this always be 1?
		if isFanfare {
			sf.Meta.CachePolicy = int8(0x1)
		} else {
			sf.Meta.CachePolicy = int8(0x2)
		}
		if err != nil {
			return 0, fmt.Errorf("Couldn't generate %s's soundfont. Skipping\n", archiveFilename)
		}
		if err := swriter.WriteZipEntry(sf, bufferedLW, bufferedCW, lw.GetLength()); err != nil {
			return 0, err
		}
		filesWritten++
	} else {
		seqFilename := seqEntry.Name
		seqExt := filepath.Ext(seqFilename)
		seqBasename := seqFilename[0 : len(seqFilename)-len(seqExt)]
		newBank, err := strconv.ParseUint(seqBasename, 16, 16)
		bankId = uint32(newBank)
		if err != nil {
			return 0, fmt.Errorf("Could not parse bank information")
		}
	}
	fSeq, err := seqEntry.Open()
	if err != nil {
		return 0, err
	}
	sequenceName := strings.ReplaceAll(archiveBasename, "_", " ")
	sequenceName += "_" + sequenceSuffix
	banks := []byte{uint8(bankId)}
	seq, err := seq.NewSequenceFromStream(fSeq, sequenceName, &banks)
	if err := swriter.WriteZipEntry(seq, bufferedLW, bufferedCW, lw.GetLength()); err != nil {
		return 0, err
	}
	filesWritten++
	lw.CopyFrom(bufferedLW)
	cw.CopyFrom(bufferedCW)
	return filesWritten, nil
}

func getCategoriesFromArchive(src *zip.File) *[]string {
	fSrc, err := src.Open()
	if err != nil {
		return nil
	}
	defer fSrc.Close()
	catText := readIntoString(fSrc, src.FileInfo().Size())
	return seqcat.GetCategoriesFromString(catText)
}

func readIntoString(r io.Reader, len int64) string {
	buf := make([]byte, len)
	r.Read(buf)
	return string(buf[:])
}
