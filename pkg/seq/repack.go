package seq

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/frogssoldseparately/shippacker/pkg/maps"
	"github.com/frogssoldseparately/shippacker/pkg/seqcat"
	"github.com/frogssoldseparately/simpleseek/swriter"
)

func RepackSequence(paths maps.Paths, file os.DirEntry, lw *swriter.SimpleWriter, cw *swriter.SimpleWriter) (uint16, error) {
	filename := filepath.Base(file.Name())
	ext := filepath.Ext(filename)
	basename := filename[0 : len(filename)-len(ext)]
	seqPath := filepath.Join(paths.MSrc, filename)
	metaPath := filepath.Join(paths.MSrc, basename+".meta")
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
	categories := seqcat.GetCategoriesFromString(nameParts[len(nameParts)-1])
	bankHex := nameParts[len(nameParts)-2]
	moddedName := strings.Join(nameParts[0:len(nameParts)-2], " ") + "_" + strings.Join(*categories, "-")
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
