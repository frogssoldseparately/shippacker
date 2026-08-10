package seq

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/frogssoldseparately/shippacker/pkg/globals"
	"github.com/frogssoldseparately/shippacker/pkg/seqcat"
	"github.com/frogssoldseparately/simpleseek/swriter"
)

func RepackSequence(musicSrcPath string, file os.DirEntry, lw *swriter.SimpleWriter, cw *swriter.SimpleWriter) (uint16, error) {
	filename := filepath.Base(file.Name())
	ext := filepath.Ext(filename)
	basename := filename[0 : len(filename)-len(ext)]
	seqPath := filepath.Join(musicSrcPath, filename)
	metaPath := filepath.Join(musicSrcPath, basename+".meta")
	seq, err := NewSequenceFromFileWithMeta(seqPath, metaPath)
	if err != nil {
		return 0, err
	}
	if err := swriter.WriteZipEntry(seq, lw, cw, 0x0); err != nil {
		return 0, err
	}
	return 1, nil
}

func ExtractInformationFromPath(path string) (string, *[]byte, error) {
	filename := filepath.Base(path)
	ext := filepath.Ext(filename)
	basename := filename[0 : len(filename)-len(ext)]
	nameParts := strings.Split(basename, "_")
	if len(nameParts) < 3 {
		return "", nil, fmt.Errorf("Name lacks bank and category information")
	}
	bankHex := nameParts[len(nameParts)-2]
	catStr := nameParts[len(nameParts)-1]
	seqSuffix := "bgm"
	if globals.UseNumericCategories {
		seqSuffix = strings.Join(*seqcat.GetCategoriesFromString(catStr), "-")
	} else if seqcat.IsFanfareString(catStr) {
		seqSuffix = "fanfare"
	}
	moddedName := strings.Join(nameParts[0:len(nameParts)-2], " ") + "_" + seqSuffix
	bankIds, err := parseBankHex(bankHex)
	if err != nil {
		return "", nil, err
	}
	return moddedName, bankIds, nil
}

func parseBankHex(hex string) (*[]byte, error) {
	v, err := strconv.ParseUint(hex, 16, 16)
	if err != nil {
		return nil, err
	}
	buf := []byte{uint8(v)}
	return &buf, nil
}
